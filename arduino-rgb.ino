// all communications will be as long as the macro below
#define MESSAGE_SIZE 	32
// feel free to change the baud rate as desired
#define BAUDRATE 		115200
// this pin will be used to blink error messages
#define ERROR_PIN 		2
// the number of milliseconds for delay between error blinks
#define BLINK_TIME 		200
// blink counts for error messages
#define ERR_BAD_READ 	1
#define ERR_BAD_CMD		2
#define ERR_BAD_PARAMS	3
#define ERR_UNKWN_CMD 	4
// RGB modes
#define STATIC 			"static"
#define RAINBOW 		"rainbow"
#define FADE 			"fade"
#define STROBE 			"strobe"
// RGB pins
#define RED 			9
#define GREEN			10
#define BLUE			11

char* buffer;
char* rgb_mode;

int color[3];
int rainbow[3];
int rainbow_index;
int strobe;
float fade;
int fade_dir;
int speed;
unsigned long then;

int doCommand(char* params, int (*function)(int*), int argc);
int changeMode(char* params);
int changeSpeed(char* params);
void success(int v1, int v2);
void error(char* error_message, int blinks);
void blink(int blinks);
char* nextWS(char* s);

int digitalWriteWrapper(int* args);
int analogWriteWrapper(int* args);
int digitalReadWrapper(int* args);
int analogReadWrapper(int* args);
int pinModeWrapper(int* args);

void doRainbow();
void doStrobe();
void doFade();

void setup() {
	// initialize any pins as needed here
	// setting pins to output may help if you are experiencing irregular voltage floating on startup
	pinMode(ERROR_PIN, OUTPUT);
	pinMode(RED, OUTPUT);
	pinMode(GREEN, OUTPUT);
	pinMode(BLUE, OUTPUT);

	buffer = (char*)malloc(sizeof(char)*(MESSAGE_SIZE + 1));
	buffer[MESSAGE_SIZE] = 0;

	rgb_mode = (char*)malloc(sizeof(char)*(MESSAGE_SIZE + 1));
	memset(rgb_mode, 0, MESSAGE_SIZE+1);
	strcpy(rgb_mode, RAINBOW);

	memset(color, 0, sizeof(int)*3);
	memset(rainbow, 0, sizeof(int)*3);
	rainbow[0] = 255;
	rainbow_index = 0;
	speed = 5;
	strobe = 1;
	fade = 1;
	fade_dir = -1;
	then = millis();

	Serial.begin(BAUDRATE);
}

void loop() {
	if(Serial.available() > 0) {
		int bytesRead = Serial.readBytes(buffer, MESSAGE_SIZE);
		if(bytesRead != MESSAGE_SIZE)  {
			error("Serial read failure", ERR_BAD_READ);
			return;
		}

		char* ws = nextWS(buffer);
		int dif = ws - buffer;
		if(ws == NULL || dif == 0 || dif >= MESSAGE_SIZE-2)  {
			error("Bad command", ERR_BAD_CMD);
			return;
		}

		char command[dif + 1];
		char params[MESSAGE_SIZE + 1 - (dif + 1)];
		memcpy(command, buffer, dif);
		memcpy(params, buffer + dif + 1, MESSAGE_SIZE + 1 - (dif + 1));
		command[dif] = 0;

		int result = 0;
		if(strcmp("digital_write", command) == 0) {
			result = doCommand(params, &digitalWriteWrapper, 2);
		} else if(strcmp("analog_write", command) == 0) {
			result = doCommand(params, &analogWriteWrapper, 2);
		} else if(strcmp("digital_read", command) == 0) {
			result = doCommand(params, &digitalReadWrapper, 1);
		} else if(strcmp("analog_read", command) == 0) {
			result = doCommand(params, &analogReadWrapper, 1);
		} else if(strcmp("set_pin_mode", command) == 0) {
			result = doCommand(params, &pinModeWrapper, 2);
		} else if(strcmp("set_rgb_mode", command) == 0) {
			result = changeMode(params);
		} else if(strcmp("set_speed", command) == 0) {
			result = changeSpeed(params);
		} else {
			error("Unknown command", ERR_UNKWN_CMD);
			return;
		}

		if(result < 0) {
			error("Bad parameters", ERR_BAD_PARAMS);
			return;
		}
	} else if(strcmp(STATIC, rgb_mode) != 0) {

		if(millis() - then > speed) {
			if(strcmp(RAINBOW, rgb_mode) == 0) {
				doRainbow();
			} else if(strcmp(STROBE, rgb_mode) == 0) {
				doStrobe();
			} else if(strcmp(FADE, rgb_mode) == 0) {
				doFade();
			}
			then = millis();
		}
	}
}

