package tool

type Token interface {
	Generate(userID uint64) (string, error)

	Parse(token string) (uint64, error)
}
