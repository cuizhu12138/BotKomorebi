package main

import (
	"EutopiaQQBot/database"
	"EutopiaQQBot/receive"
	_ "EutopiaQQBot/send"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go receive.InitRoute()
	wg.Add(1)
	go database.InitDatabase()
	wg.Wait()
}
