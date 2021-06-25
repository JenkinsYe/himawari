package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLegalAndFormatGroundTimeParam(t *testing.T) {
	s, err := LegalAndFormatGroundTimeParam("2021030412", "hour")
	assert.Nil(t, err)
	assert.Equal(t, "20210304_12", s)

	s, err = LegalAndFormatGroundTimeParam("20210304", "day")
	assert.Nil(t, err)
	assert.Equal(t, "20210304", s)

	s, err = LegalAndFormatGroundTimeParam("202103", "month")
	assert.Nil(t, err)
	assert.Equal(t, "202103", s)
}

func TestGetyyyyMMdd(t *testing.T) {
	formatTime := GetyyyyMMdd(time.Now())
	fmt.Printf("time: %v", formatTime)
}