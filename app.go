package main

import (
	"os"
	"time"

	"github.com/ninjasphere/app-security-light/service"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/support"
)

var info = ninja.LoadModuleInfo("./package.json")

// This object is persisted by HomeCloud, and provided when the app starts.
type Config struct {
	Lights []service.SecurityLightConfig `json:"lights"`
}

type App struct {
	support.AppSupport
}

func (a *App) Start(config *Config) error {
	return service.Start(config.Lights, a.Conn, func(config []service.SecurityLightConfig) {
		a.SendEvent("config", Config{config})
	})
}

// Stop the security light app.
func (a *App) Stop() error {
	// Can't really stop at the moment, so just bomb out.
	a.Log.Infof("Stop called. Quitting in 3 sec...")
	go func() {
		time.Sleep(time.Second * 3)
		os.Exit(0)
	}()
	return nil
}

func main() {
	app := &App{}
	err := app.Init(info)
	if err != nil {
		app.Log.Fatalf("failed to initialize app: %v", err)
	}

	err = app.Export(app)
	if err != nil {
		app.Log.Fatalf("failed to export app: %v", err)
	}

	support.WaitUntilSignal()
}
