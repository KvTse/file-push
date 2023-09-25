package netcall

// FtpMessage 请求消息
type FtpMessage struct {
	// 消息id
	MessageId string `json:"messageId"`
	// ftp远程存储的路径
	RemoteStorePath string `json:"remoteStorePath"`
	// 本地文件的路径
	LocalFilePath string `json:"localFilePath"`
	// ftp地址
	FtpHost string `json:"ftpHost"`
	// ftp端口号
	FtpPort string `json:"ftpPort"`
	// ftp baseDir
	FtpBaseDir string `json:"ftpBaseDir"`
	// ftp用户名
	FtpUser string `json:"ftpUser"`

	// ftp密码
	FtpPassword string `json:"ftpPassword"`
	// 字段不为空的情况下,发送会掉信息
	CallbackUrl string `json:"callbackUrl"`
}

// ResponseVo 返回体 包括接口请求返回和回调接口返回
type ResponseVo struct {
	// 消息id 业务方的唯一id键
	MessageId string `json:"messageId"`
	// 返回码 0 成功
	Code int `json:"code"`
	// 是否成功
	IsSuccess bool `json:"isSuccess"`
	// 信息 失败原因
	Msg string `json:"msg"`
}

func Success(messageId string) ResponseVo {
	vo := ResponseVo{MessageId: messageId, Code: 0, IsSuccess: true, Msg: "success"}
	return vo
}
func SuccessWithMsg(messageId string, msg string) ResponseVo {
	vo := ResponseVo{MessageId: messageId, Code: 0, IsSuccess: true, Msg: msg}
	return vo
}
func Failed(messageId string) ResponseVo {
	vo := ResponseVo{MessageId: messageId, Code: 1, IsSuccess: false, Msg: "failed"}
	return vo
}
func FailedWithMsg(messageId string, msg string) ResponseVo {
	vo := ResponseVo{MessageId: messageId, Code: 1, IsSuccess: false, Msg: msg}
	return vo
}
func FailedWithCodeMsg(messageId string, code int, msg string) ResponseVo {
	vo := ResponseVo{MessageId: messageId, Code: code, IsSuccess: false, Msg: msg}
	return vo
}
