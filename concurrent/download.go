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

// workState stores the state of all the jobs and
// is shared across workers
type workState struct {
	Total         int
	Completed     int
	Successes     int
	Failures      int
	SaveDirectory string
}

// BeginDownload takes a slice of posts, a directory to save them in, and a
// number of concurrent workers to make. It blocks until all the post have
// been processed. It returns the number of successes, failures, and the total
// amount of posts.
func BeginDownload(posts *[]e621.Post, saveDirectory *string, maxConcurrents *int) (*int, *int, *int) {
	// Channel for main goroutine to give workers a post when they are done downloading one
	wc := make(chan *e621.Post)

	var current int

	total := len(*posts)

	state := workState{
		Total:         total,
		SaveDirectory: *saveDirectory,
	}

	// If we have more workers than posts, then we don't need all of them
	if *maxConcurrents > total {
		*maxConcurrents = total
	}

	for i := 0; i < *maxConcurrents; i++ {
		// Create our workers
		go work(i+1, &state, wc)

		// Give them their initial posts
		wc <- &(*posts)[current]
		current++

		time.Sleep(time.Millisecond * 50)
	}

	for {
		// Wait for a worker to be done (they send nil to wc)
		<-wc

		// If we finished downloading all posts, break out of the loop
		if state.Successes+state.Failures == total {
			break
		}

		// If there's no more posts to give, stop the worker
		if current >= total {
			wc <- nil
			continue
		}

		// Give the worker the next post in the array
		wc <- &(*posts)[current]
		current++
	}

	return &state.Successes, &state.Failures, &total
}

func work(wn int, state *workState, wc chan *e621.Post) {
	for {
		state.Completed++

		// Wait for a post from main
		post := <-wc
		if post == nil { // nil means there aren't any more posts, so we're OK to break
			return
		}

		progress := aurora.Sprintf(aurora.Green("[%d/%d]"), state.Completed, state.Total)
		workerText := aurora.Sprintf(aurora.Cyan("[w%d]"), wn)

		fmt.Println(aurora.Sprintf(
			"%s %s Downloading post %d (%s) -> %s...",
			progress,
			workerText,
			post.ID,
			humanize.Bytes(uint64(post.FileSize)),
			getSavePath(post, &state.SaveDirectory),
		))

		err := downloadPost(post, state.SaveDirectory)
		if err != nil {
			fmt.Printf("[w%d] Failed to download post %d: %v\n", wn, post.ID, err)
			state.Failures++
		} else {
			state.Successes++
		}

		wc <- nil
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
