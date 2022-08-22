package server

import (
	"os"
	"rss/app/handlers"
	"rss/app/models"
	"strings"
)

type TaskBaseHandler struct {
	TaskBaseHandler models.ArticlesRepo
}

func TaskNewBaseHandle(task models.ArticlesRepo) *TaskBaseHandler {
	return &TaskBaseHandler{
		TaskBaseHandler: task,
	}
}
func (t *TaskBaseHandler) TaskBbc() {

	urls := strings.Split(os.Getenv("BBC_URLS"), ",")
	// Create a channel to process the feeds
	feedc := make(chan models.RssBbc, len(urls))

	// Start a goroutine for each feed url
	for _, u := range urls {
		go handlers.GetRssFeedBbc(u, feedc)
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

func (t *TaskBaseHandler) TaskSky() {

	urls := strings.Split(os.Getenv("SKY_URLS"), ",")
	// Create a channel to process the feeds
	feedc := make(chan models.RssSky, len(urls))

	// Start a goroutine for each feed url
	for _, u := range urls {
		go handlers.GetRssFeedSky(u, feedc)
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
