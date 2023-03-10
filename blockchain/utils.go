package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

// Add new element in queue and return the queue
func enqueue[T comparable](queue []*T, element *T) []*T {
	queue = append(queue, element)
	return queue
}

// Remove the first element from the queue and return it.
func dequeue[T comparable](queue []*T) (*T, []*T, error) {
	if len(queue) == 0 {
		return nil, queue, errors.New("cannot dequeue from an empty queue")
	}
	removed := queue[0]
	return removed, queue[1:], nil
}

/*
This function converts an int64 (a 64-bit signed integer) to a slice of bytes in big-endian order.

bytes := ToHex(1234567890)
fmt.Println(bytes) // Output: [72, 144, 173, 205, 1, 0, 0, 0]
*/
func toHex(num int64) []byte {
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.BigEndian, num); err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

// MarshalJSON : Converting the transaction object to a JSON object.
func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender"`
		Recipient string  `json:"recipient"`
		Value     float32 `json:"Value"`
	}{
		Sender:    transaction.SenderWalletAddress,
		Recipient: transaction.RecipientWalletAddress,
		Value:     transaction.Value,
	})
}

// The Serialize method is a way to convert the data stored in a Transaction object into a slice of bytes that can be stored or transmitted. The serialized data contains all the information needed to recreate the original Transaction object, including the values of its fields.
func (transaction *Transaction) Serialize() []byte {
	var buffer bytes.Buffer

	prevKeyBytes, _ := x509.MarshalECPrivateKey(transaction.SenderPrivateKey)
	err := binary.Write(&buffer, binary.BigEndian, prevKeyBytes)
	if err != nil {
		return nil
	}

	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(transaction.SenderPublicKey)
	err = binary.Write(&buffer, binary.BigEndian, pubKeyBytes)
	if err != nil {
		return nil
	}

	err = binary.Write(&buffer, binary.BigEndian, []byte(transaction.SenderWalletAddress))
	if err != nil {
		return nil
	}

	err = binary.Write(&buffer, binary.BigEndian, []byte(transaction.RecipientWalletAddress))
	if err != nil {
		return nil
	}

	err = binary.Write(&buffer, binary.BigEndian, transaction.Value)
	if err != nil {
		return nil
	}

	return buffer.Bytes()
}

func DecodeStringToECDSA(encodedPriv string, encodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	encodedPriv = strings.TrimSpace(encodedPriv)
	encodedPub = strings.TrimSpace(encodedPub)

	privateKeyBytes, err1 := base64.StdEncoding.DecodeString(encodedPriv)
	if err1 != nil {
		return nil, nil, err1
	}

	privateKey, err1 := x509.ParseECPrivateKey(privateKeyBytes)
	if err1 != nil {
		return nil, nil, err1
	}

	publicKeyBytes, err2 := base64.StdEncoding.DecodeString(encodedPub)
	if err2 != nil {
		return nil, nil, err2
	}
	publicKeyInterface, err2 := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err2 != nil {
		return nil, nil, err2
	}

	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("error casting public key to ECDSA public key type")
	}

	return privateKey, publicKey, nil
}
