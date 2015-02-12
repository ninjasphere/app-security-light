package service

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ninjasphere/astrotime"
	"github.com/ninjasphere/go-ninja/bus"
	"github.com/ninjasphere/go-ninja/logger"
)

type securityLight struct {
	config        SecurityLightConfig
	log           *logger.Logger
	subscriptions map[string]*bus.Subscription
	timeout       *time.Timer
}

func newSecurityLight(config SecurityLightConfig) (*securityLight, error) {

	l := &securityLight{
		config:        config,
		log:           logger.GetLogger("[Name: " + config.Name + "]"),
		subscriptions: make(map[string]*bus.Subscription),
	}

	l.timeout = time.AfterFunc(time.Minute, l.turnOffLights)
	l.timeout.Stop()

	l.log.Infof("start(%+v)", config)

	// start listening to sensors
	for _, sensorID := range config.Sensors {
		sub, err := listenToSensor(sensorID, l.onSensor)
		if err == nil {
			l.subscriptions[sensorID] = sub
		} else {
			l.log.Warningf("Failed to subscribe to sensor %s: %s", sensorID, err)
		}
	}

	return l, nil
}

func (l *securityLight) destroy() {
	l.log.Infof("destroy()")
	l.timeout.Stop()
	for _, s := range l.subscriptions {
		s.Cancel()
	}
	l.subscriptions = nil
}

func (l *securityLight) turnOnLights() {
	l.log.Infof("turnOnLights() timeout %s", time.Second*time.Duration(l.config.Timeout))

	// TODO: Turn on the lights
	l.timeout.Reset(time.Second * time.Duration(l.config.Timeout))
}

func (l *securityLight) turnOffLights() {
	l.log.Infof("turnOffLights()")

	// TODO: Turn off the lights
}

func (l *securityLight) onSensor(id string) {

	if !l.isActiveNow() {
		l.log.Debugf("Sensor activated: %s (but light not active)", id)
		return
	}

	l.log.Infof("Sensor activated: %s", id)
	l.turnOnLights()
	// TODO: Turn on lights
}

func (l *securityLight) isActiveNow() bool {
	if !l.config.Enabled {
		return false
	}

	if l.config.TimeStart == "" {
		return true
	}

	start, err := getTime(l.config.TimeStart)
	if err != nil {
		l.log.Warningf("Failed to parse start time (allowing): %s", err)
		return true
	}

	end, err := getTime(l.config.TimeEnd)
	if err != nil {
		l.log.Warningf("Failed to parse end time (allowing): %s", err)
		return true
	}

	spew.Dump("start", l.config.TimeStart, start, "end", l.config.TimeEnd, end)

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
		fmt.Errorf("TODO: Parse named times! : %s - %s", s, err)
	}

	// Get that time, today
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), tmp.Hour(), tmp.Minute(), 0, 0, time.Now().Location()), nil
}
