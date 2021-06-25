package pm

import (
	"encoding/json"
	"github.com/gocarina/gocsv"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type MonitorInfo struct {
	Code      string  `csv:"监测站编号"`
	Longitude float64 `csv:"经度"`
	Latitude  float64 `csv:"纬度"`
	Type      string  `csv:"站点类型"`
	Name      string  `csv:"监测站名称"`
}

const MonitorPath = "/Users/ye/Desktop/himawari/monitor/"
const MonitorInputPath = "/Users/ye/Desktop/himawari/monitor/monitor.csv"
const ServerMonitorInputPath = "/himawari/monitor/monitor.csv"


var MonitorInfos []*MonitorInfo

func MakeMonitorFile(filePath string) error {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		logrus.WithError(err).Error("open file failed")
		return err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	var info ResponseInfo
	err = json.Unmarshal(byteValue, &info)
	if err != nil {
		logrus.WithError(err).Error("unmarshal failed")
		return err
	}

	file, err := os.OpenFile(MonitorPath+"monitor.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		logrus.WithError(err).Error("create csv file failed")
		return err
	}
	defer file.Close()

	pmCSVInfos := []*MonitorInfo{}
	entryMap := make(map[string]Entry)
	for _, entry := range info.Data.Entries {
		if _, ok := entryMap[entry.Code]; !ok {
			entryMap[entry.Code] = entry
		}
	}
	for _, entry := range entryMap {
		pmCSVInfos = append(pmCSVInfos, &MonitorInfo{
			Code:      entry.Code,
			Latitude:  entry.Latitude,
			Longitude: entry.Longitude,
			Type:      entry.Sitetypename,
			Name:      entry.Name,
		})
	}

	err = gocsv.Marshal(pmCSVInfos, file)
	if err != nil {
		logrus.WithError(err).Error("write csv file failed")
		return err
	}
	return nil
}

func ReadMonitor(isLocal bool) error {
	var monitorFile *os.File
	var err error
	if isLocal {
		monitorFile, err = os.OpenFile(MonitorInputPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			logrus.WithError(err).Panic("generate monitor file failed")
			return err
		}
	} else {
		monitorFile, err = os.OpenFile(ServerMonitorInputPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			logrus.WithError(err).Panic("generate monitor file failed")
			return err
		}
	}


	if err := gocsv.Unmarshal(monitorFile, &MonitorInfos); err != nil {
		logrus.WithError(err).Panic("unmarshal monitor.csv failed")
		return err
	}

	logrus.Infof("readMonitor success, total number: %v", len(MonitorInfos))
	logrus.Infof("first monitor info: %v", MonitorInfos[0])
	logrus.Infof("last monitor info: %v", MonitorInfos[len(MonitorInfos) - 1])
	return nil
}
