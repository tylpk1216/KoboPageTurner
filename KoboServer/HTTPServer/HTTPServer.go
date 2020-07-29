package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"
    "strconv"
    "time"
)

//----------------------------------------------------------------------
/*
struct input_event {
    struct timeval time;
    unsigned short type;
    unsigned short code;
    unsigned int value;
};
*/
// little endian, Clara is 32bit OS?
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h
var gEvent = []byte{
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x01, 0x00, 0x4A, 0x01, 0x01, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x03, 0x00, 0x39, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x03, 0x00, 0x3B, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x03, 0x00, 0x35, 0x00, 0x56, 0x03, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x03, 0x00, 0x36, 0x00, 0x3F, 0x03, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x03, 0x00, 0x3A, 0x00, 0x1E, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x03, 0x00, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x03, 0x00, 0x31, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x03, 0x00, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x02, 0x00,  0x01, 0x00, 0x4A, 0x01, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x02, 0x00,  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func addTimeStamp(buf []byte) error {
    n := int32(time.Now().Unix())
    s := fmt.Sprintf("%08X(%d)", n, n);
    //fmt.Println(s)

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
    f, err := os.OpenFile("/dev/input/event1", os.O_RDWR, 0777)
    if err != nil {
        fmt.Printf("Open File Error (%v) \n", err)
        return err
    }

    defer f.Close()

    n, err := f.Write(buf)
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

func setXY(x, y int) {
    // little endian
    x1 := x % 256
    x2 := x / 256

    y1 := y % 256
    y2 := y / 256

    i := 60
    gEvent[i]   = byte(x1)
    gEvent[i+1] = byte(x2)

    i = i + 16
    gEvent[i]   = byte(y1)
    gEvent[i+1] = byte(y2)
}

func leftPage() error {
    setXY(0x0356, 0x033F)
    return TouchPage(gEvent)
}

func rightPage() error {
    setXY(0x0479, 0x013F)
    return TouchPage(gEvent)
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

func main() {
    fmt.Println("Prepare to run Server")

    m := http.NewServeMux()

    s := http.Server{Addr: ":80", Handler: m}

    m.HandleFunc("/left", left)
    m.HandleFunc("/right", right)

    m.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Clsoe Server")
        
        err := os.Remove("/mnt/onboard/.koboserver/PID")
        sendResponse(w, err)

        waitSecs := 3
        closeTimer := time.NewTimer(time.Duration(waitSecs) * time.Second)
        go func() {
            <-closeTimer.C
            fmt.Printf("After %d seconds, closeTimer fired. \n", waitSecs)
            s.Shutdown(context.Background())
        }()
    })

    if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }

    fmt.Println("Server Finished")
}
