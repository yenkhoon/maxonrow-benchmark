package lib

import (
	"crypto/rand"
)

//create the key_address

func CreateAddress(n int) (error, [][]byte) {
	accounts := make([][]byte, n)
	for i := 0; i < n; i++ {
		accounts[i] = make([]byte, 20)
		rand.Read(accounts[i])

	}
	return nil, accounts
}
