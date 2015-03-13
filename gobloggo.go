package main
// This is my first go program. I'm probably doing everything wrong, wicked sorry.
//
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
    posts []*blogpost
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

var monthindexblock string = ""
var tweetblock string = ""

var monthmap = make(map[string]*blogmonth)
var allposts = make(map[string]*blogpost)

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

//
// This is where each post is created.  For each .txt file, create a container .shtml file that is
// full page, and a .ihtml file that is the markdown-converted version of the .txt
// 
func postprocess (monthdir string, year string, month string, postfile string) {
    var thismonth *blogmonth
    thismonth = monthmap[monthdir]
    if thismonth == nil {
        thismonth = &blogmonth{year, month, make([]*blogpost, 0)}
        monthmap[monthdir] = thismonth
    }
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
        shtmlpath := fmt.Sprintf("%s/%s.shtml", monthdir, slug)
        shtmlbasepath := fmt.Sprintf("%s.shtml", slug)
        plainhtmlbasepath := fmt.Sprintf("%s.html", slug)
        url := fmt.Sprintf("%s/%s/%s.html", year, month, slug)

        post := &blogpost{title, url, plainhtmlbasepath, monthdir, postfile, slug, year,
            month, shtmlbasepath, preview, datestamp}

        thismonth.posts = append(thismonth.posts, post)
        allposts[datestamp] = post

        shtmlf, err := os.Create(shtmlpath)
        if err != nil {
            panic(err)
        }
        defer shtmlf.Close()

        writer := bufio.NewWriter(shtmlf)
        _, err = writer.WriteString("<!-- made by go blog go! -->\n")
        _, err = writer.WriteString(fmt.Sprintf("<!--#set var=\"title\" value=\"%s\" -->\n", title))
        _, err = writer.WriteString("<!--#include virtual=\"/header.shtml\" -->\n")
        _, err = writer.WriteString(fmt.Sprintf("<!--#include virtual=\"%s.ihtml\" -->\n", slug))
        _, err = writer.WriteString("<!--#include virtual=\"/footer.shtml\" -->\n")
        writer.Flush()

        mdcmd := fmt.Sprintf("Markdown.pl --html4tags %s > %s", filepath, htmlpath)
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

    
    if len(thismonth.posts) > 0 {
        fmt.Printf("posts len = %d\n", len(thismonth.posts))
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

func reversekeys (m map[string]*blogmonth) []string {
    keys := make([]string, 0, len(m))
    for key := range m {
        keys = append(keys, key)
    }
    sort.Sort(sort.Reverse(sort.StringSlice(keys)))
    
    return keys
}

func reversepostkeys (m map[string]*blogpost) []string {
    keys := make([]string, 0, len(m))
    for key := range m {
        keys = append(keys, key)
    }
    sort.Sort(sort.Reverse(sort.StringSlice(keys)))
    
    return keys
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// 
// Assembles master index, month indicies and RSS feed.
//
func postdump () {
    fmt.Println("***POST DUMP***")
    dat, err := ioutil.ReadFile(masterdir + "/" + "monthindex.shtml")
    check(err)
    monthindexblock = string(dat)

    dat, err = ioutil.ReadFile(masterdir + "/" + "tweet.shtml")
    check(err)
    tweetblock = string(dat)


    sidebarf, err := os.Create(fmt.Sprintf("%s/sidebar.shtml", blogdir))
    if err != nil {
        panic(err)
    }
    defer sidebarf.Close()

    sidebarf.WriteString("<div class=\"sidebar\">\n")
    sidebarf.WriteString("<p>Archive</p>\n<ul class=\"sidebarlist\">")

    reversemonths := reversekeys(monthmap)
    for _,k := range reversemonths {
        month := monthmap[k]
        posts := month.posts
        monthurl := strings.Replace(k, blogdir + "/", "", 1)
        link := fmt.Sprintf("<li><a href=\"/%s}\">%s(%d)</a></li>\n", monthurl, 
            monthurl, len(posts))
        sidebarf.WriteString(link)
        //fmt.Println(posts)

        // write the month index here
        contentpath := k + "/contents.shtml"
        indexpath := k + "/index.shtml"
        indexf, err := os.Create(indexpath)
        if err != nil {
            panic(err)
        }
        defer indexf.Close()
        titleset := fmt.Sprintf("<!--#set var=\"title\" value=\"%s\" -->\n", monthurl)
        indexf.WriteString(titleset)
        indexf.WriteString(monthindexblock)

        mcontentf, err := os.Create(contentpath)
        if err != nil {
            panic(err)
        }
        defer mcontentf.Close()
        for _,p := range posts {
            entry := fmt.Sprintf("<div class=\"monthindexentry\"><a href=\"%s\">%s</a></div>\n",
                p.baseurl, p.title)
            mcontentf.WriteString(entry)
        }
        

    }
    sidebarf.WriteString("</ul></div>\n")

    contentf, err := os.Create(fmt.Sprintf("%s/contents.shtml", blogdir))
    if err != nil {
        panic(err)
    }
    defer contentf.Close()
    reversestamps := reversepostkeys(allposts)
    count := 0
    for _,k := range reversestamps {
        post := allposts[k]
        if count < 5 {
            contentf.WriteString("<div class=\"frontindexentry\">\n")
            ps := fmt.Sprintf("<!--#include virtual=\"%s/%s/%s.ihtml\" -->\n", 
                post.year, post.month, post.name)
            contentf.WriteString(ps)
            ps = fmt.Sprintf("<a href=\"%s\">Permalink</a>\n", post.url)
            contentf.WriteString(ps)
            contentf.WriteString("</div>\n")
            
        }
        count++
    }
    
}


func main() {
    options()
    blogscan(blogdir)
    postdump()
}

