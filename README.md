# e6dl

This is a command line tool for downloading posts that match certain tags on e621 or e926.

It does basically the same thing as [this tool](https://www.npmjs.com/package/e6dl) except it was written in Go and the output is a lot less pretty.

It supports concurrently downloading posts, and you can set a maximum number of workers that should be downloading at a time.

I made this because I wanted to rewrite one of my previous projects in Go, so I decided to start with this one since it's a pretty small and simple command line tool that encapsulates a lot of concepts in the language (type definition, imports, goroutines, standard libraries, slices, to name a few).

## Installing, Building, etc.

See [here](https://github.com/tjhorner/nplcsv/blob/master/README.md) since it uses the same Makefile.

## Usage

```
Usage of e6dl:
  --concurrents int
    	Maximum amount of concurrent downloads (default 5)
  --limit int
    	Maximum amount of posts to grab from e621 (default 10)
  --out string
    	The directory to write the downloaded posts to (default "dl")
  --sfw
    	Download posts from e926 instead of e621.
  --tags string
    	Tags to search for
```