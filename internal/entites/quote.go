package entites

type Quote struct {
	ID    int    `json:"-" db:"q_id"`
	Quote string `json:"quote" db:"quote"`
	By    string `json:"by" db:"by"`
}
