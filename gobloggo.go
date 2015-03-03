package main

import "fmt"
import "os"
import "strings"

var forcemarkdown bool = false;
var blogdir string = "/users/billo/sites/egopoly.com/blog"
var masterdir string = "/users/billo/sites/egopoly.com/"


func options() {
    for i := 0; i < len(os.Args); {
        arg := os.Args[i]
        switch {
        case strings.HasPrefix(arg, "-force"):
            fmt.Println("FORCE\n")
            forcemarkdown = true; 
        case strings.HasPrefix(arg, "-blog"):
            i++
            blogdir = os.Args[i]
            fmt.Printf("blogdir=%s\n", blogdir)
        case strings.HasPrefix(arg, "-master"):
            i++
            masterdir = os.Args[i]
            fmt.Printf("masterdir=%s\n", masterdir)
        }
        i++
    }
}

func blogscan (path string) {
    dir, err := os.Open(path)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer dir.Close()
    files, err := dir.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _,f := range files {
        fmt.Println(f.Name())
    }
}

func main() {
    options()
    blogscan(blogdir)
}

