# KVAss

a small function that using AWS S3 as key/value storage via HTTP API

### Requirements

- Golang
- [Apex](http://apex.run/)
- [AWS credentials](http://apex.run/)

## Init

```shell
$ make init
```

## Deploy

```shell
$ make deploy
```

I recommanded a tool, [direnv](https://direnv.net/) to localize and automatically set the variables when you're working on a project.

# Invoke

### PUT

PUT 'Invoke URL'/__KEY__

request body: __BODY__

response status code: 201

### GET

GET 'Invoke URL'/__KEY__

response content type: text/plain

response status code: 200
