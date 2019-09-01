
# URL Shortener

## Running

```bash
$ docker build . -t shorter:latest
$ docker run shorter:latest -p 8000:8000
```

## Configuration

Environment variables for configuring the service.

| Name             | Description                         | Example                    |
| ---------------- | ----------------------------------- | -------------------------- |
| SHORTER_SQLITE   | File name and path to SQLite file   | `/var/db/shorter.sqlite3`  |
| SHORTER_BASEURL  | Base URL for short URLs             | `https://sho.rt/`          |
| SHORTER_PORT     | Port to listen on                   | `8000`                     |


## TODO

* [ ] Dockerfile
* [ ] URL validation and javascript injection security
* [ ] UI for creating URLs in browser
* [ ] Expiry of URL records
* [ ] Request scoped correlation GUID
* [ ] User sign-up and sign-in
* [ ] PostgreSQL