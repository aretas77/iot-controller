package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func StripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

func SplitTopic4(topic string) (string, string, string, string) {
	tok := strings.Split(topic, "/")
	return tok[0], tok[1], tok[2], tok[3]
}

// SplitDataReadLine splits a read data line from the sensor data.
func SplitDataReadLine(line string) (error, string, float64, float64, float64) {
	tok := strings.Split(line, ";")

	var err error
	var temperature, pressure, meters float64
	if temperature, err = strconv.ParseFloat(tok[1], 32); err != nil {
		return errors.New("failed to parse temperature"), tok[0], 0, 0, 0
	}

	if pressure, err = strconv.ParseFloat(tok[2], 32); err != nil {
		return errors.New("failed to parse pressure"), tok[0], 0, 0, 0
	}

	if meters, err = strconv.ParseFloat(tok[3], 32); err != nil {
		return errors.New("failed to parse meters"), tok[0], 0, 0, 0
	}

	return nil, tok[0], temperature, pressure, meters
}

// HashAndSalt will take a password as an input and extract its hash
func HashAndSalt(password string) string {
	bytePsw := []byte(password)

	// Use GenerateFromPassword to hash and salt password.
	hash, err := bcrypt.GenerateFromPassword(bytePsw, bcrypt.MinCost)
	if err != nil {
		logrus.Error(err)
	}

	// GenerateFromPassword returns a byte slice so we convert the bytes to a
	// string and return the hash.
	return string(hash)
}

func ComparePasswords(hashedPsw string, plainPsw string) bool {
	byteHash := []byte(hashedPsw)
	bytePsw := []byte(plainPsw)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePsw)
	if err != nil {
		logrus.Error(err)
		return false
	}

	return true
}
