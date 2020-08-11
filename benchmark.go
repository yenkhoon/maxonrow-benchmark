package main

import (
	bankSend "github.com/maxonrow/maxonrow-benchmark/bank"
	"github.com/maxonrow/maxonrow-go/app"
	//rpc "github.com/maxonrow/maxonrow-go/tests"
)

func main() {

	app.MakeDefaultCodec()

	bankSend.BankSend()
}
