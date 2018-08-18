package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/zserge/webview"
)

type RGB struct {
	Color 	[]int	`json:"color"`
	Mode	string	`json:"mode"`
	Index 	int 	`json:"index"`
	Name 	string 	`json:"name"`
	Speed 	int 	`json:"speed"`
}

type Profiles struct {
	List	[]RGB 	`json:"list"`
	Current	int 	`json:"current"`
}

func initProfiles() Profiles {
	raw, err := ioutil.ReadFile(os.Getenv("RGB_PROFILES"))
	errorHandler("Failed to read profiles", err, false)

	var p Profiles
	err = json.Unmarshal(raw, &p)
	errorHandler("Failed to unmarshal profiles", err, false)
	
	return p
}

func (p *Profiles) save() {
	f, err := os.Create(os.Getenv("RGB_PROFILES"))
	errorHandler("Failed to open profiles", err, false)
	defer f.Close()

	content, err := json.MarshalIndent(p, "", "\t")
	errorHandler("Failed to marshal profiles", err, false)

	_, err = f.Write(content)
	errorHandler("Failed to write profiles", err, false)
}

func (p *Profiles) send(w webview.WebView) {
	msg, err := json.Marshal(p)
	errorHandler("Failed to marshal profiles", err, true)
	w.Eval("rgb_controller.$emit('loadProfiles', " + string(msg) + " )")
}

func (p *Profiles) getColor() []int {
	return p.List[p.Current].Color
}

func (p *Profiles) getMode() string {
	return p.List[p.Current].Mode
}

func (p *Profiles) setColor(color []int) {
	copy(p.List[p.Current].Color, color)
}

func (p *Profiles) setColorChannel(channel, val int) {
	p.List[p.Current].Color[channel] = val
}

func (p *Profiles) setCurrent(index int) {
	p.Current = index
}

func (p *Profiles) addProfile(profile RGB) {
	p.List = append(p.List, profile)
	p.Current = profile.Index
}

func (p *Profiles) deleteProfile(index int) {
	p.List = append(p.List[:index], p.List[index+1:]...)
	for i := index; i < len(p.List); i++ {
		p.List[i].Index--
	}
}

func (p *Profiles) setName(index int, name string) {
	p.List[index].Name = name
}

func (p *Profiles) setMode(index int, mode string) {
	p.List[index].Mode = mode
}

func (p *Profiles) setSpeed(index, speed int) {
	p.List[index].Speed = speed
}