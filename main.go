package main

//go:generate go-bindata -o assets.go assets/vue/... assets/styles.css

import  (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/ContinentalBreakfast17/seriard"
	"github.com/zserge/webview"
)

const (
	SINGLE_CHANNEL_CHANGE = "channel"
	ALL_CHANNEL_CHANGE = "full"
	PROFILE_CHANGE = "profile_change"
	ADD_PROFILE = "add_profile"
	DELETE_PROFILE = "delete_profile"
	NAME_CHANGE = "name_change"
	MODE_CHANGE = "mode_change"
	SPEED_CHANGE = "speed_change"
)

type Message struct {
	Color 	[]int 	`json:"color,omitempty"`
	Index 	int 	`json:"index,omitempty"`
	Value 	int 	`json:"value,omitempty"`
	StrVal	string	`json:"strVal,omitempty"`
	Profile RGB		`json:"profile"`
}

type MessageHolder struct {
	Type 	string 	`json:"type"`
	Data 	Message	`json:"data"`
}

func handle(w webview.WebView, msg string) {
	log.Print(msg)
	var message MessageHolder
	err := json.Unmarshal([]byte(msg), &message)
	errorHandler("Failed to read message", err, true)

	switch message.Type {
		case SINGLE_CHANNEL_CHANGE:
			profiles.setColorChannel(message.Data.Index, message.Data.Value)
			writeColor(arduino, message.Data.Index, message.Data.Value)
		case ALL_CHANNEL_CHANGE:
			profiles.setColor(message.Data.Color)
			writeColors(arduino, message.Data.Color)
		case PROFILE_CHANGE:
			profiles.setCurrent(message.Data.Index)
			writeColors(arduino, message.Data.Color)
		case ADD_PROFILE:
			profiles.addProfile(message.Data.Profile)
		case DELETE_PROFILE:
			profiles.deleteProfile(message.Data.Index)
		case NAME_CHANGE:
			profiles.setName(message.Data.Index, message.Data.StrVal)
		case MODE_CHANGE:
			profiles.setMode(message.Data.Index, message.Data.StrVal)
			writeMode(arduino, message.Data.StrVal)
		case SPEED_CHANGE:
			profiles.setSpeed(message.Data.Index, message.Data.Value)
			writeSpeed(arduino, message.Data.Value)
		default:
			errorHandler("Failed to read message", errors.New("Unrecognized message type"), true)
	}
}

var profiles Profiles
var arduino *seriard.Arduino

func main() {
	// Read RGB Profiles and connect to arduino
	profiles = initProfiles()
	arduino = initArduino()

	// Read html
	html, err := ioutil.ReadFile("assets/web.html")
	errorHandler("Failed to read web.html", err, false)

	// Intialize Webview
	w := webview.New(webview.Settings{
		Title: "Arduino RGB Controller",
		URL: `data:text/html,` + url.PathEscape(string(html)),
		ExternalInvokeCallback: handle,
		Debug: true,
	})
	defer w.Exit()
	defer profiles.save()
	defer arduino.Disconnect()

	// Dispatch Webview
	w.Dispatch(func() {
		// Inject CSS
		w.InjectCSS(string(MustAsset("assets/styles.css")))

		// Inject Vue.js
		w.Eval(string(MustAsset("assets/vue/vendor/vue.min.js")))
		// Inject app code
		w.Eval(string(MustAsset("assets/vue/app.js")))

		profiles.send(w)	
		writeColors(arduino, profiles.getColor())
		writeMode(arduino, profiles.getMode())
	})
	w.Run()
	
}

func errorHandler(msg string, err error, shouldSave bool) {
	if err != nil {
		log.Printf("%s", msg)
		if shouldSave {
			profiles.save()
		}
		panic(err)
	}
}