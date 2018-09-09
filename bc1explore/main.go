package main

//go:generate go run templates/gen.go

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const name = "Dumb Block Explorer (now in Go, at least)"

type (
	ChainInfo struct {
		Blocks   int    `json:"blocks"`
		BestHash string `json:"bestblockhash"`
	}

	Block struct {
		Height        int64    `json:"height"`
		TxCount       int64    `json:"nTx"`
		Size          int64    `json:"size"`
		Weight        int64    `json:"weight"`
		Hash          string   `json:"hash"`
		PrevHash      string   `json:"previousblockhash"`
		Confirmations int      `json:"confirmations"`
		NextHash      string   `json:"nextblockhash"`
		Ts            int64    `json:"mediantime"`
		Txs           []string `json:"tx"`
	}

	Vout struct {
		N            int     `json:"n"`
		Value        float64 `json:"value"`
		ScriptPubKey struct {
			Addresses []string `json:"addresses"`
			Type      string   `json:"type"`
		} `json:"scriptPubKey"`

		Spent bool
	}
	Vin struct {
		Txid string `json:"txid"`
		Vout int    `json:"vout"`

		PrevOut Vout
	}
	Tx struct {
		Id        string `json:"txid"`
		Hash      string `json:"hash"`
		Size      int    `json:"size"`
		VSize     int    `json:"vsize"`
		Weight    int    `json:"weight"`
		BlockHash string `json:"blockhash"`
		Vins      []Vin  `json:"vin"`
		Vouts     []Vout `json:"vout"`

		TotalIn,
		TotalOut,
		Fee,
		FeeRate float64

		JSON string
	}

	PageData struct {
		Testnet bool
		HtmlTitle,
		Title,
		BaseUrl string

		// overview
		Blocks []Block

		// block
		Block    Block
		ActiveTx string

		// transaction
		Tx Tx
		N  int64

		BreadCrumbs template.HTML
		Footer      string
	}
)

var (
	testnetEndpoint = flag.String("testnet", "http://127.0.0.1:18332", "Point to the REST interface of your testnet Bitcoin node")
	mainnetEndpoint = flag.String("mainnet", "http://127.0.0.1:8332", "Point to the REST interface of your mainnet Bitcoin node")
	port            = flag.Int("port", 8080, "What port should this blockchain explorer work on")

	baseUrl = fmt.Sprintf("http://127.0.0.1:%d", *port)

	templ, blockTempl, txTempl *template.Template

	defaultPageData = PageData{
		HtmlTitle: name,
		Title:     name,
		Testnet:   false,
		BaseUrl:   baseUrl,
	}

	funcMap = template.FuncMap{
		"add": func(i, j int) int { return i + j },
		"tounit": func(v float64, unit string) string {
			//TODO: distinction between "not set" & "zero"
			//if v == 0 {
			//	return "N/A"
			//}

			return fmt.Sprintf("%.08f %s", v, unit)
		},
	}
)

func (b Block) Url() string         { return fmt.Sprintf("%s/block/%s", baseUrl, b.Hash) }
func (b Block) Age() string         { return fmt.Sprintf("%d", b.Ts) }
func (b Block) Time() string        { return fmt.Sprintf("%d", b.Ts) }
func (vo Vout) Addresses() []string { return vo.ScriptPubKey.Addresses }
func (vo Vout) Type() string        { return vo.ScriptPubKey.Type }

func init() {
	flag.Parse()

	templ = template.Must(template.New("overview").Parse(_escFSMustString(false, "/templates/overview.html")))

	var err error
	blockTempl, err = template.Must(templ.Clone()).Funcs(funcMap).Parse(_escFSMustString(false, "/templates/block.html"))
	if err != nil {
		panic(err)
	}

	txTempl, err = template.Must(templ.Clone()).Funcs(funcMap).Parse(_escFSMustString(false, "/templates/tx.html"))
	if err != nil {
		panic(err)
	}
}

func getNodeUrl(testnet bool, path string) (url string) {
	url = *mainnetEndpoint
	if testnet {
		url = *testnetEndpoint
	}
	return fmt.Sprintf("%s/rest/%s", url, path)
}

func getInfo(testnet bool) (ci ChainInfo, err error) {
	url := getNodeUrl(testnet, "chaininfo.json")
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&ci)
	return
}

