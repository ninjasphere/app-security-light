package service

import (
	"fmt"
	"math/rand"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/bus"
	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/logger"
)

var mocking = config.Bool(false, "mock")

// This isn't an example of good code. Fix it or ignore it, but don't copy it.

var lightsConfig map[string]SecurityLightConfig
var conn *ninja.Connection
var saveConfig func()

var log = logger.GetLogger("service")
var lights = make(map[string]*securityLight)

var started bool

var latitude, longitude float64

type sensor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type light struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func Start(config []SecurityLightConfig, conn1 *ninja.Connection, saveConfig1 func([]SecurityLightConfig)) error {

	if lightsConfig != nil {
		return fmt.Errorf("Already started!")
	}

	lightsConfig = make(map[string]SecurityLightConfig)

	for _, c := range config {
		onConfigUpdated(c)
	}

	// If you're ugly and you know it clap your hands
	conn = conn1
	// *clap*
	// *clap*

	if mocking {
		latitude, longitude = -33.86, -151.20
	}

	saveConfig = func() {

		cfg := []SecurityLightConfig{}

		for _, c := range lightsConfig {
			cfg = append(cfg, c)
		}

		saveConfig1(cfg)
	}

	log.Infof("start()")

	go func() {
		err := startRestServer()
		if err != nil {
			log.Fatalf("Failed to start rest server: %s", err)
		}
	}()

	return nil
}

func onConfigUpdated(cfg SecurityLightConfig) {
	log.Infof("Config updated: " + cfg.ID)

	lightsConfig[cfg.ID] = cfg

	// If we have previously created this security light, destroy it.
	if light, ok := lights[cfg.ID]; ok {
		light.destroy()
	}

	// Create and start the new security light
	light, err := newSecurityLight(cfg)

	if err != nil {
		log.Fatalf("Failed to create security light %s: %s", cfg.ID, err)
	}
	lights[cfg.ID] = light

	saveConfig()
}

func saveSecurityLight(cfg *SecurityLightConfig) error {
	if cfg.ID == "" {
		// TODO: Fix me
		cfg.ID = fmt.Sprintf("%d", rand.Intn(99999))

		cfg.Enabled = true
	}

	onConfigUpdated(*cfg)
	return nil
}

func deleteSecurityLight(id string) error {
	l, ok := lights[id]
	if ok {
		l.destroy()
		delete(lightsConfig, id)
		delete(lights, id)
		saveConfig()
	}

	return nil
}

func getSensors() ([]sensor, error) {
	if mocking {
		return mockSensors, nil
	}
	return nil, nil
}

func getLights() ([]light, error) {
	if mocking {
		return mockLights, nil
	}

	return nil, nil
}

func listenToSensor(thingID string, callback func(thingID string)) (*bus.Subscription, error) {

	if mocking {
		return mockSensor(thingID, callback)
	}
	return nil, nil
}
