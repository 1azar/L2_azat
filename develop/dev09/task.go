package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type config struct {
	url    *url.URL
	depth  int
	folder string
}

var lg *log.Logger

func main() {
	lg = log.New(os.Stdout, "", 0)

	urlRaw := flag.String("url", "", "website url: https://hackerspaces.org")
	depth := flag.Int("depth", 0, "limits the recursion depth of visited URLs")
	folderPath := flag.String("folder", ".\\", "path to the save folderPath")

	flag.Parse()

	webUrl, err := url.Parse(*urlRaw)
	if err != nil {
		lg.Fatal("invalid url:", *urlRaw)
	}

	//if *depth < 1 {
	//	lg.Fatal("depth value must be more than 0")
	//}

	_, err = os.Stat(*folderPath)
	if os.IsNotExist(err) {
		lg.Fatal("such folder does not exists:", folderPath)
	}

	cfg := config{
		url:    webUrl,
		depth:  *depth,
		folder: *folderPath,
	}

	gatherWebSite(&cfg)
}

func gatherWebSite(cfg *config) {

	pageVisited := make(map[string]string) // web url: local path

	rootDir, err := filepath.Abs(filepath.Join(cfg.folder, cfg.url.Host))
	pageVisited["/"] = filepath.Join(rootDir, cfg.url.Host+".htm")
	pageVisited[""] = filepath.Join(rootDir, cfg.url.Host+".htm")
	if err != nil {
		lg.Fatal("failed to create folder for website", err)
	}
	err = os.Mkdir(rootDir, os.ModePerm)
	if !errors.Is(err, os.ErrExist) && err != nil {
		lg.Fatal("could not create folder:", rootDir, err)
	}

	c := colly.NewCollector(
		colly.AllowedDomains(cfg.url.Host),
		colly.IgnoreRobotsTxt(),
		colly.MaxDepth(cfg.depth),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		lg.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		lg.Println("Visiting " + r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		// check if this page was visited before
		//filePath, visited := pageVisited[GetFullUrl(r.Request.URL)]

		pageName := GetLastEntryInUrl(r.Request.URL)

		pageExtensions, err := mime.ExtensionsByType(r.Headers.Get("Content-Type"))
		pageExtension := pageExtensions[0]
		if err != nil {
			lg.Println("could not determine extension for content type", pageExtension, err)
		}

		if !in(pageExtension, []string{"html", ".htm"}) {
			return
		}

		fileName := pageName + pageExtension

		dir := filepath.Join(rootDir, GetUpperEntry(r.Request.URL))

		err = os.MkdirAll(dir, os.ModePerm)
		if !errors.Is(err, os.ErrExist) && err != nil {
			lg.Fatal("could not create folder:", dir, err)
		}

		filePath := filepath.Join(dir, fileName)

		pageVisited[r.Request.URL.Path] = filePath

		r.Body, err = replaceLinks(r, pageVisited)
		if err != nil {
			lg.Println("Failed to replace href in page:", err)
		}

		lg.Println("Saving page:", r.Request.URL.Path, "->", filePath)

		// save imgs, css, etc - пока не работает. поэтому без изображений, стилей, жс
		//if pageExtension == ".html" || pageExtension == ".htm" {
		//	newBody, err := DownloadImgsAndReplaceSrc(r.Body, dir, GetFullUrl(r.Request.URL))
		//	if err != nil {
		//		lg.Println("failed to process html body", err)
		//	}
		//	r.Body = newBody
		//}

		// a href replace

		r.Save(filePath)

	})

	c.Visit(GetFullUrl(cfg.url))

	//c.Request()

	fmt.Print(len(pageVisited))
}

func GetFullUrl(u *url.URL) string {
	return u.Scheme + "://" + u.Host + u.Path
}

func GetLastEntryInUrl(u *url.URL) string {
	u2 := u.Host + u.Path
	t := strings.Split(u2, "/")
	return t[len(t)-1]
}

func GetUpperEntry(u *url.URL) string {
	//t := strings.Split(u.Host+u.Path, "/")
	t := strings.Split(u.Path, "/")
	return strings.Join(t[:len(t)-1], "/")
}

func in[T comparable](e T, sl []T) bool {
	for _, i := range sl {
		if e == i {
			return true
		}
	}
	return false
}

func replaceLinks(r *colly.Response, linkMap map[string]string) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(r.Body))
	if err != nil {
		return nil, err
	}

	var replace func(*html.Node)
	replace = func(n *html.Node) {
		//if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "link" || n.Data == "script" || n.Data == "img") {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i, attr := range n.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					if replacement, ok := linkMap[attr.Val]; ok {
						n.Attr[i].Val = replacement
					} else {
						err = r.Request.Visit(attr.Val)
						if err == nil {
							n.Attr[i].Val = linkMap[attr.Val]
						}
						if err != nil {
							fmt.Println("could not visit:", attr.Val, err)
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			replace(c)
		}
	}

	replace(doc)

	var result bytes.Buffer
	err = html.Render(&result, doc)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

func DownloadImgsAndReplaceSrc(body []byte, dir string, url string) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return []byte{}, err
	}

	var processNode func(*html.Node)
	processNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					imageURL := attr.Val
					localPath, err := downloadImage(dir, url+imageURL)
					if err != nil {
						lg.Print("failed to load img:", imageURL, err)
					}
					lg.Println("succeeded download img:", path.Base(localPath))
					attr.Val = localPath
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processNode(c)
		}
	}

	processNode(doc)

	var result bytes.Buffer
	err = html.Render(&result, doc)
	if err != nil {
		lg.Println("failed to render html", err)
	}

	return result.Bytes(), nil
}

func downloadImage(dir string, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	localPath := filepath.Join(dir, filepath.Base(url))
	err = os.WriteFile(localPath, body, 0644)
	if err != nil {
		return "", err
	}

	return localPath, nil
}
