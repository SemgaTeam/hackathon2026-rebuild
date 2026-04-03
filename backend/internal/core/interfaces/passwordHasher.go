package interfaces

type IPasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}
