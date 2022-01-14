package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var txnInfoList *TxnInfo = &TxnInfo{Txn: make(map[int64]*TransactionInfo)}
var txnIdListByType *Txns_per_type = &Txns_per_type{TxnIdList: make(map[string][]int64)}
var parentChildIdList *Parent_child_txn_map = &Parent_child_txn_map{TxnIdMap: make(map[int64][]int64)}

func main() {

	// creates a new instance of a mux router
	router := mux.NewRouter()

	router.HandleFunc("/transactionservice/transaction/{transaction_id}", handleTxn) // GET, PUT, UPDATE, DELETE (PUT and UPDATE are handled in same way)
	router.HandleFunc("/transactionservice/types/{type}", listTxnByTxnType)          // GET
	router.HandleFunc("/transactionservice/sum/{transaction_id}", calSumOfAmount)    // GET
	router.HandleFunc("/transactionservice/transaction/list/", listAllTxn)           // GET

	http.ListenAndServe(":8080", router)

}

// handle PUT, GET and DELETE of transactions
func handleTxn(w http.ResponseWriter, r *http.Request) {

	// get transaction_id from path
	m := mux.Vars(r)
	txn_id, err := strconv.Atoi(m["transaction_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// request will be processed based on method passed
	if r.Method == "PUT" || r.Method == "POST" {

		// store/update transaction info
		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		txn := &TransactionInfo{}
		err = json.Unmarshal(data, txn)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = storeTxnInfo(txn, int64(txn_id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("status: OK"))

	} else if r.Method == "GET" {

		// get transactio info
		txnInfo := getTxnInfo(int64(txn_id))
		if txnInfo != nil {
			data, err := json.Marshal(txnInfo)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusFound)
			w.Write(data)

		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("transaction info not found"))
		}

	} else if r.Method == "DELETE" {

		err = deleteTxn(int64(txn_id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Status: OK"))

	} else {

		w.WriteHeader(http.StatusBadRequest)
		s := r.Method + " not supported"
		w.Write([]byte(s))

	}
}

func storeTxnInfo(txn *TransactionInfo, txn_id int64) error {

	// if txn have parent_id and that id does not exist, should not allow this kind of request
	if txn.Parent_id != 0 && txnInfoList.GetTxnInfo(txn.Parent_id) == nil {

		return errors.New("invalid parent_id passed in transaction")

	}

	txnInfoList.SetTxnInfo(txn_id, txn)

	txnIdListByType.SetTxnByTxntype(txn.Txn_type, txn_id)

	if txn.Parent_id != 0 {
		parentChildIdList.SetParentChildMap(txn.Parent_id, txn_id)
	}

	return nil

}

func getTxnInfo(txn_id int64) *TransactionInfo {

	return txnInfoList.GetTxnInfo(txn_id)
}

// 1. check if given txn_id have any child, if yes deletion will not be allowed
// 2. check txn_type and parent_id from txnInfo and delete the entry from txns_per_type and parent_child_txn_map as well
func deleteTxn(txn_id int64) error {

	txn := txnInfoList.GetTxnInfo(txn_id)

	if txn == nil {
		return errors.New("transaction not found")
	}
	if len(parentChildIdList.GetParentChildMap(txn_id)) != 0 {
		return errors.New("transaction already has children, deletion not allowed")
	}

	// delete from tnxInfoList
	delete(txnInfoList.Txn, txn_id)

	// delete from listByTxntype append(slice[:s], slice[s+1:]...)
	i := findIndex(txnIdListByType.TxnIdList[txn.Txn_type], txn_id)
	txnIdListByType.TxnIdList[txn.Txn_type] = append(txnIdListByType.TxnIdList[txn.Txn_type][:i], txnIdListByType.TxnIdList[txn.Txn_type][i+1:]...)

	// delete from parent_child_map list
	if txn.Parent_id != 0 {
		delete(parentChildIdList.TxnIdMap, txn_id)
		i = findIndex(parentChildIdList.TxnIdMap[txn.Parent_id], txn_id)
		parentChildIdList.TxnIdMap[txn.Parent_id] = append(parentChildIdList.TxnIdMap[txn.Parent_id][:i], parentChildIdList.TxnIdMap[txn.Parent_id][i+1:]...)
	}

	return nil

}

// find index of an element in given slice
func findIndex(slice []int64, txn_id int64) int {

	for k, v := range slice {
		if v == txn_id {
			return k
		}
	}
	return -1
}

// list transactions based on transaction type
func listTxnByTxnType(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {

		w.WriteHeader(http.StatusBadRequest)
		s := r.Method + " not supported"
		w.Write([]byte(s))
		return

	}

	// read txn_type from path
	v := mux.Vars(r)
	txnType := v["type"]

	// take list of txn_id for given txn_type
	txnIdList := txnIdListByType.GetTxnByTxntype(txnType)

	data, err := json.Marshal(txnIdList)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusFound)
	w.Write(data)

}

// calculate sum of amount for given transaction_id
// queue will be used to process this request
func calSumOfAmount(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		s := r.Method + " not supported"
		w.Write([]byte(s))
		return
	}

	// get transaction id from path
	m := mux.Vars(r)
	txn_id, err := strconv.Atoi(m["transaction_id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sum := calSum(int64(txn_id))

	w.WriteHeader(http.StatusFound)
	w.Write([]byte(fmt.Sprintf("%g", sum)))

}

func calSum(txn_id int64) float64 {

	idList := make([]int64, 0)
	idList = append(idList, int64(txn_id))

	var sum float64
	for len(idList) > 0 {

		// calculate sum
		sum = sum + (txnInfoList.GetTxnInfo(idList[0])).Amount
		idList = append(idList, (parentChildIdList.GetParentChildMap(idList[0]))...)
		idList = idList[1:]
	}
	return sum

}

// List all available transactions
func listAllTxn(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		s := r.Method + " not supported"
		w.Write([]byte(s))
		return
	}

	txnList := txnInfoList.GetAllTxnInfo()

	if txnList == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No transaction found"))
		return
	}

	data, err := json.Marshal(txnList)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusFound)
	w.Write(data)
}
