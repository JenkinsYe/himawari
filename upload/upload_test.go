package upload

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestUpload(t *testing.T) {
	client := ResultFtpClient{}
	client.Init(true)
	err := client.Upload("20210629", logrus.WithField("test", true))
	if err != nil {
		logrus.WithError(err).Error("upload failed")
	}
}

