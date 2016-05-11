//-*- mode: go -*-

package main

import (
    "io"
    "os"
    "os/user"
    "flag"
    "log"
    "path"
    "strings"
    "net/url"
    "gopkg.in/gcfg.v1"
    "github.com/mitchellh/goamz/aws"
    "github.com/mitchellh/goamz/s3"
)

type Options struct {
    Runtime struct {
        Mountpoint string
        Bucket string
    }
}

var configFile string
var options Options
var s3Client *s3.S3
var s3Bucket *s3.Bucket

var uri string

var app string
var bucketName string
var filePath string
var fullPath string
var mediaId string
var mount string

var Logger *log.Logger;

func init(){
    flag.StringVar(&configFile, "file", "/etc/threepio.ini", "ConfigFile file for threepio")
    flag.StringVar(&configFile, "f", "/etc/threepio.ini", "ConfigFile file for threepio (Shorthand)")

    flag.StringVar(&uri, "uri", "threepio+prelude:///some/path?id=12345", "Project URI; see docs")
    flag.StringVar(&uri, "u", "threepio+prelude:///some/path?id=12345", "Project URI; see docs (shorthand)")

    usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }

    auth, err := aws.EnvAuth()
    if err != nil {
        log.Fatal(err)
    }

    s3Client = s3.New(auth, aws.EUWest)

    file, _ := os.OpenFile(path.Join(usr.HomeDir, ".threepio.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    Logger = log.New(file,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func readOptions() (string, string){
    err := gcfg.ReadFileInto(&options, configFile)
    if err != nil {
        log.Fatalf("Failed to parse gcfg data: %s", err)
    }

    return options.Runtime.Mountpoint, options.Runtime.Bucket
}

func parseUri(uri string) (string, string, string) {
    urlObj, err := url.Parse(uri)
    if err != nil {
        log.Fatal( err )
    }

    queryObj, err := url.ParseQuery(urlObj.RawQuery)
    if err != nil {
        log.Fatal( err )
    }

    schemeSplit := strings.Split(urlObj.Scheme, "+")
    return schemeSplit[len(schemeSplit)-1], urlObj.Path, queryObj.Get("id")
}

func createDirIfMissing(){
    err := os.MkdirAll(fullPath, 0755)
    if err != nil {
        Logger.Fatal( err )
    }
}

func syncAssets(){
    auth, err := aws.EnvAuth()
    if err != nil {
        Logger.Fatal(err)
    }
    s3Client = s3.New(auth, aws.EUWest)
    s3Bucket = s3Client.Bucket(bucketName)

    prefix := mediaId + "/"
    resp, err := s3Bucket.List(prefix, "/", "", 1000)

    if err != nil {
        Logger.Fatal(err)
    }

    for _,c := range resp.Contents {
        filename := strings.TrimPrefix(c.Key, prefix)

        if filename == "" {
            continue
        }

        outFile, err := os.Create(path.Join(fullPath, filename))
        rc,err := s3Bucket.GetReader(c.Key)

        if err != nil {
            Logger.Fatal(err)
        }

        defer outFile.Close()
        _, err = io.Copy(outFile, rc)

        if err != nil {
            Logger.Fatal(err)
        }


        Logger.Println(c)
    }
}

func launch(){
}

func main(){
    flag.Parse()

    app, filePath, mediaId = parseUri(uri)
    mount, bucketName = readOptions()

    fullPath = path.Join(mount, filePath)
    Logger.Printf("Launching %s on path %s to edit project %s with assets from %s", app, fullPath, mediaId, bucketName)

    // Lets go to work
    createDirIfMissing()
    syncAssets()

}
