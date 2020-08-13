package main

import (
	"github.com/maxonrow/maxonrow-benchmark-go/bank"
	"github.com/maxonrow/maxonrow-benchmark-go/lib"
	"github.com/maxonrow/maxonrow-go/app"
)

var receiverList = 30

func main() {
	_, receiverAccList := lib.CreateAddress(receiverList)
	app.MakeDefaultCodec()
	bank.BankSend(receiverAccList)
}
