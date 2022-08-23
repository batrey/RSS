package server

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"rss/app/models"
	"strings"
	"time"
)

type TaskBaseHandler struct {
	TaskBaseHandler models.ArticlesRepo
}

func TaskNewBaseHandle(task models.ArticlesRepo) *TaskBaseHandler {
	return &TaskBaseHandler{
		TaskBaseHandler: task,
	}
}

// Gets BBC articles
func (t *TaskBaseHandler) TaskBbc() {

	urls := strings.Split(os.Getenv("BBC_URLS"), ",")
	// Create a channel to process the feeds
	feedc := make(chan models.RssBbc, len(urls))

	// Start a goroutine for each feed url
	for _, u := range urls {
		go GetRssFeedBbc(u, feedc)
	}

	// Wait for the goroutines to write their results to the channel
	var feeds []models.RssBbc
	for i := 0; i < len(urls); i++ {
		res := <-feedc
		feeds = append(feeds, res)
	}

	//loop over each sites articles and add them to the database
	for _, feed := range feeds {
		t.TaskBaseHandler.AddArticles(feed.Channel.Title, feed)
	}
}

// Gets Sky Articles
func (t *TaskBaseHandler) TaskSky() {

	urls := strings.Split(os.Getenv("SKY_URLS"), ",")
	// Create a channel to process the feeds
	feedc := make(chan models.RssSky, len(urls))

	// Start a goroutine for each feed url
	for _, u := range urls {
		go GetRssFeedSky(u, feedc)
	}

	// Wait for the goroutines to write their results to the channel
	var feeds []models.RssSky
	for i := 0; i < len(urls); i++ {
		res := <-feedc
		feeds = append(feeds, res)
	}

	//loop over each sites articles and add them to the database
	for _, feed := range feeds {
		t.TaskBaseHandler.AddArticles(feed.Channel.Title, feed)
	}
}

// Gets Articles from BBC
func GetRssFeedBbc(url string, feedc chan models.RssBbc) {
	// Create a client with a default timeout
	net := &http.Client{
		Timeout: time.Second * 10,
	}
	//GET request for the feed
	res, err := net.Get(url)
	// If there was an error write that to the channel and return immediately
	if err != nil {
		feedc <- models.RssBbc{}
		return
	}
	defer res.Body.Close()
	// Read the body of the request and parse the feed
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		feedc <- models.RssBbc{}
		return
	}
	feed, err := parseFeedBbc(body)
	if err != nil {
		feedc <- models.RssBbc{}
		return
	}
	feedc <- *feed
}

// Gets Articles from sky
func GetRssFeedSky(url string, feedc chan models.RssSky) {
	// Create a client with a default timeout
	net := &http.Client{
		Timeout: time.Second * 10,
	}
	//GET request for the feed
	res, err := net.Get(url)
	// If there was an error write that to the channel and return immediately
	if err != nil {
		feedc <- models.RssSky{}
		return
	}
	defer res.Body.Close()
	// Read the body of the request and parse the feed
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		feedc <- models.RssSky{}
		return
	}
	feed, err := parseFeedSky(body)
	if err != nil {
		feedc <- models.RssSky{}
		return
	}
	feedc <- *feed
}

// Parses the BBC feeds
func parseFeedBbc(body []byte) (*models.RssBbc, error) {
	feed := models.RssBbc{}
	err := xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

// Parses the Sky feed
func parseFeedSky(body []byte) (*models.RssSky, error) {
	feed := models.RssSky{}
	err := xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}
