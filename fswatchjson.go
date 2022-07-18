package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
)

const (
	fileName = "alert.log"
	maxSize  = 65536
)

func main() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(-1)
	}
	defer file.Close()

	file.Seek(0, io.SeekEnd)
	filedata := make([]byte, maxSize)

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("Event: %#v\n", event)
				n1, err := io.ReadAtLeast(file, filedata, 1)
				if err != nil {
					fmt.Println("ERROR: ", err)
					os.Exit(-1)
				}
				fmt.Printf("{Content: %d bytes: %s\n", n1, string(filedata[:n1]))

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
				os.Exit(-1)
			}
		}
	}()

	if err := watcher.Add(fileName); err != nil {
		fmt.Println("ERROR: ", err)
		os.Exit(-1)
	}

	<-done
	os.Exit(0)
}
