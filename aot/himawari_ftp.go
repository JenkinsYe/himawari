package aot

import (
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sort"
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

func (client *HimawariFtpClient) Init(local bool) {
	client.Account = "www.875773677_qq.com"
	client.Password = "SP+wari8"
	client.Host = "ftp.ptree.jaxa.jp:21"
	if !local {
		client.AOTWorkDir = "/himawari/AOT/"
		client.AHIWorkDir = "/himawari/AHI/"
	} else {
		client.AOTWorkDir = "/Users/ye/Desktop/himawari/AOT/"
		client.AHIWorkDir = "/Users/ye/Desktop/himawari/AHI/"
	}
}

const FTP_AOT_DIR = "/pub/himawari/L3/ARP/031/"
const FTP_AHI_DIR = "/jma/netcdf/"
const MAX_AOT_TOLERANCE = 2
const MAX_AHI_TOLERANCE = 0

func (client *HimawariFtpClient) DownloadAOT(logger *logrus.Entry) error {
	con, err := client.getConnection()
	if err != nil {
		logger.WithError(err).Error("get connection failed")
		return err
	}
	aotPath, fileNameSub := GetAOTPath()
	err = con.ChangeDir(FTP_AOT_DIR + aotPath)
	if err != nil {
		logger.WithError(err).Error("change dir failed")
		return err
	}
	entries, err := con.List("")
	if err != nil {
		logger.WithError(err).Error("")
		return err
	}

	// check
	filePath := client.AOTWorkDir
	if _, err := os.Stat(filePath); err != nil {
		logger.Infof("path %v not exist, create it", filePath)
		err := os.MkdirAll(filePath, 0711)
		if err != nil {
			logger.WithError(err).Error("create path failed")
		}
	}

	fileMap, fileSize, err := BuildFileMap(filePath)
	if err != nil {
		logger.WithError(err).Error("build File map failed")
		return err
	}
	for _, entry := range entries {
		tooSmall := false
		// 本地有而且本地大小与服务端不一致
		if fileMap[entry.Name] && uint64(fileSize[entry.Name]) < entry.Size {
			tooSmall = true
			logger.Warnf("unexpected size file: %v", entry.Name)
			err = os.Remove(filePath + entry.Name)
			if err != nil {
				logger.WithError(err).Error("remove file failed")
			}
		}
		if strings.Contains(entry.Name, ".nc") && (!fileMap[entry.Name] || tooSmall) && strings.Contains(entry.Name, fileNameSub) {
			logger.Infof("downloading %v", entry.Name)
			r, err := con.Retr(entry.Name)
			if err != nil {
				logger.WithError(err).Errorf("retr file: %v failed", entry.Name)
				return nil
			}

			file, err := os.Create(filePath + entry.Name)
			if err != nil {
				logger.WithError(err).Error("create file failed")
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

func (client *HimawariFtpClient) DownloadAHI(logger *logrus.Entry) error {
	con, err := client.getConnection()
	if err != nil {
		logger.WithError(err).Error("get connection failed")
		return err
	}

	for tolerance := 0; tolerance <= MAX_AHI_TOLERANCE; tolerance++ {
		AHIPath := GetAHIPath(tolerance)
		err = con.ChangeDir(FTP_AHI_DIR + AHIPath)
		if err != nil {
			logger.WithError(err).Error("change dir failed")
			continue
		}
		entries, err := con.List("")
		if err != nil {
			logger.WithError(err).Error("")
			continue
		}

		// check
		filePath := client.AHIWorkDir + AHIPath
		if _, err := os.Stat(filePath); err != nil {
			logger.Infof("path %v not exist, create it", filePath)
			err := os.MkdirAll(filePath, 0711)
			if err != nil {
				logger.WithError(err).Error("create path failed")
				continue
			}
		}

		fileMap, fileSize, err := BuildFileMap(filePath)
		if err != nil {
			logger.WithError(err).Error("build File map failed")
			continue
		}
		downloadCount := 0
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name > entries[j].Name
		})
		for _, entry := range entries {
			tooSmall := false
			// 本地有而且本地大小与服务端不一致
			if fileMap[entry.Name] && uint64(fileSize[entry.Name]) * 100 < entry.Size * 70 {
				tooSmall = true
				logger.Warnf("unexpected size file: %v", entry.Name)
				err = os.Remove(filePath + entry.Name)
				if err != nil {
					logger.WithError(err).Error("remove file failed")
				}
			}
			if strings.Contains(entry.Name, "R21_FLDK.06001_06001.nc") && (!fileMap[entry.Name] || tooSmall) {
				logger.Infof("downloading %v", entry.Name)
				r, err := con.Retr(entry.Name)
				if err != nil {
					logger.WithError(err).Errorf("retr file: %v failed", entry.Name)
					return nil
				}

				file, err := os.Create(filePath + entry.Name)
				if err != nil {
					logger.WithError(err).Error("create file failed")
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
				downloadCount++
				if downloadCount >= 2 {
					break //AHI文件太大，单次任务最多下载两个
				}
			}
		}
	}
	con.Logout()
	con.Quit()
	return nil
}

func (client *HimawariFtpClient) ClearJob(logger *logrus.Entry) error {

	for tolerance := 3; tolerance <= 5; tolerance++ {
		AHIPath := GetAHIPath(tolerance)
		filePath := client.AHIWorkDir + AHIPath
		if _, err := os.Stat(filePath); err != nil {
			logger.Infof("path %v not exist, there is no need to clear", filePath)
		}
		err := os.RemoveAll(filePath)
		if err != nil {
			logger.WithError(err).Errorf("remove directory %v failed", filePath)
		}
	}

	aotFiles, err := ioutil.ReadDir(client.AOTWorkDir)
	if err != nil {
		logger.WithError(err).Error("readDir failed")
		return err
	}
	for _, f := range aotFiles {
		if NeedDeleteAOT(f.Name()) {
			logger.Warnf("expired aot file: %v", f.Name())
			err = os.Remove(client.AOTWorkDir + f.Name())
			if err != nil {
				logger.WithError(err).Errorf("remove file %v failed", f.Name())
			}
		}
	}
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
