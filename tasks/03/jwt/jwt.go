package jwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash"
	"strings"
	"time"
)

type SignMethod string

const (
	HS256 SignMethod = "HS256"
	HS512 SignMethod = "HS512"
)

var (
	ErrInvalidSignMethod      = errors.New("invalid sign method")
	ErrSignatureInvalid       = errors.New("signature invalid")
	ErrTokenExpired           = errors.New("token expired")
	ErrSignMethodMismatched   = errors.New("sign method mismatched")
	ErrConfigurationMalformed = errors.New("configuration malformed")
	ErrInvalidToken           = errors.New("invalid token")
)

type header struct {
	Alg string `json:"alg"`
	Type string `json:"typ"`
}

type payload struct {
	Data interface{} `json:"d"`
	Exp *int64 `json:"exp,omitempty"`
}



func Encode(data interface{}, opts ...Option) ([]byte, error) {
	var c config

	for _, opt := range opts {
		opt(&c)
	}

	if c.TTL != nil && c.Expires != nil {
		return nil, ErrConfigurationMalformed
	}

	if c.Expires != nil && c.Expires.Before(timeFunc()) {
		return nil, ErrConfigurationMalformed
	}
	headjs, _ := json.Marshal(header{
		Alg: string(c.SignMethod),
		Type: "JWT",
	})
	headerPld := base64.RawURLEncoding.EncodeToString(headjs)

	// Body
	var exp *int64

	if c.Expires != nil {
		exp = func(x int64) *int64 { return &x} (c.Expires.Unix())
	}

	if c.TTL != nil {
		exp = func(x int64) *int64 { return &x} (timeFunc().Add(*c.TTL).Unix())
	}

	pldJson, _ := json.Marshal(payload{data, exp})
	pldPld := base64.RawURLEncoding.EncodeToString(pldJson)

	var b bytes.Buffer
	b.WriteString(headerPld + "." + pldPld)

	ssum, err := getSum(b.Bytes(), &c)
	if err != nil {
		return nil, err
	}
	b.WriteString(".")
	b.Write(ssum)
	return b.Bytes(), nil
}

func Decode(token []byte, data interface{}, opts ...Option) error {
	var c config
	for _, opt := range opts {
		opt(&c)
	}

	prts := strings.Split(string(token), ".")
	if len(prts) != 3 {
		return ErrInvalidToken
	}

	var hdr header
	if err := decodeStr(prts[0], &hdr); err != nil {
		return err
	}

	pld := payload{
		data,
		nil,
	}
	if err := decodeStr(prts[1], &pld); err != nil {
		return err
	}

	lastdotidx := bytes.LastIndex(token, []byte("."))
	sum := token[lastdotidx + 1:]
	toverify := token[:lastdotidx]
	expSum, err := getSum(toverify, &c)
	if err != nil {
		return err
	}
	if len(sum) != len(expSum) {
		return ErrSignMethodMismatched
	}
	if !bytes.Equal(sum, expSum) {
		return ErrSignatureInvalid
	}

	if pld.Exp != nil && timeFunc().After(time.Unix(*pld.Exp, 0)) {
		return ErrTokenExpired
	}

	return nil
}

// To mock time in tests
var timeFunc = time.Now

func getSum(b []byte, c *config) ([]byte, error) {
	var method func() hash.Hash

	switch c.SignMethod {
	case HS256:
		method = sha256.New
	case HS512:
		method = sha512.New
	default:
		return nil, ErrInvalidSignMethod
	}
	h := hmac.New(method, []byte(c.Key))
	h.Write(b)
	bs := h.Sum(nil)
	sum := make([]byte, base64.RawStdEncoding.EncodedLen(len(bs)))
	base64.RawURLEncoding.Encode(sum, bs)
	return sum, nil
}

func decodeStr(p string, d interface{}) error {
	pJson, err := base64.RawURLEncoding.DecodeString(p)
	if err != nil {
		return ErrInvalidToken
	}

	err = json.Unmarshal(pJson, &d)
	if err != nil {
		return ErrInvalidToken
	}
	return nil
}