package servers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/martinsaporiti/blockchain-sample/internal/blkcrypto"
	"github.com/martinsaporiti/blockchain-sample/internal/dto"
	"github.com/martinsaporiti/blockchain-sample/internal/wallet"
)

var tempDir1 = "../../internal/wallet/templates"
var tempDir2 = "./internal/wallet/templates"

type Server struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *Server {
	return &Server{port, gateway}
}

func (ws *Server) Port() uint16 {
	return ws.port
}

func (ws *Server) Gateway() string {
	return ws.gateway
}

func (ws *Server) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println(os.Getwd())
	currentPath, _ := os.Getwd()
	var tempDir string
	if strings.Contains(currentPath, "cmd") {
		tempDir = tempDir1
	} else {
		tempDir = tempDir2
	}

	switch r.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		t.Execute(w, "")
	default:
		log.Println("ERROR: Invalid request method")
	}
}

func (ws *Server) Wallet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.New()
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid request method")
	}
}

func (ws *Server) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t wallet.TransactionRequest
		if err := decoder.Decode(&t); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("ERROR: %s", err)
			io.WriteString(w, string(dto.JsonStatus("fail")))
			return
		}
		if !t.IsValid() {
			log.Println("ERROR: Invalid missing fields")
			io.WriteString(w, string(dto.JsonStatus("fail")))
			return
		}

		publicKey := blkcrypto.PublicKeyFromString(*t.SenderPrivateKey)
		privateKey := blkcrypto.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			io.WriteString(w, string(dto.JsonStatus("fail")))
			return
		}

		value32 := float32(value)
		timestamp := time.Now().Unix()
		transaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress, value32, timestamp)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()
		bt := &dto.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value32,
			Timestamp:                  &timestamp,
			Signature:                  &signatureStr,
		}

		m, _ := json.Marshal(bt)
		buf := bytes.NewBuffer(m)
		resp, err := http.Post(ws.Gateway()+"/transactions", "application/json", buf)

		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			io.WriteString(w, string(dto.JsonStatus("fail")))
			return
		}

		if resp.StatusCode == http.StatusCreated {
			log.Println("Transaction created")
			io.WriteString(w, string(dto.JsonStatus("Money sent")))
			return
		}

		io.WriteString(w, string(dto.JsonStatus("fail")))
		return

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid request method")
	}
}

func (ws *Server) WalletAmountHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		blockchainAddress := r.URL.Query().Get("blockchain_address")
		endpoint := fmt.Sprintf("%s/amount", ws.Gateway())
		client := http.Client{}
		bcsReq, _ := http.NewRequest(http.MethodGet, endpoint, nil)
		q := bcsReq.URL.Query()
		q.Add("blockchain_address", blockchainAddress)
		bcsReq.URL.RawQuery = q.Encode()

		resp, err := client.Do(bcsReq)

		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			io.WriteString(w, string(dto.JsonStatus("fail")))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if resp.StatusCode == http.StatusOK {
			decoder := json.NewDecoder(resp.Body)

			var bar dto.AmountResponse
			if err := decoder.Decode(&bar); err != nil {
				io.WriteString(w, string(dto.JsonStatus("fail")))
				return
			}

			m, _ := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			})

			io.WriteString(w, string(m[:]))
			return
		}

		io.WriteString(w, string(dto.JsonStatus("fail")))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid request method")
	}
}

func (ws *Server) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransactionHandler)
	http.HandleFunc("/wallet/amount", ws.WalletAmountHandler)
	log.Printf("Listening on port %d", ws.port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.port)), nil))
}
