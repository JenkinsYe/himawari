package schedule

import (
	"github.com/JenkinsYe/himawari/aot"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"time"
)

func RunScheduledJobs() {
	var ftpClient aot.HimawariFtpClient
	ftpClient.Init()

	cronjobs := cron.New()
	cronjobs.AddFunc("@every 10m", func() {
		logger := logrus.WithField("job id", time.Now().Unix())
		logger.Infof("----- CronJob DownloadAOT Start -----")
		err := ftpClient.DownloadAOT()
		if err != nil {
			logger.WithError(err).Error("----- CronJob DownloadAOT Failed -----")
		}
	})

	cronjobs.AddFunc("@every 25m", func() {
		logger := logrus.WithField("job id", time.Now().Unix())
		logger.Infof("----- CronJob DownloadAHI Start -----")
		err := ftpClient.DownloadAHI()
		if err != nil {
			logger.WithError(err).Error("----- CronJob DownloadAHI Failed -----")
		}
	})
	cronjobs.Run()
}
