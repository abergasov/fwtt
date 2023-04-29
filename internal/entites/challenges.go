package entites

type Challenges struct {
	Challenges []string      `json:"challenges"`
	Difficulty uint32        `json:"difficulty"`
	Algorithm  string        `json:"algorithm"`
	AlgoParams *ScryptConfig `json:"algo_params,omitempty"`
}
