package main

import (
	"context"
	"encoding/json"
	"errors"
	ftpClient "file-push/ftp"
	"file-push/global"
	"file-push/log"
	"file-push/netcall"
	"file-push/redis"
	"file-push/redis/redis_lock"
	redmq "file-push/redis/redismq"
	"file-push/tool"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	global.GVA_VP = tool.Viper()

	//redis.InitRedis()
	//go generateMessage()
	//go mq.ListenerSignal()
	//mq.ConsumeMessageFromKafka(limit.New(2))
	go startConsumer()
	runServer()
}
func startConsumer() {
	client := redis.NewClient("tcp",
		global.GVA_CONFIG.RedisConfig.Addr, global.GVA_CONFIG.RedisConfig.Password)
	// 接收到消息后的处理函数
	callbackFunc := func(ctx context.Context, msg *redis.MsgEntity) error {
		log.Infof("receive msg, msg id: %s, msg key: %s, msg val: %s", msg.MsgID, msg.Key, msg.Val)

		ftpMessage := netcall.FtpMessage{}
		if err := json.Unmarshal([]byte(msg.Val), &ftpMessage); err != nil {
			log.Errorf("message is not a json string %s", msg.Val)
			// 消息算是处理成功了, 是格式不对,在接口侧校验,理论上不应该走到这里来
			return nil
		}
		lockCxt := context.Background()
		client := redis.NewClient("tcp",
			global.GVA_CONFIG.RedisConfig.Addr, global.GVA_CONFIG.RedisConfig.Password)
		lock := redis_lock.NewRedisLock(ftpMessage.MessageId, client, redis_lock.WithExpireSeconds(30))
		// 加锁成功,正在处理任务
		if err := lock.Lock(lockCxt); err == nil {
			business := ftpClient.DoConsumeFtpPushBusiness(ftpMessage)
			if business {
				log.Infof("message %s handle success ", msg.Val)
				return nil
			}
			// 处理完成释放锁
			defer lock.Unlock(lockCxt)
		} else {
			// 回调 正在处理,请勿重复提交
			vo := netcall.FailedWithMsg(ftpMessage.MessageId, "重复任务...")
			netcall.FtpReqCallbackIfNecessary(vo, ftpMessage.CallbackUrl)
			return errors.New("repeat wrong")
		}
		return errors.New("DoConsumeFtpPushBusiness wrong")
	}

	// 自定义实现的死信队列
	demoDeadLetterMailbox := NewDemoDeadLetterMailbox(func(msg *redis.MsgEntity) {
		log.Infof("receive dead letter, msg id: %s, msg key: %s, msg val: %s", msg.MsgID, msg.Key, msg.Val)
	})

	// 构造并启动消费者
	_, err := redmq.NewConsumer(client, "myTopic", "consumerGroup", "consumerID", callbackFunc,
		// 每条消息最多重试 2 次
		redmq.WithMaxRetryLimit(2),
		// 每轮接收消息的超时时间为 2 s
		redmq.WithReceiveTimeout(2*time.Second),
		// 注入自定义实现的死信队列
		redmq.WithDeadLetterMailbox(demoDeadLetterMailbox))
	if err != nil {
		log.Errorf("", err)
		return
	}
}

// 自定义实现的死信队列
type DemoDeadLetterMailbox struct {
	do func(msg *redis.MsgEntity)
}

func NewDemoDeadLetterMailbox(do func(msg *redis.MsgEntity)) *DemoDeadLetterMailbox {
	return &DemoDeadLetterMailbox{
		do: do,
	}
}

// Deliver 死信队列接收消息的处理方法
func (d *DemoDeadLetterMailbox) Deliver(ctx context.Context, msg *redis.MsgEntity) error {
	d.do(msg)
	return nil
}

func runServer() {
	log.Infof("start server ...")
	http.HandleFunc("/", pushFileByFtp)      // 设置访问的路由
	err := http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatalf("ListenAndServe: ", err)
	}
}
func pushFileByFtp(w http.ResponseWriter, r *http.Request) {
	ftpParams, _ := io.ReadAll(r.Body)
	ftpMessage := netcall.FtpMessage{}
	if err := json.Unmarshal(ftpParams, &ftpMessage); err != nil {
		log.Errorf("message is not a json string %s", string(ftpParams))
		vo := netcall.ResponseVo{MessageId: ftpMessage.MessageId, Code: 1, IsSuccess: false, Msg: "request param is not a json string"}
		responseJson, _ := json.Marshal(vo)
		fmt.Fprintf(w, string(responseJson))
		return
	}

	client := redis.NewClient("tcp",
		global.GVA_CONFIG.RedisConfig.Addr, global.GVA_CONFIG.RedisConfig.Password)
	// 最多保留十条消息
	producer := redmq.NewProducer(client, redmq.WithMsgQueueLen(10))
	ctx := context.Background()
	msgID, err := producer.SendMsg(ctx, "myTopic", "test_kk", string(ftpParams))
	if err != nil {
		vo := netcall.ResponseVo{MessageId: ftpMessage.MessageId, Code: 500, IsSuccess: false, Msg: "server error please contact the person in charge."}
		responseJson, _ := json.Marshal(vo)
		fmt.Fprintf(w, string(responseJson))
		log.Errorf("send to redis error ", err)
		return
	}
	log.Infof(msgID)

	responseVo := netcall.Success(ftpMessage.MessageId)
	marshal, _ := json.Marshal(responseVo)
	fmt.Fprintf(w, string(marshal))
}
