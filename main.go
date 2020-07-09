package main

import (
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/fsnotify/fsnotify"
	"github.com/otiai10/copy"
)

const SUPAI = `
  _____ _    _ _____        _____ 
 / ____| |  | |  __ \ /\   |_   _|
| (___ | |  | | |__) /  \    | |  
 \___ \| |  | |  ___/ /\ \   | |  
 ____) | |__| | |  / ____ \ _| |_ 
|_____/ \____/|_| /_/    \_\_____|				 
`


const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func main() {
	customExitHanlder()

	fmt.Printf(WarningColor, SUPAI+"\n")

	arguments := checkArguments()
	
	watchFolder(arguments[0], arguments[1])
}

// Check if arguments are 2 (watch, outdir)
// and check if both are valid directories
// returns array of arguments
func checkArguments()(hola []string)  {
	arguments := os.Args[1:]

	if len(arguments) < 2 {
		log.Fatalf(ErrorColor,"Invalid number of arguments, must be 2!\n")
	}

	for _, path := range arguments {
		dirExists(path)
	}

	return arguments
}

// Check if the provided path is a valid directory
// if the provided path is not a valid one it will exit
func dirExists(path string) {
	_, err := os.Stat(path)
	
    if err != nil {
      log.Fatalf(ErrorColor, path+" is not a valid path!")
    }
}

// watch directory files changes
func watchFolder(folderToWatch string, folderToCopy string)  {
	watcher, err := fsnotify.NewWatcher()
	
	if err != nil {
		log.Fatal(err)
	}
	
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf(InfoColor, event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf(DebugColor, "modified file: "+ event.Name)
					copyError := copy.Copy(folderToWatch, folderToCopy)
					if copyError != nil {
						log.Printf(ErrorColor, copyError)
					}
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf(DebugColor, "created file: "+ event.Name)
					copyError := copy.Copy(folderToWatch, folderToCopy)
					if copyError != nil {
						log.Printf(ErrorColor, copyError)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf(ErrorColor, err)
			}
		}
	}()

	err = watcher.Add(folderToWatch)
	
	if err != nil {
		log.Fatal(err)
	}

	<-done
}

func customExitHanlder() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf(WarningColor,"\r- made by x0s3\n")
		os.Exit(0)
	}()
}