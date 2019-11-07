package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getcasa/plugin-philipshue/devices"
	"github.com/getcasa/sdk"
)

func main() {}

var Config = sdk.Configuration{
	Name:        "philipshue",
	Version:     "1.0.0",
	Author:      "casa",
	Description: "Control Philips Hue ecosystem",
	Main:        "",
	FuncData:    "",
	Discover:    true,
	Triggers:    []sdk.Trigger{},
	Actions: []sdk.Action{
		sdk.Action{
			Name: "switchLight",
			Fields: []sdk.Field{
				sdk.Field{
					Name:   "On",
					Type:   "bool",
					Config: true,
				},
				sdk.Field{
					Name:   "Sat",
					Type:   "int",
					Config: true,
				},
				sdk.Field{
					Name:   "Bri",
					Type:   "int",
					Config: true,
				},
				sdk.Field{
					Name:   "Hue",
					Type:   "int",
					Config: true,
				},
			},
		},
		sdk.Action{
			Name:   "toggleLight",
			Fields: []sdk.Field{},
		},
	},
}

// State define each element of the global state
type State struct {
	Bridge   Bridge
	Device   devices.LCT0152A19ECLv5
	DeviceID int
}

type savedConfig struct {
	BridgeID string
	Username string
}

// States is the global state of plugin
var States []State
var client http.Client

// Init plugin config
func Init() []byte {
	bridges := DiscoverBridges()
	var savedConfigs []savedConfig

	for _, bridge := range bridges {
		savedConfigs = append(savedConfigs, savedConfig{
			BridgeID: bridge.ID,
			Username: bridge.CreateUser(),
		})
	}

	config, _ := json.Marshal(savedConfigs)

	return config
}

// OnStart create http client
func OnStart(config []byte) {
	var savedConfigs []savedConfig

	if err := json.Unmarshal(config, &savedConfigs); err != nil {
		panic(err)
	}

	bridges := DiscoverBridges()

	for _, savedConfig := range savedConfigs {
		for i, bridge := range bridges {
			if savedConfig.BridgeID != bridge.ID {
				continue
			}
			bridges[i].Username = savedConfig.Username
			for j, light := range bridges[i].GetLights() {
				States = append(States, State{
					Bridge:   bridges[i],
					Device:   light,
					DeviceID: j + 1,
				})
			}
		}
	}
}

// Params define actions parameters available
type Params struct {
	ID  string
	On  bool
	Sat int
	Bri int
	Hue int
}

// CallAction call functions from actions
func CallAction(name string, params []byte, config []byte) {
	if string(params) == "" {
		fmt.Println("Params must be provided")
		return
	}

	// declare parameters
	var req Params

	// unmarshal parameters to use in actions
	err := json.Unmarshal(params, &req)
	if err != nil {
		fmt.Println(err)
	}

	bridge := GetBridge(req.ID)
	if bridge.ID == "" {
		return
	}

	// use name to call actions
	switch name {
	case "switchLight":
		bridge.SwitchLight(req)
	case "toggleLight":
		bridge.ToggleLight(req)
	default:
		return
	}
}

// OnStop close connection
func OnStop() {
}
