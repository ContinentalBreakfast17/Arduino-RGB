package main

//go:generate go-bindata -o assets.go assets/vue/... assets/styles.css

import  (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/ContinentalBreakfast17/seriard"
	"github.com/zserge/webview"
)

type Message struct {
	Color 	[]int 	`json:"color,omitempty"`
	Index 	int 	`json:"index,omitempty"`
	Value 	int 	`json:"value,omitempty"`
	StrVal	string	`json:"strVal,omitempty"`
}

type MessageHolder struct {
	Type 	string 	`json:"type"`
	Data 	Message	`json:"data"`
}

func handle(w webview.WebView, msg string) {
	log.Print(msg)
	var message MessageHolder
	err := json.Unmarshal([]byte(msg), &message)
	if err != nil {
		log.Print(err)
		return
	}

	switch message.Type {
		case SINGLE_CHANNEL_CHANGE:
			profiles.List[profiles.Current].Color[message.Data.Index] = message.Data.Value
		case ALL_CHANNEL_CHANGE:
			copy(profiles.List[profiles.Current].Color, message.Data.Color)
		default:
			log.Print("Bad message")
	}
}

func sendProfile(w webview.WebView, profile RGB) error {
	msg, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	w.Eval("rgb_controller.$emit('set', " + string(msg) + " )")
	return nil
}

var profiles Profiles
var arduino *seriard.Arduino

func main() {
	// Read RGB Profiles and connect to arduino
	profiles = initProfiles()
	arduino = initArduino()

	// Read html
	html, err := ioutil.ReadFile("assets/web.html")
	if err != nil {
		panic(err)
	}

	// Intialize Webview
	w := webview.New(webview.Settings{
		Title: "Arduino RGB Controller",
		URL: `data:text/html,` + url.PathEscape(string(html)),
		ExternalInvokeCallback: handle,
		Debug: true,
	})
	defer w.Exit()
	defer saveProfiles(profiles)
	defer arduino.Disconnect()

	// Dispatch Webview
	w.Dispatch(func() {
		// Inject CSS
		w.InjectCSS(string(MustAsset("assets/styles.css")))

		// Inject Vue.js
		w.Eval(string(MustAsset("assets/vue/vendor/vue.min.js")))
		// Inject app code
		w.Eval(string(MustAsset("assets/vue/app.js")))

		err = sendProfile(w, profiles.List[profiles.Current])
		if err != nil {
			panic(err)
		}		
	})
	w.Run()
	
}