package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"
)

// BeginDownload is meant to be called as a goroutine and begins the post download process.
func BeginDownload(posts *[]Post, saveDirectory *string, maxConcurrents *int) {
	var wg sync.WaitGroup
	var completed int

	total := len(*posts)

	// Distribute the posts based on the number of workers
	ppw := len(*posts) / *maxConcurrents
	mod := len(*posts) % *maxConcurrents

	for i := 0; i < *maxConcurrents; i++ {
		postsLower := i * ppw
		postsUpper := i*ppw + ppw

		if i == *maxConcurrents-1 {
			// Give the last worker the remaining posts
			// TODO: compensate it for labor
			postsUpper += mod
		}

		wg.Add(1)
		go work(i+1, (*posts)[postsLower:postsUpper], *saveDirectory, &completed, &total, &wg)
	}

	wg.Wait()
}

func work(wn int, posts []Post, directory string, completed *int, total *int, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, post := range posts {
		*completed++
		fmt.Printf(
			"[%d/%d] [w%d] Downloading post %d (%s) -> %s...\n",
			*completed,
			*total,
			wn,
			post.ID,
			humanize.Bytes(uint64(post.FileSize)),
			getSavePath(&post, &directory),
		)
		downloadPost(&post, directory)
	}
}

func getSavePath(post *Post, directory *string) string {
	pathSliced := strings.Split(post.FileURL, ".")
	extension := pathSliced[len(pathSliced)-1]

	savePath := path.Join(*directory, strconv.Itoa(post.ID)+"."+extension)

	return savePath
}

func downloadPost(post *Post, directory string) {
	savePath := getSavePath(post, &directory)

	resp, err := HTTPGet(post.FileURL)
	if err != nil {
		fmt.Println("Unable to download, skipping...")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read post response body, skipping...")
		return
	}

	err = ioutil.WriteFile(savePath, body, 0755)
	if err != nil {
		fmt.Printf("Error: could not write to file: %v\n", err)
		return
	}
}
