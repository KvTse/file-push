package ftpClient

import (
	"file-push/log"
	"file-push/netcall"
)

func StoreFile2Ftp(ftpMessage netcall.FtpMessage) (string, error) {
	// new a ftp connect
	log.Infof("start to connect ftp ...")
	conn, _ := NewFtpConn(ftpMessage)
	log.Infof("start to change to work dir ...")
	// change to work dir
	desDir := ftpMessage.RemoteStorePath
	err := Change2WorkDir(conn, desDir, ftpMessage.FtpBaseDir)
	if err != nil {
		// TODO send a message to mq file push error
		return "change work dir error", err
	}

	// upload file to ftp
	result, err := UploadFile(conn, ftpMessage.LocalFilePath)
	if err != nil {
		// TODO
		log.Errorf("upload file error %s,%v", result, err)
	}
	defer Quit(conn)
	return result, err
}

func DoConsumeFtpPushBusiness(ftpMessage netcall.FtpMessage) bool {

	// TODO 需要处理两个问题
	// 1.重启之后redisKey存在,导致业务无法进行
	// 2.如果两条消息写的ftp路径一致,有问题.
	//if redis.SetNxWithExp(ftpMessage.MessageId, "1", 1*time.Minute) {
	//	StoreFile2Ftp(ftpMessage)
	//	// 业务处理成功了,删除锁
	//	redis.Delete(ftpMessage.MessageId)
	//} else {
	//	log.Printf("repeat commit ...")
	//	return false
	//}
	res, _ := StoreFile2Ftp(ftpMessage)

	vo := netcall.ResponseVo{MessageId: ftpMessage.MessageId}
	if res == "" {
		vo.Code = 0
		vo.Msg = "success"
		vo.IsSuccess = true
	} else {
		vo.Code = 1
		vo.Msg = res
		vo.IsSuccess = false
	}

	cbRes, _ := netcall.FtpReqCallbackIfNecessary(vo, ftpMessage.CallbackUrl)
	return cbRes == ""

}
