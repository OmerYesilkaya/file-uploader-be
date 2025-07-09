package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	return salt, err
}

func HashPassword(password string) (string, error) {
	// Parameters
	salt, err := generateSalt(16)
	if err != nil {
		return "", err
	}
	time := uint32(1)
	memory := uint32(64 * 1024) // 64 MB
	threads := uint8(4)
	keyLen := uint32(32)

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// Return as encoded string (include params + salt for later verification)
	encoded := fmt.Sprintf("$argon2id$v=19$t=%d$m=%d$p=%d$%s$%s",
		time, memory, threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

func CheckPasswordHash(password, hash string) error {
	// Parse the hash
	var time, memory, threads uint32
	var saltB64, keyB64 string

	parts := strings.Split(hash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return nil
	}

	// extract values from parts[2] to parts[5]
	timeStr := strings.TrimPrefix(parts[2], "t=")
	memoryStr := strings.TrimPrefix(parts[3], "m=")
	threadsStr := strings.TrimPrefix(parts[4], "p=")

	_, err := fmt.Sscanf(timeStr, "%d", &time)
	if err != nil {
		return err
	}
	_, err = fmt.Sscanf(memoryStr, "%d", &memory)
	if err != nil {
		return err
	}
	_, err = fmt.Sscanf(threadsStr, "%d", &threads)
	if err != nil {
		return err
	}

	saltB64 = parts[5]
	keyB64 = parts[6]

	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return err
	}
	key, err := base64.RawStdEncoding.DecodeString(keyB64)
	if err != nil {
		return err
	}

	// Hash the input password with the same parameters
	computedHash := argon2.IDKey([]byte(password), salt, time, memory, uint8(threads), uint32(len(key)))

	if subtle.ConstantTimeCompare(computedHash, key) == 1 {
		return nil // Password matches
	} else {
		return fmt.Errorf("password does not match") // Password does not match
	}
}

func GenerateJWT(userID string) (string, error) {
	fmt.Printf("secret: %s", secret)

	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return secret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["userId"].(string), nil
	}

	return "", err
}
