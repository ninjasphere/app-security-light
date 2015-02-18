package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/bus"
	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/logger"
	"github.com/ninjasphere/go-ninja/model"
)

var mocking = config.Bool(false, "mock")

// This isn't an example of good code. Fix it or ignore it, but don't copy it.

var lightsConfig map[string]SecurityLightConfig
var conn *ninja.Connection
var saveConfig func()
var thingModel *ninja.ServiceClient

var log = logger.GetLogger("service")
var lights = make(map[string]*securityLight)

var allThings []model.Thing

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

	// If you're ugly and you know it clap your hands
	conn = conn1
	// *clap*
	// *clap*

	conn.MustExportService(&configService{}, "$app/com.ninjablocks.securitylight/configure", &model.ServiceAnnouncement{
		Schema: "/protocol/configuration",
	})

	thingModel = conn.GetServiceClient("$home/services/ThingModel")

	lightsConfig = make(map[string]SecurityLightConfig)

	for _, c := range config {
		onConfigUpdated(c)
	}

	var err error
	allThings, err = getAllThings()
	if err != nil {
		return err
	}

	if mocking {
		latitude, longitude = -33.86, -151.20 // Sydney, AU
	} else {
		getSiteLocation()
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

func getSiteLocation() {
	siteModel := conn.GetServiceClient("$home/services/SiteModel")
	for {

		var site model.Site

		err := siteModel.Call("fetch", config.MustString("siteId"), &site, time.Second*5)

		if err == nil && site.Latitude != nil {
			latitude, longitude = *site.Latitude, *site.Longitude
			break
		}

		log.Infof("Failed to fetch siteid from sitemodel: %s", err)
		time.Sleep(time.Second * 5)
	}
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

func getAllThings() ([]model.Thing, error) {
	var things []model.Thing

	err := thingModel.Call("fetchAll", []interface{}{}, &things, time.Second*20)
	//err = client.Call("fetch", "c7ac05e0-9999-4d93-bfe3-a0b4bb5e7e78", &thing)

	if err != nil {
		return nil, fmt.Errorf("Failed to get things!: %s", err)
	}

	return things, nil
}

func getSensors() ([]sensor, error) {
	if mocking {
		return mockSensors, nil
	}

	sensors := []sensor{}

	for _, thing := range allThings {
		if thing.Device != nil {
			isSensor := len(thing.Device.GetChannelsByProtocol("motion")) > 0
			if isSensor {
				sensors = append(sensors, sensor{
					ID:   thing.ID,
					Name: thing.Name,
				})
			}
		}
	}

	return sensors, nil
}

func getLights() ([]light, error) {
	if mocking {
		return mockLights, nil
	}

	lights := []light{}

	for _, thing := range allThings {
		if thing.Type == "light" || thing.Type == "lamp" {
			hasOnOff := len(thing.Device.GetChannelsByProtocol("on-off")) > 0
			if hasOnOff {
				lights = append(lights, light{
					ID:   thing.ID,
					Name: thing.Name,
				})
			}
		}
	}

	return lights, nil
}

func listenToSensor(thingID string, callback func(thingID string)) (*bus.Subscription, error) {

	if mocking {
		return mockSensor(thingID, callback)
	}

	var sensor model.Thing

	err := thingModel.Call("fetch", []string{thingID}, &sensor, time.Second*20)

	if err != nil {
		return nil, err
	}

	//spew.Dump("got sensor", sensor, err)

	if sensor.Device == nil {
		return nil, fmt.Errorf("The sensor %s is not attached to a device", thingID)
	}

	motionChannels := sensor.Device.GetChannelsByProtocol("motion")

	if len(motionChannels) == 0 {
		return nil, fmt.Errorf("The sensor %s has no motion channels!", thingID)
	}

	// XXX: TODO: Just listen to the first for now...

	return conn.GetServiceClient(motionChannels[0].Topic).OnEvent("state", func() bool {
		callback(thingID)
		return true
	})
}

func getOnOffChannelClient(thingID string) (*ninja.ServiceClient, error) {
	var light model.Thing

	err := thingModel.Call("fetch", []string{thingID}, &light, time.Second*20)
	if err != nil {
		return nil, err
	}

	//spew.Dump("got light", light, err)

	if light.Device == nil {
		return nil, fmt.Errorf("The light %s is not attached to a device", thingID)
	}

	channels := light.Device.GetChannelsByProtocol("on-off")

	if len(channels) == 0 {
		return nil, fmt.Errorf("The light %s has no on-off channel!", thingID)
	}

	// XXX: TODO: Just use the first for now...
	return conn.GetServiceClient(channels[0].Topic), nil
}
