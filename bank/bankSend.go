package bank

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/maxonrow/maxonrow-benchmark/lib"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/maxonrow/maxonrow-go/app"
	util "github.com/maxonrow/maxonrow-go/tests"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	client "github.com/tendermint/tendermint/rpc/lib/client"
)

var tCdc *codec.Codec

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

type key struct {
	Name        string
	MasterPriv  string
	DerivedPriv string
	Address     string
	Mnemonic    string
}

var tKeys map[string]*keyInfo

func BankSend() {

	//0.1 read from keys.json of sender list
	readFile()

	//0.2 read from ArrayList of receiver list
	_, receiverAccList := lib.CreateAddress(30)

	for i, receiver := range receiverAccList {

		receiverAddress, _ := sdkTypes.AccAddressFromBech32(receiver)
		//1.
		fees, _ := types.ParseCoins("200000000cin")
		amt, _ := types.ParseCoins("1cin")
		msg := bank.NewMsgSend(tKeys["alice"].addr, receiverAddress, amt)

		//2.
		tx, bz := makeSignedTx(i, "alice", "alice", 1, 0, fees, "MEMO: P2P sending.......", msg)
		fmt.Printf("test case - (%v) with SignedTx Msg: %v\n", i+1, tx)

		//3.
		res := util.BroadcastTxAsync(bz)
		resHash := res.Hash.Bytes()

		fmt.Printf("test case - (%v) with Response.Log : %v\n", i+1, resHash)

	}

}

var store = map[string]uint64{}

func increaseSequence(accAddress string, seq uint64, acc sdkAuth.BaseAccount) uint64 {

	if seq < 1 {
		seq = acc.GetSequence()
	}

	store[accAddress] += seq
	return store[accAddress]
	//fmt.Printf("%v\n", store[accAddress])

}

//Read the all the account in keys.json file
func readFile() {

	var keys []key
	content, _ := ioutil.ReadFile("./config/keys.json")
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
func makeSignedTx(i int, sender string, signer string, seq uint64, gas uint64, fees sdkTypes.Coins, memo string, msg sdkTypes.Msg) (sdkAuth.StdTx, []byte) {

	acc := Account(tKeys[sender].addrStr)

	// require.NotNil(t, acc, "alias:%s", sender)

	//seq := increaseSequence(tKeys["alice"].addr, i, acc)
	seq = 1
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

func Account(addr string) *sdkAuth.BaseAccount {

	client := client.NewJSONRPCClient("http://localhost:26657")
	ctypes.RegisterAmino(client.Codec())
	var bg string
	_, err := client.Call("account_cdc", map[string]interface{}{"address": addr}, &bg)
	if err == nil {
		acc := new(sdkAuth.BaseAccount)
		cdc := app.MakeDefaultCodec()

		err := cdc.UnmarshalJSON([]byte(bg), acc)
		if err != nil {
			fmt.Print("Error unmarshal account", err)
		}
		return acc
	}
	return nil
}
