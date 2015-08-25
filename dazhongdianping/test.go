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
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/henrylee2cn/pholcus/common/mlog"
	"github.com/hu17889/go_spider/core/common/page"
	"github.com/hu17889/go_spider/core/pipeline"
	"github.com/hu17889/go_spider/core/spider"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

// Parse html dom here and record the parse result that we want to Page.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}
	// fmt.Println("json===\n")
	// fmt.Println(*p.GetJson())
	// fmt.Println("Rquest===\n")
	// res2B, _ := json.Marshal(*p.GetRequest())
	// fmt.Println(string(res2B))
	// fmt.Println(*p.GetRequest())

	// fmt.Println("URLTAG====\n", p.GetUrlTag())
	query := p.GetHtmlParser()

	var urls []string
	query.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		// s.
		// urls = append(urls, "http://github.com/"+href)
		// fmt.Println("======================================")
		// fmt.Println(s.Attr("href"))
		// fmt.Println(href)
		// fmt.Println(s.AttrOr(attrName, defaultValue))
		// fmt.Println("======================================")
		var absHref string
		urlHref, err := url.Parse(href)
		if err != nil {
			mlog.LogInst().LogError(err.Error())
			return
		}
		if !urlHref.IsAbs() {
			urlPrefix := p.GetRequest().GetUrl()
			absHref = urlPrefix + href
			absHref = absHref
			// fmt.Println(absHref)
			urls = append(urls, absHref)
		} else {
			// fmt.Println(href)
			urls = append(urls, href)
		}

	})
	// these urls will be saved and crawed by other coroutines.
	p.AddTargetRequests(urls, "html")
	// content, _ := query.Html()
	// fmt.Println(content)

	name := query.Find(".entry-title .author").Text()
	name = strings.Trim(name, " \t\n")
	repository := query.Find(".entry-title .js-current-repository").Text()
	repository = strings.Trim(repository, " \t\n")
	//readme, _ := query.Find("#readme").Html()
	if name == "" {
		p.SetSkip(true)
	}
	// the entity we want to save by Pipeline
	p.AddField("author", name)
	p.AddField("project", repository)
	//p.AddField("readme", readme)
}

func (this *MyPageProcesser) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}

func main() {
	// Spider input:
	//  PageProcesser ;
	//  Task name used in Pipeline for record;
	spider.NewSpider(NewMyPageProcesser(), "TaskName").
		AddUrl("http://t.dianping.com/beijing", "html"). // Start url, html is the responce type ("html" or "json" or "jsonp" or "text")
		AddPipeline(pipeline.NewPipelineConsole()).      // Print result on screen
		SetThreadnum(3).                                 // Crawl request by three Coroutines
		Run()
}
