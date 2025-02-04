package main

import (
	"github.com/dovelus/corb4n-c2/server/c2ext"
	"github.com/dovelus/corb4n-c2/server/c2int"
	"github.com/dovelus/corb4n-c2/server/db"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	// Start the internal and external HTTP servers
	go c2int.StartIntHTTPServer()
	c2ext.StartExtHTTPServer()
}
