package aot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetAOTPath() (string, string) {
	t := time.Now()
	for i := 0; i < 30; i++ {
		t = t.Add(-time.Hour * 24)
	}
	time := t.UTC()
	year := strconv.Itoa(time.Year())
	month := strconv.Itoa(int(time.Month()))
	day := strconv.Itoa(time.Day())
	if len(month) < 2 {
		month = fmt.Sprintf("0%s", month)
	}
	if len(day) < 2 {
		day = fmt.Sprintf("0%s", day)
	}
	return fmt.Sprintf("%s%s/daily", year, month), fmt.Sprintf("%s%s%s", year, month, day)
}

func GetAHIPath(tolerance int) string {
	t := time.Now()
	for i := 0; i < tolerance; i++ {
		t = t.Add(-time.Hour * 24)
	}
	time := t.UTC()
	year := strconv.Itoa(time.Year())
	month := strconv.Itoa(int(time.Month()))
	day := strconv.Itoa(time.Day())

	if len(month) < 2 {
		month = fmt.Sprintf("0%s", month)
	}
	if len(day) < 2 {
		day = fmt.Sprintf("0%s", day)
	}
	return fmt.Sprintf("%s%s/%s/", year, month, day)
}

func BuildFileMap(filePath string) (map[string]bool, map[string]int64, error) {
	fileMap := make(map[string]bool)
	fileSize := make(map[string]int64)
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		logrus.WithError(err).Error("readDir failed")
		return nil, nil, err
	}
	for _, f := range files {
		if f.Size() < 1024*4 {
			logrus.Warnf("unexpected size file: %v", f.Name())
			err = os.Remove(filePath + f.Name())
			if err != nil {
				logrus.WithError(err).Error("remove file failed")
			}
			continue
		}
		fileSize[f.Name()] = f.Size()
		fileMap[f.Name()] = true
	}
	return fileMap, fileSize, nil
}

func NeedDeleteAOT(name string) bool {
	t := time.Now()
	for i := 0; i < 90; i++ {
		t = t.Add(-time.Hour * 24)
	}
	time := t.UTC()
	year := strconv.Itoa(time.Year())
	month := strconv.Itoa(int(time.Month()))
	if len(month) < 2 {
		month = fmt.Sprintf("0%s", month)
	}
	toleranceDate := fmt.Sprintf("%s%s", year, month)
	if strings.Contains(name, toleranceDate) {
		return true
	} else {
		return false
	}
}
