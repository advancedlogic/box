package interfaces

type AuthN interface {
	Login(string, string) (interface{}, error)
	Logout(string) error

	Register(string, string) (interface{}, error)
	Delete(string) error
	Reset(string, string) (interface{}, error)
}
