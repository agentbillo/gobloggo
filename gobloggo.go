package main
// to compile and run:
// go install github.com/agentbillo/gobloggo
// ~/go/bin/gobloggo
import "fmt"
import "os"
import "strings"
import "sort"

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

func listdir (path string) []string {
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
    entries := make([]string, 0)
    for _,f := range files {
        name := f.Name()
        entries = append(entries, name)
    }
    sort.Strings(entries)
    return entries
}

func monthscan (path string, month string) {
    fmt.Println(month)
    monthdir := path + "/" + month
    entries := listdir(monthdir)
    textfiles := make([]string, 0)
    for _,e := range entries {
        if strings.HasSuffix(e, ".txt") {
            textfiles = append(textfiles, e)
        }
    }
    sort.Strings(textfiles)
    fmt.Println(textfiles)
}


func yearscan (path string, year string) {
    fmt.Println(year)
    yeardir := path + "/" + year
    entries := listdir(yeardir)
    months := make([]string, 0)
    for _,e := range entries {
        if len(e) == 2 {
            months = append(months, e)
        }
    }
    sort.Strings(months)
    for _,month := range months {
        monthscan(yeardir, month)
    }
}


func blogscan (path string) {

    entries := listdir(path)
    years := make([]string, 0)
    for _,e := range entries {
        // go blog go only worjs in the 21 century
        if strings.HasPrefix(e, "20") && len(e) == 4 {
            years = append(years, e)
        }
    }
    sort.Strings(years)
    fmt.Println(years)
    for _,year := range years {
        yearscan(path, year)
    }
}

func main() {
    options()
    blogscan(blogdir)
}

