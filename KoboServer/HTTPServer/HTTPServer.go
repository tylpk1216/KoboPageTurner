package main

import (
    "context"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
)

const PID_FILE = "/mnt/onboard/.koboserver/PID"
const CFG_FILE = "/mnt/onboard/.koboserver/koboserver.cfg"
const KOBO_INI_FILE = "/mnt/onboard/.kobo/Kobo/Kobo eReader.conf"

//----------------------------------------------------------------------
/*
struct input_event {
    struct timeval time;
    unsigned short type;
    unsigned short code;
    unsigned int value;
};

in the structre of Go.
type input_event struct {
    seconds int32
    microseconds int32

    type int16
    code int16
    value int32
}
*/
// little endian.
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h

// Clara HD, H2O
var gClaraHD = []byte{
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x01, 0x00,  0x4A, 0x01,  0x01, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x39, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x3B, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x35, 0x00,  0x56, 0x03, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x36, 0x00,  0x3F, 0x03, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x3A, 0x00,  0x1E, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x30, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x31, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x34, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x00, 0x00,  0x02, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x00, 0x00,  0x00, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x01, 0x00,  0x4A, 0x01,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x00, 0x00,  0x00, 0x00,  0x00, 0x00, 0x00, 0x00,
}

// Elipsa
var gElipsa = []byte{
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x39, 0x00,  0xCA, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x37, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x01, 0x00,  0x4A, 0x01,  0x01, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x30, 0x00,  0x60, 0x06, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x3A, 0x00,  0x60, 0x06, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x35, 0x00,  0x56, 0x03, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x03, 0x00,  0x36, 0x00,  0x3F, 0x03, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x04, 0x00,  0x00, 0x00,  0x00, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x39, 0x00,  0xCA, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x37, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x30, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x3A, 0x00,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x35, 0x00,  0x56, 0x03, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x08, 0x00,  0x03, 0x00,  0x36, 0x00,  0x3F, 0x03, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x0c, 0x00,  0x03, 0x00,  0x39, 0x00,  0xFF, 0xFF, 0xFF, 0xFF,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x0c, 0x00,  0x01, 0x00,  0x4A, 0x01,  0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x0c, 0x00,  0x00, 0x00,  0x00, 0x00,  0x00, 0x00, 0x00, 0x00,
}

type eventData struct {
    eventFile string

    rawData []byte

    leftX int
    leftY int

    rightX int
    rightY int
}

var eventItem eventData
var touchEvent *os.File

func addTimeStamp(buf []byte) error {
    n := int32(time.Now().Unix())
    s := fmt.Sprintf("%08X(%d)", n, n);
    fmt.Println(s)

    line := len(buf) / 16
    for i := 0; i < line; i++ {
        for d_i := 0; d_i < 4; d_i++ {
            index := i * 16

            s1 := s[8 - 2 * (d_i + 1) : 8 - 2 * d_i]
            n1, err := strconv.ParseUint(s1, 16, 8)
            if err != nil {
                fmt.Printf("Strconv Error (%v) \n", err)
                return fmt.Errorf("Convert Error : %s", s)
            }

            buf[index + d_i] = byte(n1)
        }
    }

    return nil
}

func debugEvent(buf []byte) {
    for i := 0; i < len(buf); i++ {
        fmt.Printf("%02X ", buf[i])

        if (i + 1) % 16 == 0 {
            fmt.Printf("\n")
        }
    }
}

func TriggerTouch(buf []byte) error {
    n, err := touchEvent.Write(buf)
    if err != nil {
        fmt.Printf("Write File Error (%v) \n", err)
        return err
    }

    fmt.Printf("Wrote %d bytes \n", n)

    debugEvent(buf)

    return nil
}

func TouchPage(buf []byte) error {
    err := addTimeStamp(buf)
    if err != nil {
        fmt.Println(err)
        return err
    }

    return TriggerTouch(buf)
}

func pixelToValue(x, y int) (int, int) {
    // The origin point is right-top.
    // x = y, y = x
    return y, x
}

