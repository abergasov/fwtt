package entites

const (
	SHA256 = "sha256"
	Scrypt = "scrypt"
)

func ValidateHashAlgo(algo string) bool {
	return algo == SHA256 || algo == Scrypt
}
