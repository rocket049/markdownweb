package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/kataras/iris"

	_ "github.com/mattn/go-sqlite3"
	md "github.com/russross/blackfriday/v2"
)

const mdTmpl = `<html>
<head>
<meta http-equiv="content-type" content="text/html;charset=utf-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<link type="text/css" rel="stylesheet" href="/style.css"/>
<title>{{.title}}</title>
<head>
<body>
{{.body}}
</body>
</html>`

var logger *log.Logger

func relatePath(items ...string) string {
	exe1, _ := os.Executable()
	base := filepath.Dir(exe1)
	paths := append([]string{base}, items...)
	res := filepath.Join(paths...)
	if strings.Contains(res, "..") {
		return "404"
	}
	return res

}

func initMdTmpl(path1 string) string {
	//get relate template
	tpl, err := ioutil.ReadFile(filepath.Join(path.Dir(path1), "md.tpl"))
	if err == nil {
		//logger.Println(filepath.Join(path.Dir(path1), "md.tpl"))
		return string(tpl)
	}

	//get top template
	tpl, err = ioutil.ReadFile(relatePath("md.tpl"))
	if err == nil {
		//logger.Println(relatePath("md.tpl"))
		return string(tpl)
	}

	//get default template
	return mdTmpl
}

func getSize(filename string) int64 {
	fh, err := os.Stat(filename)
	if err != nil {
		return 0
	}
	return fh.Size()
}

func sendFile(ctx iris.Context, filename string) error {
	fname := relatePath("files", filename)
	info, err := os.Stat(fname)
	if err != nil {
		//ctx.StatusCode(404)
		logger.Println(err)
		adRedirect(ctx)
		return err
	}

	ctx.Header("Content-Length", fmt.Sprint(info.Size()))
	if strings.HasSuffix(filename, ".wasm") {
		//ctx.Header("Content-Type", "application/wasm")
		ctx.ContentType("application/wasm; binary")
		fp, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer fp.Close()
		ctx.StatusCode(200)
		io.Copy(ctx.ResponseWriter(), fp)
	}
	ctx.SendFile(fname, filename)

	return nil
}

func getTitle(p []byte) string {
	n := bytes.IndexByte(p, '\n')
	if n <= 0 {
		return ""
	}
	line1 := string(p[0:n])
	return strings.Trim(line1, "# \n\r\t")
}

func sendMarkdown(ctx iris.Context, filename string) {
	fname := relatePath("files", filename)
	//get cache dir
	fdir := filepath.Dir(fname)
	cacheDir := filepath.Join(fdir, ".cache")
	os.MkdirAll(cacheDir, os.ModePerm)
	// is exist file
	fstat, err := os.Lstat(fname)
	if err != nil {
		//ctx.StatusCode(404)
		logger.Println(err)
		adRedirect(ctx)
		return
	}

	cacheName := filepath.Join(cacheDir, filepath.Base(filename)+".htm")
	cacheStat, err := os.Stat(cacheName)
	if err == nil {
		if cacheStat.ModTime().Unix() > fstat.ModTime().Unix() {
			ctx.ServeFile(cacheName, ctx.ClientSupportsGzip())
			return
		}
	}

	fp, err := os.Create(cacheName)
	if err != nil {
		//ctx.StatusCode(404)
		logger.Println(err)
		adRedirect(ctx)
		return
	}
	defer fp.Close()

	writer := io.MultiWriter(fp, ctx.ResponseWriter())

	data := make(map[string]interface{})

	file1, err := os.Open(fname)
	if err != nil {
		//ctx.StatusCode(404)
		logger.Println(err)
		adRedirect(ctx)
		return
	}
	stat1, _ := file1.Stat()
	var buf []byte = make([]byte, stat1.Size())
	n, _ := io.ReadFull(file1, buf)
	file1.Close()
	if int64(n) == stat1.Size() {
		data["title"] = getTitle(buf)

		body := md.Run(buf, md.WithExtensions(md.CommonExtensions))
		data["body"] = string(body)
		v, _ := getAdList()
		data["ad"] = v
		adLock.Lock()
		adList = v
		adLock.Unlock()

	} else {
		ctx.StatusCode(500)
		return
	}
	tmpl := initMdTmpl(fname)
	t := template.New("")
	t.Parse(tmpl)
	t.Execute(writer, data)
}

type AdItem struct {
	Img  string `json:"img"`
	Href string `json:"href"`
	Text string `json:"text"`
}

var adLock sync.RWMutex
var adList []AdItem

func getAdList() ([]AdItem, error) {
	res := make([]AdItem, 5)
	jsonFile := relatePath("ad.json")
	jsonData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		logger.Println(err)
		return res, err
	}
	err = json.Unmarshal(jsonData, &res)
	if err != nil {
		logger.Println(err)
		return res, err
	}

	return res, nil
}

func adRedirect(ctx iris.Context) {
	adLock.RLock()
	defer adLock.RUnlock()

	l := len(adList)
	if l == 0 {
		return
	}
	rnd := rand.Int()
	n := rnd % l
	url := adList[n].Href

	ctx.Redirect(url, 302)
}

func main() {
	var addr = flag.String("addr", ":8080", "format [IP:Port]")
	var tls = flag.Bool("tls", false, "use tls or not")
	flag.Parse()

	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	adList, _ = getAdList()
	app := iris.New()

	app.Get("/{key:path}", func(ctx iris.Context) {
		fn := strings.TrimSpace(ctx.Params().Get("key"))
		if strings.Contains(fn, "..") || strings.Contains(fn, "/.") {
			logger.Printf("ERROR %s Get /%s\n", ctx.RemoteAddr(), fn)
			//ctx.StatusCode(404)
			adRedirect(ctx)
			return
		}
		if strings.HasSuffix(strings.ToLower(fn), ".md") {
			sendMarkdown(ctx, fn)
			logger.Printf("%s Get /%s\n", ctx.RemoteAddr(), fn)
		} else if strings.HasSuffix(fn, ".tpl") {
			logger.Println(".tpl")
			adRedirect(ctx)
		} else {
			err := ctx.ServeFile(relatePath("files", fn), ctx.ClientSupportsGzip())
			if err != nil {
				logger.Printf("ERROR %s Get /%s\n", ctx.RemoteAddr(), fn)
				adRedirect(ctx)
			} else {
				logger.Printf("%s Get /%s\n", ctx.RemoteAddr(), fn)
			}
		}

	})

	app.Get("/", func(ctx iris.Context) {
		sendMarkdown(ctx, "index.md")
		logger.Printf("%s Get /\n", ctx.RemoteAddr())
	})
	if *tls == false {
		err := app.Run(iris.Addr(*addr))
		log.Println(err)
		return
	}

	runner, err := getTLSRunner()
	if err != nil {
		log.Println(err)
	} else {
		err = app.Run(runner)
		log.Println(err)
	}
}
