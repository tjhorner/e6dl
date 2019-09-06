# e6dl

**["I don't care about any of this, just take me to the downloads"](https://github.com/tjhorner/e6dl/releases/latest)**

![](https://user-images.githubusercontent.com/2646487/54313480-b1170e80-4596-11e9-811f-d73f1ea13b99.gif)

This is a command line tool for downloading posts that match certain tags on e621 or e926.

It does basically the same thing as [this tool](https://www.npmjs.com/package/e6dl) except it was written in Go and the output is a lot less pretty.

It supports concurrently downloading posts, and you can set a maximum number of workers that should be downloading at a time.

I made this because I wanted to rewrite one of my previous projects in Go, so I decided to start with this one since it's a pretty small and simple command line tool that encapsulates a lot of concepts in the language (type definition, imports, goroutines, standard libraries, slices, to name a few).

## Installing

### Prebuilt Binaries

There is an install script available for Linux/macOS [here](https://github.com/tjhorner/e6dl/blob/master/install.sh). It automatically grabs the latest release (at the time of writing) and saves to `/usr/local/bin/e6dl`.

### From Source

Clone repo.

```bash
# Build e6dl then install it to /usr/local/bin
make install
```

You can also install it to a custom path of your choosing:

```bash
make install INSTALLPATH="/bin"
```

If you wanna get rid of it:

```bash
make uninstall
```

## Building

First, install Go.

Then just:

```bash
make build
```

`bin/e6dl` will magically appear. You can also just:

```bash
go run main.go
```

To distribute:

```bash
make dist
```

3 files will be spit out in the `dist` directory - one for Windows, one for macOS, and one for Linux.

You can also run `make dist-win`, `make dist-macos`, or `make dist-linux` to build them individually.

## Example

Here's a situation that uses every flag:

If you wanted to download a maximum of 20 posts with the tag `pokemon` in random order from e926 to the directory `./posts` with a maximum of 2 downloading at a time:

```bash
e6dl --tags "pokemon order:random" --out ./posts --limit 20 --concurrents 2 --sfw
```

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
    	Download posts from e926 instead of e621
  --tags string
    	Tags to search for
  --pages int
    	Number of search result pages to download (default 1)
```