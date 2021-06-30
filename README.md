# GoGround: Open Platform

This is the repository houses of monolithic flat architecture that make up GoGroud's guest list API.

## Things you need

* [Go](https://golang.org/dl/): `brew install go`
* Install docker
* Setup `go-guest` folder in your home directory: `/Users/<user_id>/go-guest`. So, `go env` should show `GOPATH="/Users/<user_id>/go-guest"`
* Clone `HasanShahjahan/go-guest` repository in `/Users/<user_id>/`: `git clone <repo_url>`
* After cloning, the `go.mod` file should be found in this directory `/Users/<user_id>/go-guest/go.mod`
* Add `/Users/<user_id>/go-guest/bin` in environment variable `PATH` if it is not already there.

## Running the App

### Run Scripts
We use [Go Modules](https://blog.golang.org/using-go-modules) for dependency management.

Run below commands step by step:
```shell
$ go mod vendor
```
1. `go mod vendor` should create a `vendor` directory in project root directory. (`/Users/<user_id>/go-ground/vendor`)
   If for any reason, this directory is not being created then try to clear the cache and run the command again. To clear the cache run below command:

```shell
$ go clean -modcache  #clean module cache
$ go mod vendor       #setup vendor dir again
```

### Setup the Environment variable and configuration
```shell
$ export CONFIG_PATH=go-guest.json
$ export DB_DRIVER=mysql
$ export DB_USERNAME=root
$ export DB_PASSWORD=password
$ export DB_NAME=guests
```
You could also put all environment related variables in your .bash_profile, so you don't need to set them up all the time.

### Run Docker Compose

Start the **docker engine**. Next, run below command.
```shell
$ docker-compose up
```

### Run application

After docker run, make sure your mysql is running inside container and ready for connections...
```shell
$ go build
$ go run main.go
```

Now application is listing port :8080 and ready to go.

## Makefile
Use `make` command to check available targets.

## Running go-vet

Vet examines Go source code and reports suspicious constructs,
such as Printf calls whose arguments do not align with the format string.
Vet uses heuristics that do not guarantee all reports are genuine problems,
but it can find errors not caught by the compilers.

```shell
$ go vet ./.
```

## Running golint

`golint` prints out style mistakes.

```shell
$ go list ./... | grep -v /vendor/ | xargs -n 1
```