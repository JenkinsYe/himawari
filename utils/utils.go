package utils

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	HOUR  = "hour"
	DAY   = "day"
	MONTH = "month"
)

func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}


func LegalAndFormatGroundTimeParam(time string, timeType string) (string, error) {

	if timeType == HOUR { // exp 2021052021 return 20210520_21
		if len(time) != 10 {
			return "", errors.New("unexpected length")
		}
		return fmt.Sprintf("%s_%s", time[:8], time[8:]), nil
	} else if timeType == DAY { // exp 20210520 return 20210520
		if len(time) != 8 {
			return "", errors.New("unexpected length")
		}
		return time, nil
	} else if timeType == MONTH {
		if len(time) != 6 {
			return "", errors.New("unexpected length")
		}
		return time, nil
	} else {
		return "", errors.New("unexpected timeType")
	}
}

func GetyyyyMMdd(time time.Time) string {
	return time.Format("20060102")
}
