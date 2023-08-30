package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
}

func New() *Hasher{
	return &Hasher{}
}

func (h *Hasher) HashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func (h *Hasher) IsHashedPassEquals(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
