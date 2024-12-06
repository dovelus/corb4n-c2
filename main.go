package main

import (
	"github.com/dovelus/corb4n-c2/server/db"
)

func main() {

	db.InitDB()
	db.CloseDB()
}
