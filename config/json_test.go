// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"strings"
	"testing"
)

func Test_SerializeConfig_NoErrorIsReturned(t *testing.T) {
	// arrange
	writeBuffer := new(bytes.Buffer)

	config := &Config{
		Server: Server{
			ThemeFolder: "/some/folder",
			Http: Http{
				Port: 80,
			},
		},
	}

	serializer := JSONSerializer{}

	// act
	err := serializer.SerializeConfig(writeBuffer, config)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("The serialization of the config object return an error. %s", err)
	}
}

func Test_SerializeConfig_JsonContainsConfigValues(t *testing.T) {
	// arrange
	writeBuffer := new(bytes.Buffer)

	config := &Config{
		Server: Server{
			ThemeFolder: "/some/folder",
			Http: Http{
				Port: 80,
			},
		},
	}

	serializer := JSONSerializer{}

	// act
	serializer.SerializeConfig(writeBuffer, config)

	// assert
	json := writeBuffer.String()

	// assert: json contains theme folder
	if !strings.Contains(json, config.Server.ThemeFolder) {
		t.Fail()
		t.Logf("The produced json does not contain the 'ThemeFolder' value %q. The produced JSON is this: %s", config.Server.ThemeFolder, json)
	}

	// assert: json contains http port
	if !strings.Contains(json, config.Server.Http.Port.String()) {
		t.Fail()
		t.Logf("The produced json does not contain the 'Http Port' value %q. The produced JSON is this: %s", config.Server.Http.Port, json)
	}
}

func Test_SerializeConfig_JsonIsFormatted(t *testing.T) {
	// arrange
	writeBuffer := new(bytes.Buffer)

	config := &Config{
		Server: Server{
			ThemeFolder: "/some/folder",
			Http: Http{
				Port: 80,
			},
		},
	}

	serializer := JSONSerializer{}

	// act
	serializer.SerializeConfig(writeBuffer, config)

	// assert
	json := writeBuffer.String()

	// assert: json contains theme folder
	if !strings.Contains(json, "\n") {
		t.Fail()
		t.Logf("The produced json does not seem to be formatted. The produced JSON is this: %s", json)
	}
}

func Test_DeserializeConfig_EmptyObjectString_NoErrorIsReturned(t *testing.T) {
	// arrange
	json := `{}`
	jsonReader := bytes.NewBuffer([]byte(json))

	serializer := JSONSerializer{}

	// act
	_, err := serializer.DeserializeConfig(jsonReader)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("The deserialization of %q should not produce an error. But it did produce this error: %s", json, err)
	}
}

func Test_DeserializeConfig_FullConfigString_AllFieldsAreSet(t *testing.T) {
	// arrange
	json := `{
		"Server": {
			"ThemeFolder": "/some/folder",
			"Http": {
				"Port": 80
			}
		}
	}`
	jsonReader := bytes.NewBuffer([]byte(json))

	serializer := JSONSerializer{}

	// act
	config, _ := serializer.DeserializeConfig(jsonReader)

	// assert: Theme folder
	if config.Server.ThemeFolder == "" {
		t.Fail()
		t.Logf("The deserialized config object should have the %q field properly initialized. Deserialization result: %#v", "ThemeFolder", config)
	}

	// assert: http port
	if config.Server.Http.Port == 0 {
		t.Fail()
		t.Logf("The deserialized config object should have the %q field properly initialized. Deserialization result: %#v", "Http.Port", config)
	}
}

func Test_DeserializeConfig_ObjectWithDifferentFields_ConfigWithDefaultValuesIsReturned(t *testing.T) {
	// arrange
	json := `{
		"Name": "Ladi da",
		"AnotherField": {
		},
		"SomeList": [ "1", "2", "3" ]
	}
	`
	jsonReader := bytes.NewBuffer([]byte(json))

	serializer := JSONSerializer{}

	// act
	config, _ := serializer.DeserializeConfig(jsonReader)

	// assert
	emptyConfig := Config{}
	if config.Server.ThemeFolder != emptyConfig.Server.ThemeFolder {
		t.Fail()
		t.Logf("When the JSON cannot be mapped to the Config type the deserializer should return an uninitialized config object.")
	}
}

func Test_DeserializeConfig_EmptyString_ErrorIsReturned(t *testing.T) {
	// arrange
	json := ""
	jsonReader := bytes.NewBuffer([]byte(json))

	serializer := JSONSerializer{}

	// act
	_, err := serializer.DeserializeConfig(jsonReader)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("DeserializeConfig should return an error if supplied JSON is invalid")
	}
}

func Test_DeserializeConfig_InvalidJson_ErrorIsReturned(t *testing.T) {
	// arrange
	json := `dsajdklasdj/(/)(=7897402
		38748902
		;;;
	`
	jsonReader := bytes.NewBuffer([]byte(json))

	serializer := JSONSerializer{}

	// act
	_, err := serializer.DeserializeConfig(jsonReader)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("DeserializeConfig should return an error if supplied JSON is invalid")
	}
}