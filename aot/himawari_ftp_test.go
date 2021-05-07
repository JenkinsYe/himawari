package aot

import (
	"fmt"
	"testing"
)

func TestGetL2Path(t *testing.T) {
	path := getL2Path()
	fmt.Printf("path %v", path)
}

func TestDownloadL2(t *testing.T) {
	var FtpClient HimawariFtpClient
	FtpClient.Init()
	FtpClient.DownloadAOT()
}

func TestDownloadAHI(t *testing.T) {
	var FtpClient HimawariFtpClient
	FtpClient.Init()
	FtpClient.DownloadAHI()
}
