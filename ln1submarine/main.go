package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/lightningnetwork/lnd/input"
	"github.com/meeDamian/bc1toolkit/lib/common"
	"github.com/meeDamian/bc1toolkit/lib/help"
)

const (
	BinaryName = "ln1submarine"

	description = `Constructs and returns a Bitcoin submarine atomic swap script.

Setup an atomic swap unlocking BTC/on-chain payment, only upon successful LN/off-chain payment.`
)

type Swap struct {
	MyPubKey    btcec.PublicKey
	TheirPubKey btcec.PublicKey
	Hash        string
	Height      int64
}

var (
	opts struct {
		TestNet bool   `long:"testnet" short:"T" description:"Uses testnet addresses.  Connects to testnet nodes."`
		PubKey  string `long:"pubkey" description:""`
		Bitcoin string `long:"bitcoin" description:"Bitcoin's RPC interface" default:"localhost:8332" default:"localhost:18332" default-mask:"localhost:8332, or localhost:18332 for testnet"`

		Lnd struct {
			Host     string `long:"host" description:"Lnd's RPC interface'" default:"localhost:10009"`
			Cert     string `long:"cert" description:"TLS certificate that LND uses to auth RPC connections'" default-mask:"~/.lnd/tls.cert"`
			Macaroon string `long:"macaroon" description:"Path to admin.macaroon lnd generated.  See description for \"why admin?\"" default:"~/.lnd/data/chain/bitcoin/mainnet/admin.macaroon" default:"~/.lnd/data/chain/bitcoin/testnet/admin.macaroon" default-mask:"~/.lnd/data/chain/bitcoin/(main|test)net/admin.macaroon"`
		} `namespace:"lnd" group:"lnd"`
	}
	commonOpts help.Opts

	swap Swap
)

func init() {
	common.Logger.Name(BinaryName)

	help.Customize(
		"[OPTIONS] their-pubkey [invoice|hash] [your-pubkey] [height]",
		description,
		help.DisableTor,
		BinaryName, &opts,
	)

	swap = Swap{}

	// read parameters passed to a binary
	var args []string
	args, commonOpts = help.Parse()

	if len(args) == 0 {
		fmt.Println("At least the pubkey of redeemer has to be provided")
		os.Exit(1)
	}
}

func genSubmarineSwapScript(swapperPubKey, payerPubKey, hash []byte, lockHeight int64) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_HASH160)
	builder.AddData(input.Ripemd160H(hash))
	builder.AddOp(txscript.OP_EQUAL) // Leaves 0P1 (true) on the stack if preimage matches
	builder.AddOp(txscript.OP_IF)
	builder.AddData(swapperPubKey) // Path taken if preimage matches
	builder.AddOp(txscript.OP_ELSE)
	builder.AddInt64(lockHeight)
	builder.AddOp(txscript.OP_CHECKSEQUENCEVERIFY)
	builder.AddOp(txscript.OP_DROP)
	builder.AddData(payerPubKey) // Refund back to payer
	builder.AddOp(txscript.OP_ENDIF)
	builder.AddOp(txscript.OP_CHECKSIG)

	return builder.Script()
}


func main() {

	b, err := hex.DecodeString(opts.PubKey)

	key, err := btcec.ParsePubKey(b, btcec.S256())

	fmt.Println(opts.PubKey, len(opts.PubKey), key, err)
	//client, err := rpcclient.New(&rpcclient.ConnConfig{
	//	HTTPPostMode: true,
	//	DisableTLS:   true,
	//	Host:         "127.0.0.1:8332",
	//	User:         "testing",
	//	Pass:         "T4NHKPmdDKQeF1DvZLfbhP6WSNm2oQqv-lzk1WSsbvs=",
	//}, nil)
	//if err != nil {
	//	log.Fatalf("error creating new btc client: %v", err)
	//}
	//
	//// list accounts
	//info, err := client.GetPeerInfo()
	//if err != nil {
	//	log.Fatalf("error getting info: %v", err)
	//}
	//
	//log.Println(info)

}
