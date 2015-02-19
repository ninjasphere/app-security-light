package service

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/suit"
)

type configService struct {
}

func (c *configService) error(message string) (*suit.ConfigurationScreen, error) {

	return &suit.ConfigurationScreen{
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.Alert{
						Title:        "Error",
						Subtitle:     message,
						DisplayClass: "danger",
					},
				},
			},
		},
		Actions: []suit.Typed{
			suit.ReplyAction{
				Label: "Cancel",
				Name:  "list",
			},
		},
	}, nil
}
func (c *configService) list() (*suit.ConfigurationScreen, error) {

	var lights []suit.ActionListOption

	for _, light := range lightsConfig {
		lights = append(lights, suit.ActionListOption{
			Title:    light.Name,
			Subtitle: fmt.Sprintf("%d sensors, %d lights", len(light.Sensors), len(light.Lights)),
			Value:    light.ID,
		})
	}

	screen := suit.ConfigurationScreen{
		Title: "Security Lights",
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.ActionList{
						Name:    "light",
						Options: lights,
						PrimaryAction: &suit.ReplyAction{
							Name:        "edit",
							DisplayIcon: "pencil",
						},
						SecondaryAction: &suit.ReplyAction{
							Name:         "delete",
							Label:        "Delete",
							DisplayIcon:  "trash",
							DisplayClass: "danger",
						},
					},
				},
			},
		},
		Actions: []suit.Typed{
			suit.CloseAction{
				Label: "Close",
			},
			suit.ReplyAction{
				Label:        "New Security Light",
				Name:         "new",
				DisplayClass: "success",
				DisplayIcon:  "star",
			},
		},
	}

	return &screen, nil
}

func (c *configService) Configure(request *model.ConfigurationRequest) (*suit.ConfigurationScreen, error) {
	spew.Dump("configure()", request)

	switch request.Action {
	case "list":
		fallthrough
	case "":
		return c.list()
	case "new":
		return c.edit(SecurityLightConfig{
			Timeout: 5,
			Time: suit.TimeRange{
				From: "sunset",
				To:   "sunrise",
			},
		})
	case "edit":

		var vals map[string]string
		json.Unmarshal(request.Data, &vals)
		config, ok := lightsConfig[vals["light"]]

		if !ok {
			return c.error(fmt.Sprintf("Could not find light with id: %s", vals["light"]))
		}

		return c.edit(config)
	case "delete":

		var vals map[string]string
		json.Unmarshal(request.Data, &vals)
		deleteSecurityLight(vals["light"])

		return c.list()
	case "save":
		var lightConfig SecurityLightConfig
		err := json.Unmarshal(request.Data, &lightConfig)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal save config request %s: %s", request.Data, err)
		}

		spew.Dump("Got security light config", lightConfig)

		saveSecurityLight(&lightConfig)
		return c.list()

	default:
		return c.error(fmt.Sprintf("Unknown action: %s", request.Action))
	}
}

func (c *configService) edit(config SecurityLightConfig) (*suit.ConfigurationScreen, error) {

	var sensorOptions []suit.OptionGroupOption
	sensors, err := getSensors()
	if err != nil {
		return c.error(fmt.Sprintf("Could not find sensors: %s", err))
	}

	for _, s := range sensors {
		sensorOptions = append(sensorOptions, suit.OptionGroupOption{
			Title:    s.Name,
			Value:    s.ID,
			Selected: contains(config.Sensors, s.ID),
		})
	}

	var lightOptions []suit.OptionGroupOption
	lights, err := getLights()
	if err != nil {
		return c.error(fmt.Sprintf("Could not find lights: %s", err))
	}

	for _, s := range lights {
		lightOptions = append(lightOptions, suit.OptionGroupOption{
			Title:    s.Name,
			Value:    s.ID,
			Selected: contains(config.Lights, s.ID),
		})
	}

	title := "New Security Light"
	if config.ID == "" {
		title = "Edit Security Light"
	}

	screen := suit.ConfigurationScreen{
		Title: title,
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.InputHidden{
						Name:  "id",
						Value: config.ID,
					},
					suit.InputText{
						Name:        "name",
						Before:      "Name",
						Placeholder: "My Security Light",
						Value:       config.Name,
					},
					suit.OptionGroup{
						Name:           "sensors",
						Title:          "When these devices detect motion",
						MinimumChoices: 1,
						Options:        sensorOptions,
					},
					suit.OptionGroup{
						Name:           "lights",
						Title:          "Turn on these lights",
						MinimumChoices: 1,
						Options:        lightOptions,
					},
					suit.InputTimeRange{
						Name:  "time",
						Title: "When",
						Value: suit.TimeRange{
							From: config.Time.From,
							To:   config.Time.To,
						},
					},
					suit.InputText{
						Title:     "Turn off again after",
						After:     "minutes",
						Name:      "timeout",
						InputType: "number",
						Minimum:   i(0),
						Value:     config.Timeout,
					},
				},
			},
		},
		Actions: []suit.Typed{
			suit.CloseAction{
				Label: "Cancel",
			},
			suit.ReplyAction{
				Label:        "Save",
				Name:         "save",
				DisplayClass: "success",
				DisplayIcon:  "star",
			},
		},
	}

	return &screen, nil
}

func i(i int) *int {
	return &i
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
