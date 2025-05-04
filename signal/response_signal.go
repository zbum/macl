package controlsignal

type ResponseSignal struct {
	TxId      int        `json:"txId"`
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	FiveTuple *FiveTuple `json:"fiveTuple,omitempty"`
}

func NewFailResponseSignal(txId int, message string, err error) *ResponseSignal {
	return &ResponseSignal{
		TxId:    txId,
		Success: false,
		Message: message + ": " + err.Error(),
	}
}

func NewSuccessResponseSignal(txId int, fiveTuple *FiveTuple) *ResponseSignal {
	return &ResponseSignal{
		TxId:      txId,
		Success:   true,
		Message:   "",
		FiveTuple: fiveTuple,
	}
}
