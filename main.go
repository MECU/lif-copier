package main

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	// Load the key
	key, err := ioutil.ReadFile("./tractrak.key")
	if err != nil {
		panic(err)
	}
	log.Println("Got the Key")

	// Setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	log.Println("Ready to rock this meet ...")
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					log.Print("created/modified file:", event.Name)
					// See if it's a .LIF file
					// TODO: lynx.sch, .ppl, evt
					fileExtension := event.Name[len(event.Name)-3:]
					if strings.ToUpper(fileExtension) == "LIF" {
						log.Print(" ... trying to upload")
						err = Upload("https://tractrak.com/api/meet-file", event.Name, key)
						if err != nil {
							log.Println(" ... upload failed: ", err)
						} else {
							log.Println(" ... upload success")
						}
					}
				}
			case err := <-watcher.Errors:
				log.Println("error: ", err)
			}
		}
	}()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
