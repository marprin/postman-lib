package security

//go:generate mockgen -source=./contract.go -destination=./mock/contract.go -package=mock

type HashContract interface {
	CheckHashValidity(hash string, bareString string) error
	GenerateHash(bareString string) (*string, error)
}
