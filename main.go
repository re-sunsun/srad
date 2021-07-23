package main

import "SunsunSRAD/ssrad"

func main() {
	serversPool := ssrad.Default()
	serversPool.Run(":8080")
}
