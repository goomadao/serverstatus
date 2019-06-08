package auth

var (
	Password string
)

func login(key []byte) (bool, error) {
	if string(key) == Password {
		return true, nil
	}
	return false, nil
}
