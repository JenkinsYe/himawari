package aot

import (
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type HimawariFtpClient struct {
	Account    string
	Password   string
	Host       string
	AOTWorkDir string
	AHIWorkDir string
}

func (client *HimawariFtpClient) Init() {
	client.Account = "www.875773677_qq.com"
	client.Password = "SP+wari8"
	client.Host = "ftp.ptree.jaxa.jp:21"
	client.AOTWorkDir = "/home/admin/himawari/AOT/"
	client.AHIWorkDir = "/home/admin/himawari/AHI/"
}

const FTPAOTDir = "/pub/himawari/L2/ARP/030/"
const FTPAHIDir = "/jma/netcdf/"

func (client *HimawariFtpClient) Test() error {
	logrus.Info("start dialing...")
	con, err := ftp.Dial(client.Host, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		logrus.WithError(err).Error("connect to ftp server failed")
		return err
	}
	con.Login(client.Account, client.Password)
	if err != nil {
		logrus.WithError(err).Error("login to ftp server failed")
		return err
	}
	files, err := con.List("")
	if err != nil {
		logrus.WithError(err).Error("list failed")
		return err
	}
	for _, file := range files {
		logrus.Infof("files: %v\n", file.Name)
	}
	return nil
}

func (client *HimawariFtpClient) DownloadAOT() error {
	con, err := client.getConnection()
	if err != nil {
		logrus.WithError(err).Error("get connection failed")
		return err
	}
	aotPath := getAOTPath()
	err = con.ChangeDir(FTPAOTDir + aotPath)
	if err != nil {
		logrus.WithError(err).Error("change dir failed")
		return err
	}
	entries, err := con.List("")
	if err != nil {
		logrus.WithError(err).Error("")
		return err
	}

	// check
	filePath := client.AOTWorkDir + aotPath
	if _, err := os.Stat(filePath); err != nil {
		logrus.Infof("path %v not exist, create it", filePath)
		err := os.MkdirAll(filePath, 0711)
		if err != nil {
			logrus.WithError(err).Error("create path failed")
			return err
		}
	}

	fileMap, err := buildFileMap(filePath)
	if err != nil {
		logrus.WithError(err).Error("build File map failed")
		return err
	}
	for _, entry := range entries {
		if strings.Contains(entry.Name, ".nc") && !fileMap[entry.Name] {
			logrus.Infof("downloading %v", entry.Name)
			r, err := con.Retr(entry.Name)
			if err != nil {
				logrus.WithError(err).Errorf("retr file: %v failed", entry.Name)
				return nil
			}

			file, err := os.Create(filePath + entry.Name)
			if err != nil {
				logrus.WithError(err).Error("create file failed")
				return err
			}
			for {
				var buf = make([]byte, 1024)
				n, _ := r.Read(buf)
				if n == 0 {
					break
				}
				file.Write(buf[:n])
			}
			file.Close()
			r.Close()
		}
	}
	con.Logout()
	con.Quit()
	return nil
}

func (client *HimawariFtpClient) DownloadAHI() error {
	con, err := client.getConnection()
	if err != nil {
		logrus.WithError(err).Error("get connection failed")
		return err
	}

	AHIPath := getAHIPath()
	err = con.ChangeDir(FTPAHIDir + AHIPath)
	if err != nil {
		logrus.WithError(err).Error("change dir failed")
		return err
	}
	entries, err := con.List("")
	if err != nil {
		logrus.WithError(err).Error("")
		return err
	}

	// check
	filePath := client.AHIWorkDir + AHIPath
	if _, err := os.Stat(filePath); err != nil {
		logrus.Infof("path %v not exist, create it", filePath)
		err := os.MkdirAll(filePath, 0711)
		if err != nil {
			logrus.WithError(err).Error("create path failed")
			return err
		}
	}

	fileMap, err := buildFileMap(filePath)
	if err != nil {
		logrus.WithError(err).Error("build File map failed")
		return err
	}
	for _, entry := range entries {
		if strings.Contains(entry.Name, "R21_FLDK.06001_06001.nc") && !fileMap[entry.Name] {
			logrus.Infof("downloading %v", entry.Name)
			r, err := con.Retr(entry.Name)
			if err != nil {
				logrus.WithError(err).Errorf("retr file: %v failed", entry.Name)
				return nil
			}

			file, err := os.Create(filePath + entry.Name)
			if err != nil {
				logrus.WithError(err).Error("create file failed")
				return err
			}
			for {
				var buf = make([]byte, 1024)
				n, _ := r.Read(buf)
				if n == 0 {
					break
				}
				file.Write(buf[:n])
			}
			file.Close()
			r.Close()
		}
	}
	con.Logout()
	con.Quit()
	return nil
}

func (client *HimawariFtpClient) getConnection() (*ftp.ServerConn, error) {
	//logrus.Infof("time %v, start dialing", time.Now().String())
	con, err := ftp.Dial(client.Host, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		logrus.WithError(err).Error("connect to ftp server failed")
		return nil, err
	}
	con.Login(client.Account, client.Password)
	if err != nil {
		logrus.WithError(err).Error("login to ftp server failed")
		return nil, err
	}
	return con, nil
}
