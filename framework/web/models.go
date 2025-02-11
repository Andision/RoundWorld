package web

type tableMessage struct {
	sender  *Player
	message []byte
}

type tableResponse struct {
	TableId string      `json:"table_id"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	tableResponseCodeUpdate = iota
	tableResponseCodeError
	tableResponseCodeEnd
	tableResponseCodeSuccess
)
