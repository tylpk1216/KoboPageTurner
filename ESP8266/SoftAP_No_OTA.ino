#include <ESP8266WiFi.h>

#define AP_SSID    "KoboPageTurner"
#define AP_PWD     "123456"

#define HOST       "192.168.4.2"
#define PORT       80

#define LED_PIN     2
#define LEFT_PIN   14
#define RIGHT_PIN   4

#define PRESSED_STATUS  LOW

#define LEFT_PRESSED   1
#define RIGHT_PRESSED  2
#define DUAL_PRESSED   3

String LEFT_API = "GET /left";
String RIGHT_API = "GET /right";
String EXIT_API = "GET /exit";

void initPIN()
{
    pinMode(LED_PIN, OUTPUT);
    pinMode(LEFT_PIN, INPUT_PULLUP);
    pinMode(RIGHT_PIN, INPUT_PULLUP);
}

void initAP()
{
    Serial.print("Setting soft-AP ... ");

    // AP Mode
    boolean result = WiFi.softAP(AP_SSID, AP_PWD, 1, false, 1);
    if (result == true) {
        Serial.println("Ready");
    } else {
        Serial.println("Failed!");
    }
}

void setup()
{
    Serial.begin(74880);
    Serial.println();

    initPIN();
    initAP();

    Serial.println("After setup");

    showLED();
}

void showLED()
{
    digitalWrite(LED_PIN, LOW);
    delay(2000);
    digitalWrite(LED_PIN, HIGH);
}

bool isJustCheckPressed(volatile char pin)
{
    if (digitalRead(pin) == PRESSED_STATUS) {
        return true;
    }
    return false;
}

void releaseButton(volatile char pin)
{
    while (digitalRead(pin) == PRESSED_STATUS) {
        delay(20);
    }
}

char getButtonStatus()
{
    char status = 0;

    bool isLeft = isJustCheckPressed(LEFT_PIN);
    bool isRight = isJustCheckPressed(RIGHT_PIN);

    if (isLeft) {
        status += LEFT_PRESSED;
        releaseButton(LEFT_PIN);
    }

    if (isRight) {
        status += RIGHT_PRESSED;
        releaseButton(RIGHT_PIN);
    }

    return status;
}

void loop()
{
    if (WiFi.softAPgetStationNum() >= 1) {
        char status = getButtonStatus();

        if (status == DUAL_PRESSED) {
            turnPage(EXIT_API);
            showLED();
        } else if (status == LEFT_PRESSED) {
            turnPage(LEFT_API);
        } else if (status == RIGHT_PRESSED) {
            turnPage(RIGHT_API);
        }
    }
}

void turnPage(String api)
{
    WiFiClient client;

    if (!client.connect(HOST, PORT)) {
        Serial.println("connection failed");
    } else {
        String getStr = api + " HTTP/1.1\r\n" +
                              "Host: 192.168.4.2\n" +
                              "Connection: close\r\n\r\n";

        client.print(getStr);
        delay(10);

        client.stop();
    }
}
