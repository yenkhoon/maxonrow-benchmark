package main

import (
	"sync"

	"github.com/maxonrow/maxonrow-benchmark-go/bank"
	"github.com/maxonrow/maxonrow-benchmark-go/lib"
	"github.com/maxonrow/maxonrow-go/app"
)

var receiverList = 5000

func main() {

	// added the go-routine
	var wg sync.WaitGroup
	wg.Add(1)
	_, receiverAccList := lib.CreateAddress(receiverList)
	app.MakeDefaultCodec()

	//bank.BankSend([]string{"jeansoon"}, receiverAccList)

	go func() {
		sender := []string{"jeansoon"}
		bank.BankSend(sender, receiverAccList)
		wg.Done()
	}()

	go func() {
		sender := []string{"yk"}
		bank.BankSend(sender, receiverAccList)
		wg.Done()
	}()

	go func() {
		sender := []string{"alice"}
		bank.BankSend(sender, receiverAccList)
		wg.Done()
	}()

	go func() {
		sender := []string{"bob"}
		bank.BankSend(sender, receiverAccList)
		wg.Done()
	}()

	go func() {
		sender := []string{"eve"}
		bank.BankSend(sender, receiverAccList)
		wg.Done()
	}()

	go func() {
		sender := []string{"maintainer-1"}
		bank.BankSend(sender, receiverAccList)
		wg.Done()
	}()
	go func() {
		sender := []string{"maintainer-2"}
		bank.BankSend(sender, receiverAccList)
		wg.Done()
	}()
	go func() {
		sender := []string{"maintainer-3"}
		bank.BankSend(sender, receiverAccList)
		wg.Done()
	}()
	wg.Wait()
}
