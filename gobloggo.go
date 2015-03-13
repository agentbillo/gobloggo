package main
// to compile and run:
// go install github.com/agentbillo/gobloggo
// ~/go/bin/gobloggo
//
// to run test:
// gobloggo -blog ~/Sites/blog -master ~/Sites/blog
//
import "fmt"
import "os"
import "os/exec"
import "strings"
import "sort"
import "regexp"
import "bufio"
import "io/ioutil"

//import "time"

var forcemarkdown bool = false;
var blogdir string = "/users/billo/sites/egopoly.com/blog"
var masterdir string = "/users/billo/sites/egopoly.com/"

type blogmonth struct {
    year string
    month string
    posts []string
}

type blogpost struct {
    title string
    url string
    baseurl string
    monthdir string
    postfile string
    name string
    year string
    month string
    shtmlfile string
    preview string
    stamp string
}

var monthmap = make(map[string]*blogmonth)

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

// returns true if path1 ctime is older than path2 ctime
// if either path is unstat-able then return false
func isolder(path1 string, path2 string) bool {
    info1,err := os.Stat(path1)
    if err == nil {
        info2,err := os.Stat(path2)
        if err == nil {
            //fmt.Printf("time comparison %s %s %s %s\n", path1, info1.ModTime(), path2, info2.ModTime())
            if info1.ModTime().Before(info2.ModTime()) {
                return true
            }
            return false
        }
    }
    return false
}

func pathexists(path string) bool {
    if _, err := os.Stat(path); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func postprocess (monthdir string, year string, month string, postfile string) {
    var thismonth *blogmonth
    thismonth = monthmap[monthdir]
    if thismonth == nil {
        thismonth = &blogmonth{year, month, make([]string, 0)}
    }
    posts := thismonth.posts
    re := regexp.MustCompile("^(.*)\\.txt")
    matches := re.FindStringSubmatch(postfile)
    slug := ""
    if len(matches) > 1 {
        slug = matches[1]
        //fmt.Printf("slug is %s\n", slug)
    }
    title := "wut"
    datestamp := "1970-01-01 00:00"
    preview := ""

    f,err := os.Open(monthdir + "/" + postfile)
    if err == nil {
        defer f.Close()
        
        scanner := bufio.NewScanner(f)
        lc := 0
        for scanner.Scan() {
            line := scanner.Text()
            if lc == 0 {
                title = line
            } else if lc == 1 {
            } else if lc == 2 {
                datestamp = line
            } else {
                preview += line + "\n"
            }
            lc++
        }
        fmt.Printf("title = %s\n", title)
        //fmt.Printf("stamp = %s\n", datestamp)
        //fmt.Println(preview)

        filepath := fmt.Sprintf("%s/%s.txt", monthdir, slug)
        htmlpath := fmt.Sprintf("%s/%s.ihtml", monthdir, slug)
        //shtmlpath := fmt.Sprintf("%s/%s.shtml", monthdir, slug)
        shtmlbasepath := fmt.Sprintf("%s.shtml", slug)
        plainhtmlbasepath := fmt.Sprintf("%s.html", slug)
        url := fmt.Sprintf("%s/%s/%s.html", year, month, slug)

        post := &blogpost{title, url, plainhtmlbasepath, monthdir, postfile, slug, year,
            month, shtmlbasepath, preview, datestamp}

        mdcmd := fmt.Sprintf("Markdown.pl --html4tags %s --output %s", filepath, htmlpath)
        if forcemarkdown || !pathexists(htmlpath) || isolder(htmlpath, filepath) {
            cmd := exec.Command("Markdown.pl", "--html4tags", filepath)
            output, err := cmd.Output()
            if err != nil {
                fmt.Printf("***Markdown failed: %s\n", err)
                return
            }
            err = ioutil.WriteFile(htmlpath, output, 0644)
            if err != nil {
                panic(err)
            }
            fmt.Printf("%s\n", post.title)
            //fmt.Printf("%s\n", shtmlpath)
            fmt.Printf("%s\n", mdcmd)
        }

        


    }

    
    if len(posts) > 0 {
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

func monthscan (path string, year string,  month string) {
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
    for _,tf := range textfiles {
        postprocess(monthdir, year, month, tf)
    }
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
        monthscan(yeardir, year, month)
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

