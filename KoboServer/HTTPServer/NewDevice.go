package main

import (
    "fmt"
    "os"
    "time"
)

var touchEvent *os.File

func cc() {
    touchEvent.Close()
}

func main() {
    fmt.Println("Prepare to run Server", time.Now())

    touchEvent, err := os.OpenFile("/dev/input/event2", os.O_RDWR, 0777)
    if err != nil {
        panic(fmt.Errorf("Open File Error (%v) \n", err))
    }

    //defer touchEvent.Close()

    time.AfterFunc(60 * time.Second, cc)

    buf := make([]byte, 1024)

    for {
        if touchEvent == nil {
            break
        }

        count, err := touchEvent.Read(buf)
        if err != nil {
            touchEvent.Close()
            break
        }

        for i := 0; i < count; i++ {
            fmt.Printf("%02X ", buf[i])

            if (i + 1) % 16 == 0 {
                fmt.Printf("\n")
            }
        }
    }

    fmt.Println("Done")
}
