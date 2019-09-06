package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/tjhorner/e6dl/concurrent"
	"github.com/tjhorner/e6dl/e621"
)

func main() {
	// define cmd line flags
	tags := flag.String("tags", "", "Tags to search for")
	maxConcurrents := flag.Int("concurrents", 5, "Maximum amount of concurrent downloads")
	postLimit := flag.Int("limit", 10, "Maximum amount of posts to grab from e621")
	saveDirectory := flag.String("out", "dl", "The directory to write the downloaded posts to")
	sfw := flag.Bool("sfw", false, "Download posts from e926 instead of e621")
	pages := flag.Int("pages", 1, "Number of search result pages to download")

	flag.Parse()

	fmt.Printf("Fetching posts for \"%s\" (limit=%d, pages=%d)\n", *tags, *postLimit, *pages)

	var allPosts []e621.Post

	for i := 1; i <= *pages; i++ {
		fmt.Printf("Fetching page %d/%d...", i, *pages)

		posts, err := e621.GetPostsForTags(*tags, *postLimit, *sfw, i)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf(" fetched %d posts\n", len(posts))

		allPosts = append(allPosts, posts...)
	}

	fmt.Printf("Found %d posts. Starting download with %d workers...\n\n", len(allPosts), *maxConcurrents)

	cwd, _ := os.Getwd()
	absSaveDir := path.Join(cwd, *saveDirectory)

	err := os.MkdirAll(absSaveDir, 0755)
	if err != nil {
		fmt.Printf("Cannot create output directory (%s). Do you have the right permissions?\n", absSaveDir)
		os.Exit(1)
	}

	successes, failures, _ := concurrent.BeginDownload(&allPosts, saveDirectory, maxConcurrents)

	fmt.Printf("\nAll done! %d posts downloaded and saved. (%d failed to download)\n", *successes, *failures)
}
