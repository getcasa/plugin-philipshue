package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/getcasa/plugin-philipshue/devices"
	"github.com/getcasa/sdk"
)

type Context struct {
	Bridge Bridge
	User   User
}

type Bridge struct {
	ID                string
	InternalIPAddress string
	Username          string
}

func Discover() []Bridge {
	res, err := http.Get("https://discovery.meethue.com/")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var bridges []Bridge
	if err := json.Unmarshal(body, &bridges); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return bridges
}

type User struct {
	Username string
}

type UserResponse struct {
	Success User
}

func (bridge *Bridge) CreateUser() string {
	data := []byte(`{
		"devicetype": "casa#plugin"
	}`)
	res, err := http.Post("http://"+bridge.InternalIPAddress+"/api", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var responses []UserResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return responses[0].Success.Username
}

func (bridge *Bridge) getLights() []devices.LCT0152A19ECLv5 {
	res, err := http.Get("http://" + bridge.InternalIPAddress + "/api/" + bridge.Username + "/lights")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var result map[string]devices.LCT0152A19ECLv5
	json.Unmarshal([]byte(body), &result)

	var lights []devices.LCT0152A19ECLv5
	for i := 1; result[strconv.Itoa(i)].Name != ""; i++ {
		lights = append(lights, result[strconv.Itoa(i)])
	}

	return lights
}

func main() {}

var Config = sdk.Configuration{
	Name:        "philips hue",
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
	},
}

type state struct {
	Bridge Bridge
	Device devices.LCT0152A19ECLv5
}

// Params define actions parameters available
type Params struct {
	ID  string
	On  bool
	Sat int
	Bri int
	Hue int
}

var states []state
var client http.Client

// Init plugin config
func Init() []byte {
	bridges := Discover()

	for i, bridge := range bridges {
		bridges[i].Username = bridge.CreateUser()
	}

	config, _ := json.Marshal(bridges)

	return config
}

// OnStart create http client
func OnStart(config []byte) {
	var bridges []Bridge

	if err := json.Unmarshal(config, &bridges); err != nil {
		panic(err)
	}

	for _, bridge := range bridges {
		for _, light := range bridge.getLights() {
			states = append(states, state{
				Bridge: bridge,
				Device: light,
			})
		}
	}
}

// CallAction call functions from actions
func CallAction(name string, params []byte) {
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

	// use name to call actions
	switch name {
	case "switchLight":
		// TODO: add call
	default:
		return
	}
}

// OnStop close connection
func OnStop() {
}
