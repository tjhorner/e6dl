package e621

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type PostFile struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Ext    string `json:"ext"`
	Size   int    `json:"size"`
	Md5    string `json:"md5"`
	URL    string `json:"url"`
}

type PostPreview struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type PostSample struct {
	Has    bool   `json:"has"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

type PostScore struct {
	Up    int `json:"up"`
	Down  int `json:"down"`
	Total int `json:"total"`
}

type PostTags struct {
	General   []string `json:"general"`
	Species   []string `json:"species"`
	Character []string `json:"character"`
	Copyright []string `json:"copyright"`
	Artist    []string `json:"artist"`
	Invalid   []string `json:"invalid"`
	Lore      []string `json:"lore"`
	Meta      []string `json:"meta"`
}

func (t *PostTags) All() []string {
	allTags := []string{}

	groups := [][]string{
		t.General,
		t.Species,
		t.Character,
		t.Copyright,
		t.Artist,
		t.Invalid,
		t.Lore,
		t.Meta,
	}

	for _, g := range groups {
		allTags = append(allTags, g...)
	}

	return allTags
}

type PostFlags struct {
	Pending      bool `json:"pending"`
	Flagged      bool `json:"flagged"`
	NoteLocked   bool `json:"note_locked"`
	StatusLocked bool `json:"status_locked"`
	RatingLocked bool `json:"rating_locked"`
	Deleted      bool `json:"deleted"`
}

type PostRelationships struct {
	ParentID          int   `json:"parent_id"`
	HasChildren       bool  `json:"has_children"`
	HasActiveChildren bool  `json:"has_active_children"`
	Children          []int `json:"children"`
}

type Post struct {
	ID            int               `json:"id"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
	File          PostFile          `json:"file"`
	Preview       PostPreview       `json:"preview"`
	Sample        PostSample        `json:"sample"`
	Score         PostScore         `json:"score"`
	Tags          PostTags          `json:"tags"`
	LockedTags    []string          `json:"locked_tags"`
	ChangeSeq     int               `json:"change_seq"`
	Flags         PostFlags         `json:"flags"`
	Rating        string            `json:"rating"`
	FavCount      int               `json:"fav_count"`
	Sources       []string          `json:"sources"`
	Pools         []int             `json:"pools"`
	Relationships PostRelationships `json:"relationships"`
	ApproverID    int               `json:"approver_id"`
	UploaderID    int               `json:"uploader_id"`
	Description   string            `json:"description"`
	CommentCount  int               `json:"comment_count"`
	IsFavorited   bool              `json:"is_favorited"`
}

// SerializedDate represents a serialized date passed via JSON
type SerializedDate struct {
	JSONClass   string `json:"json_class"`
	Seconds     int64  `json:"s"`
	Nanoseconds int64  `json:"n"`
}

// Time returns a time.Time object representing the SerializedDate
func (date *SerializedDate) Time() time.Time {
	return time.Unix(0, date.Nanoseconds)
}

// GetPostsForTags gets a list of e621 Posts
func GetPostsForTags(tags string, limit int, sfw bool, page int) ([]Post, error) {
	client := &http.Client{}

	var domain string

	if sfw {
		domain = "e926.net"
	} else {
		domain = "e621.net"
	}

	req, _ := http.NewRequest("GET", "https://"+domain+"/posts.json", nil)
	req.Header.Set("User-Agent", "e6dl: go edition (@tjhorner on Telegram)")

	qs := req.URL.Query()
	qs.Add("tags", tags)
	qs.Add("limit", strconv.Itoa(limit))
	qs.Add("page", strconv.Itoa(page))

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

	var respBody struct {
		Posts []Post
	}

	json.Unmarshal(body, &respBody)

	return respBody.Posts, nil
}
