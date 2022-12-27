package blockchain

import (
	"bytes"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"errors"
)

func enqueue[T comparable](queue []*T, element *T) []*T {
	queue = append(queue, element)
	return queue
}

func dequeue[T comparable](queue []*T) (*T, []*T, error) {
	if len(queue) == 0 {
		return nil, queue, errors.New("cannot dequeue from an empty queue")
	}
	removed := queue[0]
	return removed, queue[1:], nil
}

func peakQueue[T comparable](queue []*T) *T {
	return queue[0]
}

func isQueueEmpty[T comparable](queue []*T) bool {
	return len(queue) == 0
}

// Converting the transaction object to a JSON object.
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
