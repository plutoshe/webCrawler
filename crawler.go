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
	"flag"
	"fmt"
	"log"
	"net/url"
	"regexp"

	"./urlstore"
	"github.com/PuerkitoBio/goquery"
	"github.com/henrylee2cn/pholcus/common/mlog"
	"github.com/hu17889/go_spider/core/common/page"
	"github.com/hu17889/go_spider/core/pipeline"
	"github.com/hu17889/go_spider/core/spider"
	"github.com/plutoshe/webCrawler/repetition"
	"github.com/plutoshe/webCrawler/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/redis.v3"
)

var (
	exitChan       chan struct{}
	getURL         chan string
	releaseSlot    chan int
	collection     *mgo.Collection
	rep            repetition.RepetitionJudgement
	urlstr         urlstore.URLCrawlerStore
	base           = "www.dazhongdainping.com"
	databaseURL    = flag.String("dbURL", storage.MONGODB_URL, "Identify the linked db adress")
	databaseName   = flag.String("dbname", storage.MONGODB_DB, "Denote the database name of mongodb")
	databaseAuth   = flag.Bool("dbauth", false, "Denote whether link datatbase by authentication or not")
	databaseUser   = flag.String("dbuser", storage.MONGODB_USER, "Denote the user name of db")
	databasePwd    = flag.String("dbpwd", storage.MONGODB_PWD, "Regarding the user name, denote corresponding password")
	collectionName = flag.String("collection", storage.MONGODB_COLLECTION, "Denote the coolection name to operate")
	threadNum      = flag.Int("threadNum", 4, "Specify the thread number to crawl. The default value is 4.")
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

func checkMatchPattern(base, href string) bool {
	// inDomain := regexp.MustCompile("dianping\\.com")
	inDomain := regexp.MustCompile("dianping\\.com(\\/[^\\/\n]+){2,2}$")
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
			visited, _ := rep.CheckIfVisited(href)
			if !visited {
				rep.VisitedNewNode(href)
				// urls = append(urls, href)
				urlstr.UploadURL(href)
			}
		}
	})

	// store content to db

	fmt.Printf("====store & commit : %s====\n\n\n", currentUrl)
	content, _ := query.Html()
	// content := ""
	storage.StoreInsert(collection, storage.StoreFormat{currentUrl, content})
	urlstr.CommitURL(currentUrl)
	releaseSlot <- 1

	url := GetOneURL()
	if url != "" {
		urls = append(urls, url)
	}

	p.AddTargetRequests(urls, "html")

}

func (this *MyPageProcesser) Finish() {
	// fmt.Printf("TODO:before end spider \r\n")

}

func distributeURL(threadNum int, urlstr urlstore.URLCrawlerStore) {
	var remain int = 0
	for {
		select {
		case <-releaseSlot:
			remain++
			oldRemain := remain
			for i := 0; i < oldRemain; i++ {
				url, err := urlstr.GetOneNeedCrawlerURL()
				if err != nil {
					log.Printf("Distribute URL error, error Msg = \"%v\"\n", err)
					continue
				}
				if url == "" {
					break
				}
				remain--
				getURL <- url
			}
			if oldRemain == threadNum && remain == threadNum {
				close(exitChan)
				break
			}
		}
	}
}

func GetOneURL() string {
	for {
		select {
		case URL := <-getURL:
			return URL
		case <-exitChan:
			return ""
		}
	}
}

func main() {
	flag.Parse()
	// chan init
	exitChan = make(chan struct{})
	getURL = make(chan string, *threadNum)
	releaseSlot = make(chan int, *threadNum)

	// repetition and urlstor initialization
	// TODO:
	// Add flag configuration of redis
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rep = repetition.RepetitionJudgement{}
	err := rep.InitializeVisited(c, "repetition")
	if err != nil {
		log.Fatal(err)
	}

	urlstr = urlstore.URLCrawlerStore{}
	_, err = urlstr.InitialURLsStore(c, "colNeedCrawl", "colNeedCommit", "colNeedCrawl", "colNeedCommit")
	visited, _ := rep.CheckIfVisited("http://www.dianping.com/")
	if !visited {
		rep.VisitedNewNode("http://www.dianping.com/")
		urlstr.UploadURL("http://www.dianping.com/")
	}
	if err != nil {
		log.Fatal(err)
	}

	// db initilization
	dbSession, err := storage.Link2Db(*databaseURL)
	defer dbSession.Close()
	if err != nil {
		log.Fatal(err)
	}
	collection = storage.Link2Collection(dbSession, *databaseName, *databaseUser, *databasePwd, *collectionName, *databaseAuth)
	go distributeURL(*threadNum, urlstr)
	// url initilziation
	for i := 0; i < *threadNum; i++ {
		releaseSlot <- 1
	}
	rootURL := GetOneURL()

	// Spider input:
	//  PageProcesser ;
	//  Task name used in Pipeline for record;
	// TO-DO :
	// Change to goroutine mechanism, use channel to get url
	// Able to start serveral goroutine to crawler a same website.
	spider.NewSpider(NewMyPageProcesser(), "TaskName").
		AddUrl(rootURL, "html").                    // Start url, html is the responce type ("html" or "json" or "jsonp" or "text")
		AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
		SetThreadnum((uint)(*threadNum)).           // Crawl request by three Coroutines
		Run()
}
