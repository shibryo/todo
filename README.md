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

## Onion-Architecture

### オニオンアーキテクチャとは
オニオンアーキテクチャはJeffry Palermoさんによって考案されたアーキテクチャのパターン。
特に、DDDのような大規模で、長期的な運用が想定されるアプリケーションに適したアプリケーションに向いたアーキテクチャのパターンです。

### レイヤードアーキテクチャとの比較
レイヤードアーキテクチャと比較して、依存の向きに大きく違いがある。
レイヤードアーキテクチャはプレゼンテーション層、アプリケーション層、業務ロジック層、データアクセス層（インフラ層）というレイヤーに分けて関心ごとを分離してる。
依存の向きとしては、　プレゼンテーション層→アプリケーション層→業務ロジック層→データアクセス層（インフラ層）といった形で、データアクセス層に向いて依存している。
レイヤードアーキテクチャによって責務の分離ができたものの、データアクセス層は３〜５年でDBの変更やスキーマの変更などでデータにアクセスする方法が変わります。
データアクセス層に依存したアプリケーションは業務ロジック以下を変更することになるため、多くの工数をかけて修正するか変更を断念することになります。
基本的に、機能追加に当たらない範囲での変更には大きな予算付きづらいため放置されることになるでしょう。

### 依存性逆転の法則
この問題を解決するために、オニオンアーキテクチャは業務ロジック層を中心とし、データアクセス層が一番外にくるようなアーキテクチャとなっています。
依存関係としては、データアクセス層（インフラ層）、プレゼンテーション層→アプリケーション層→業務ロジック層のようになります。
※インフラ層が一番依存としては外になり、業務ロジック層が一番依存の中心になる。
この依存関係を実現するために”依存性逆転の法則”を利用しています。
以下の例1ではapp層のserviceがinfra層のrepositoryに依存してしまっています。
例2では、app層にinterfaceを定義することで、infra層でinterfaceに沿って実装を作成する形になりapp←infraという依存になります。

例1）app層のserviceがinfra層のrepositoryに依存してしまっている様子
app層

type Service：app層でserviceを定義
  infra.repository.create()：infra層のrepositoryを呼び出してしまっている。
------ ↓依存の方向 -----
type Repository：infra層でrepositoryを定義

infra層


例2）infra層が依存性逆転の法則でapp層のインターフェースに依存している様子
app層

type repositoryInterface：app層でinterfaceを定義
------ ↑依存の方向 -----
type repositoryImpl：app層のinterfaceを実装(implementantion)

infra層

さらにいうと、app層にrepositoryのinterfaceがあることで、interfaceの関数を満たすオブジェクトであれば
何を実装しても良いため、テストの際にはモックを実装しても良いし、別のDBを利用しても良い。

依存性逆転の法則を実装するにあたっては依存性注入ができるオブジェクトとする必要もある。

例3）ServiceにRepositoryInterfaceの依存性を注入できるようにする。
```
type Service struct { repo RepositoryInterface }
func NewService(repo RepositoryInterface) Service{
    return &Service{ repo = repo}
}
```

### オニオンアーキテクチャの主要な原則
発案者によると以下の原則を守っているとオニオンアーキテクチャと呼べるとのこと。
オニオンアーキテクチャの主要な原則:
* アプリケーションは独立したオブジェクトモデルに基づいて構築されている
* 内側の層はインターフェースを定義します。外側の層はインターフェースを実装します。
* 結合方向は中心に向かう
* すべてのアプリケーションコアコードはインフラストラクチャとは別にコンパイルおよび実行できます。

### 各レイヤーの概要
「ドメイン駆動設計をはじめよう」によると以下の通りである。
主に、DDDの流れを汲んでいると考えられる。
##### 業務ロジック層
プログラムの業務ロジックを実装しカプセル化する役割。
業務上の意思決定が行われる場所

#### サービス層（アプリケーション層）
サービス層はプレゼンテーション層と業務ロジック層との仲介役として機能します。
サービス層は、業務ロジック層のファサード（窓口）として機能します。つまり、公開するメソッドのインタフェースを定義し、下位レイヤー（業務ロジック層）で実装された振る舞いをカプセル化します。

#### プレゼンテーション層
エンドユーザーとの対話のためのユーザーインタフェースを実装します。
もともとは、このレイヤーはブラウザインタフェースやデスクトップアプリケーションのようなグラフィカルインタフェースを表します。
厳密に言えば、プレゼンテーション層はプログラムの公開インタフェースです。

#### データアクセス層（インフラ層）
データアクセス層は、永続化の仕組みと接続する機能を提供します。
このレイヤーでは、プログラムの機能を実現するために、外部のさまざまな情報源と連携します。

## Reference

### Onion-Architecture
The Onion Architecture : part 1
https://jeffreypalermo.com/2008/07/the-onion-architecture-part-1/
※オニオンアーキテクチャの原文

The Onion Architecture : part 3
https://jeffreypalermo.com/2008/08/the-onion-architecture-part-3/
※オニオンアーキテクチャの原則を参照


### explanation Onion-Architecture
オニオンアーキテクチャとは何か
https://qiita.com/cocoa-maemae/items/e3f2eabbe0877c2af8d0

### explanation Onion-Architecture and DDD(Domain Driven Design)
ドメイン駆動設計をはじめよう
https://www.oreilly.co.jp/books/9784814400737/
※オニオンアーキテクチャについての説明 p.152
※各アーキテクチャのレイヤーの対応は以下のように説明されている p.147
プレゼンテーション層＝ユーザーインターフェース層
サービス層＝アプリケーション層
業務ロジック層＝ドメイン層＝モデル層
データアクセス層＝インフラストラクチャ層

### architecting domain
現場で役立つシステム設計の原則 ~変更を楽で安全にするオブジェクト指向の実践技法
https://gihyo.jp/book/2017/978-4-7741-9087-7
https://masuda220.hatenablog.com/entry/2020/06/26/182317


### relation each layers
データ詰め替え戦略
https://scrapbox.io/kawasima/%E3%83%87%E3%83%BC%E3%82%BF%E8%A9%B0%E3%82%81%E6%9B%BF%E3%81%88%E6%88%A6%E7%95%A5