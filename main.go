package main

import (
	"github.com/dovelus/corb4n-c2/server/c2"
	"github.com/dovelus/corb4n-c2/server/db"
)

func main() {
	db.InitDB()
	defer db.CloseDB()
	c2.StartHTTPServer()
}
