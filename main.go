package main

import (
	"sync"

	"github.com/maxonrow/maxonrow-benchmark-go/bank"
	"github.com/maxonrow/maxonrow-benchmark-go/lib"
	"github.com/maxonrow/maxonrow-go/app"
)

var receiverList = 30

func main() {

	// added the go-routine
	var wg sync.WaitGroup
	wg.Add(1)
	_, receiverAccList := lib.CreateAddress(receiverList)
	app.MakeDefaultCodec()

	go func() {
		bank.BankSend(receiverAccList)
		wg.Done()
	}()
	wg.Wait()
}
