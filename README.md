# GoGround: Open Platform

This is the repository houses of monolithic flat architecture that make up GoGroud's guest list API.

## Things you need

* [Go](https://golang.org/dl/): `brew install go`
* Install docker
* Setup `go-ground` folder in your home directory: `/Users/<user_id>/go-ground`. So, `go env` should show `GOPATH="/Users/<user_id>/go-ground"`
* Clone `HasanShahjahan/go-ground` repository in `/Users/<user_id>/`: `git clone <repo_url>`
* After cloning, the `go.mod` file should be found in this directory `/Users/<user_id>/go-ground/go.mod`
* Add `/Users/<user_id>/go-ground/bin` in environment variable `PATH` if it is not already there.

## Running the App

### Run Scripts
We use [Go Modules](https://blog.golang.org/using-go-modules) for dependency management.

Run below commands step by step:
```shell
$ go mod vendor
$ go build
$ go run main.go
```
1. `go mod vendor` should create a `vendor` directory in project root directory. (`/Users/<user_id>/go-ground/vendor`)
   If for any reason, this directory is not being created then try to clear the cache and run the command again. To clear the cache run below command:

```shell
$ go clean -modcache  #clean module cache
$ go mod vendor       #setup vendor dir again
```

### Run docker.compose

Start the **docker engine**. Next, run below command.
```shell
$ 
```


## Running DB migrations

```shell
$ 
```