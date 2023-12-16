package encryptor

import (
	"github.com/0loff/grade_gophermart/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("Cannot create hash", zap.Error(err))
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func Compare(hash, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
