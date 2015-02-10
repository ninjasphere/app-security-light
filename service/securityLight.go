package service

import (
	"github.com/ninjasphere/go-ninja/bus"
	"github.com/ninjasphere/go-ninja/logger"
)

type securityLight struct {
	config        SecurityLightConfig
	log           *logger.Logger
	subscriptions map[string]*bus.Subscription
}

func (l *securityLight) updateConfig(config SecurityLightConfig) error {

	l.log = logger.GetLogger("[Name: " + config.Name + "]")
	l.log.Infof("start(%+v)", config)
	l.config = config

	if l.subscriptions == nil {
		l.subscriptions = make(map[string]*bus.Subscription)
	}

	// start listening to sensors

	return nil
}

func (l *securityLight) onSensor(id string) {
	// When the sensor goes off...

}

func (l *securityLight) isActiveNow() bool {
	if !l.config.Enabled {
		return false
	}

	if l.config.TimeStart == "" {
		return true
	}

	l.log.Warningf("TODO: Check the time!")
	return true
}
