package server

import (
	"time"

	"github.com/fatih/color"
)

var (
	blue   = color.New(color.FgHiBlue).SprintFunc()
	red    = color.New(color.FgHiRed).SprintFunc()
	green  = color.New(color.FgHiGreen).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
)

func getTimeSt() string {
	return time.Now().Format("2006-01-02_15:04:05")
}
