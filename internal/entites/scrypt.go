package entites

type ScryptConfig struct {
	N      int `json:"n"`
	R      int `json:"r"`
	P      int `json:"p"`
	KeyLen int `json:"key_len"`
}
