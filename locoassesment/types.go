package main

import "sync"

// transaction information
type TransactionInfo struct {
	Amount    float64
	Txn_type  string
	Parent_id int64
}

// TxnInfo will have all transaction info per transaction id
type TxnInfo struct {
	Txn     map[int64]*TransactionInfo
	TxnLock sync.RWMutex
}

// SetTxnInfo will store transaction info
func (txn *TxnInfo) SetTxnInfo(id int64, txnInfo *TransactionInfo) {

	txn.TxnLock.Lock()
	txn.Txn[id] = txnInfo
	txn.TxnLock.Unlock()

}

// GetTxnInfo will fetch transaction info and return pointer
func (txn *TxnInfo) GetTxnInfo(id int64) *TransactionInfo {

	txn.TxnLock.RLock()
	t := txn.Txn[id]
	txn.TxnLock.RUnlock()
	return t

}

// GetAllTxnInfo will fetch all available transaction info
func (txn *TxnInfo) GetAllTxnInfo() []TransactionInfo {

	t := []TransactionInfo{}
	txn.TxnLock.RLock()
	for _, v := range txn.Txn {
		t = append(t, *v)
	}
	txn.TxnLock.RUnlock()
	return t

}

// Txn_per_type will have map of txn_id's per txn_type
type Txns_per_type struct {
	TxnIdList   map[string][]int64
	TxnTypeLock sync.RWMutex
}

// SetTxnByTxntype will store map of txn type and all txn_id of that txn type
func (t *Txns_per_type) SetTxnByTxntype(txnType string, id int64) {

	t.TxnTypeLock.Lock()
	l := t.TxnIdList[txnType]
	l = append(l, id)
	t.TxnIdList[txnType] = l
	t.TxnTypeLock.Unlock()

}

// GetTxnByTxntype will return list of txn_id of same txn_type
func (t *Txns_per_type) GetTxnByTxntype(txnType string) []int64 {

	t.TxnTypeLock.RLock()
	l := t.TxnIdList[txnType]
	t.TxnTypeLock.RUnlock()
	return l

}

// Parent_child_txn_map will have map of parent and child txn_id, key = parent_id and value = list of child_id
type Parent_child_txn_map struct {
	TxnIdMap           map[int64][]int64
	ParentChildTxnlock sync.RWMutex
}

// SetParentChildMap will store parent-child id in Parent_child_txn_map
func (p *Parent_child_txn_map) SetParentChildMap(parent_id, child_id int64) {

	p.ParentChildTxnlock.Lock()
	l := p.TxnIdMap[parent_id]
	l = append(l, child_id)
	p.TxnIdMap[parent_id] = l
	p.ParentChildTxnlock.Unlock()

}

// GetParentChildMap will return list of all children for given parent's txn_id
func (p *Parent_child_txn_map) GetParentChildMap(parent_id int64) []int64 {

	p.ParentChildTxnlock.RLock()
	l := p.TxnIdMap[parent_id]
	p.ParentChildTxnlock.RUnlock()
	return l

}
