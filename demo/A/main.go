package main

import (
	"SunsunSRAD/ssradc"
	"log"
	"time"
)

func main() {
	clint, err := ssradc.Default("localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer clint.Close()

	clint.Register("LogIn", "sunsun", "159.75.14.159:8080")

	time.Sleep(time.Duration(30) * time.Second)
}
