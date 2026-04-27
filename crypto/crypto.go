package crypto

import (
	"crypto/aes"
	gocipher "crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"

	"github.com/tengolang/tengo/v3"
)

// Module is the Tengo "crypto" module.
//
//	crypto := import("crypto")
//	crypto.sha256(data bytes) => bytes
//	crypto.sha256_hex(data bytes) => string
//	crypto.sha512(data bytes) => bytes
//	crypto.sha512_hex(data bytes) => string
//	crypto.hmac_sha256(key bytes, data bytes) => bytes
//	crypto.hmac_sha256_hex(key bytes, data bytes) => string
//	crypto.aes_encrypt(key bytes, plaintext bytes) => bytes | error
//	crypto.aes_decrypt(key bytes, ciphertext bytes) => bytes | error
var Module = map[string]tengo.Object{
	"sha256": &tengo.UserFunction{
		Name: "sha256",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 1); err != nil {
				return nil, err
			}
			b, err := tengo.ArgBytes(args, 0, "data")
			if err != nil {
				return nil, err
			}
			sum := sha256.Sum256(b)
			return &tengo.Bytes{Value: sum[:]}, nil
		},
	},

	"sha256_hex": &tengo.UserFunction{
		Name: "sha256_hex",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 1); err != nil {
				return nil, err
			}
			b, err := tengo.ArgBytes(args, 0, "data")
			if err != nil {
				return nil, err
			}
			sum := sha256.Sum256(b)
			return &tengo.String{Value: hex.EncodeToString(sum[:])}, nil
		},
	},

	"sha512": &tengo.UserFunction{
		Name: "sha512",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 1); err != nil {
				return nil, err
			}
			b, err := tengo.ArgBytes(args, 0, "data")
			if err != nil {
				return nil, err
			}
			sum := sha512.Sum512(b)
			return &tengo.Bytes{Value: sum[:]}, nil
		},
	},

	"sha512_hex": &tengo.UserFunction{
		Name: "sha512_hex",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 1); err != nil {
				return nil, err
			}
			b, err := tengo.ArgBytes(args, 0, "data")
			if err != nil {
				return nil, err
			}
			sum := sha512.Sum512(b)
			return &tengo.String{Value: hex.EncodeToString(sum[:])}, nil
		},
	},

	"hmac_sha256": &tengo.UserFunction{
		Name: "hmac_sha256",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 2); err != nil {
				return nil, err
			}
			key, err := tengo.ArgBytes(args, 0, "key")
			if err != nil {
				return nil, err
			}
			data, err := tengo.ArgBytes(args, 1, "data")
			if err != nil {
				return nil, err
			}
			mac := hmac.New(sha256.New, key)
			mac.Write(data)
			return &tengo.Bytes{Value: mac.Sum(nil)}, nil
		},
	},

	"hmac_sha256_hex": &tengo.UserFunction{
		Name: "hmac_sha256_hex",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 2); err != nil {
				return nil, err
			}
			key, err := tengo.ArgBytes(args, 0, "key")
			if err != nil {
				return nil, err
			}
			data, err := tengo.ArgBytes(args, 1, "data")
			if err != nil {
				return nil, err
			}
			mac := hmac.New(sha256.New, key)
			mac.Write(data)
			return &tengo.String{Value: hex.EncodeToString(mac.Sum(nil))}, nil
		},
	},

	// aes_encrypt uses AES-GCM. key must be 16, 24, or 32 bytes (AES-128/192/256).
	// Returns nonce+ciphertext as a single byte slice.
	"aes_encrypt": &tengo.UserFunction{
		Name: "aes_encrypt",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 2); err != nil {
				return nil, err
			}
			key, err := tengo.ArgBytes(args, 0, "key")
			if err != nil {
				return nil, err
			}
			plaintext, err := tengo.ArgBytes(args, 1, "plaintext")
			if err != nil {
				return nil, err
			}
			block, e := aes.NewCipher(key)
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			gcm, e := gocipher.NewGCM(block)
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			nonce := make([]byte, gcm.NonceSize())
			if _, e = io.ReadFull(rand.Reader, nonce); e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			ct := gcm.Seal(nonce, nonce, plaintext, nil)
			return &tengo.Bytes{Value: ct}, nil
		},
	},

	// aes_decrypt expects the nonce+ciphertext format produced by aes_encrypt.
	"aes_decrypt": &tengo.UserFunction{
		Name: "aes_decrypt",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 2); err != nil {
				return nil, err
			}
			key, err := tengo.ArgBytes(args, 0, "key")
			if err != nil {
				return nil, err
			}
			ciphertext, err := tengo.ArgBytes(args, 1, "ciphertext")
			if err != nil {
				return nil, err
			}
			block, e := aes.NewCipher(key)
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			gcm, e := gocipher.NewGCM(block)
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			ns := gcm.NonceSize()
			if len(ciphertext) < ns {
				return &tengo.Error{Value: &tengo.String{Value: errors.New("ciphertext too short").Error()}}, nil
			}
			nonce, ct := ciphertext[:ns], ciphertext[ns:]
			pt, e := gcm.Open(nil, nonce, ct, nil)
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			return &tengo.Bytes{Value: pt}, nil
		},
	},
}
