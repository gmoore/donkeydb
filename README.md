### Donkey DB

A terrible log-driven key-value store.

This was a lunchtime hack project based on the ideas presented in Chapter 3 of **Designing Data-Intensive Applications** by Martin Kleppmann

#### Quick Start

You'll need Go. I used 1.8 but any version is probably fine.  

```
git clone git@github.com:gmoore/donkeydb.git
cd donkeydb
go build donkeydb.go
./donkeydb
```

#### Usage

Insert some data  

`./donkeydb insert (key) (value)`

Read some data

`./donkeydb insert (key)`

Delete some data

`./donkeydb delete (key)`

#### Features

* Slow
* Unstable 
* Not distributed
* Not concurrent
* Not durable (If you delete the data file, your data is gone forever)
* Keys can't be deleted