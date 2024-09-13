# todo

## How to run server
first, you should copy .env.sample and rename to .env
```
$cp ./devcontainer/.env.sample ./devcontainer/.env
$cp .env.sample .env
```

second, install package.
```
$ go mod vedor
```


```
$ make run
```

and access to
http://localhost:8080/swagger/index.html

## How to generate openapi

update openapis

```
$ make swagger
```
