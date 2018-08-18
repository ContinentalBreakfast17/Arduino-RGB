var rgb_controller = new Vue({
	el: "#rgb",
	data: {
		hex_val: "#000000",
		profiles: {
			current: 0,
			list: [
				{
					name: "",
					mode: "static",
					speed: 5,
					color: [0, 0, 0]
				}
			]
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
			this.profiles.list[this.profiles.current].color = color.slice();

			var msg = {
				"type": "full",
				"data": {
					"color": color
				}
			};
			window.external.invoke(JSON.stringify(msg));
		},
		setChannel: function(channel, val) {
			this.profiles.list[this.profiles.current].color[channel] = parseInt(val, 10);
			this.hex_val = rgbToHex(this.profiles.list[this.profiles.current].color);

			var msg = {
				"type": "channel",
				"data": {
					"index": channel,
					"value": parseInt(val,10)
				}
			};
			window.external.invoke(JSON.stringify(msg));
		},
		setName: function() {
			var msg = {
				"type": "name_change",
				"data": {
					"index": this.profiles.current,
					"strVal": this.profiles.list[this.profiles.current].name
				}
			};
			window.external.invoke(JSON.stringify(msg));
		},
		setMode: function() {
			var msg = {
				"type": "mode_change",
				"data": {
					"index": this.profiles.current,
					"strVal": this.profiles.list[this.profiles.current].mode
				}
			};
			window.external.invoke(JSON.stringify(msg));
			this.setSpeed();
		},
		setSpeed: function() {
			var msg = {
				"type": "speed_change",
				"data": {
					"index": this.profiles.current,
					"value": parseInt(this.profiles.list[this.profiles.current].speed, 10)
				}
			};
			window.external.invoke(JSON.stringify(msg));
		},
		setCurrentProfile: function(index) {
			this.profiles.current = index;
			color = this.profiles.list[index].color;
			this.hex_val = rgbToHex(color);

			var msg = {
				"type": "profile_change",
				"data": {
					"index": index,
					"color": color,
					"strVal": this.profiles.list[index].mode
				}
			};
			window.external.invoke(JSON.stringify(msg));
			this.setMode();
		},
		newProfile: function() {
			profile = {
				color: [
					Math.floor(Math.random() * 256),
					Math.floor(Math.random() * 256),
					Math.floor(Math.random() * 256)
				],
				index: this.profiles.list.length,
				mode: "static",
				name: "Profile " + (this.profiles.list.length + 1).toString(),
				speed: 5
			};
			this.profiles.list.push(profile);

			var msg = {
				"type": "add_profile",
				"data": {
					"profile": profile
				}
			};
			window.external.invoke(JSON.stringify(msg));
			this.setCurrentProfile(profile.index);
		},
		deleteProfile: function() {
			if(this.profiles.list.length == 1) {
				return
			}

			index = this.profiles.current;
			for(i = this.profiles.current; i < this.profiles.list.length-1; i++) {
				this.profiles.list[i] = this.profiles.list[i+1];
				this.profiles.list[i].index--;
			}
			this.profiles.list.pop();

			var msg = {
				"type": "delete_profile",
				"data": {
					"index": index
				}
			};
			window.external.invoke(JSON.stringify(msg));
			if(this.profiles.current == this.profiles.list.length) this.setCurrentProfile(this.profiles.current - 1);
			else this.setCurrentProfile(this.profiles.current);
		},
		loadProfiles: function(profiles) {
			this.profiles.list.pop();
			this.profiles.current = profiles.current;
			for (var i = 0; i < profiles.list.length; i++) {
				this.profiles.list.push(profiles.list[i]);
			}
			color = this.profiles.list[this.profiles.current].color;
			this.hex_val = rgbToHex(color);
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

function rgbToHex(color) {
	return "#" + componentToHex(color[0]) + componentToHex(color[1]) + componentToHex(color[2]);
}