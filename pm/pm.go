package pm

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gocarina/gocsv"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type ResponseInfo struct {
	Code int `json:"code"`
	Data struct {
		Entries []Entry `json:"entries"`
	} `json:"data"`
	Message string `json:"message"`
}

type Entry struct {
	Sitetypename string  `json:"siteTypeName"`
	Sitetypecode string  `json:"siteTypeCode"`
	Name         string  `json:"name"`
	Code         string  `json:"code"`
	Gridname     string  `json:"gridName"`
	Gridcode     string  `json:"gridCode"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Pm25         float32 `json:"PM2.5"`
	Pm10         float32 `json:"PM10"`
	Time         string  `json:"time"`
}

type PMCSVInfo struct {
	Code string  `csv:"监测站编号"`
	PM25 float32 `csv:"PM2.5"`
	PM10 float32 `csv:"PM10"`
}

type ParticleInfoClient struct {
	Host           string
	Port           string
	Url            string
	MonitorDataDir string
}

func (client *ParticleInfoClient) Init(isLocal bool) {
	if !isLocal {
		client.Host = "http://10.229.7.117"
		client.Port = strconv.FormatInt(8899, 10)
		client.Url = "/aims-server/web/api/v1/external/site-data/hour"
		client.MonitorDataDir = "/himawari/monitor/pmdata/"
	} else {
		client.Host = "http://127.0.0.1"
		client.Port = strconv.FormatInt(8080, 10)
		client.Url = "/ping"
		client.MonitorDataDir = "/Users/ye/Desktop/himawari/monitor/pmdata/"
	}
}

func (client *ParticleInfoClient) GetInfo(logger *logrus.Entry) error {
	var err error
	httpClient := resty.New()

	param := GetTimeParam()

	resp, err := httpClient.R().
		SetQueryParam("startTime", param).
		EnableTrace().
		Get(client.Host + ":" + client.Port + client.Url)
	if err != nil {
		logger.WithError(err).Error("GetPMInfo failed")
		return err
	}
	var respInfo ResponseInfo

	if err = json.Unmarshal(resp.Body(), &respInfo)
		err != nil {
		logger.WithError(err).Errorf("Unmarshal json failed")
		return err
	}

	if _, err := os.Stat(client.MonitorDataDir); err != nil {
		logger.Infof("path %v not exist, create it", client.MonitorDataDir)
		err := os.MkdirAll(client.MonitorDataDir, 0711)
		if err != nil {
			logger.WithError(err).Error("create path failed")
			return err
		}
	}
	if len(respInfo.Data.Entries) == 0 {
		logger.Info("data not ready yet")
		return nil
	}
	clientFile, err := os.OpenFile(client.MonitorDataDir+param+".csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		logger.WithError(err).Error("create csv file failed")
		return err
	}
	defer clientFile.Close()

	codeToEntryMap := make(map[string]Entry)
	pmCSVInfos := []*PMCSVInfo{}
	for _, entry := range respInfo.Data.Entries {
		codeToEntryMap[entry.Code] = entry
	}
	var invalid int64
	for _, monitor := range MonitorInfos {
		entry, ok := codeToEntryMap[monitor.Code]
		if ok {
			pmCSVInfos = append(pmCSVInfos, &PMCSVInfo{
				Code: entry.Code,
				PM25: entry.Pm25,
				PM10: entry.Pm10,
			})
		} else {
			invalid++
			pmCSVInfos = append(pmCSVInfos, &PMCSVInfo{
				Code: monitor.Code,
				PM25: -1,
				PM10: -1,
			})
		}
	}

	err = gocsv.Marshal(pmCSVInfos, clientFile)
	if err != nil {
		logger.WithError(err).Error("write csv file failed")
		return err
	}
	logger.Infof("total invalid info: %v", invalid)
	return nil
}

func (client *ParticleInfoClient) GetInfoFromInternal(start string, end string) error {
	params := GetParamListByTimeInternal(start, end)
	httpClient := resty.New()
	for _, param := range params {
		logrus.Infof("doing %v", param)
		resp, err := httpClient.R().
			SetQueryParam("startTime", param).
			EnableTrace().
			Get(client.Host + ":" + client.Port + client.Url)
		if err != nil {
			logrus.WithError(err).Error("GetPMInfo failed")
			return err
		}
		var respInfo ResponseInfo

		if err = json.Unmarshal(resp.Body(), &respInfo)
			err != nil {
			logrus.WithError(err).Errorf("Unmarshal json failed")
			return err
		}

		if _, err := os.Stat(client.MonitorDataDir); err != nil {
			logrus.Infof("path %v not exist, create it", client.MonitorDataDir)
			err := os.MkdirAll(client.MonitorDataDir, 0711)
			if err != nil {
				logrus.WithError(err).Error("create path failed")
				return err
			}
		}
		if len(respInfo.Data.Entries) == 0 {
			logrus.Info("data not ready yet")
			return nil
		}
		clientFile, err := os.OpenFile(client.MonitorDataDir+param+".csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			logrus.WithError(err).Error("create csv file failed")
			return err
		}
		defer clientFile.Close()

		codeToEntryMap := make(map[string]Entry)
		pmCSVInfos := []*PMCSVInfo{}
		for _, entry := range respInfo.Data.Entries {
			codeToEntryMap[entry.Code] = entry
		}
		var invalid int64
		for _, monitor := range MonitorInfos {
			entry, ok := codeToEntryMap[monitor.Code]
			if ok {
				pmCSVInfos = append(pmCSVInfos, &PMCSVInfo{
					Code: entry.Code,
					PM25: entry.Pm25,
					PM10: entry.Pm10,
				})
			} else {
				invalid++
				pmCSVInfos = append(pmCSVInfos, &PMCSVInfo{
					Code: monitor.Code,
					PM25: -1,
					PM10: -1,
				})
			}
		}

		err = gocsv.Marshal(pmCSVInfos, clientFile)
		if err != nil {
			logrus.WithError(err).Error("write csv file failed")
			return err
		}
		logrus.Infof("total invalid info: %v", invalid)
	}
	return nil
}

func (client *ParticleInfoClient) GetInfoFromFile(value []byte) error {
	var respInfo ResponseInfo
	var err error
	if err = json.Unmarshal(value, &respInfo)
		err != nil {
		logrus.WithError(err).Errorf("Unmarshal json failed")
		return err
	}
	param := GetTimeParam()

	if _, err := os.Stat(client.MonitorDataDir); err != nil {
		logrus.Infof("path %v not exist, create it", client.MonitorDataDir)
		err := os.MkdirAll(client.MonitorDataDir, 0711)
		if err != nil {
			logrus.WithError(err).Error("create path failed")
			return err
		}
	}
	if len(respInfo.Data.Entries) == 0 {
		logrus.Info("data not ready yet")
		return nil
	}
	clientFile, err := os.OpenFile(client.MonitorDataDir+param+".csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		logrus.WithError(err).Error("create csv file failed")
		return err
	}
	defer clientFile.Close()

	codeToEntryMap := make(map[string]Entry)
	pmCSVInfos := []*PMCSVInfo{}
	for _, entry := range respInfo.Data.Entries {
		codeToEntryMap[entry.Code] = entry
	}
	var invalid int64
	for _, monitor := range MonitorInfos {
		entry, ok := codeToEntryMap[monitor.Code]
		if ok {
			pmCSVInfos = append(pmCSVInfos, &PMCSVInfo{
				Code: entry.Code,
				PM25: entry.Pm25,
				PM10: entry.Pm10,
			})
		} else {
			invalid++
			pmCSVInfos = append(pmCSVInfos, &PMCSVInfo{
				Code: monitor.Code,
				PM25: -1,
				PM10: -1,
			})
		}
	}

	err = gocsv.Marshal(pmCSVInfos, clientFile)
	if err != nil {
		logrus.WithError(err).Error("write csv file failed")
		return err
	}
	logrus.Infof("total invalid info: %v", invalid)
	return nil
}

func GetFileName(time time.Time) string {
	return time.Format("2006010215") + ".csv"
}

func GetTimeParam() string {
	nowT := time.Now()
	paramT := nowT.Add(-time.Hour * 2)
	param := paramT.Format("2006010215")
	return param
}

func GetParamListByTimeInternal(start string, end string) []string {
	var result []string
	startPrefix := start[:8]
	endPrefix := end[:8]
	startSuffix := start[8:]
	endSuffix := end[8:]
	startT, _ := strconv.ParseInt(startSuffix, 10, 64)
	endT, _ := strconv.ParseInt(endSuffix, 10, 64)
	if startPrefix == endPrefix {
		for i := startT; i <= endT; i++ {
			if i < 10 {
				result = append(result, fmt.Sprintf("%s0%d", startPrefix, i))
			} else {
				result = append(result, fmt.Sprintf("%s%d", startPrefix, i))
			}
		}
	} else {
		for i := startT; i < 24; i++ {
			if i < 10 {
				result = append(result, fmt.Sprintf("%s0%d", startPrefix, i))
			} else {
				result = append(result, fmt.Sprintf("%s%d", startPrefix, i))
			}
		}

		m, _ := strconv.ParseInt(startPrefix, 10, 64)
		n, _ := strconv.ParseInt(endPrefix, 10, 64)
		for i := m + 1; i < n; i++ {
			for j := 0; j < 24; j++ {
				if j < 10 {
					result = append(result, fmt.Sprintf("%d0%d", i, j))
				} else {
					result = append(result, fmt.Sprintf("%d%d", i, j))
				}
			}
		}

		for i := 0; i <= int(endT); i++ {
			if i < 10 {
				result = append(result, fmt.Sprintf("%s0%d", endPrefix, i))
			} else {
				result = append(result, fmt.Sprintf("%s%d", endPrefix, i))
			}
		}
	}
	return result
}