int doCommand(char* params, int (*function)(int*), int argc) {
	int args[argc];
	int i;

	for(i = 0; i < argc; i++) {
		int n = sscanf(params, "%u", &args[i]);
		if(n != 1) return -1;
		params += (int)args[i]/10 + 2; // + 1 for ceil, + 1 to skip space
	}

	return (*function)(args);
}

int changeMode(char* params) {
	char tmp[MESSAGE_SIZE+1];
	memset(tmp, 0, MESSAGE_SIZE+1);
	int n = sscanf(params, "%s ", tmp);
	if(n != 1) return -1;
	if(strcmp(STATIC, tmp) == 0 || strcmp(RAINBOW, tmp) == 0 || 
		strcmp(FADE, tmp) == 0 || strcmp(STROBE, tmp) == 0) {
		if(strcmp(RAINBOW, tmp) != 0) {
			analogWrite(RED, color[0]);
			analogWrite(GREEN, color[1]);
			analogWrite(BLUE, color[2]);
		}
		memset(rgb_mode, 0, MESSAGE_SIZE+1);
		strcpy(rgb_mode, tmp);
		success(-1, -1);
		return 0;
	} else return -1;
}

int changeSpeed(char* params) {
	int tmp;
	int n = sscanf(params, "%d", &tmp);
	if(n != 1) return -1;
	if(tmp <= 0) return -1;
	speed = tmp;
	success(-1, -1);
	return 0;
}

void success(int v1, int v2) {
	char msg[MESSAGE_SIZE+1];
	memset(msg, 32, MESSAGE_SIZE);
	sprintf(msg, "%d %d", v1, v2);
	msg[strlen(msg)] = 32;
	msg[MESSAGE_SIZE] = 0;
	Serial.print(msg);
}

void error(const char* error_message, int blinks) {
	char msg[MESSAGE_SIZE];
	memset(msg, 0, MESSAGE_SIZE);
	memcpy(msg, error_message, strlen(error_message));

	Serial.print(msg);
	blink(blinks);
}

void blink(int blinks) {
	int i;
	int val = digitalRead(ERROR_PIN);
	digitalWrite(ERROR_PIN, LOW);
	delay(BLINK_TIME);
	for(i = 0; i < blinks; i++) {
		digitalWrite(ERROR_PIN, HIGH);
		delay(BLINK_TIME);
		digitalWrite(ERROR_PIN, LOW);
		delay(BLINK_TIME);
	}
	digitalWrite(ERROR_PIN, val);
}

// returns a pointer to the next white space in s, or NULL if there is no more white space
// assumes string is null terminated
char* nextWS(char* s) {
	char* p = s;
	for(; *p != 0 && *p != ' ' && *p != '\t'; p++);
	if(*p) return p;
	return NULL;
}

// wrapper functions to make c work

int digitalWriteWrapper(int* args) {
	digitalWrite(args[0], args[1]);
	success(args[0], args[1]);
	return 0; 
}

int analogWriteWrapper(int* args) {
	analogWrite(args[0], args[1]);
	// update color
	if(args[0] == RED) color[0] = args[1];
	else if(args[0] == GREEN) color[1] = args[1];
	else if(args[0] == BLUE) color[2] = args[1];
	success(args[0], args[1]);
	return 0; 
}

int digitalReadWrapper(int* args) {
	int val = digitalRead(args[0]);
	success(args[0], val);
	return 0;
}

int analogReadWrapper(int* args) {
	int val = analogRead(args[0]);
	success(args[0], val);
	return 0;
}

int pinModeWrapper(int* args) {
	pinMode(args[0], args[1]);
	success(args[0], args[1]);
	return 0; 
}

// RGB Modes

void doRainbow() {
	rainbow[rainbow_index]--;
	rainbow[(rainbow_index == 2) ? (0) : (rainbow_index + 1)]++;
	if(rainbow[rainbow_index] == 0) {
		rainbow_index = (rainbow_index == 2) ? (0) : (rainbow_index + 1);
	}
	analogWrite(RED, rainbow[0]);
	analogWrite(GREEN, rainbow[1]);
	analogWrite(BLUE, rainbow[2]);
}

void doStrobe() {
	strobe = !strobe;
	analogWrite(RED, color[0]*strobe);
	analogWrite(GREEN, color[1]*strobe);
	analogWrite(BLUE, color[2]*strobe);
}

void doFade() {
	fade += fade_dir;
	if(fade >= 255) fade_dir = -1;
	else if(fade <= 0) fade_dir = 1;
	analogWrite(RED, color[0]*(fade/255));
	analogWrite(GREEN, color[1]*(fade/255));
	analogWrite(BLUE, color[2]*(fade/255));
}
