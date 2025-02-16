package responses

type WisdomResponse struct {
	Quote string `json:"quote"`
}

func (w *WisdomResponse) Encode() ([]byte, error) {
	buff := make([]byte, len(w.Quote))
	copy(buff, w.Quote)
	return buff, nil
}

func (w *WisdomResponse) Decode(buff []byte) error {
	w.Quote = string(buff)
	return nil
}
