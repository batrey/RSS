package handlers

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"rss/app/models"
	"strconv"
	"time"
)

const AddForm = `
<form method="GET" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>`

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

func parseFeedBbc(body []byte) (*models.RssBbc, error) {
	feed := models.RssBbc{}
	err := xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

func parseFeedSky(body []byte) (*models.RssSky, error) {
	feed := models.RssSky{}
	err := xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func CategoryCheck(category string) string {
	switch category {
	case "bbc":
		return "BBC News - UK"
	case "bbc-tech":
		return "BC News - Technology"
	case "sky":
		return "K News - The latest headlines from the UK | Sky News"
	case "sky-tech":
		return "Tech News - Latest Technology and Gadget News | Sky News"
	default:
		return "No category found"
	}
}
