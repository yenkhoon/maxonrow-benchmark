package bank

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	util "github.com/maxonrow/maxonrow-go/tests/"
)

type bankInfo struct {
	from   string
	to     string
	amount string
}

type keyInfo struct {
	addr    sdkTypes.AccAddress
	priv    tmCrypto.PrivKey
	pub     tmCrypto.PubKey
	addrStr string
}

var tKeys map[string]*keyInfo

func processBankSend() {

	//0. read from keys.json of sender list
	readFileKeyJson()

	// Test-cases example :
	var caseDesc := "alice sending 1 cin to Bob"

	//1.
	var fees, _ = types.ParseCoins("200000000cin")
	var amt, _ = types.ParseCoins("1cin")
	var msg = bank.NewMsgSend(tKeys["alice"].addr, tKeys["bob"].addr, amt)

	//2.
	var tx, bz = makeSignedTx("alice", "alice", 0, 0, fees, "MEMO: P2P sending.......", msg)

	//3.
	var res = util.BroadcastTxAsync(bz)
	var resHash = res.Hash.Bytes()

	fmt.Printf("test case (%v) with CheckTx.Log : %v\n", caseDesc, res.CheckTx.Log)
	fmt.Printf("test case (%v) with DeliverTx.Log : %v\n", caseDesc, res.DeliverTx.Log)

}

func readFileForSenderKeyJson() {

	type key struct {
		Name        string
		MasterPriv  string
		DerivedPriv string
		Address     string
		Mnemonic    string
	}

	var keys []key
	content, _ := ioutil.ReadFile("../config/keys.json")
	json.Unmarshal(content, &keys)
	tKeys = make(map[string]*keyInfo)

	for _, k := range keys {
		bz, _ := hex.DecodeString(k.DerivedPriv)
		var priv [32]byte
		copy(priv[:], bz)
		addr, _ := sdkTypes.AccAddressFromBech32(k.Address)

		tKeys[k.Name] = &keyInfo{
			addr,
			secp256k1.PrivKeySecp256k1(priv),
			secp256k1.PrivKeySecp256k1(priv).PubKey(),
			k.Address,
		}

	}
}

// for most of transactions, sender is same as signer.
// only for multi-sig transactions sender and signer are different.
func makeSignedTx(sender string, signer string, seq uint64, gas uint64, fees sdkTypes.Coins, memo string, msg sdkTypes.Msg) (sdkAuth.StdTx, []byte) {
	acc := util.Account(tKeys[sender].addrStr)
	// require.NotNil(t, acc, "alias:%s", sender)

	seq = acc.GetSequence() //no need get from stateDB, directly +1 base from current

	signMsg := authTypes.StdSignMsg{
		AccountNumber: acc.GetAccountNumber(),
		ChainID:       "maxonrow-chain",
		Fee:           authTypes.NewStdFee(gas, fees),
		Memo:          memo,
		Msgs:          []sdkTypes.Msg{msg},
		Sequence:      seq,
	}

	signBz, signBzErr := tCdc.MarshalJSON(signMsg)
	if signBzErr != nil {
		panic(signBzErr)
	}

	sig, err := tKeys[signer].priv.Sign(sdkTypes.MustSortJSON(signBz))
	if err != nil {
		panic(err)
	}

	pub := tKeys[signer].priv.PubKey()
	stdSig := sdkAuth.StdSignature{
		PubKey:    pub,
		Signature: sig,
	}

	sdtTx := authTypes.NewStdTx(signMsg.Msgs, signMsg.Fee, []authTypes.StdSignature{stdSig}, signMsg.Memo)
	bz, err := tCdc.MarshalBinaryLengthPrefixed(sdtTx)
	if err != nil {
		panic(err)
	}
	return sdtTx, bz
}