func getBlock(testnet bool, blockHash string, headerOnly bool) (block Block, err error) {
	urlTemplate := "block/notxdetails/%s.json"
	if headerOnly {
		urlTemplate = "headers/1/%s.json"
	}

	url := getNodeUrl(testnet, fmt.Sprintf(urlTemplate, blockHash))
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if !headerOnly {
		err = json.NewDecoder(res.Body).Decode(&block)
		return
	}

	var blockArr []Block
	err = json.NewDecoder(res.Body).Decode(&blockArr)
	return blockArr[0], err
}
func getBlocks(testnet bool, lastBlock string, count int) (blocks []Block, err error) {
	block, err := getBlock(testnet, lastBlock, false)
	if err != nil {
		block, err = getBlock(testnet, lastBlock, true)
		if err != nil {
			return
		}
	}

	if count > 1 {
		blocks, err = getBlocks(testnet, block.PrevHash, count-1)
	}

	return append([]Block{block}, blocks...), nil
}

func getUTXO(testnet bool, txid string, n int) (spent bool, err error) {
	url := getNodeUrl(testnet, fmt.Sprintf("getutxos/%s-%d.json", txid, n))
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	x := struct{ Bitmap string `json:"bitmap"` }{}
	err = json.NewDecoder(res.Body).Decode(&x)
	if x.Bitmap == "0" {
		return true, nil
	}

	return
}
func getTx(testnet bool, txid string, complete bool) (tx Tx, err error) {
	url := getNodeUrl(testnet, fmt.Sprintf("tx/%s.json", txid))
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &tx)
	if err != nil {
		return
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, data, "", "  ")

	if complete {
		tx.JSON = prettyJSON.String()

		for i, vin := range tx.Vins {
			if len(vin.Txid) == 0 {
				continue
			}

			prevTx, err := getTx(testnet, vin.Txid, false)
			if err != nil {
				continue
			}

			tx.Vins[i].PrevOut = prevTx.Vouts[vin.Vout]
			tx.TotalIn += prevTx.Vouts[vin.Vout].Value
		}

		for i, vout := range tx.Vouts {
			tx.TotalOut += vout.Value
			tx.Vouts[i].Spent, err = getUTXO(testnet, txid, vout.N)
		}

		tx.Fee = tx.TotalIn - tx.TotalOut
		tx.FeeRate = tx.Fee * 1e8 / float64(tx.VSize)
	}

	return
}

func overview(w http.ResponseWriter, testnet bool) {
	ci, err := getInfo(testnet)
	if err != nil {
		http.Error(w, "unable to get chaindata", 500)
		return
	}

	blocks, err := getBlocks(testnet, ci.BestHash, 10)
	if err != nil {
		http.Error(w, "unable to get recent blocks", 500)
		return
	}

	pd := defaultPageData
	pd.Testnet = testnet
	pd.Blocks = blocks

	err = templ.Execute(w, pd)
	if err != nil {
		http.Error(w, "unable to get render overview view", 500)
		return
	}
}

func block(w http.ResponseWriter, testnet bool, hash, activeTxId string) {
	blocks, err := getBlocks(testnet, hash, 1)
	if err != nil {
		http.Error(w, "unable to get block", 500)
		log.Println(err)
		return
	}

	block := blocks[0]
	pd := defaultPageData
	pd.HtmlTitle = fmt.Sprintf("%s, height %d", name, block.Height)
	pd.Testnet = testnet
	pd.Block = block

	if len(activeTxId) > 0 {
		pd.ActiveTx = activeTxId
	}

	err = blockTempl.Execute(w, pd)
	if err != nil {
		http.Error(w, "unable to get render block view", 500)
		log.Println(err)
		return
	}
}

func tx(w http.ResponseWriter, testnet bool, txid, n string) {
	tx, err := getTx(testnet, txid, true)
	if err != nil {
		http.Error(w, "unable to get transaction", 500)
		log.Println(err)
		return
	}

	pd := defaultPageData
	pd.Testnet = testnet
	pd.Tx = tx

	// NOTE: discard errors - missing data here is ok
	pd.N, _ = strconv.ParseInt(n, 10, 64)
	pd.Block, _ = getBlock(testnet, tx.BlockHash, true)

	err = txTempl.Execute(w, pd)
	if err != nil {
		http.Error(w, "unable to get render transaction view", 500)
		log.Println(err)
		return
	}
}

func simpleRouter(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")

	var testnet bool
	if strings.ToLower(parts[0]) == "testnet" {
		testnet, parts = true, parts[1:]
	}

	if len(parts) == 1 {
		overview(w, testnet)
		return
	}

	route, id, parts := parts[0], parts[1], parts[2:]

	switch strings.ToLower(route) {
	case "block":
		var txid string
		if len(parts) >= 2 && strings.ToLower(parts[0]) == "tx" {
			txid = parts[1]
		}
		block(w, testnet, id, txid)

	case "tx":
		var n string
		if len(parts) >= 2 && strings.ToLower(parts[0]) == "n" {
			n = parts[1]
		}
		tx(w, testnet, id, n)

	case "search":
		return

	default:
		//TODO: 404 not found
	}

	return
}

func main() {
	http.HandleFunc("/", simpleRouter)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
