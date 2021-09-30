# my-yarn-api

## Running locally
Before running this api, start your local MongoDB.

MacOS
```sh
$ brew services start mongodb-community@5.0
```

Linux
```sh
$ sudo systemctl start mongod
```
or
```sh
sudo service mongod start
```

Then run the api
```sh
$ go run main.go
```

If you use the VSCode launcher, change the `program` property to
```json
"program": "${workspaceFolder}"
```

The api will start on `http://localhost:8080`
