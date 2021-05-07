package aot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"time"
)

func getAOTPath() string {
	time := time.Now().UTC()

	year := strconv.Itoa(time.Year())
	month := strconv.Itoa(int(time.Month()))
	day := strconv.Itoa(time.Day())
	hour := strconv.Itoa(time.Hour() - 1)
	if len(month) < 2 {
		month = fmt.Sprintf("0%s", month)
	}
	if len(day) < 2 {
		day = fmt.Sprintf("0%s", day)
	}
	if len(hour) < 2 {
		hour = fmt.Sprintf("0%s", hour)
	}
	return fmt.Sprintf("%s%s/%s/%s/", year, month, day, hour)
}

func getAHIPath() string {
	time := time.Now().UTC()

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

func buildFileMap(filePath string) (map[string]bool, error) {
	fileMap := make(map[string]bool)
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		logrus.WithError(err).Error("readDir failed")
		return nil, err
	}
	for _, f := range files {
		fileMap[f.Name()] = true
	}
	return fileMap, nil
}
