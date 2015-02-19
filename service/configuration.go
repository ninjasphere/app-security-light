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

func (c *configService) Configure(request *model.ConfigurationRequest) (*suit.ConfigurationScreen, error) {
	spew.Dump("configure()", request)

	if request.Action == "save" {
		var lightConfig SecurityLightConfig
		err := json.Unmarshal(request.Data, &lightConfig)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal save config request %s: %s", request.Data, err)
		}

		spew.Dump("Got security light config", lightConfig)
	}

	screen := suit.ConfigurationScreen{
		Title: "Edit Security Light",
		Sections: []suit.Section{
			suit.Section{
				Contents: []suit.Typed{
					suit.InputText{
						Name:        "name",
						Before:      "Name",
						Placeholder: "My Security Light",
						Value:       "Front door light",
					},
					suit.OptionGroup{
						Name:           "sensors",
						Title:          "When these devices detect motion",
						MinimumChoices: 1,
						Options: []suit.OptionGroupOption{
							suit.OptionGroupOption{
								Title:    "Front Door Motion",
								Subtitle: "Motion",
								Value:    "fd",
							},
							suit.OptionGroupOption{
								Title:    "Back Door 1",
								Subtitle: "Presence",
								Value:    "bd",
								Selected: true,
							},
						},
					},
					suit.OptionGroup{
						Name:           "lights",
						Title:          "Turn on these lights",
						MinimumChoices: 1,
						Options: []suit.OptionGroupOption{
							suit.OptionGroupOption{
								Title:    "Front Door",
								Subtitle: "Lamp in Hallway",
								Value:    "fd",
							},
							suit.OptionGroupOption{
								Title:    "Front Door Spotlight",
								Subtitle: "Light in Front Step",
								Value:    "fds",
								Selected: true,
							},
						},
					},
					suit.InputTimeRange{
						Name:  "time",
						Title: "When",
						Value: suit.TimeRange{
							From: "10:00",
							To:   "sunset",
						},
					},
					suit.InputText{
						Title:     "Turn off again after",
						After:     "minutes",
						Name:      "timeout",
						InputType: "number",
						Minimum:   i(0),
						Value:     5,
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
