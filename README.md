# Inventory API

API to manage your simple inventory

.
## Features

Get my postman collection from [here](https://www.getpostman.com/collections/fde01ca3ea1ba772b6cf)

| Method   | Endpoints | Header | Notes |
|:-------- |:--------- | ------ | ----- |
| `POST`   | */product*  | `Content-Type: application/json` | |
| `PATCH`  | */product/:id* | `Content-Type: application/json` | Only update name | 
| `DELETE` | */product_variant/:id* | - | Delete per variant using product variant id |
| `POST`   | */purchase_order* | `Content-Type: application/json` | |
| `DELETE` | */purchase_order/:id* | - | |
| `POST`   | */stock_in* | `Content-Type: application/json` | |
| `DELETE` | */stock_in/:id* | - | |
| `POST`   | */sales_order* | `Content-Type: application/json` | |
| `DELETE` | */sales_order/:id* | - | |
| `POST`   | */stock_out* | `Content-Type: application/json` | |
| `DELETE` | */stock_out/:id* | - | |
| `GET`    | */actual_stock* | - | |
| `GET`    | */assets_report* | - | Use `until_date` parameter which used as before date on `stock-in` & `stock-out`  |
| `GET`    | */sales_report* | - | Use `start_date` and `end_date` parameters used on sales order time |
| `GET`    | */export/product* | - | CSV export for product |
| `GET`    | */export/stock_in* | - | CSV export for stock-in |


.
## Database

Save as `data.db` file which generated automatically when running app at first

.
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

.
## Development

### Add endpoint

Go to `route/public_routes.go` to add another endpoint

.
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