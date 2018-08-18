var rgb_controller = new Vue({
	el: "#rgb",
	data: {
		hex_val: "#000000",
		r: 0,
		g: 0,
		b: 0,
		profiles: {
			current: 0,
			list: []
		}
	},
	components: {
		'rgb-input': {
			props: {
				channel: 	String,
				val: 		Number
			},
			template: `
				<span><h class="channel_label">{{ channel }}:</h>
					<input type="number" v-model="val" v-on:change="$emit('color-change', val)" min="0" max="255" class="channel_textbox"></input>
					<input type="range" v-model="val" v-on:change="$emit('color-change', val)" min="0" max="255" class="slider"></input> 
				</span>
			`
		},
		'profile-button': {
			props: {
				name: 		String,
				index: 		Number
			},
			template: `
				<button v-on:click="$emit('profile-click', index)" class="profile_button"> {{ name }} </button>
			`
		}
	},
	methods: {
		hex: function() {
			this.hex_val = document.getElementById("hex_color").value;
			var color = parseColor(this.hex_val);
			this.r = color[0]; this.g = color[1]; this.b = color[2];
			this.profiles.list[this.profiles.current].color = color.slice();

			var msg = {
				"type": "full",
				"data": {
					"color": [this.r, this.g, this.b]
				}
			}
			window.external.invoke(JSON.stringify(msg))
		},
		color: function(channel, val) {
			switch(channel) {
				case 0:
					this.r = parseInt(val,10);
					break;
				case 1:
					this.g = parseInt(val,10);
					break;
				case 2:
					this.b = parseInt(val,10);
					break;
				default:
					console.log("Invalid color channel");
					return;
			}
			this.hex_val = rgbToHex(this.r, this.g, this.b);

			var msg = {
				"type": "channel",
				"data": {
					"index": channel,
					"value": parseInt(val,10)
				}
			}
			window.external.invoke(JSON.stringify(msg))
		},
		setCurrentProfile: function(index) {
			this.profiles.current = index;
			color = this.profiles.list[this.profiles.current].color;
			this.r = color[0]; this.g = color[1]; this.b = color[2];
			this.hex_val = rgbToHex(this.r, this.g, this.b);

			var msg = {
				"type": "profile_change",
				"data": {
					"index": index
				}
			}
			window.external.invoke(JSON.stringify(msg))
		},
		loadProfiles: function(profiles) {
			this.profiles.current = profiles.current;
			for (var i = 0; i < profiles.list.length; i++) {
				this.profiles.list.push(profiles.list[i])
			}
			color = this.profiles.list[this.profiles.current].color;
			this.r = color[0]; this.g = color[1]; this.b = color[2];
			this.hex_val = rgbToHex(this.r, this.g, this.b);
		}
	},
	created: function() {
		this.$on('loadProfiles', this.loadProfiles);
	}
})

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