package data

import (
	"crypto/rand"
	"crypto/subtle"

	"golang.org/x/crypto/argon2"
)

// Argon2 parameters (Tuned for good security vs speed balance)
const (
	memory      = 64 * 1024 // 64MB
	iterations  = 3
	parallelism = 2
	keyLength   = 32
	saltLength  = 16
)

// hashPassword takes a plain text password and returns a unique salt and hash
func hashPassword(password string) (hash []byte, salt []byte, err error) {
	salt = make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	hash = argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	return hash, salt, nil
}

// VerifyPassword compares a plain text password against the stored hash and salt.
func VerifyPassword(password string, storedHash, storedSalt []byte) bool {
	computedHash := argon2.IDKey([]byte(password), storedSalt, iterations, memory, parallelism, keyLength)

	return subtle.ConstantTimeCompare(storedHash, computedHash) == 1
}
