package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

func main() {
	// define cmd line flags
	tags := flag.String("tags", "", "Tags to search for")
	maxConcurrents := flag.Int("concurrents", 5, "Maximum amount of concurrent downloads")
	postLimit := flag.Int("limit", 10, "Maximum amount of posts to grab from e621")
	saveDirectory := flag.String("out", "dl", "The directory to write the downloaded posts to")
	sfw := flag.Bool("sfw", false, "Download posts from e926 instead of e621.")

	flag.Parse()

	fmt.Printf("Fetching posts for \"%v\" (limit %v)\n", *tags, *postLimit)

	posts, err := GetPostsForTags(*tags, *postLimit, *sfw)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Found %d posts. Starting download with %d workers...\n\n", len(posts), *maxConcurrents)

	cwd, _ := os.Getwd()
	absSaveDir := path.Join(cwd, *saveDirectory)

	err = os.MkdirAll(absSaveDir, 0755)
	if err != nil {
		fmt.Printf("Cannot create output directory (%s). Do you have the right permissions?\n", absSaveDir)
		os.Exit(1)
	}

	BeginDownload(&posts, saveDirectory, maxConcurrents)

	fmt.Printf("All done! %d posts downloaded and saved.\n", len(posts))
}
