package aot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"time"
)

func GetAOTPath(tolerance int) string {
	t := time.Now()
	for i := 0; i < tolerance; i++ {
		t = t.Add(-time.Hour)
	}
	time := t.UTC()
	year := strconv.Itoa(time.Year())
	month := strconv.Itoa(int(time.Month()))
	day := strconv.Itoa(time.Day())
	hour := strconv.Itoa(time.Hour())
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
