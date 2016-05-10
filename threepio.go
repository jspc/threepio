//-*- mode: go -*-

package main

import (
    "os"
    "os/user"
    "flag"
    "log"
    "path"
    "strings"
    "net/url"
    "gopkg.in/gcfg.v1"
)

type Options struct {
    Runtime struct {
        Mountpoint string
    }
}

var configFile string
var options Options

var uri string

var app string
var filePath string
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
        "THREPIO: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func readOptions() string{
    err := gcfg.ReadFileInto(&options, configFile)
    if err != nil {
        log.Fatalf("Failed to parse gcfg data: %s", err)
    }

    return options.Runtime.Mountpoint
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

func main(){
    flag.Parse()

    app, filePath, mediaId = parseUri(uri)
    mount = readOptions()

    fullPath := path.Join(mount, filePath)
    Logger.Printf("Launching %s on path %s to edit project %s", app, fullPath, mediaId)
}
