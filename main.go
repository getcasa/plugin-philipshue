package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getcasa/plugin-philipshue/devices"
	"github.com/getcasa/sdk"
)

func main() {}

// Config define the casa plugin configuration
var Config = sdk.Configuration{
	Name:        "philipshue",
	Version:     "1.0.0",
	Author:      "casa",
	Description: "Control Philips Hue ecosystem",
	Discover:    true,
	Devices: []sdk.Device{
		sdk.Device{
			Name:           "Philips-LCT015-2-A19ECLv5",
			DefaultTrigger: "",
			DefaultAction:  "toggleLight",
			Triggers:       []sdk.Trigger{},
			Actions:        []string{"switchLight", "toggleLight"},
		},
	},
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
var configPlugin []savedConfig

// Init plugin config
func Init() []byte {
	res, _ := json.Marshal([]savedConfig{})
	return res
}

//UpdateConfig update plugin's config if necessary
func UpdateConfig(config []byte) []byte {
	var savedConfigs []savedConfig

	if err := json.Unmarshal(config, &savedConfigs); err != nil {
		panic(err)
	}

	bridges := DiscoverBridges()

	for _, bridge := range bridges {
		indexBridge := findBridgeFromID(savedConfigs, bridge.ID)
		if indexBridge < 0 {
			savedConfigs = append(savedConfigs, savedConfig{
				BridgeID: bridge.ID,
				Username: bridge.CreateUser(),
			})
			continue
		} else {
			if savedConfigs[indexBridge].Username == "" {
				username := bridge.CreateUser()
				savedConfigs[indexBridge].Username = username
			}
		}
	}
	configPlugin = savedConfigs

	marshalConfig, _ := json.Marshal(savedConfigs)

	return marshalConfig
}

func findBridgeFromID(config []savedConfig, ID string) int {
	for ind, conf := range config {
		if conf.BridgeID == ID {
			return ind
		}
	}
	return -1
}

// OnStart discover brdiges and create the global state
func OnStart(config []byte) {
	if err := json.Unmarshal(config, &configPlugin); err != nil {
		panic(err)
	}

	discover()
}

func findStateFromID(arrayStates []State, ID string) bool {
	for _, state := range arrayStates {
		if state.Device.UniqueID == ID {
			return true
		}
	}
	return false
}

// Discover return array of all found devices
func Discover() []sdk.DiscoveredDevice {
	var discovered []sdk.DiscoveredDevice

	discover()

	for _, state := range States {
		discovered = append(discovered, sdk.DiscoveredDevice{
			Name:         state.Device.Name,
			PhysicalID:   state.Device.UniqueID,
			PhysicalName: state.Device.ProductID, // strings.ToLower(reflect.TypeOf(state.Device).Name()),
			Plugin:       Config.Name,
		})
	}

	return discovered
}

func discover() {
	bridges := DiscoverBridges()

	// create global state to store bridges and lights
	for _, savedConfig := range configPlugin {
		for i, bridge := range bridges {
			if savedConfig.BridgeID != bridge.ID {
				continue
			}
			bridges[i].Username = savedConfig.Username
			for j, light := range bridges[i].GetLights() {
				if !findStateFromID(States, light.UniqueID) {
					States = append(States, State{
						Bridge:   bridges[i],
						Device:   light,
						DeviceID: j + 1,
					})
				}
			}
		}
	}
}

// Params define actions parameters available
type Params struct {
	On  bool `json:"on"`
	Sat int  `json:"sat"`
	Bri int  `json:"bri"`
	Hue int  `json:"hue"`
}

// CallAction call functions from actions
func CallAction(physicalID string, name string, params []byte, config []byte) {
	if string(params) == "" {
		fmt.Println("Params must be provided")
	}

	// declare parameters
	var req Params

	// unmarshal parameters to use in actions
	err := json.Unmarshal(params, &req)
	if err != nil {
		fmt.Println(err)
	}

	// get the light's bridge
	bridge := GetBridge(physicalID)
	if bridge.ID == "" {
		return
	}

	// use name to call actions
	switch name {
	case "switchLight":
		bridge.SwitchLight(physicalID, req)
	case "toggleLight":
		bridge.ToggleLight(physicalID)
	default:
		return
	}
}

// OnStop close connection
func OnStop() {
}
