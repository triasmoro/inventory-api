# Inventory API

API to manage your simple inventory

## Database

Save as `data.db` file which generated automatically when running app at first

## Run

### Default port 8080

```shell
$ go run main.go
```

### Custom Port

Change `8787` with your choice.

```shell
$ go run main.go -p=8787
```

## Build

If you use Windows OS, please use bash console to build binary such as Git Bash.
Use `... GOARCH=386 ...` if you want to build for 32-bit machine.

Don't forget to bring `db-sqlite.sql` with your binary and `data.db` if you want to bring your existing data also.

Read [Build Go Executables from DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04) for more information.

### Linux64 builder

```shell
$ GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o main
```

### Windows64 builder

Set filename as `main.exe` if you expect to run it via Command Prompt. Since the app will not run if you build without `.exe` extension.

```shell
$ GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o main.exe
```

### MacOS64 builder

```shell
$ GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o main
```
