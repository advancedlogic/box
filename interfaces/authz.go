package interfaces

type AuthZ interface {
	NewToken(string) (string, error)
	RefreshToken(string) (string, error)
	RevokeToken(string) error
	CheckToken(string) error
}

