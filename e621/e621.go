package e621

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Post represents an e621 post object returned by the e621 API.
type Post struct {
	ID             int            `json:"id"`             // The ID of the post
	Tags           string         `json:"tags"`           // Space-separated list of tags attached to this post
	LockedTags     bool           `json:"locked_tags"`    // (undocumented)
	Description    string         `json:"description"`    // The post's description
	CreatedAt      SerializedDate `json:"created_at"`     // When the post was uploaded
	CreatorID      int            `json:"creator_id"`     // User ID of the user who uploaded the post
	Author         string         `json:"author"`         // Username of the user who uploaded the post
	Change         int            `json:"change"`         // (undocumented)
	Source         string         `json:"source"`         // URL that the source for this post can be found at
	Score          int            `json:"score"`          // The post's score (upvotes - downvotes)
	FavoritesCount int            `json:"fav_count"`      // Amount of users that favorited this post
	MD5Hash        string         `json:"md5"`            // MD5-sum of the post's file's content
	FileSize       int            `json:"file_size"`      // Size of the post's file
	FileURL        string         `json:"file_url"`       // URL to the full-sized file
	FileExt        string         `json:"file_ext"`       // File extension
	PreviewURL     string         `json:"preview_url"`    // URL to preview-sized version of the file
	PreviewHeight  int            `json:"preview_height"` // Height of the preview
	PreviewWidth   int            `json:"preview_width"`  // Width of the preview
	Rating         string         `json:"rating"`         // Rating of the file ("safe", "questionable", "explicit")
	Status         string         `json:"status"`         // Moderation status ("active" or "pending")
	Width          int            `json:"width"`          // Width of the original file
	Height         int            `json:"height"`         // Height of the original file
	HasComments    bool           `json:"has_comments"`   // True if post has comments
	HasNotes       bool           `json:"has_notes"`      // True if post has notes
	HasChildren    bool           `json:"has_children"`   // True if post has children
	Children       string         `json:"children"`       // Comma-separated list of children post IDs
	ParentID       int            `json:"parent_id"`      // ID of the parent post
	Artist         []string       `json:"artist"`         // Slice of artist names
	Sources        []string       `json:"sources"`        // Slice of source URLs
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
