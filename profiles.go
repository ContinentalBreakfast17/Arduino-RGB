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

	content, err := json.Marshal(p)
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

func (p *Profiles) setColor(color []int) {
	copy(p.List[p.Current].Color, color)
}

func (p *Profiles) setColorChannel(channel, val int) {
	p.List[p.Current].Color[channel] = val
}

func (p *Profiles) setCurrent(index int) {
	p.Current = index
}