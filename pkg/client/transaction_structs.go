package client

import (
	"encoding/json"
	"fmt"
	"github.com/eteu-technologies/borsh-go"

	"github.com/maogongyin/near-api-go/pkg/types"
	"github.com/maogongyin/near-api-go/pkg/types/action"
	"github.com/maogongyin/near-api-go/pkg/types/hash"
	"github.com/maogongyin/near-api-go/pkg/types/signature"
)

type OutcomeStatus struct {
	SuccessValue     string          `json:"SuccessValue"`
	SuccessReceiptID string          `json:"SuccessReceiptId"`
	Failure          json.RawMessage `json:"Failure"` // TODO
}

type TransactionStatus struct {
	Enum         borsh.Enum `borsh_enum:"true"`
	NotStarted   StatusNotStarted
	Started      StatusStarted
	Failure      StatusFailure
	SuccessValue StatusSuccessValue
}

type StatusNotStarted struct {
}

type StatusStarted struct {
}

type StatusFailure struct {
	json.RawMessage
}

type StatusSuccessValue struct {
	json.RawMessage
}

const (
	ordNotStarted uint8 = iota
	ordStarted
	ordFailure
	ordSuccessValue
)

var (
	ordMappings = map[string]uint8{
		"NotStarted":   ordNotStarted,
		"Started":      ordStarted,
		"Failure":      ordFailure,
		"SuccessValue": ordSuccessValue,
	}

	simpleStatus = map[string]bool{
		"NotStarted": true,
		"Started":    true,
	}
)

func (t *TransactionStatus) UnderlyingValue() interface{} {
	switch uint8(t.Enum) {
	case ordNotStarted:
		return &t.NotStarted
	case ordStarted:
		return &t.Started
	case ordFailure:
		return &t.Failure
	case ordSuccessValue:
		return &t.SuccessValue
	}

	panic("unreachable")
}

func (t *TransactionStatus) UnmarshalJSON(b []byte) (err error) {
	var obj map[string]json.RawMessage

	// status can be either strings, or objects, so try deserializing into string first
	var statusType string
	if len(b) > 0 && b[0] == '"' {
		if err = json.Unmarshal(b, &statusType); err != nil {
			return
		}

		if _, ok := simpleStatus[statusType]; !ok {
			err = fmt.Errorf("Status '%s' had no body", statusType)
			return
		}

		obj = map[string]json.RawMessage{
			statusType: json.RawMessage(`{}`),
		}
	} else {
		if err = json.Unmarshal(b, &obj); err != nil {
			return
		}
	}

	if l := len(obj); l > 1 {
		err = fmt.Errorf("status object contains invalid amount of keys (expected: 1, got: %d)", l)
		return
	}

	for k := range obj {
		statusType = k
		break
	}

	ord := ordMappings[statusType]
	*t = TransactionStatus{Enum: borsh.Enum(ord)}
	ul := t.UnderlyingValue()

	if err = json.Unmarshal(obj[statusType], ul); err != nil {
		return
	}

	return nil
}

type SignedTransactionView struct {
	SignerID   types.AccountID           `json:"signer_id"`
	Nonce      types.Nonce               `json:"nonce"`
	ReceiverID types.AccountID           `json:"receiver_id"`
	Actions    []action.Action           `json:"actions"`
	Signature  signature.Base58Signature `json:"signature"`
	Hash       hash.CryptoHash           `json:"hash"`
}

type FinalExecutionOutcomeView struct {
	Status             TransactionStatus            `json:"status"`
	Transaction        SignedTransactionView        `json:"transaction"`
	TransactionOutcome ExecutionOutcomeWithIdView   `json:"transaction_outcome"`
	ReceiptsOutcome    []ExecutionOutcomeWithIdView `json:"receipts_outcome"`
}

type FinalExecutionOutcomeWithReceiptView struct {
	FinalExecutionOutcomeView
	Receipts []ReceiptView `json:"receipts"`
}

type ReceiptView struct {
	PredecessorID types.AccountID `json:"predecessor_id"`
	ReceiverID    types.AccountID `json:"receiver_id"`
	ReceiptID     hash.CryptoHash `json:"receipt_id"`
	Receipt       Receipt         `json:"receipt"`
}

type Receipt struct {
	Action struct {
		SignerId string `json:"signer_id"`
		//SignerPublicKey     string        `json:"signer_public_key"`
		GasPrice string `json:"gas_price"`
		//OutputDataReceivers []interface{} `json:"output_data_receivers"`
		//InputDataIds        []interface{} `json:"input_data_ids"`
		Actions []action.Action `json:"actions"`
	} `json:"Action"`
}

type ExecutionOutcomeView struct {
	Logs        []string          `json:"logs"`
	ReceiptIDs  []hash.CryptoHash `json:"receipt_ids"`
	GasBurnt    types.Gas         `json:"gas_burnt"`
	TokensBurnt types.Balance     `json:"tokens_burnt"`
	ExecutorID  types.AccountID   `json:"executor_id"`
	Status      OutcomeStatus     `json:"status"`
}

type MerklePathItem struct {
	Hash      hash.CryptoHash `json:"hash"`
	Direction string          `json:"direction"` // TODO: enum type, either 'Left' or 'Right'
}

type MerklePath = []MerklePathItem

type ExecutionOutcomeWithIdView struct {
	Proof     MerklePath           `json:"proof"`
	BlockHash hash.CryptoHash      `json:"block_hash"`
	ID        hash.CryptoHash      `json:"id"`
	Outcome   ExecutionOutcomeView `json:"outcome"`
}
