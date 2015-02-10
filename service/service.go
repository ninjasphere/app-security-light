package service

import (
	"fmt"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/bus"
	"github.com/ninjasphere/go-ninja/logger"
)

var config []SecurityLightConfig
var conn *ninja.Connection
var saveConfig func([]SecurityLightConfig)

var log = logger.GetLogger("service")
var lights map[string]*securityLight

var started bool

func Start(config1 []SecurityLightConfig, conn1 *ninja.Connection, saveConfig1 func([]SecurityLightConfig)) error {

	if config != nil {
		return fmt.Errorf("Already started!")
	}

	// If you're ugly and you know it clap your hands
	config = config1
	conn = conn1
	saveConfig = saveConfig1
	// *clap*
	// *clap*

	log.Infof("start()")
	return nil
}

func listenToSensor(deviceId string, callback func(string)) (*bus.Subscription, error) {

}
