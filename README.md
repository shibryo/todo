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


## Reference

MVC
https://ja.wikipedia.org/wiki/Model_View_Controller
https://github.com/system-sekkei/isolating-the-domain

architecting Model
現場で役立つシステム設計の原則 ~変更を楽で安全にするオブジェクト指向の実践技法
https://gihyo.jp/book/2017/978-4-7741-9087-7
https://masuda220.hatenablog.com/entry/2020/06/26/182317


relation between repository and model
データ詰め替え戦略
https://scrapbox.io/kawasima/%E3%83%87%E3%83%BC%E3%82%BF%E8%A9%B0%E3%82%81%E6%9B%BF%E3%81%88%E6%88%A6%E7%95%A5