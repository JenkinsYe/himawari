package aot

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"sort"
	"testing"
	"time"
)

func TestGetL2Path(t *testing.T) {
	path, fileNameSub := GetAOTPath()
	fmt.Printf("path %v fileNameSub: %v", path, fileNameSub)
}

func TestDownloadL2(t *testing.T) {
	var FtpClient HimawariFtpClient
	FtpClient.Init(true)
	FtpClient.DownloadAOT()
}

func TestDownloadAHI(t *testing.T) {
	var FtpClient HimawariFtpClient
	FtpClient.Init(true)
	FtpClient.DownloadAHI(logrus.WithField("test", true))
}

func TestGetAHIPath(t *testing.T) {
	path := GetAHIPath(0)
	fmt.Printf("path %v", path)
}

func TestBuildFileMap(t *testing.T) {
	fileMap, _, _ := BuildFileMap("/Users/ye/Desktop/test/")
	logrus.Infof("map: %v", fileMap)
}

func TestNeedDeleteAOT(t *testing.T) {
	need := NeedDeleteAOT("NC_H08_20210513_1110_L2ARP030_FLDK.02401_02401.nc")
	logrus.Infof("need: %v", need)
}

func TestHimawariFtpClient_ClearJob(t *testing.T) {
	var FtpClient HimawariFtpClient
	FtpClient.Init(true)
	FtpClient.ClearJob()
}

func TestSortEntry(t *testing.T) {
	var entries = []*ftp.Entry{
		{
			Name:   "asd0020",
			Target: "",
			Type:   0,
			Size:   0,
			Time:   time.Time{},
		},
		{
			Name:   "asd0030",
			Target: "",
			Type:   0,
			Size:   0,
			Time:   time.Time{},
		},
	}
	logrus.Infof("%v %v", entries[0], entries[1])
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name > entries[j].Name
	})
	logrus.Infof("%v %v", entries[0], entries[1])
}
