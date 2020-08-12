package main

import (
	"github.com/maxonrow/maxonrow-benchmark-go/bank"
	"github.com/maxonrow/maxonrow-go/app"
)

func main() {

	app.MakeDefaultCodec()
	bank.BankSend()
}
