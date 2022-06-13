package servers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/martinsaporiti/blockchain-sample/internal/blockchain"
	"github.com/martinsaporiti/blockchain-sample/internal/config"
	"github.com/martinsaporiti/blockchain-sample/internal/controller"
	"github.com/martinsaporiti/blockchain-sample/internal/dto"
)

type BlockchainServer struct {
	config     config.Config
	controller controller.Controller
}

func NewBlockchainServer(config config.Config, controller controller.Controller) *BlockchainServer {
	return &BlockchainServer{
		config:     config,
		controller: controller,
	}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.config.Port
}

func (bcs *BlockchainServer) GetChainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.controller.GetBlockchain()
		m, _ := json.Marshal(bc)
		io.WriteString(w, string(m[:]))
	default:
		log.Println("ERROR: Invalid request method")
		io.WriteString(w, string(dto.JsonStatus("fail")))
	}
}

func (bcs *BlockchainServer) TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		transactions := bcs.controller.GetTransactions()
		m, _ := json.Marshal(struct {
			Transactions []*blockchain.Transaction `json:"transactions"`
			Length       int                       `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		io.WriteString(w, string(m[:]))
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t dto.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %bcs", err)
			io.WriteString(w, string(dto.JsonStatus("fail")))
			return
		}

		if !t.Validate() {
			log.Println("ERROR: missing field(bcs)")
			io.WriteString(w, string(dto.JsonStatus("fail")))
			return
		}

		isCreated := bcs.controller.CreateTransaction(&t)
		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			m = dto.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = dto.JsonStatus("success")
		}
		io.WriteString(w, string(m))

	case http.MethodPut:
		decoder := json.NewDecoder(r.Body)
		var t dto.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %bcs", err)
			io.WriteString(w, string(dto.JsonStatus("fail")))
			return
		}

		if !t.Validate() {
			log.Println("ERROR: missing field(bcs)")
			io.WriteString(w, string(dto.JsonStatus("fail")))
			return
		}

		isUpdated := bcs.controller.AddTransaction(&t)
		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isUpdated {
			w.WriteHeader(http.StatusBadRequest)
			m = dto.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusOK)
			m = dto.JsonStatus("success")
		}
		io.WriteString(w, string(m))
	default:
		log.Println("ERROR: Invalid request method")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (bcs *BlockchainServer) AmountHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		blockchainAddress := r.URL.Query().Get("blockchain_address")
		amount := bcs.controller.CalculateTotalAmount(blockchainAddress)
		ar := &dto.AmountResponse{Amount: amount}
		w.Header().Add("Content-Type", "application/json")
		m, _ := json.Marshal(ar)
		w.Write(m)

	default:
		log.Println("ERROR: Invalid request method")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (bcs *BlockchainServer) AddNewBlockHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var proposedBlock blockchain.Block
		err := decoder.Decode(&proposedBlock)
		if err != nil {
			log.Printf("ERROR: %bcs", err)
			io.WriteString(w, string(dto.JsonStatus("fail")))
			return
		}
		bcs.controller.AddProposedBlockFromNetwork(&proposedBlock)
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(dto.JsonStatus("success")))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Start() {
	http.HandleFunc("/", bcs.GetChainHandler)
	http.HandleFunc("/transactions", bcs.TransactionsHandler)
	http.HandleFunc("/amount", bcs.AmountHandler)
	http.HandleFunc("/block", bcs.AddNewBlockHandler)
	log.Printf("Listening on port %d", bcs.config.Port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.config.Port)), nil))
}
