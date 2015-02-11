package service

import (
	"fmt"
	"math/rand"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/bus"
	"github.com/ninjasphere/go-ninja/logger"
)

// This isn't an example of good code. Fix it or ignore it, but don't copy it.

var lightsConfig map[string]SecurityLightConfig
var conn *ninja.Connection
var saveConfig func()

var log = logger.GetLogger("service")
var lights = make(map[string]*securityLight)

var started bool

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

	light, ok := lights[cfg.ID]
	if !ok {
		light = &securityLight{}
		lights[cfg.ID] = light
	}

	light.updateConfig(cfg)

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
	return []sensor{
		sensor{"mysensor1", "My Sensor1"},
		sensor{"mysensor2", "My Sensor2"},
		sensor{"mysensor3", "My Sensor3"},
		sensor{"mysensor4", "My Sensor4"},
		sensor{"mysensor5", "My Sensor5"},
		sensor{"mysensor6", "My Sensor6"},
	}, nil
}

func getLights() ([]light, error) {
	return []light{
		light{"mylight1", "My Light1"},
		light{"mylight2", "My Light2"},
		light{"mylight3", "My Light3"},
		light{"mylight4", "My Light4"},
		light{"mylight5", "My Light5"},
		light{"mylight6", "My Light6"},
	}, nil
}

func listenToSensor(deviceId string, callback func(string)) (*bus.Subscription, error) {
	return nil, nil
}
