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
	}

	var bridges []Bridge
	if err := json.Unmarshal(body, &bridges); err != nil {
		fmt.Println(err)
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
	}

	var responses []UserResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		fmt.Println(err)
	}

	return responses[0].Success.Username
}

// GetLights list all lights from bridge
func (bridge *Bridge) GetLights() []devices.Hue {
	res, err := http.Get("http://" + bridge.InternalIPAddress + "/api/" + bridge.Username + "/lights")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var result map[string]devices.Hue
	json.Unmarshal([]byte(body), &result)

	var lights []devices.Hue
	for i := 1; i <= 50; i++ {
		if result[strconv.Itoa(i)].Name != "" {
			lights = append(lights, result[strconv.Itoa(i)])
		}
	}
	return lights
}

// GetLight get state light from bridge
func (bridge *Bridge) GetLight(id int) devices.Hue {
	res, err := http.Get("http://" + bridge.InternalIPAddress + "/api/" + bridge.Username + "/lights/" + strconv.Itoa(id))
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var light devices.Hue
	json.Unmarshal(body, &light)

	return light
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

// SwitchLight send user params to the light
func (bridge *Bridge) SwitchLight(physicalID string, params Params) {
	byteParams, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	id := GetLightID(physicalID)
	if id == -1 {
		fmt.Println("Wrong id")
		return
	}

	req, err := http.NewRequest(http.MethodPut, "http://"+bridge.InternalIPAddress+"/api/"+bridge.Username+"/lights/"+strconv.Itoa(id)+"/state", bytes.NewBuffer(byteParams))
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

// ToggleLight get the current state light to toggle on/off
func (bridge *Bridge) ToggleLight(physicalID string) {
	id := GetLightID(physicalID)
	if id == -1 {
		fmt.Println("Wrong id")
		return
	}

	light := bridge.GetLight(id)
	on := !light.State.ON

	lightReq := Params{
		On:  on,
		Sat: light.State.Sat,
		Bri: light.State.Bri,
		Hue: light.State.Hue,
	}
	byteParams, err := json.Marshal(lightReq)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, "http://"+bridge.InternalIPAddress+"/api/"+bridge.Username+"/lights/"+strconv.Itoa(id)+"/state", bytes.NewBuffer(byteParams))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("http://"+bridge.InternalIPAddress+"/api/"+bridge.Username+"/lights/"+strconv.Itoa(id)+"/state", bytes.NewBuffer(byteParams))

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	return
}
