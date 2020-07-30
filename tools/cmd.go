package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {

	os.Chdir("../")

	LoadConfig()
	InitDb()
	LoadTemplates(os.Args[1])
	log.Println("Server exiting")
}
