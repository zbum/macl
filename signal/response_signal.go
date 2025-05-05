package controlsignal

import "fmt"

type ResponseSignal struct {
	TxId      string     `json:"txId"`
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	FiveTuple *FiveTuple `json:"fiveTuple,omitempty"`
}

func (r ResponseSignal) String() string {
	return fmt.Sprintf("txId: %s, success: %t, message: %s, fiveTuple: %s", r.TxId, r.Success, r.Message, r.FiveTuple)
}

func NewFailResponseSignal(txId string, message string, err error) *ResponseSignal {
	return &ResponseSignal{
		TxId:    txId,
		Success: false,
		Message: message + ": " + err.Error(),
	}
}

func NewSuccessResponseSignal(txId string, fiveTuple *FiveTuple) *ResponseSignal {
	return &ResponseSignal{
		TxId:      txId,
		Success:   true,
		Message:   "",
		FiveTuple: fiveTuple,
	}
}
