package main

import (
    "fmt"
    "time"
    "os"
    "flag"
    "path/filepath"
)

import "github.com/gosexy/exif"

// todo:
// - take arguments for input dir and output dir
// - fallback to file date/time if no exif data found
// - warn and skip when non-images are found
// - support for movies with hash tags
// - color support for terminal
// - multithreading

func main() {
    // inputFile := "./_input/test.jpg"
    inputDir := flag.Arg(0)
    outputDir := flag.Arg(1)
    fmt.Printf("Vars: %s %s\n\n", inputDir, outputDir)
    filepath.Walk("./_input", fileScanFunc)
    fmt.Println("Done.")
}

func fileScanFunc(fileName string, _ os.FileInfo, _ error) (err error) {
    fmt.Printf("[.] Trying %s\n", fileName)
    if err != nil { return err }

    stat, err := os.Stat(fileName)
    if stat.IsDir() { return nil }

    reader := readExifData(fileName)
    parsedDate, _ := getDateTime(reader)
    newDir := constructPath(parsedDate)

    fmt.Println("[+] creating dir ", newDir)
    if os.MkdirAll(newDir, 0777) != nil {
        panic(fmt.Sprintf("[!] Unable to create directory: %s\n", newDir))
    }

    destinationFile := fmt.Sprint(newDir, constructFileName(fileName, parsedDate))
    fmt.Printf("[+] copying file %s to %s\n", fileName, destinationFile)
    copyFile(fileName, destinationFile)

    return nil
}

func readExifData(inputFile string) (*exif.Data) {
    // open an exif reader and read the file from system
    reader := exif.New()
    err := reader.Open(inputFile)

    // check to make sure we can read the exif data
    if err != nil {
        fmt.Println("[!] Error:", err.Error())
    }

    return reader
}

func constructPath(parsedDate time.Time) string {
    return fmt.Sprintf("./_output/%d/%02d/", parsedDate.Year(), int(parsedDate.Month()))
}

func constructFileName(inputFile string, parsedDate time.Time) string {
    fileHash, _ := getHash(inputFile)
    return fmt.Sprintf("%d-%02d-%02d-%02d-%02d.jpg", parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), parsedDate.Minute(), fileHash)
}

func getDateTime(reader *exif.Data) (time.Time, error) {
    exifDateLayout := "2006:01:02 15:04:05"
    parsedDate := reader.Tags["Date and Time"]
    return time.Parse(exifDateLayout, parsedDate)
}
