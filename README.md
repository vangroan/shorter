
# URL Shortener

## Getting Started

Build and run the Docker container.

```bash
$ docker build . -t shorter:latest
$ sudo docker run shorter:latest -it \
  -p 4000:8000 \
  -e SHORTER_BASEURL='https://localhost:4000/'
```

## Usage

Create a short URL.

```bash
$ curl --data '{"url":"https://www.github.com"}' https://localhost:4000/
```

## Configuration

Environment variables for configuring the service.

| Name             | Description                         | Example                    |
| ---------------- | ----------------------------------- | -------------------------- |
| SHORTER_SQLITE   | File name and path to SQLite file   | `/var/db/shorter.sqlite3`  |
| SHORTER_BASEURL  | Base URL for short URLs             | `https://sho.rt/`          |
| SHORTER_PORT     | Port to listen on                   | `8000`                     |


## TODO

* [x] Dockerfile
* [x] URL validation and javascript injection security
* [x] UI for creating URLs in browser
* [x] Expiry of URL records
* [ ] Request scoped correlation GUID
* [ ] User sign-up and sign-in
* [ ] PostgreSQL