package common

// FtpMessage 消息
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
}
