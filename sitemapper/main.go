package main

import (
	"container/list"
	"encoding/xml"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"../html_link_parser/parse"
)

var urlFlag string
var domainURL string

func main() {
	flag.StringVar(&urlFlag, "urlLink", "https://gophercises.com", "url you want the sitemap for")
	depth := flag.Int("depth", 5, "depth till which links should be retrieved from the domain")
	xmlFile := flag.String("xml", "sitemap.xml", "XML filename which will contain the sitemap")

	flag.Parse()
	getDomain("")
	links := traverseDomain(urlFlag, *depth)

	f, err := os.Create(*xmlFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writeToXMLFile(links, f)
}

type UrlSet struct {
	XMLName   xml.Name `xml:"urlset"`
	NameSpace string   `xml:"xmlns,attr"`
	Links     *[]Url   `xml:"url"`
}

type Url struct {
	XMLName  xml.Name `xml:"url"`
	Location string   `xml:"loc"`
}

func writeToXMLFile(links []string, r io.Writer) {
	log.Println("Writing to XML file...")
	linkSlice := make([]Url, 0)
	urlset := &UrlSet{NameSpace: "http://www.sitemaps.org/schemas/sitemap/0.9"}

	for _, link := range links {
		linkStruct := Url{Location: link}
		linkSlice = append(linkSlice, linkStruct)
	}

	urlset.Links = &linkSlice
	r.Write([]byte(xml.Header))
	enc := xml.NewEncoder(r)
	enc.Indent("", "  ")
	if err := enc.Encode(urlset); err != nil {
		panic(err)
	}
	log.Println("...writing done.")
}

func getDomain(urlLink string) {
	log.Println("Getting domain...")
	if urlLink == "" {
		urlLink = urlFlag
	}
	resp, err := http.Head(urlLink)
	if err != nil {
		panic(err)
	}

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	domainURL = baseUrl.String()
	log.Println("...Domain got.")
}

func traverseDomain(root string, maxDepth int) []string {
	log.Println("Traversing domain...")
	visited := make(map[string]bool)
	all_links := make([]string, 0)

	queue := list.New()
	queue.PushBack(root)

	curDepth := 0
	elementsToDepthIncrease := 1
	nextElementsToDepthIncrease := 0

	visited[root] = true
	for queue.Len() > 0 {
		link := queue.Front()
		linkStr := (*link).Value.(string)

		all_links = append(all_links, linkStr)
		links := getLinksFromPage(linkStr)
		nextElementsToDepthIncrease += len(links)
		elementsToDepthIncrease--
		if elementsToDepthIncrease == 0 {
			curDepth++
			if curDepth > maxDepth {
				break
			}
			elementsToDepthIncrease = nextElementsToDepthIncrease
			nextElementsToDepthIncrease = 0
		}

		for _, urlLink := range links {
			if _, ok := visited[urlLink]; !ok {
				queue.PushBack(urlLink)
				visited[urlLink] = true
			}
		}
		queue.Remove(link)
	}

	return all_links
}

func getLinksFromPage(urlLink string) []string {
	log.Println("Gettling links from", urlLink)
	resp, err := http.Get(urlLink)
	if err != nil {
		panic(err)
	}

	links := parse.GetLinksStrings(resp.Body)
	links = buildUrls(links)
	return links

}

func buildUrls(urls []string) []string {
	retUrls := make([]string, 0)
	for _, link := range urls {
		switch {
		case strings.Contains(link, domainURL):
			retUrls = append(retUrls, link)

		case strings.HasPrefix(link, "/"):
			retUrls = append(retUrls, domainURL+link)
		}
	}
	return retUrls
}
