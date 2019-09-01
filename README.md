
# URL Shortener

## Running

```bash
$ docker build . -t shorter:latest
$ docker run shorter:latest -p 8000:8000
```

## TODO

* [ ] Dockerfile
* [ ] URL validation and javascript injection security
* [ ] UI for creating URLs in browser
* [ ] Expiry of URL records
* [ ] Request scoped correlation GUID
* [ ] User sign-up and sign-in
* [ ] PostgreSQL