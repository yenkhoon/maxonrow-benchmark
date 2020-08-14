package bank

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/maxonrow/maxonrow-go/app"
	"github.com/maxonrow/maxonrow-go/x/bank"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	clientrpc "github.com/tendermint/tendermint/rpc/lib/client"
)

var tCdc *codec.Codec

var client = clientrpc.NewJSONRPCClient("http://192.168.20.219:26657")
//var client = clientrpc.NewJSONRPCClient("http://localhost:26657")

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

func BankSend(senders []string, receiverAccList [][]byte) {

	tKeys := readFile()

	//senders := []string{"gohck", "carlo", "mostafa", "nago", "jeansoon", "yk"}

	var txs [][]byte

	if len(senders) > 0 {

		for _, sender := range senders {
			acc := Account(tKeys[sender].addrStr)
			accNum := acc.GetAccountNumber()
			seq := acc.GetSequence()
			for i, receiver := range receiverAccList {

				receiverAddress := sdkTypes.AccAddress(receiver)
				//1.
				fees, _ := types.ParseCoins("800400000cin")
				amt, _ := types.ParseCoins("1cin")
				msg := bank.NewMsgSend(tKeys[sender].addr, receiverAddress, amt)

				if i > 0 {
					seq += uint64(1)
				}

				_, bz := makeSignedTx(sender, sender, seq, accNum, 0, fees, "", msg)

				txs = append(txs, bz)

				fmt.Printf(".")
			}
			//fmt.Printf("test case - (%v) with SignedTx Msg: %v\n", i+1, tx)
			//}()
		}
		if len(txs) > 0 {
			result := new(ctypes.ResultBroadcastTx)

			for _, tx := range txs {
				client.Call("broadcast_tx_async", map[string]interface{}{"tx": tx}, result)
			}
		}
	}

}

var store = map[string]uint64{}

func increaseSequence(accAddress string, seq uint64, acc sdkAuth.BaseAccount) uint64 {

	if seq < 1 {
		seq = acc.GetSequence()
	}

	store[accAddress] += seq
	return store[accAddress]

}

//Read the all the account in keys.json file
func readFile() map[string]*keyInfo {

	var keys []key
	content, _ := ioutil.ReadFile("./config/keys.json")
	json.Unmarshal(content, &keys)
	tKeys := make(map[string]*keyInfo)

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

	return tKeys
}

// for most of transactions, sender is same as signer.
// only for multi-sig transactions sender and signer are different.
func makeSignedTx(sender string, signer string, seq, accNum uint64, gas uint64, fees sdkTypes.Coins, memo string, msg sdkTypes.Msg) (sdkAuth.StdTx, []byte) {
	tKeys := readFile()

	//acc := Account(tKeys[sender].addrStr)
	// require.NotNil(t, acc, "alias:%s", sender)

	tCdc = app.MakeDefaultCodec()

	signMsg := authTypes.StdSignMsg{
		AccountNumber: accNum,
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

	//pub := tKeys[signer].priv.PubKey()
	stdSig := sdkAuth.StdSignature{
		//PubKey:    pub,
		Signature: sig,
	}

	sdtTx := authTypes.NewStdTx(signMsg.Msgs, signMsg.Fee, []authTypes.StdSignature{stdSig}, signMsg.Memo)

	bz, err := tCdc.MarshalBinaryLengthPrefixed(sdtTx)
	// fmt.Println("sdtTx [MarshalBinaryLengthPrefixed] : ", string(bz))
	if err != nil {
		panic(err)
	}
	return sdtTx, bz
}

func Account(addr string) *sdkAuth.BaseAccount {
	acc := new(sdkAuth.BaseAccount)

	var bg string
	_, err := client.Call("account_cdc", map[string]interface{}{"address": addr}, &bg)
	if err == nil {
		cdc := app.MakeDefaultCodec()
		err := cdc.UnmarshalJSON([]byte(bg), acc)
		if err != nil {
			fmt.Print("Error unmarshal account", err)
		}
		return acc
	}
	return acc
}

func BroadcastTxAsync(tx []byte) *ctypes.ResultBroadcastTx {
	result := new(ctypes.ResultBroadcastTx)
	_, err := client.Call("broadcast_tx_async", map[string]interface{}{"tx": tx}, result)
	if err == nil {
		return result
	}
	panic(err)
}

func BroadcastTxSync(tx []byte) *ctypes.ResultBroadcastTx {
	result := new(ctypes.ResultBroadcastTx)
	_, err := client.Call("broadcast_tx_sync", map[string]interface{}{"tx": tx}, result)
	if err == nil {
		return result
	}
	panic(err)
}

func BroadcastTxCommit(tx []byte) *ctypes.ResultBroadcastTxCommit {

	result := new(ctypes.ResultBroadcastTxCommit)
	_, err := client.Call("broadcast_tx_commit", map[string]interface{}{"tx": tx}, result)
	if err == nil {
		return result
	}
	panic(err)

}
