package concurrent

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/logrusorgru/aurora"
	"github.com/tjhorner/e6dl/e621"
)

// BeginDownload takes a slice of posts, a directory to save them in, and a
// number of concurrent workers to make. It blocks until all the post have
// been processed. It returns the number of successes, failures, and the total
// amount of posts.
func BeginDownload(posts *[]e621.Post, saveDirectory *string, maxConcurrents *int) (*int, *int, *int) {
	// Channel for worker goroutines to notify the main goroutine that it is done
	// downloading a post
	done := make(chan interface{})

	// Channel for main goroutine to give workers a post when they are done downloading
	// one
	pc := make(chan *e621.Post)

	var completed int
	var successes int
	var failures int
	var current int

	total := len(*posts)

	// If we have more workers than posts, then we don't need all of them
	if *maxConcurrents > total {
		*maxConcurrents = total
	}

	for i := 0; i < *maxConcurrents; i++ {
		// Create our workers
		go work(i+1, *saveDirectory, &completed, &total, &successes, &failures, done, pc)

		// Give them their initial posts
		pc <- &(*posts)[current]
		current++

		time.Sleep(time.Millisecond * 50)
	}

	for {
		// Wait for a worker to be done
		<-done

		// If we finished downloading all posts, break out of the loop
		if successes+failures == total {
			break
		}

		// If there's no more posts to give, stop the worker
		if current >= total {
			pc <- nil
			continue
		}

		// Give the worker the next post in the array
		pc <- &(*posts)[current]
		current++
	}

	return &successes, &failures, &total
}

func work(wn int, directory string, completed *int, total *int, successes *int, failures *int, done chan interface{}, pc chan *e621.Post) {
	for {
		*completed++

		// Wait for a post from main
		post := <-pc
		if post == nil { // nil means there aren't any more posts, so we're OK to break
			return
		}

		progress := aurora.Sprintf(aurora.Green("[%d/%d]"), *completed, *total)
		workerText := aurora.Sprintf(aurora.Cyan("[w%d]"), wn)

		fmt.Println(aurora.Sprintf(
			"%s %s Downloading post %d (%s) -> %s...",
			progress,
			workerText,
			post.ID,
			humanize.Bytes(uint64(post.FileSize)),
			getSavePath(post, &directory),
		))

		err := downloadPost(post, directory)
		if err != nil {
			fmt.Printf("[w%d] Failed to download post %d: %v\n", wn, post.ID, err)
			*failures++
		} else {
			*successes++
		}

		done <- nil
	}
}

func getSavePath(post *e621.Post, directory *string) string {
	savePath := path.Join(*directory, strconv.Itoa(post.ID)+"."+post.FileExt)
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
