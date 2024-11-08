package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
	"winqroo/config"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

type HashingSystemUtils struct {
	SecretKey string
}

func NewHashingSystemUtils() *HashingSystemUtils {
	return &HashingSystemUtils{
		SecretKey: config.GetHashingSecretKey(),
	}
}

// GenerateSalt generates a random salt.
func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

// HashPassword hashes a password securely using Argon2id and bcrypt.
func (h *HashingSystemUtils) HashPassword(password string) (string, error) {
	// Step 1: Generate a unique salt for this password
	salt, err := GenerateSalt(16)
	if err != nil {
		return "", err
	}

	// Step 2: Use Argon2id for key derivation
	argonHash := argon2.IDKey([]byte(password+h.SecretKey), []byte(salt), 1, 64*1024, 4, 32)

	// Step 3: Hash the Argon2id output with bcrypt for additional security
	finalHash, err := bcrypt.GenerateFromPassword(argonHash, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Combine the salt and the bcrypt hash
	fullHash := salt + ":" + base64.RawStdEncoding.EncodeToString(finalHash)
	return fullHash, nil
}

// ComparePassword compares a hashed password with a plain text password.
func (h *HashingSystemUtils) ComparePassword(storedHash, password string) (bool, error) {
	// Split the stored hash to retrieve the salt and the bcrypt hash
	parts := strings.Split(storedHash, ":")
	if len(parts) != 2 {
		return false, errors.New("invalid hash format")
	}
	salt := parts[0]
	bcryptHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	// Recreate the Argon2id hash
	argonHash := argon2.IDKey([]byte(password+h.SecretKey), []byte(salt), 1, 64*1024, 4, 32)

	// Compare the bcrypt hash of the Argon2id result with the stored hash
	err = bcrypt.CompareHashAndPassword(bcryptHash, argonHash)
	return err == nil, nil
}
