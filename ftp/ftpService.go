package ftpClient

import (
	"encoding/json"
	"file-push/common"
	"file-push/redis"
	"time"

	"log"
)

func StoreFile2Ftp(ftpMessage common.FtpMessage) {
	// new a ftp connect
	conn, _ := NewFtpConn(ftpMessage)

	// change to work dir
	desDir := ftpMessage.RemoteStorePath
	err := Change2WorkDir(conn, desDir, ftpMessage.FtpBaseDir)
	if err != nil {
		// TODO send a message to mq file push error
		return
	}

	// upload file to ftp
	result, err := UploadFile(conn, ftpMessage.LocalFilePath)
	if err != nil {
		// TODO
		log.Fatalf("upload file error %s,%v", result, err)
	}
	defer Quit(conn)
}

func DoConsumeFtpPushBusiness(mqMessage []byte) bool {
	ftpMessage := common.FtpMessage{}
	if err := json.Unmarshal(mqMessage, &ftpMessage); err != nil {
		log.Printf("message is not a json string %s", string(mqMessage))
		return false
	}
	if redis.SetNxWithExp(ftpMessage.MessageId, "1", 1*time.Hour) {
		StoreFile2Ftp(ftpMessage)
		// 业务处理成功了,删除锁
		redis.Delete(ftpMessage.MessageId)
	} else {
		log.Printf("repeat commit ...")
		return false
	}
	return true

}