func setXY(x, y int, buf []byte) {
    xValue, yValue := pixelToValue(x, y)

    //fmt.Println("After pixelToValue", xValue, yValue)

    // little endian
    x1 := xValue % 256
    x2 := xValue / 256

    y1 := yValue % 256
    y2 := yValue / 256

    i := 10
    for i  < len(buf) {
        // code = ABS_MT_POSITION_X
        if buf[i] == 0x35 && buf[i+1] == 0x00 {
            // value
            buf[i+2] = byte(x1)
            buf[i+3] = byte(x2)
        }

        // code = ABS_MT_POSITION_Y
        if buf[i] == 0x36 && buf[i+1] == 0x00 {
            // value
            buf[i+2] = byte(y1)
            buf[i+3] = byte(y2)
        }

        i += 16
    }
}

func leftPage() error {
    setXY(eventItem.leftX, eventItem.leftY, eventItem.rawData)
    return TouchPage(eventItem.rawData)
}

func rightPage() error {
    setXY(eventItem.rightX, eventItem.rightY, eventItem.rawData)
    return TouchPage(eventItem.rawData)
}

func left(w http.ResponseWriter, r *http.Request) {
    fmt.Println("/left")
    sendResponse(w, leftPage())
}

func right(w http.ResponseWriter, r *http.Request) {
    fmt.Println("/right")
    sendResponse(w, rightPage())
}

func sendResponse(w http.ResponseWriter, err error) {
    t := time.Now()
    n := int32(t.Unix())
    io.WriteString(w, fmt.Sprintf("%v(%d) (%v)", t, n, err))
}

func deletePID() {
    err := os.Remove(PID_FILE)
    if err != nil {
        fmt.Println(err)
    }
}

func getData(content, key, value string) string {
    i := strings.Index(content, key)
    if i == -1 {
        fmt.Printf("%s=%s(default)\n", key, value)
        return value
    }

    s := content[i+len(key)+1:]
    res := ""
    for i := 0; i < len(s); i++ {
        if s[i] == '\r' || s[i] == '\n' {
            break
        }
        res += string(s[i])
    }

    if len(s) == 0 {
        return value
    }

    return res
}

func atoi(s string) int {
    n, err := strconv.Atoi(s)
    if err != nil {
        return 0
    }
    return n
}

func getEventData() error {
    content, err := ioutil.ReadFile(CFG_FILE)
    if err != nil {
        return err
    }

    s := string(content)

    eventItem.eventFile = getData(s, "eventFile", "/dev/input/event1")

    // for the raw data of different devices.
    bufStr := getData(s, "rawData", "ClaraHD")
    switch bufStr {
    case "ClaraHD":
        eventItem.rawData = gClaraHD
        //eventItem.setXY = setXYClara
    case "H2O":
        eventItem.rawData = gClaraHD
        //eventItem.setXY = setXYClara
    case "Elipsa":
        eventItem.rawData = gElipsa
        //eventItem.setXY = setXYElipsa
    default:
        return fmt.Errorf("device is not supproted")
    }

    eventItem.leftX = atoi(getData(s, "leftX", "800"))
    eventItem.leftY = atoi(getData(s, "leftY", "500"))
    eventItem.rightX = atoi(getData(s, "rightX", "100"))
    eventItem.rightY = atoi(getData(s, "rightY", "500"))

    fmt.Printf("eventFile=%s \n", eventItem.eventFile)
    fmt.Printf("rawData=%s, len=%d \n", bufStr, len(eventItem.rawData))
    fmt.Printf("leftX=%d, leftY=%d \n", eventItem.leftX, eventItem.leftY)
    fmt.Printf("rightX=%d, rightY=%d \n", eventItem.rightX, eventItem.rightY)

    return nil
}

func main() {
    fmt.Println("Prepare to run Server", time.Now())

    defer deletePID()

    var err error
    err = getEventData()
    if err != nil {
        panic(fmt.Errorf("getEventData Error (%v) \n", err))
    }

    touchEvent, err = os.OpenFile(eventItem.eventFile, os.O_RDWR, 0777)
    if err != nil {
        panic(fmt.Errorf("Open File Error (%v) \n", err))
    }

    defer touchEvent.Close()

    m := http.NewServeMux()
    s := http.Server{Addr: ":80", Handler: m}

    m.HandleFunc("/left", left)
    m.HandleFunc("/right", right)

    m.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Clsoe Server")
        sendResponse(w, nil)

        waitSecs := 3
        closeTimer := time.NewTimer(time.Duration(waitSecs) * time.Second)
        go func() {
            <-closeTimer.C
            fmt.Printf("After %d seconds, closeTimer fired. \n", waitSecs)
            s.Shutdown(context.Background())
        }()
    })

    if err = s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }

    fmt.Println("Server Finished", time.Now())
}
