package main

import "fmt"
import "time"
import "os"
import "io/ioutil"
import "io"
import "github.com/gosexy/exif"
import "hash/crc32"

// todo:
// - recursively parse directory
// - warn and skip when non-images are found
// - take arguments for input dir and output dir
// - fallback to file date/time if no exif data found
// - support for movies with hash tags
// - color support for terminal

func main() {
    inputFile := "./_input/test.jpg"
    // open an exif reader and read the file from system
    reader := exif.New()
    err := reader.Open(inputFile)

    // check to make sure we can read the exif data
    if err != nil {
        fmt.Printf("Error: %s", err.Error())
    }

    parsedDate, _ := getDateTime(reader)
    newDir := constructPath(parsedDate)

    fmt.Println("[+] creating dir ", newDir)
    if os.MkdirAll(newDir, 0777) != nil {
        panic(fmt.Sprintf("Unable to create directory: %s\n", newDir))
    }

    destinationFile := fmt.Sprint(newDir, constructFileName(inputFile, parsedDate))
    fmt.Println("[+] creating file ", destinationFile)
    copyFile(inputFile, destinationFile)
}

func getHash(filename string) (uint32, error) {
    bs, err := ioutil.ReadFile(filename)
    if err != nil {
        return 0, err
    }
    h := crc32.NewIEEE()
    h.Write(bs)
    return h.Sum32(), nil
}

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
