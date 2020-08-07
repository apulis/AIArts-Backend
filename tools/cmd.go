package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {

	os.Chdir("../")
	log.Println("Server exiting")
}
