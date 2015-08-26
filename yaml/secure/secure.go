package secure

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"

	"github.com/square/go-jose"
	"gopkg.in/yaml.v2"
)

// Parse parses and returns the secure section of the
// yaml file as plaintext parameters.
func Parse(raw, privKey string) (map[string]string, error) {
	secrets, err := parseSecure(raw)
	if err != nil {
		return nil, err
	}

	// unarmshal the private key from PEM
	rsaPrivKey, err := decodePrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	// decrypt each item in the secure section,
	// extract the named parameter, and return
	// as a map
	params := map[string]string{}
	for _, secret := range secrets {
		plain, err := decrypt(secret, rsaPrivKey)
		if err != nil {
			return nil, err
		}
		key, value := parseVariable(plain)
		params[key] = value
	}

	return params, nil
}

// encrypt encrypts a plaintext variable using JOSE with
// RSA_OAEP and A128GCM algorithms.
func encrypt(text string, pubKey *rsa.PublicKey) (string, error) {
	var encrypted string
	var plaintext = []byte(text)

	// Creates a new encrypter using defaults
	encrypter, err := jose.NewEncrypter(jose.RSA_OAEP, jose.A128GCM, pubKey)
	if err != nil {
		return encrypted, err
	}
	// Encrypts the plaintext value and serializes
	// as a JOSE string.
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		return encrypted, err
	}
	return object.CompactSerialize()
}

// decrypt decrypts a JOSE string and returns the
// plaintext value.
func decrypt(secret string, privKey *rsa.PrivateKey) (string, error) {
	var plaintext string

	// parses the encrypted JSON JOSE string
	object, err := jose.ParseEncrypted(secret)
	if err != nil {
		return plaintext, err
	}

	// decrypts the JOSE object
	decrypted, err := object.Decrypt(privKey)
	if err != nil {
		return plaintext, err
	}
	plaintext = string(decrypted)

	return plaintext, nil
}

// parseSecure is helper function to parse the variable
// declaration in KEY=VALUE format, and return the
// individual parts.
func parseVariable(param string) (string, string) {
	var key string
	var val string
	parts := strings.SplitN(param, "=", 2)
	if len(parts) == 2 {
		key = parts[0]
		val = parts[1]
	}
	return key, val
}

// parseSecure is helper function to parse the Secure data from
// the raw yaml file.
func parseSecure(raw string) ([]string, error) {
	data := struct {
		Secure []string
	}{}
	err := yaml.Unmarshal([]byte(raw), &data)

	return data.Secure, err
}

// decodePrivateKey is a helper function that unmarshals a PEM
// bytes to an RSA Private Key
func decodePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	derBlock, _ := pem.Decode([]byte(privateKey))
	return x509.ParsePKCS1PrivateKey(derBlock.Bytes)
}

// encodePrivateKey is a helper function that marshals an RSA
// Private Key to a PEM encoded file.
func encodePrivateKey(privkey *rsa.PrivateKey) string {
	privateKeyMarshaled := x509.MarshalPKCS1PrivateKey(privkey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Headers: nil, Bytes: privateKeyMarshaled})
	return string(privateKeyPEM)
}
