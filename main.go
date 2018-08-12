package main

//go:generate go-bindata -o assets.go assets/vue/... assets/styles.css

import  (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/zserge/webview"
)

type Tester struct {
	Value string `json:"value"`
}

func (t *Tester) Hey(s string) {
	t.Value = s
	fmt.Println(t.Value)
}

func main() {
	html, err := ioutil.ReadFile("web.html")
	if err != nil {
		panic(err)
	}

	w := webview.New(webview.Settings{
		Title: "Arduino RGB Controller",
		URL: `data:text/html,` + url.PathEscape(html),
	})
	defer w.Exit()

	w.Dispatch(func() {
		// Inject tester
		w.Bind("tester", &Tester{})

		// Inject CSS
		w.InjectCSS(string(MustAsset("assets/styles.css")))

		// Inject Vue.js
		w.Eval(string(MustAsset("assets/vue/vendor/vue.min.js")))
		// Inject app code
		w.Eval(string(MustAsset("assets/vue/app.js")))
	})
	w.Run()
	
}