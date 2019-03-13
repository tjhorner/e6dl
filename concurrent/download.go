package concurrent

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/tjhorner/e6dl/e621"
)

// BeginDownload takes a slice of posts, a directory to save them in, and a
// number of concurrent workers to make. It blocks until all the post have
// been processed. It returns the number of successes, failures, and the total
// amount of posts.
func BeginDownload(posts *[]e621.Post, saveDirectory *string, maxConcurrents *int) (*int, *int, *int) {
	var wg sync.WaitGroup
	var completed int
	var successes int
	var failures int

	total := len(*posts)

	// Distribute the posts based on the number of workers
	ppw := len(*posts) / *maxConcurrents // ppw: posts per worker
	mod := len(*posts) % *maxConcurrents // mod: remainder of posts

	for i := 0; i < *maxConcurrents; i++ {
		postsLower := i * ppw
		postsUpper := i*ppw + ppw

		if i == *maxConcurrents-1 {
			// Give the last worker the remaining posts
			// TODO: compensate it for labor
			postsUpper += mod
		}

		wg.Add(1)
		go work(i+1, (*posts)[postsLower:postsUpper], *saveDirectory, &completed, &successes, &failures, &total, &wg)
		// Spawn workers with a little bit of a delay so as to not DDOS e621
		// but also make the initial numbers show up correctly
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()

	return &successes, &failures, &total
}

func work(wn int, posts []e621.Post, directory string, completed *int, successes *int, failures *int, total *int, wg *sync.WaitGroup) {
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

		err := downloadPost(&post, directory)
		if err != nil {
			fmt.Printf("[w%d] Failed to download post %d: %v\n", wn, post.ID, err)
			*failures++
		} else {
			*successes++
		}
	}
}

func getSavePath(post *e621.Post, directory *string) string {
	pathSliced := strings.Split(post.FileURL, ".")
	extension := pathSliced[len(pathSliced)-1]

	savePath := path.Join(*directory, strconv.Itoa(post.ID)+"."+extension)

	return savePath
}

func downloadPost(post *e621.Post, directory string) error {
	savePath := getSavePath(post, &directory)

	resp, err := e621.HTTPGet(post.FileURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(savePath, body, 0755)
	if err != nil {
		return err
	}

	return nil
}
