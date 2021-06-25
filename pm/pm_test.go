package pm

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestGetFileName(t *testing.T) {
	file := GetFileName(time.Now())
	logrus.Infof("file: %v", file)
}

func TestMakeFile(t *testing.T) {
	err := MakeMonitorFile("/Users/ye/Desktop/葵花卫星/2021052018.json")
	assert.Nil(t, err)
}

func TestReadMonitor(t *testing.T) {
	ReadMonitor(true)
}

func TestGetPMInfo(t *testing.T) {
	client := ParticleInfoClient{}
	client.Init(true)
	ReadMonitor(true)

	jsonFile, err := os.Open("/Users/ye/Desktop/葵花卫星/2021052818.json")
	if err != nil {
		logrus.WithError(err).Error("open file failed")
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	client.GetInfoFromFile(byteValue)
}

func TestGetParamListByTimeInternal(t *testing.T) {
	start1 := "2021053110"
	end1 := "2021053123"
	start2 := "2021060100"
	end2 := "2021060821"
	result := GetParamListByTimeInternal(start1, end1)
	logrus.Info(result)
	result = GetParamListByTimeInternal(start2, end2)
	logrus.Info(result)
}
