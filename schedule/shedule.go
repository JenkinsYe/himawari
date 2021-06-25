package schedule

import (
	"github.com/JenkinsYe/himawari/aot"
	"github.com/JenkinsYe/himawari/pm"
	"github.com/JenkinsYe/himawari/upload"
	"github.com/JenkinsYe/himawari/utils"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var aotLock sync.Mutex
var ahiLock sync.Mutex
var pmLock sync.Mutex
var uploadLock sync.Mutex
var resultClearLock sync.Mutex

func RunScheduledJobs(isLocal bool) {
	var ftpClient aot.HimawariFtpClient
	ftpClient.Init(isLocal)

	var particleInfoClient pm.ParticleInfoClient
	particleInfoClient.Init(isLocal)

	var uploadClient upload.ResultFtpClient
	uploadClient.Init(isLocal)

	cronjobs := cron.New()
	// aot数据下载任务
	cronjobs.AddFunc("@every 3h", func() {
		aotLock.Lock()
		defer aotLock.Unlock()
		logger := logrus.WithField("aot job id", time.Now().Unix())
		logger.Infof("----- CronJob DownloadAOT Start -----")
		err := ftpClient.DownloadAOT(logger)
		if err != nil {
			logger.WithError(err).Error("----- CronJob DownloadAOT Fail  -----")
		} else {
			logger.Infof("----- CronJob DownloadAOT Done  -----")
		}
	})

	// ahi数据下载任务
	cronjobs.AddFunc("@every 10m", func() {
		ahiLock.Lock()
		defer ahiLock.Unlock()
		logger := logrus.WithField("ahi job id", time.Now().Unix())
		logger.Infof("----- CronJob DownloadAHI Start -----")
		err := ftpClient.DownloadAHI(logger)
		if err != nil {
			logger.WithError(err).Error("----- CronJob DownloadAHI Fail  -----")
		} else {
			logger.Infof("----- CronJob DownloadAHI Done  -----")
		}
	})

	// aot/ahi清理任务
	cronjobs.AddFunc("@every 6h", func() {
		logger := logrus.WithField("clear job id", time.Now().Unix())
		logger.Infof("----- CronJob ClearJob Start -----")
		err := ftpClient.ClearJob(logger)
		if err != nil {
			logger.WithError(err).Error("----- CronJob ClearJob Fail  -----")
		} else {
			logger.Infof("----- CronJob ClearJob Done  -----")
		}
	})

	// pm数据获取任务
	cronjobs.AddFunc("@every 10m", func() {
		pmLock.Lock()
		defer pmLock.Unlock()
		logger := logrus.WithField("pm job id", time.Now().Unix())
		logger.Infof("----- CronJob GetPmInfo Start -----")
		err := particleInfoClient.GetInfo(logger)
		if err != nil {
			logger.WithError(err).Error("----- CronJob GetPmInfo Fail  -----")
		} else {
			logger.Infof("----- CronJob GetPmInfo Done  -----")
		}
	})

	// 结果文件上传任务
	cronjobs.AddFunc("@every 10m", func() {
		uploadLock.Lock()
		defer uploadLock.Unlock()
		logger := logrus.WithField("upload job id", time.Now().Unix())
		logger.Infof("----- CronJob Upload Start -----")
		nowTime := time.Now()
		yesterdayTime := nowTime.Add(-time.Hour * 24)
		err1 := uploadClient.Upload(utils.GetyyyyMMdd(yesterdayTime), logger)
		err2 := uploadClient.Upload(utils.GetyyyyMMdd(nowTime), logger)
		if err1 != nil {
			logger.WithError(err1).Error("----- CronJob Upload Fail  -----")
			return
		}
		if err2 != nil {
			logger.WithError(err2).Error("----- CronJob Upload Fail  -----")
			return
		}
		logger.Infof("----- CronJob Upload Done  -----")
	})

	// 结果数据清理任务
	cronjobs.AddFunc("@every 24h", func() {
		resultClearLock.Lock()
		defer resultClearLock.Unlock()
		logger := logrus.WithField("clear result job id", time.Now().Unix())
		logger.Infof("----- CronJob ClearResult Start -----")
		nowTime := time.Now()
		deleteTime := nowTime.Add(-time.Hour * 24 * 15)
		err := uploadClient.ClearResult(utils.GetyyyyMMdd(deleteTime), logger)
		if err != nil {
			logger.WithError(err).Error("---- CronJob ClearResult Fail ----")
			return
		}
		logger.Infof("---- CronJob ClearResult Fail ----")
	})
	cronjobs.Run()
}
