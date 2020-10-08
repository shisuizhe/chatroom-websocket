package util

import (
	"fmt"
	"time"
)

func LogInfo(msg string) {
	fmt.Printf("[INFO][%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}

func LogErr(trace string, err error) {
	fmt.Printf("[ERROR][%s] %s - %s\n", time.Now().Format("2006-01-02 15:04:05"), trace, err.Error())
}
