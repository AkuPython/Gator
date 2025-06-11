package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var body io.Reader
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, body)
	if err != nil {
		return nil, fmt.Errorf("Could not create Context: %v", err)
	}
	client := &http.Client{}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do error: %v", err)
	}
	defer resp.Body.Close()

	xmlBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("client.Do error: %v", err)
	}
	var feed RSSFeed
	err = xml.Unmarshal(xmlBody, &feed)
	if err != nil {
		return nil, fmt.Errorf("Could not Unmarshal: %v", err)
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
	fmt.Println(feed)
	return &feed, nil
}


