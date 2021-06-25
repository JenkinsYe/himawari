package server

import (
	"github.com/JenkinsYe/himawari/errex"
	"github.com/JenkinsYe/himawari/upload"
	"github.com/JenkinsYe/himawari/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	HOUR  = "hour"
	DAY   = "day"
	MONTH = "month"
)

var AHI_OUT_PATH string
var PM25_GROUND_HOUR_PATH string
var PM25_GROUND_DAY_PATH string
var PM25_GROUND_MON_PATH string
var PM10_GROUND_HOUR_PATH string
var PM10_GROUND_DAY_PATH string
var PM10_GROUND_MON_PATH string
var AOT_HOUR_PATH string
var AOT_DAY_PATH string
var AOT_MON_PATH string
var PM25_HOTSPOT_HOUR_PATH string
var PM25_HOTSPOT_DAY_PATH string
var PM25_HOTSPOT_MON_PATH string
var PM10_HOTSPOT_HOUR_PATH string
var PM10_HOTSPOT_DAY_PATH string
var PM10_HOTSPOT_MON_PATH string

var IsLocal bool
func Init(isLocal bool) {
	IsLocal = isLocal
	if isLocal {
		AHI_OUT_PATH = "/Users/ye/Desktop/himawari/AHI_out/"
	} else {
		AHI_OUT_PATH = "/himawari/AHI_out/"

		// ground path
		PM25_GROUND_HOUR_PATH = "/himawari/spq/output/ground2d/output/hrly/"
		PM25_GROUND_DAY_PATH = "/himawari/spq/output/ground2d/output/daily/"
		PM25_GROUND_MON_PATH = "/himawari/spq/output/ground2d/output/monthly/"
		PM10_GROUND_HOUR_PATH = "/himawari/spq/output/pm10/ground2d/output/hrly/"
		PM10_GROUND_DAY_PATH = "/himawari/spq/output/pm10/ground2d/output/daily/"
		PM10_GROUND_MON_PATH = "/himawari/spq/output/pm10/ground2d/output/monthly/"

		// aot path
		AOT_HOUR_PATH = "/himawari/spq/output/AOTreference/"
		AOT_DAY_PATH = "/himawari/spq/output/AOTowndaily/"
		AOT_MON_PATH = "/himawari/spq/output/AOTownmonthly/"

		// hotspot path
		PM25_HOTSPOT_HOUR_PATH = "/himawari/spq/output/preDrawing/hrly/"
		PM25_HOTSPOT_DAY_PATH = "/himawari/spq/output/preDrawing/daily/"
		PM25_HOTSPOT_MON_PATH = "/himawari/spq/output/preDrawing/monthly/"
		PM10_HOTSPOT_HOUR_PATH = "/himawari/spq/output/pm10/preDrawing/hrly/"
		PM10_HOTSPOT_DAY_PATH = "/himawari/spq/output/pm10/preDrawing/daily/"
		PM10_HOTSPOT_MON_PATH = "/himawari/spq/output/pm10/preDrawing/monthly/"
	}
}

func RunServer() {
	r := gin.Default()

/*	r.GET("/api/ahi", errex.ErrorWrapper(DownloadAHIFile))
	r.GET("/api/ground", errex.ErrorWrapper(DownloadGroundFile))
	r.GET("/api/aot", errex.ErrorWrapper(DownloadAOTFile))
	r.GET("/api/hotspot", errex.ErrorWrapper(DownloadHotSpotFile))*/
    r.GET("/api/reupload", errex.ErrorWrapper(Reload))
	r.Run(":8088")
	logrus.Infof("server start")
}

func Reload(c *gin.Context) error {
	logger := logrus.New()
	logger.WithField("method", "DownloadHotSpotFile")

	time := c.Query("time")
	if len(time) == 0 {
		logger.Errorf("time is missing")
		return errex.MissingTimeError
	}
	var uploadClient upload.ResultFtpClient
	uploadClient.Init(IsLocal)
	err := uploadClient.Upload(time, logrus.WithField("web", true))
	if err != nil {
		return err
	}
	c.JSON(200, gin.H{
		"code": 0,
		"message": "success",
	})
	return nil
}


