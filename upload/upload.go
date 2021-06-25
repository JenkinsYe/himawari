package upload

import (
	"github.com/JenkinsYe/himawari/utils"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

type ResultFtpClient struct {
	Account         string
	Password        string
	Host            string
	LocalResultDir  string
	RemoteResultDir string
}

func (client *ResultFtpClient) Init(isLocal bool) {
	if isLocal {
		client.Host = "10.69.95.35:30121"
		client.Account = "ftpuser"
		client.Password = "ftpuser@123"
		client.LocalResultDir = "/Users/ye/Desktop/himawari/out/"
		client.RemoteResultDir = "/himawari/out/"
	} else {
		client.Host = "10.229.7.117:30121"
		client.Account = "ftpuser"
		client.Password = "ftpuser@123"
		client.LocalResultDir = "/himawari/out/"
		client.RemoteResultDir = "/himawari/out/"
	}
}

func (client *ResultFtpClient) Upload(timeString string, logger *logrus.Entry) error {
	con, err := client.getConnection()
	if err != nil {
		logger.WithError(err).Error("get connection failed")
		return err
	}
	defer func() {
		con.Logout()
		con.Quit()
	}()

	if err = client.changeToTimeDir(con, timeString); err != nil {
		logger.WithError(err).Errorf("changeToTimeDir failed, timeString: %v", timeString)
		return err
	}
	condir, _ := con.CurrentDir()
	logger.Infof("pwd: %v", condir)

	// 1. 获取远端ftp_server的文件列表
	remoteFileList, err := con.List("")
	if err != nil {
		logger.WithError(err).Errorf("Get RemoteFileList Failed")
		return err
	}
	remoteFileMap := make(map[string]bool, 0)
	for _, file := range remoteFileList {
		remoteFileMap[file.Name] = true
	}

	// 2. 获取本地文件列表，将未上传的上传
	fileDir := client.LocalResultDir + timeString + "/"
	if !utils.Exists(fileDir) {
		logger.Infof("local dir %v not exists", fileDir)
		return nil
	}
	fileList, err := ioutil.ReadDir(fileDir)
	if err != nil {
		logger.WithError(err).Errorf("open localDir failed, path: %v", fileDir)
		return err
	}
	var uploadFailedList []string
	for _, file := range fileList {
		if !remoteFileMap[file.Name()] {
			err = client.storeFile(con, fileDir, file.Name())
			if err != nil {
				logger.WithError(err).Errorf("upload file failed, fileName: %v", file.Name())
				uploadFailedList = append(uploadFailedList, file.Name())
				continue
			}
			logger.Infof("upload file success, fileName: %v", file.Name())
		}
	}
	logger.Infof("upload failed len : %v", len(uploadFailedList))
	if len(uploadFailedList) > 0 {
		logger.Infof("failed list: %v", uploadFailedList)
	}
	return nil
}

func (client *ResultFtpClient) ClearResult(timeString string, logger *logrus.Entry) error {
	fileDir := client.LocalResultDir + timeString + "/"
	if !utils.Exists(fileDir) {
		logger.Infof("local dir %v not exists", fileDir)
		return nil
	}
	err := os.RemoveAll(fileDir)
	if err != nil {
		logger.WithError(err).Error("remove result: %v failed", fileDir)
		return err
	}
	logger.Infof("remove result: %v success", fileDir)
	return nil
}

func (client *ResultFtpClient) storeFile(con *ftp.ServerConn, filePath, fileName string) error {
	file, err := os.Open(filePath + fileName)
	if err != nil {
		logrus.WithError(err).Errorf("open file failed, filePath: %v", filePath + fileName)
		return err
	}
	defer file.Close()

	err = con.Stor(fileName, file)
	if err != nil {
		logrus.WithError(err).Errorf("stor file failed, filePath: %v", filePath)
		return err
	}
	return nil
}

func (client *ResultFtpClient) changeToTimeDir(con *ftp.ServerConn, timeString string) error {
	err := con.ChangeDir(client.RemoteResultDir)
	if err != nil {
		logrus.WithError(err).Errorf("change to root dir failed")
		return err
	}

	rootDir, err := con.List("")
	if err != nil {
		logrus.WithError(err).Errorf("get dir list failed")
		return err
	}
	for _, dir := range rootDir {
		if dir.Name == timeString {
			err = con.ChangeDir(dir.Name)
			if err != nil {
				logrus.WithError(err).Errorf("change to dir failed, dir %v", dir)
				return err
			}
			return nil
		}
	}
	con.MakeDir(timeString)
	err = con.ChangeDir(timeString)
	if err != nil {
		logrus.WithError(err).Errorf("change to time dir failed, dir: %v", timeString)
		return err
	}
	return nil
}

func (client *ResultFtpClient) getConnection() (*ftp.ServerConn, error) {
	con, err := ftp.Dial(client.Host, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		logrus.WithError(err).Error("connect to ftp server failed")
		return nil, err
	}
	err = con.Login(client.Account, client.Password)
	if err != nil {
		logrus.WithError(err).Error("login to ftp server failed")
		return nil, err
	}
	return con, nil
}
