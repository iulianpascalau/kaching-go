package main

import (
	"fmt"
	"os"
	"os/exec"
)

const filename = "./kaching.mp3"

func main() {
	if !fileExists(filename) {
		fmt.Printf("file %s not found!\n", filename)
	}

	playSound(filename)
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
	fmt.Println("playing kaching...")
	app := "mpg123"
	arg := filename
	cmd := exec.Command(app, arg)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
