package main

import (
    "os"
)

func main() {
    img, err := OpenImage(os.Args[1])
    if err != nil {
        panic(err)
    }
    defer img.Close()

    img.Show()
}

