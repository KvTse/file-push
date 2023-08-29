package main

import (
	"context"
	"encoding/json"
	"file-push/common"
	"file-push/global"
	"file-push/limit"
	mq "file-push/mq"
	"file-push/redis"
	"file-push/tool"
	"fmt"
	"time"
)

func main() {
	global.GVA_VP = tool.Viper()
	redis.InitRedis()
	//go generateMessage()
	go mq.ListenerSignal()
	mq.ConsumeMessageFromKafka(limit.New(2))
}
func generateMessage() {
	time.Sleep(10 * time.Second)
	// get message from mq
	message := common.FtpMessage{
		MessageId:       "messageId",
		RemoteStorePath: "/test1/test2/test3",
		LocalFilePath:   "D:\\javaTest\\tif\\H1D_OPER_CZI_L1C_20221115T095502_20221115T095557_12729_10.tiff",
		FtpUser:         "ftpuser",
		FtpPort:         "21",
		FtpPassword:     "ftpuser",
	}
	cxt := context.Background()
	jsonMessage, _ := json.Marshal(message)
	fmt.Printf("%v", global.GVA_CONFIG.RedisConfig.Addr)
	mq.SendMessage2Kafka(cxt, string(jsonMessage))
	mq.SendMessage2Kafka(cxt, string(jsonMessage))
	time.Sleep(10 * time.Second)
	message1 := common.FtpMessage{
		MessageId:       "messageId1",
		RemoteStorePath: "/test1/test2/test3",
		LocalFilePath:   "D:\\javaTest\\tif\\H1D_OPER_CZI_L1C_20221115T095502_20221115T095557_12729_10-1.tiff",
		FtpUser:         "ftpuser",
		FtpPort:         "21",
		FtpPassword:     "ftpuser",
	}
	jsonMessage1, _ := json.Marshal(message1)
	mq.SendMessage2Kafka(cxt, string(jsonMessage1))
	time.Sleep(10 * time.Second)
	message2 := common.FtpMessage{
		MessageId:       "messageId2",
		RemoteStorePath: "/test1/test2/test3",
		LocalFilePath:   "D:\\javaTest\\tif\\H1D_OPER_CZI_L1C_20221115T095502_20221115T095557_12729_10-2.tiff",
		FtpUser:         "ftpuser",
		FtpPort:         "21",
		FtpPassword:     "ftpuser",
	}
	jsonMessage2, _ := json.Marshal(message2)
	mq.SendMessage2Kafka(cxt, string(jsonMessage2))
}
