package main

import (
	"github.com/JenkinsYe/himawari/pm"
	"github.com/JenkinsYe/himawari/schedule"
	"github.com/JenkinsYe/himawari/server"
	"github.com/sirupsen/logrus"
)

func main() {
	isLocal := false
	isCronjob := true
	// init
	pm.ReadMonitor(isLocal)

	debug := false
	start1 := "2021053110"
	end1 := "2021053123"
	start2 := "2021060100"
	end2 := "2021060821"

	logrus.Infof("debug: %v, start: %v, end: %v", debug, start1, end1)
	if debug {
		var particleInfoClient pm.ParticleInfoClient
		particleInfoClient.Init(isLocal)
		logrus.Info("is debug mode")
		particleInfoClient.GetInfoFromInternal(start1, end1)
		particleInfoClient.GetInfoFromInternal(start2, end2)
		return
	}

	// cronjobs
	if isCronjob {
		schedule.RunScheduledJobs(isLocal)
	} else {
		// http server
		server.Init(isLocal)
		server.RunServer()
	}
}
