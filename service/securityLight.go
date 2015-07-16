package service

import (
	"fmt"
	"time"

	"github.com/ninjasphere/astrotime"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/bus"
	"github.com/ninjasphere/go-ninja/logger"
)

type securityLight struct {
	config        SecurityLightConfig
	log           *logger.Logger
	subscriptions map[string]*bus.Subscription
	timeout       *time.Timer
	onOffServices []*ninja.ServiceClient
}

func newSecurityLight(config SecurityLightConfig) (*securityLight, error) {

	l := &securityLight{
		config:        config,
		log:           logger.GetLogger("[Name: " + config.Name + "]"),
		subscriptions: make(map[string]*bus.Subscription),
		onOffServices: []*ninja.ServiceClient{},
	}

	l.timeout = time.AfterFunc(time.Minute, l.turnOffLights)
	l.timeout.Stop()

	l.log.Infof("start(%+v)", config)

	// start listening to sensors
	for _, sensorID := range config.Sensors {
		sub, err := listenToSensor(sensorID, l.onSensor)
		if err == nil {
			if sub == nil {
				l.log.Fatalf("Got a nil subscription back for sensor id: %s", sensorID)
			} else {
				l.subscriptions[sensorID] = sub
			}
		} else {
			l.log.Warningf("Failed to subscribe to sensor %s: %s", sensorID, err)
		}
	}

	// Grab the on-off channels for the lights
	for _, lightID := range config.Lights {
		service, err := getOnOffChannelClient(lightID)
		if err == nil {
			l.onOffServices = append(l.onOffServices, service)
		} else {
			l.log.Warningf("Failed to get light %s on-off service: %s", lightID, err)
		}
	}

	return l, nil
}

func (l *securityLight) destroy() {
	l.log.Infof("destroy()")
	l.timeout.Stop()
	if l.subscriptions != nil {
		for _, s := range l.subscriptions {
			s.Cancel()
		}
	}
	l.subscriptions = nil
}

func (l *securityLight) turnOnLights() {
	l.log.Infof("turnOnLights()")

	l.timeout.Reset(time.Minute * time.Duration(l.config.Timeout))

	for _, channel := range l.onOffServices {
		channel.Call("turnOn", nil, nil, 0)
	}
}

func (l *securityLight) turnOffLights() {
	l.log.Infof("turnOffLights()")

	for _, channel := range l.onOffServices {
		channel.Call("turnOff", nil, nil, 0)
	}
}

func (l *securityLight) onSensor(id string) {

	if !l.isActiveNow() {
		l.log.Debugf("Sensor activated: %s (but light not active)", id)
		return
	}

	l.log.Infof("Sensor activated: %s", id)
	l.turnOnLights()
}

func (l *securityLight) isActiveNow() bool {
	if !l.config.Enabled {
		return false
	}

	if l.config.Time.From == "" {
		return true
	}

	if l.config.Time.From == l.config.Time.To {
		return true
	}

	start, err := getTime(l.config.Time.From)
	if err != nil {
		l.log.Warningf("Failed to parse start time (allowing): %s", err)
		return true
	}

	end, err := getTime(l.config.Time.To)
	if err != nil {
		l.log.Warningf("Failed to parse end time (allowing): %s", err)
		return true
	}

	//spew.Dump("start", l.config.TimeStart, start, "end", l.config.TimeEnd, end)

	if end.Before(start) {
		return time.Now().After(start) || time.Now().Before(end)
	}

	return time.Now().After(start) && time.Now().Before(end)
}

func getTime(s string) (time.Time, error) {

	midnight, _ := parseTimeToday("00:00")

	switch s {
	case "dawn":
		return astrotime.NextDawn(midnight, latitude, longitude, astrotime.CIVIL_DAWN), nil
	case "sunrise":
		return astrotime.NextSunrise(midnight, latitude, longitude), nil
	case "sunset":
		return astrotime.NextSunset(midnight, latitude, longitude), nil
	case "dusk":
		return astrotime.NextDusk(midnight, latitude, longitude, astrotime.CIVIL_DUSK), nil
	}

	return parseTimeToday(s)
}

func parseTimeToday(s string) (time.Time, error) {
	tmp, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("Couldn't parse time! : %s - %s", s, err)
	}

	// Get that time, today
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), tmp.Hour(), tmp.Minute(), 0, 0, time.Now().Location()), nil
}
