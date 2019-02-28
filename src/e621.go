package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Post represents an e621 post object returned by the e621 API.
type Post struct {
	ID             int            `json:"id"`
	Tags           string         `json:"tags"`
	LockedTags     bool           `json:"locked_tags"`
	Description    string         `json:"description"`
	CreatedAt      SerializedDate `json:"created_at"`
	CreatorID      int            `json:"creator_id"`
	Author         string         `json:"author"`
	Change         int            `json:"change"`
	Source         string         `json:"source"`
	Score          int            `json:"score"`
	FavoritesCount int            `json:"fav_count"`
	MD5Hash        string         `json:"md5"`
	FileSize       int            `json:"file_size"`
	FileURL        string         `json:"file_url"`
	FileExt        string         `json:"file_ext"`
	PreviewURL     string         `json:"preview_url"`
	PreviewHeight  int            `json:"preview_height"`
	PreviewWidth   int            `json:"preview_width"`
	Rating         string         `json:"rating"`
	Status         string         `json:"status"`
	Width          int            `json:"width"`
	Height         int            `json:"height"`
	HasComments    bool           `json:"has_comments"`
	HasNotes       bool           `json:"has_notes"`
	HasChildren    bool           `json:"has_children"`
	Children       string         `json:"children"`
	ParentID       int            `json:"parent_id"`
	Artist         []string       `json:"artist"`
	Sources        []string       `json:"sources"`
}

// SerializedDate represents a serialized date passed via JSON
type SerializedDate struct {
	JSONClass   string `json:"json_class"`
	Seconds     int    `json:"s"`
	Nanoseconds int    `json:"n"`
}

// GetPostsForTags gets a list of e621 Posts
func GetPostsForTags(tags string, limit int, sfw bool) ([]Post, error) {
	client := &http.Client{}

	var domain string

	if sfw {
		domain = "e926.net"
	} else {
		domain = "e621.net"
	}

	req, _ := http.NewRequest("GET", "https://"+domain+"/post/index.json", nil)
	req.Header.Set("User-Agent", "e6dl: go edition (@tjhorner on Telegram)")

	qs := req.URL.Query()
	qs.Add("tags", tags)
	qs.Add("limit", strconv.Itoa(limit))

	req.URL.RawQuery = qs.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var posts []Post
	json.Unmarshal(body, &posts)

	return posts, nil
}
