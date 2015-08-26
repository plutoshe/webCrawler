//
package main

/*
Packages must be imported:
    "core/common/page"
    "core/spider"
Pckages may be imported:
    "core/pipeline": scawler result persistent;
    "github.com/PuerkitoBio/goquery": html dom parser.
*/
import (
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/henrylee2cn/pholcus/common/mlog"
	"github.com/hu17889/go_spider/core/common/page"
	"github.com/hu17889/go_spider/core/pipeline"
	"github.com/hu17889/go_spider/core/spider"
	"github.com/plutoshe/webCrawler/storage"
	"gopkg.in/mgo.v2"
)

var collection *mgo.Collection

var base string = "www.dazhongdainping.com"

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

func checkMatchPattern(base, href string) bool {
	inDomain := regexp.MustCompile("dianping\\.com")
	if inDomain.MatchString(href) {
		// fmt.Println("in  ", href)
		return true
	}
	return false
}

// Parse html dom here and record the parse result that we want to Page.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}
	query := p.GetHtmlParser()
	currentUrl := p.GetRequest().GetUrl()
	var urls []string
	query.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		urlHref, err := url.Parse(href)
		if err != nil {
			mlog.LogInst().LogError(err.Error())
			return
		}
		if !urlHref.IsAbs() {
			href = currentUrl + href
		}
		// Temporarily check in crawler.go, it will be implemented in pattern package.
		if checkMatchPattern(base, href) {
			urls = urls
			// urls = append(urls, href)
		}
	})

	// store content to db

	fmt.Println("==store==", currentUrl)
	content, _ := query.Html()
	storage.StoreInsert(collection, storage.StoreFormat{currentUrl, content})

}

func (this *MyPageProcesser) Finish() {
	// fmt.Printf("TODO:before end spider \r\n")

}

func main() {
	// db initilization
	dbSession, err := storage.Link2DbByDefault()
	defer dbSession.Close()
	if err != nil {
		log.Fatal(err)
	}
	collection = storage.Link2CollectionByDefault(dbSession)

	// Spider input:
	//  PageProcesser ;
	//  Task name used in Pipeline for record;
	// TO-DO :
	// Change to goroutine mechanism, use channel to get url
	// Able to start serveral goroutine to crawler a same website.
	spider.NewSpider(NewMyPageProcesser(), "TaskName").
		AddUrl("http://www.dianping.com", "html").  // Start url, html is the responce type ("html" or "json" or "jsonp" or "text")
		AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
		SetThreadnum(3).                            // Crawl request by three Coroutines
		Run()
}
