package main

import (
	leveldb "github.com/syndtr/goleveldb/leveldb"
	"log"
)

func main() {
	db, err := leveldb.OpenFile("/Users/lg/.neutrinod-liveness/data/", nil)
	defer db.Close()

	if err != nil {
		log.Panicf("Unable to open DB")
	}
}
