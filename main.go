package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("Unexpected arguments %v\n", args)
	}

	displayedContent, err := getDisplayedContent(args[0])
	if err != nil {
		log.Fatalf("%v", err)
	}

	keys := keyboard.NewDriver()

	itr := 0
	length := len(displayedContent)
	work := func() {
		keys.On(keyboard.Key, func(data interface{}) {
			fmt.Printf(displayedContent[itr : itr+1])
			if itr == length-1 {
				itr = 0
			} else {
				itr++
			}
		})
	}

	log.SetOutput(ioutil.Discard)
	robot := gobot.NewRobot("keyboardbot",
		[]gobot.Connection{},
		[]gobot.Device{keys},
		work,
	)

	robot.Start()
	log.SetOutput(os.Stderr)
}

func getDisplayedContent(src string) (content string, err error) {
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		content, err = getHttpResponseBody(src)
		if err != nil {
			return
		}
	} else {
		content, err = readFile(src)
		if err != nil {
			return
		}
	}
	return
}

func getHttpResponseBody(url string) (content string, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	content = string(body)
	return
}

func readFile(path string) (content string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	content = string(b)

	return
}
