package main

import "github.com/grepsd/knowledge-database/pkg"

func main() {
	db := pkg.NewDB("user=postgres password=tpassword database=knowledge-database sslmode=disable")
	s := pkg.NewServer(db)
	s.Init()
}
