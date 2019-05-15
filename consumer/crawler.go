package consumer

import (
	"crawler/config"
	"crawler/model"
	repoCrawler "crawler/repo/crawler"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type gCrawler struct{}

var Crawler gCrawler
var cfg = config.GetConfig()
var wg sync.WaitGroup

func getBody(url string) (doc *goquery.Document, err error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func getInfo(s *goquery.Selection, i int) {
	link, _ := s.Find(".title_news a:first-child").Attr("href")

	// Check this link is existing in db or not
	exist, _ := repoCrawler.New().GetByLink(link)
	if exist.Link == "" {
		// Link is not exist in db, save this link
		title := s.Find(".title_news a:first-child").Text()
		content := s.Find(".description").Text()

		data := model.Crawler{
			Title:   title,
			Content: content,
			Link:    link,
		}

		err := repoCrawler.New().Create(data)
		if err != nil {
			log.Print(err)
		}

		// Find link's relation
		// Currently, can not run
		// getRelation(link)
	}

	wg.Done()
}

func getRelation(url string) {
	doc, err := getBody(url)
	if err == nil {
		doc.Find("#box_morelink_detail ul li").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Find("a:first-child").Attr("href")

			// Check link is existing in db or not
			exist, _ := repoCrawler.New().GetByLink(link)
			if exist.Link == "" {
				log.Print(link)
				title := s.Find("a:first-child").Text()

				data := model.Crawler{
					Title: title,
					Link:  link,
				}

				err := repoCrawler.New().Create(data)
				if err != nil {
					log.Print(err)
				} else {
					getRelation(link)
				}
			}
		})
	} else {
		log.Print(err)
	}
}

func (gCrawler) Start() error {
	// Init Monggo
	cfg.Mongo.Get("crawler").Init()

	log.Print("=========== GETTING DATA")
	log.Print(cfg.Url)

	doc, err := getBody(cfg.Url)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	doc.Find(".sidebar_1 .list_news").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go getInfo(s, i)
	})

	wg.Wait()

	log.Print("=========== DONE")

	return nil
}
