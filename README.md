# lwnfeed
lwnfeed is a simple daemon that creates a full-text RSS feed from lwn.net by
scraping the content from the website.

You must have a LWN subscription to use this tool. It is intended for personal
use only, please do not publish the generated feeds.

## Getting started
0. Install [Go](https://golang.org/dl/) 1.22 or higher.
1. Clone the repository (`git clone https://github.com/tulir/lwnfeed.git`) and enter the directory (`cd lwnfeed`).
2. Build the program (`./build.sh`).
3. Use `./lwnfeed login` to log in and store the auth cookie into a file.
4. Start the server with `./lwnfeed start`.
5. Open `http://localhost:8080/feed.rss` to view the feed. `feed.atom` and `feed.json` also work.

### Docker
1. Run `docker run --rm -itv $(pwd):/data dock.mau.dev/tulir/lwnfeed login` to log in and store the auth cookie into a file.
2. Start the container normally with your preferred method (see examples below), keeping the same bind mount as before.

#### docker start
```
$ docker start -p 8080:8080 --name lwnfeed -v $(pwd):/data dock.mau.dev/tulir/lwnfeed
```

#### docker-compose
```yaml
version: "3.7"

services:
  lwnfeed:
    container_name: lwnfeed
    image: dock.mau.dev/tulir/lwnfeed
    ports:
    - 8080:8080
    volumes:
    - .:/data
```
