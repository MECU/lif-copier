package main

import (
	"log"
	"fmt"
	"net/http"
	"github.com/go-fsnotify/fsnotify"
	"strings"
	"bytes"
	"mime/multipart"
	"io/ioutil"
	"time"
)

func OpenFile(file string) (fileContents []byte) {
	// Keep trying to read the file until there is not an error
	for ok := true; ok == true; ok = true {
		fileContents, err := ioutil.ReadFile(file)
		if err == nil {
			return fileContents
		}
		log.Println(" ... read file failed (try again in half-a-second): ", err)
		time.Sleep(500)
	}

	return fileContents
}

func Upload(url, file string, key []byte) (err error) {
	// Prepare a form that you will submit to that URL.

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	if err != nil {
		return
	}

	// Filename

	fw, err := w.CreateFormField("filename")
	if err != nil {
		return
	}
	if _, err = fw.Write([]byte(file)); err != nil {
		return
	}

	// File contents
	if fw, err = w.CreateFormField("file"); err != nil {
		return
	}

	if _, err = fw.Write(OpenFile(file)); err != nil {
		return
	}

	// Add the other fields
	if fw, err = w.CreateFormField("key"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(key)); err != nil {
		return
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(" ... upload failed: ", err)
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

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
				if event.Op & fsnotify.Write == fsnotify.Write || event.Op & fsnotify.Create == fsnotify.Create {
					log.Print("created/modified file:", event.Name)
					// See if it's a .LIF file
					// TODO: lynx.sch, .ppl, evt
					fileExtension := event.Name[len(event.Name) - 3:]
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