package lib

import (
	"fmt"

	cliKeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
)

//create the key_address

func CreateAddress(n int) (error, []string) {
	var accounts []string
	keybase, kbErr := cliKeys.NewKeyBaseFromHomeFlag()
	if kbErr != nil {
		return kbErr, nil
	}
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("acc-%v", i+1)
		info, mnemonic, err := keybase.CreateMnemonic(name, keys.English, "12345678", keys.Secp256k1)

		if err != nil {
			return fmt.Errorf("Unable to create new account: %v", err), nil
		}
		fmt.Printf("Create new account. name: %v, address: %s, mnemonic:%s\n", name, info.GetAddress(), mnemonic)

		addr := info.GetAddress()

		accounts = append(accounts, addr.String())

	}
	return nil, accounts

}
