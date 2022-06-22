package main

import (
	"fmt"
	"github.com/xukgo/gphoneAttr"
	"os"
	"time"
)

func main() {
	err := initPhoneAttr()
	if err != nil {
		os.Exit(-1)
	}

	fmt.Println("init done")
	time.Sleep(time.Hour)
}

func initPhoneAttr() error {
	filePath := "/home/hermes/work/prefix.csv"
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = gphoneAttr.InitFromReader(file)
	if err != nil {
		return err
	}
	return nil
}
