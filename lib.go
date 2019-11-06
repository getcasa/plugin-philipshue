package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/getcasa/plugin-philipshue/devices"
)

// Bridge define the physical Philips Hue bridge
type Bridge struct {
	ID                string
	InternalIPAddress string
	Username          string
}

// DiscoverBridges list all bridge on local network
func DiscoverBridges() []Bridge {
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

// User define token for requests
type User struct {
	Username string
}

// UserResponse define user API response from bridge
type UserResponse struct {
	Success User
}

// CreateUser register an user on bridge to authenticate requests
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

// GetLights list all lights from bridge
func (bridge *Bridge) GetLights() []devices.LCT0152A19ECLv5 {
	res, err := http.Get("http://" + bridge.InternalIPAddress + "/api/" + bridge.Username + "/lights")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var result map[string]devices.LCT0152A19ECLv5
	json.Unmarshal([]byte(body), &result)

	var lights []devices.LCT0152A19ECLv5
	for i := 1; result[strconv.Itoa(i)].Name != ""; i++ {
		lights = append(lights, result[strconv.Itoa(i)])
	}

	return lights
}

// SwitchLight send user params to the light
func (bridge *Bridge) SwitchLight(params Params) {
	byteParams, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	id := GetLightID(params.ID)
	if id == -1 {
		fmt.Println("Wrong id")
		return
	}

	params.ID = strconv.Itoa(id)

	req, err := http.NewRequest(http.MethodPut, "http://"+bridge.InternalIPAddress+"/api/"+bridge.Username+"/lights/"+params.ID, bytes.NewBuffer(byteParams))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	return
}

// GetLightID return int id of the light
func GetLightID(uid string) int {
	for _, state := range States {
		if state.Device.UniqueID == uid {
			return state.DeviceID
		}
	}
	return -1
}

// GetBridge return the bridge of the light
func GetBridge(uid string) Bridge {
	for _, state := range States {
		if state.Device.UniqueID == uid {
			return state.Bridge
		}
	}
	return Bridge{}
}