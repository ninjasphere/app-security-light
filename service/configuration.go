package service

import (
	"encoding/json"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/ninjasphere/app-security-light/service"
	"github.com/ninjasphere/go-ninja/model"
)

var configService = config{}

type config struct {
}

func (c *config) Configure(request *model.ConfigurationRequest) (*map[string]interface{}, error) {
	spew.Dump("configure()", request)

	if request.Action == "save" {
		var lightConfig service.SecurityLightConfig
		err := json.Unmarshal(request.Data, &lightConfig)
		if err != nil {
			log.Fatalf("Failed to unmarshal save config request %s: %s", request.Data, err)
		}

		spew.Dump("Got security light config", lightConfig)
	}

	var reply map[string]interface{}

	err := json.Unmarshal([]byte(`
    {
      "title": "Create new security light",
      "sections": [
        {
          "contents": [
            {
              "type": "inputText",
              "name": "name",
              "before": "Name",
              "placeholder": "My Security Light",
              "value": "Front door light"
            },
            {
              "type": "optionGroup",
              "name": "sensors",
              "title": "When these devices detect motion",
              "minimumChoices": 1,
              "options": [
                {
                  "title": "Front Door Motion",
                  "subtitle": "Motion",
                  "value": "fd"
                },
                {
                  "title": "Back Door 1",
                  "subtitle": "Presence",
                  "value": "bd",
                  "selected": true
                },
                {
                  "title": "Back Door Webcam",
                  "subtitle": "Camera",
                  "value": "bdc",
                  "selected": true
                }
              ]
            },
            {
              "type": "optionGroup",
              "name": "lights",
              "title": "Turn on these lights",
              "minimumChoices": 1,
              "options": [
                {
                  "title": "Front Door",
                  "subtitle": "Lamp in Hallway",
                  "value": "fd"
                },
                {
                  "title": "Front Door Spotlight",
                  "subtitle": "Light in Front Step",
                  "value": "fds"
                },
                {
                  "title": "Above Fridge",
                  "subtitle": "Lamp in Kitchen",
                  "value": "kl"
                },
                {
                  "title": "Broken",
                  "subtitle": "Light in Backyard",
                  "value": "bdf",
                  "selected": true
                }
              ]
            },
            {
              "type": "inputTimeRange",
              "name": "time",
              "title": "When",
              "value": {
                "from": "10:00",
                "to": "sunset"
              }
            },
            {
              "title": "Turn off again after",
              "type": "inputText",
              "after": "minutes",
              "name": "timeout",
              "inputType": "number",
              "minimum": 0,
              "value": 5
            }
          ]
        }
      ],
      "actions": [
        {
          "type": "close",
          "label": "Cancel"
        },
        {
          "type": "reply",
          "label": "Save",
          "name": "save",
          "displayClass": "success",
          "displayIcon": "star"
        }
      ]
    }
  `), &reply)

	return &reply, err
}
