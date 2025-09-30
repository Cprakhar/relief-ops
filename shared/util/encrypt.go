package util

import "golang.org/x/crypto/bcrypt"

// EncryptPassword hashes the given password using bcrypt.
func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ValidatePassword compares a hashed password with its possible plaintext equivalent.
func ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateToken is a placeholder function for generating tokens (e.g., JWT).
func GenerateToken(str string) (string, error) {
	// Placeholder implementation, replace with actual token generation logic
	return "token_for_" + str, nil
}
