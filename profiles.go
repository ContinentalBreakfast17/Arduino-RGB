package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	SINGLE_CHANNEL_CHANGE = "channel"
	ALL_CHANNEL_CHANGE = "full"
)

type RGB struct {
	Color 	[]int	`json:"color"`
	Mode	string	`json:"mode"`
	Index 	int 	`json:"index"`
	Name 	string 	`json:"name"`
}

type Profiles struct {
	List	[]RGB 	`json:"profiles"`
	Current	int 	`json:"current"`
}

func initProfiles() Profiles {
	raw, err := ioutil.ReadFile(os.Getenv("RGB_PROFILES"))
	if err != nil {
		panic(err)
	}

	var p Profiles
	err = json.Unmarshal(raw, &p)
	if err != nil {
		panic(err)
	}
	
	return p
}

func saveProfiles(p Profiles) {
	f, err := os.Create(os.Getenv("RGB_PROFILES"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	content, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(content)
	if err != nil {
		panic(err)
	}
}