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

// Discover list all bridge on local network
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
