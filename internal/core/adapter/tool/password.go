package tool

type Password interface {
	Hash(password string) (string, error)

	Verify(password string, hash string) bool
}
