package main

import (
    "fmt"
    "os"
    "time"
)

func catchWorker(touchEvent *os.File, c chan<- []byte ) {
    stop := make(chan bool, 1)

    time.AfterFunc(30 * time.Second, func(){
        stop <- true
    })

    buf := make([]byte, 1024)

    exit := false
    for {
        select {
        case <-stop:
            exit = true
        default:
            count, err := touchEvent.Read(buf)
            if err != nil {
                exit = true
            }

            if count > 0 {
                res := make([]byte, count)
                res = buf[0:count]
                c <- res
            }
        }

        if exit {
            break
        }
    }

    touchEvent.Close()
    close(c)
}

func main() {
    f := "/dev/input/event2"

    if len(os.Args) >= 2 {
        f = os.Args[1]
    }

    fmt.Printf("Prepare to catch input data(%s) %v \n", f, time.Now())

    touchEvent, err := os.OpenFile(f, os.O_RDWR, 0777)
    if err != nil {
        fmt.Printf("Open File Error (%v) \n", err)
        return
    }

    time.Sleep(10 * time.Second)

    c := make(chan []byte)
    go catchWorker(touchEvent, c)

    for buf := range c {
        count := len(buf)
        for i := 0; i < count; i++ {
            fmt.Printf("%02X ", buf[i])

            if (i + 1) % 16 == 0 {
                fmt.Printf("\n")
            }
        }
    }

    fmt.Printf("\n\n")
    fmt.Println("Done")
}
