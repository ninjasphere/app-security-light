package service

import (
	"math/rand"
	"time"

	"github.com/ninjasphere/go-ninja/bus"
)

var mockSensors = []sensor{
	sensor{"mysensor1", "My Sensor1"},
	sensor{"mysensor2", "My Sensor2"},
	sensor{"mysensor3", "My Sensor3"},
	sensor{"mysensor4", "My Sensor4"},
	sensor{"mysensor5", "My Sensor5"},
	sensor{"mysensor6", "My Sensor6"},
}

var mockLights = []light{
	light{"mylight1", "My Light1"},
	light{"mylight2", "My Light2"},
	light{"mylight3", "My Light3"},
	light{"mylight4", "My Light4"},
	light{"mylight5", "My Light5"},
	light{"mylight6", "My Light6"},
}

// Randomly fires the sensors every 3-8 seconds
func mockSensor(thingID string, callback func(thingID string)) (*bus.Subscription, error) {

	var cancelled bool
	sub := &bus.Subscription{
		Cancel: func() {
			cancelled = true
		},
	}

	go func() {
		for {
			if cancelled {
				break
			}
			time.Sleep(time.Second * time.Duration((rand.Intn(5) + 3)))
			callback(thingID)
		}
	}()

	return sub, nil
}
