package ftpClient

import (
	"bytes"
	"file-push/models"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// NewFtpConn 连接ftp并登录
func NewFtpConn(ftpMessage models.FtpMessage) (*ftp.ServerConn, error) {
	c, err := ftp.Dial(fmt.Sprintf("%s:%s", ftpMessage.FtpHost, ftpMessage.FtpPort),
		ftp.DialWithTimeout(5*time.Second)) // , ftp.DialWithDebugOutput(os.Stdout)
	if err != nil {
		log.Fatalf("dial ftp %s err %v", ftpMessage.FtpUser, err)
	}
	err = c.Login(ftpMessage.FtpUser, ftpMessage.FtpPassword)
	if err != nil {
		log.Fatal(err)
	}
	return c, err
}

// Change2WorkDir change to work dir if not exist then create dir
func Change2WorkDir(conn *ftp.ServerConn, workDir string, baseDir string) error {
	return ChangeDirAndMakeDirIfNotExist(conn, workDir, baseDir)
}
func ChangeDirAndMakeDirIfNotExist(conn *ftp.ServerConn, workDir string, baseDir string) error {
	pathSlash := filepath.ToSlash(workDir)
	pathSplits := strings.Split(pathSlash, "/")
	if baseDir != "" {
		err := conn.ChangeDir(baseDir)
		return err
	}
	for i := range pathSplits {
		currentPath := pathSplits[i]
		if currentPath == "" {
			continue
		}
		err := conn.ChangeDir(currentPath)
		if err == nil {
			continue
		}
		err = conn.MakeDir(currentPath)
		if err != nil {
			return err
		}
		if err = conn.ChangeDir(currentPath); err != nil {
			return err
		}
	}
	return nil
}
func UploadFile(conn *ftp.ServerConn, locateFilePath string) (uploadResult string, err error) {
	file, err := os.Open(locateFilePath)
	if err != nil {
		log.Printf("open file error %v", err)
		return "文件不存在", err
	}
	fileName := filepath.Base(file.Name())
	uploadingFile := fileName + ".temp"

	// 远程目录下已经存在文件
	fileExists := FtpFileExists(conn, fileName)
	if fileExists {
		log.Printf("file has already existed.")
		// TODO 文件存在怎么处理? 重传覆盖还是略过
		return "file has already existed", nil
	}
	// 新开始或断点续传
	var remoteSize int64 = 0
	uploadingFileExist := FtpFileExists(conn, uploadingFile)
	if uploadingFileExist {
		remoteSize, err = conn.FileSize(uploadingFile)
	}
	stat, _ := file.Stat()
	fileSize := stat.Size()
	onceReadLength := 10 * 1024 * 1024
	data := make([]byte, onceReadLength)
	n := -1
	total := 0
	for {
		n, err = file.ReadAt(data, remoteSize)
		if err == io.EOF || n == 0 {
			if n != 0 {
				total += n
				err = conn.StorFrom(uploadingFile, bytes.NewReader(data[:n]), uint64(remoteSize))
			}
			err := conn.Rename(uploadingFile, fileName)
			if err != nil {
				log.Printf("rename file err %v", err)
			}
			log.Printf("file upload completed...%s size %d", fileName, total)
			break
		} else if err != nil {
			log.Printf("read local file error %v", err)
			break
		}
		total += n
		for {
			err = conn.StorFrom(uploadingFile, bytes.NewReader(data), uint64(remoteSize))
			if err == nil {
				break
			} else {
				time.Sleep(500 * time.Millisecond)
				log.Print(err)
			}
		}
		//time.Sleep(10 * time.Second)
		remoteSize += int64(n)
		log.Printf("file [%s] uploaded %d  total size %d process %s", uploadingFile, remoteSize, fileSize, fmt.Sprintf("%.2f", float64(remoteSize)/float64(fileSize)))
		if err != nil {
			log.Fatalf("cannot Store file, fileName %s , err %v", file.Name(), err)
		}
	}
	return "", nil
}

// FtpFileExists judge the file weather exist on the ftp
func FtpFileExists(conn *ftp.ServerConn, fileName string) bool {
	exist := false
	dir, _ := conn.CurrentDir()
	entries, _ := conn.List(dir + "/" + fileName)
	if len(entries) > 0 {
		exist = true
	}
	return exist
}

// Quit quit ftp conn
func Quit(conn *ftp.ServerConn) {
	if err := conn.Quit(); err != nil {
		log.Fatal(err)
	}
}
