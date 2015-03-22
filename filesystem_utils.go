package main

import (
    "os"
    "io"
)

func copyFile(src string, dest string) {
    reader, err := os.Open(src)
    if err != nil {
        panic(err)
    }
    defer reader.Close()

    writer, err := os.Create(dest)
    if err != nil {
        panic(err)
    }
    defer writer.Close()

    _, err = io.Copy(writer, reader)
    if err != nil {
        panic(err)
    }
}

// func forEachFileInDir(targetFn func(string) error, path string) {
//     filepath.Walk(path, )
//     targetFn(path)
// }
