package chequebook

import (
	"os"
)


func GetTransactionHash()  (string) {


	TransactionHash := os.Getenv("TransactionHash")
	if len(TransactionHash) >0 {
		return TransactionHash
	}
	return ""
}
