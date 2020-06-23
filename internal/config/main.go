package config

import (
	"encoding/json"
	"errors"
	"github.com/codemicro/nota/internal/models"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
)

var (
	Settings models.Settings
)

func LoadSettings() error {
	settingsFromFile, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		return err
	}
	var unmarshalled map[string]interface{}
	err = json.Unmarshal(settingsFromFile, &unmarshalled)
	if err != nil {
		return err
	}

	// Check for missing items from the settings struct
	fields := Settings.GetFields()
	for _, v := range fields {
		if _, ok := unmarshalled[v]; !ok {
			// field not in map
			return errors.New("field " + v + " not in settings JSON.")
		}
	}

	// Cast to struct
	err = mapstructure.Decode(unmarshalled, &Settings)
	return err
}
