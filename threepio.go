//-*- mode: go -*-

package main

import (
    "os"
    "os/user"
    "flag"
    "log"
    "path"
)

var app string
var uri string

var Logger *log.Logger;

func init(){
    flag.StringVar(&app, "app", "prelude", "Application to execute")
    flag.StringVar(&app, "a", "prelude", "Application to execute (shorthand)")

    flag.StringVar(&uri, "uri", "/some/path?id=12345", "Project URI; see docs")
    flag.StringVar(&uri, "u", "/some/path?id=12345", "Project URI; see docs (shorthand)")

    usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }

    file, _ := os.OpenFile(path.Join(usr.HomeDir, ".threepio.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    Logger = log.New(file,
        "THREPIO: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func main(){
    flag.Parse()
    Logger.Printf("Loading %s in %s\n", uri, app)

}