func DownloadHotSpotFile(c *gin.Context) error {
	logger := logrus.New()
	logger.WithField("method", "DownloadHotSpotFile")

	time := c.Query("time")
	if len(time) == 0 {
		logger.Errorf("time is missing")
		return errex.MissingTimeError
	}
	fileType := c.Query("type")
	if len(fileType) == 0 {
		logger.Errorf("type is missing")
		return errex.MissingTypeError
	}
	threshold := c.Query("threshold")
	if len(threshold) == 0 {
		logger.Errorf("threshold is missing")
		return errex.MissingThresholdError
	}

	// pm类型
	pmType := c.Query("pmType")
	if len(pmType) == 0 {
		pmType = "1"
	}

	fileDir := getHotSpotFileDir(fileType, pmType)
	if len(fileDir) == 0 {
		logger.Errorf("unexpected ground file type, type: %v", fileType)
		return errex.InvalidParamError
	}

	// 时间格式校验
	formatTime, err := utils.LegalAndFormatGroundTimeParam(time, fileType)
	if err != nil {
		logger.Errorf("invalid time param, time: %v, type: %v", time, fileType)
		return errex.InvalidParamError
	}

	if !invalidThreshold(threshold) {
		logger.Errorf("invalid threshold, threshold: %v", threshold)
	}

	fileName := formatTime + threshold + ".nc"
	if !utils.Exists(fileDir + threshold + "/" + fileName) {
		logrus.Errorf("file: %v not exists", fileName)
		return errex.NotFoundError
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(fileDir + threshold + "/" + fileName)
	return nil
}

func DownloadAOTFile(c *gin.Context) error {
	logger := logrus.New()
	logger.WithField("method", "DownloadAOTFile")

	time := c.Query("time")
	if len(time) == 0 {
		logger.Errorf("time is missing")
		return errex.MissingTimeError
	}

	// 获取文件类型的路径
	groundType := c.Query("type")
	if len(groundType) == 0 {
		logger.Errorf("type is missing")
		return errex.MissingTypeError
	}

	fileDir := getAOTPathByType(groundType)
	if len(fileDir) == 0 {
		logger.Errorf("unexpected ground file type, type: %v", groundType)
		return errex.InvalidParamError
	}

	// 时间格式校验
	formatTime, err := utils.LegalAndFormatGroundTimeParam(time, groundType)
	if err != nil {
		logger.Errorf("invalid time param, time: %v, type: %v", time, groundType)
		return errex.InvalidParamError
	}

	fileName := formatTime + ".nc"
	if !utils.Exists(fileDir + fileName) {
		logrus.Errorf("file: %v not exists", fileName)
		return errex.NotFoundError
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(fileDir + fileName)
	return nil
}

func DownloadGroundFile(c *gin.Context) error {
	logger := logrus.New()
	logger.WithField("method", "DownloadGroundFile")
	time := c.Query("time")
	if len(time) == 0 {
		logger.Errorf("time is missing")
		return errex.MissingTimeError
	}

	// 获取文件类型的路径
	groundType := c.Query("type")
	if len(groundType) == 0 {
		logger.Errorf("type is missing")
		return errex.MissingTypeError
	}

	// pm类型
	pmType := c.Query("pmType")
	if len(pmType) == 0 {
		pmType = "1"
	}

	fileDir := getGroundPathByType(groundType, pmType)
	if len(fileDir) == 0 {
		logger.Errorf("unexpected ground file type, type: %v", groundType)
		return errex.InvalidParamError
	}

	// 时间格式校验
	formatTime, err := utils.LegalAndFormatGroundTimeParam(time, groundType)
	if err != nil {
		logger.Errorf("invalid time param, time: %v, type: %v", time, groundType)
		return errex.InvalidParamError
	}

	fileName := formatTime + ".nc"
	if !utils.Exists(fileDir + fileName) {
		logrus.Errorf("file: %v not exists", fileName)
		return errex.NotFoundError
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(fileDir + fileName)
	return nil
}

func DownloadAHIFile(c *gin.Context) error {
	logger := logrus.New()
	logger.WithField("method", "DownloadAHIFile")
	time := c.Query("time")
	if len(time) == 0 {
		logger.Errorf("time is missing")
		return errex.MissingTimeError
	}
	fileType := c.Query("type")
	if len(fileType) == 0 {
		logger.Errorf("fileType is missing")
		return errex.MissingTypeError
	}

	if fileType != "real" && fileType != "false" && fileType != "dust" {
		logger.Errorf("fileType is unexpected, fileType: %v", fileType)
		return errex.InvalidParamError
	}
	if !legalTimeParam(time) {
		logger.Errorf("time is unexpected, time: %v", time)
		return errex.InvalidParamError
	}

	fileName := time + fileType + ".nc"
	if !utils.Exists(AHI_OUT_PATH + fileName) {
		logrus.Errorf("file: %v not exists", fileName)
		return errex.NotFoundError
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(AHI_OUT_PATH + fileName)
	return nil
}

func legalTimeParam(time string) bool {
	if len(time) != 12 {
		return false
	}
	for _, q := range time {
		if q > '9' || q < '0' {
			return false
		}
	}
	return true
}

func getGroundPathByType(groundType string, pmType string) string {
	if pmType == "1" {
		switch groundType {
		case HOUR:
			return PM25_GROUND_HOUR_PATH
		case DAY:
			return PM25_GROUND_DAY_PATH
		case MONTH:
			return PM25_GROUND_MON_PATH
		default:
			return ""
		}
	} else {
		switch groundType {
		case HOUR:
			return PM10_GROUND_HOUR_PATH
		case DAY:
			return PM10_GROUND_DAY_PATH
		case MONTH:
			return PM10_GROUND_MON_PATH
		default:
			return ""
		}
	}
}

func getAOTPathByType(fileType string) string {
	switch fileType {
	case HOUR:
		return AOT_HOUR_PATH
	case DAY:
		return AOT_DAY_PATH
	case MONTH:
		return AOT_MON_PATH
	default:
		return ""
	}
}

func getHotSpotFileDir(fileType string, pmType string) string {
	if pmType == "1" {
		switch fileType {
		case HOUR:
			return PM25_HOTSPOT_HOUR_PATH
		case DAY:
			return PM25_HOTSPOT_DAY_PATH
		case MONTH:
			return PM25_HOTSPOT_MON_PATH
		default:
			return ""
		}
	} else {
		switch fileType {
		case HOUR:
			return PM10_HOTSPOT_HOUR_PATH
		case DAY:
			return PM10_HOTSPOT_DAY_PATH
		case MONTH:
			return PM10_HOTSPOT_MON_PATH
		default:
			return ""
		}
	}
}

func invalidThreshold(threshold string) bool {
	if threshold == "50" || threshold == "75" || threshold == "100" || threshold == "150" || threshold == "200" {
		return true
	}
	return false
}
