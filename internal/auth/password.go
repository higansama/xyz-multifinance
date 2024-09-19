package auth

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/xdg-go/pbkdf2"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type HashOptions struct {
	Salt string
	Type string
}

type CompareOptions struct {
	Hash string
	Salt string
	Type string
}

func HashPassword(password string, opts HashOptions) string {
	var res string
	switch opts.Type {
	case "bcrypt":
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			panic(err)
		}
		res = string(hash)
	case "pbkdf2":
		bhash := pbkdf2.Key([]byte(password), []byte(opts.Salt), 10000, 512, sha512.New)
		res = hex.EncodeToString(bhash)
	default:
		panic(fmt.Sprintf("unsupported hash type '%s'", opts.Type))
	}

	return res
}

func VerifyPassword(password string, opts CompareOptions) bool {
	switch opts.Type {
	case "bcrypt":
		err := bcrypt.CompareHashAndPassword([]byte(opts.Hash), []byte(password))
		return err == nil
	case "pbkdf2":
		bhash := pbkdf2.Key([]byte(password), []byte(opts.Salt), 10000, 512, sha512.New)
		return opts.Hash == hex.EncodeToString(bhash)
	default:
		panic(fmt.Sprintf("unsupported hash type '%s'", opts.Type))
	}
}

func GenerateSalt() string {
	saltBytes := make([]byte, 16)
	rand.Seed(time.Now().UnixNano())
	rand.Read(saltBytes)
	return hex.EncodeToString(saltBytes)
}
