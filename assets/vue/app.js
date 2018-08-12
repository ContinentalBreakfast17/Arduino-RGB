var rgb_controller = new Vue({
	el: "#rgb",
	data: {
		hex_val: "#000000",
		r: 0,
		g: 0,
		b: 0,
		test: tester
	},
	components: {
		'rgb-input': {
			props: {
				channel: 	String,
				val: 		Number
			},
			template: `
				<span><h class="channel_label">{{ channel }}:</h>
					<input type="number" v-model="val" min="0" max="255" class="channel_textbox"></input>
					<input type="range" v-model="val" v-on:change="$emit('color-change', val)" min="0" max="255" class="slider"></input> 
				</span>
			`
		}
	},
	methods: {
		hex: function() {
			this.hex_val = document.getElementById("hex_color").value;
			var color = parseColor(this.hex_val);
			this.r = color[0]; this.g = color[1]; this.b = color[2];
			recolor_sliders(this.hex_val);
		},
		color: function(channel, val) {
			switch(channel) {
				case 0:
					this.r = val;
					break;
				case 1:
					this.g = val;
					break;
				case 2:
					this.b = val;
					break;
				default:
					console.log("Invalid color channel");
					return;
			}
			this.hex_val = rgbToHex(this.r, this.g, this.b);
			recolor_sliders(this.hex_val);
			this.testem(this.hex_val);
		},
		testem: function(val) {
			test.hey(val);
		}
	}
})

function recolor_sliders(color) {
	var sliders = document.getElementsByClassName("slider-thumb");
	for (var i = 0; i < sliders.length; i++) {
		console.log(sliders[i].style);
		sliders[i].style.background = color;
	}
}

function parseColor(input) {
	var m = input.match(/^#([0-9a-f]{6})$/i)[1];
	if(m) {
		return [
			parseInt(m.substr(0,2),16),
			parseInt(m.substr(2,2),16),
			parseInt(m.substr(4,2),16)
		];
	}
}

function componentToHex(c) {
	if(typeof c === "string") {
		c = parseInt(c, 10);
	}
	var hex = c.toString(16);
	return hex.length == 1 ? "0" + hex : hex;
}

function rgbToHex(r, g, b) {
	return "#" + componentToHex(r) + componentToHex(g) + componentToHex(b);
}