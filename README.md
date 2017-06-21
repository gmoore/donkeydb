### Donkey DB

A terrible log-driven key-value store.

This was a lunchtime hack project based on the ideas presented in Chapter 3 of **Designing Data-Intensive Applications** by Martin Kleppmann

#### Quick Start

You'll need Go. I used 1.8 but any version is probably fine.  

```
git clone git@github.com:gmoore/donkeydb.git
cd donkeydb
make
./bin/donkeydb
```

#### Usage

Insert some data  

`./bin/donkeyclient insert (key) (value)`

Read some data

`./bin/donkeyclient insert (key)`

Delete some data

`./bin/donkeyclient delete (key)`

#### Features

* Slow
* Unstable 
* Not distributed
* Not concurrent
* Not durable
* Error prone
* Not guaranteed
* Unacceptable