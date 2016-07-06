//-*- mode: go -*-

package main

import (
    "flag"
    "github.com/mitchellh/goamz/aws"
    "github.com/mitchellh/goamz/s3"
    "gopkg.in/gcfg.v1"
    "io"
    "log"
    "net/url"
    "os"
    "os/user"
    "path"
    "strings"

    "syscall"
    "os/exec"
)

type Options struct {
    Runtime struct {
        Mountpoint string
        Bucket string
    }
}

type AppContext struct {
    Application string
    Path string
    Id string
    AWS struct {
        AccessKey string
        SecretKey string
        Token string
    }
}

var appContext AppContext
var configFile string
var options Options
var s3Bucket *s3.Bucket
var s3Client *s3.S3

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

func parseUri(uri string) {
    urlObj, err := url.Parse(uri)
    if err != nil {
        log.Fatal( err )
    }

    queryObj, err := url.ParseQuery(urlObj.RawQuery)
    if err != nil {
        log.Fatal( err )
    }

    schemeSplit := strings.Split(urlObj.Scheme, "+")

    appContext.Application = schemeSplit[len(schemeSplit)-1]
    appContext.Path = urlObj.Path
    appContext.Id = queryObj.Get("id")

    appContext.AWS.AccessKey = queryObj.Get("accessKey")
    appContext.AWS.SecretKey = queryObj.Get("secretKey")
    appContext.AWS.Token = queryObj.Get("sessionToken")
}

func createDirIfMissing(){
    err := os.MkdirAll(fullPath, 0755)
    if err != nil {
        Logger.Fatal( err )
    }
}

func syncAssets(){
    var auth aws.Auth
    auth.AccessKey = appContext.AWS.AccessKey
    auth.SecretKey = appContext.AWS.SecretKey
    auth.Token = appContext.AWS.Token

    s3Client = s3.New(auth, aws.EUWest)
    s3Bucket = s3Client.Bucket(bucketName)

    prefix := appContext.Id + "/"
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
    binary, lookErr := exec.LookPath("open")
    if lookErr != nil {
        panic(lookErr)
    }

    args := []string{"open", path.Join(fullPath, "my_project.plproj")}

    env := os.Environ()

    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        Logger.Fatal(execErr)
    }
}

func main(){
    flag.Parse()

    parseUri(uri)
    mount, bucketName = readOptions()

    fullPath = path.Join(mount, filePath)
    Logger.Printf("Launching %s on path %s to edit project %s with assets from %s",
        appContext.Application, fullPath, appContext.Id, bucketName)

    // Lets go to work
    createDirIfMissing()
    syncAssets()
    launch()
}
