package main

import (
	"Main/bowling"
	"Main/web"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"
)

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("")
	}
	if err != nil {
		log.Printf("%v", err)
	}
}

func main() {
	cfg := bowling.Config{
		NumLanes:          3,
		NumClients:        20,
		MaxClientWaitTime: 3 * time.Second,

		ClientArrival: bowling.DurationConfig{
			Base:    500 * time.Millisecond,
			Variant: 1000 * time.Millisecond,
		},

		GameDuration: bowling.DurationConfig{
			Base:    4 * time.Second,
			Variant: 2 * time.Second,
		},
	}

	state := bowling.NewState(cfg.NumLanes)
	center := bowling.New(cfg, state)
	server := web.New(state)

	go center.Run()
	go openBrowser("http://localhost:8080")

	if err := server.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
