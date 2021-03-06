package main

import (
	"fmt"
	"iulianpascalau/kaching-go/blockchain"
	"iulianpascalau/kaching-go/logging"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const filename = "./kaching.mp3"
const address = "https://api.elrond.com"
const poolInterval = time.Second * 6 //round time

func main() {
	defer logging.MainLogger.CloseLogFile()

	if !fileExists(filename) {
		logging.MainLogger.Log(fmt.Sprintf("file %s not found!", filename))
	}

	logging.MainLogger.Log("application started")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	chPlaySound := make(chan struct{}, 1)

	ew := blockchain.NewEpochWatcher(address, poolInterval, chPlaySound)
	defer ew.Close()

	for {
		select {
		case <-sigs:
			logging.MainLogger.Log("terminating at user's signal...")
			return
		case <-chPlaySound:
			playSound(filename)
		}
	}
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func playSound(filename string) {
	logging.MainLogger.Log("playing kaching...")
	app := "mpg123"
	arg := filename
	cmd := exec.Command(app, arg)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
