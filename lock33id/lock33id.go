package lock33id

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	cosineEpsilon     = 1e-12
	identityTolerance = 1e-9
)

func IdentityWitness(x float64) (float64, error) {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0, errors.New("x must be a finite number")
	}

	cosine := math.Cos(x)
	if math.Abs(cosine) < cosineEpsilon {
		return 0, errors.New("cos(x) ≈ 0; identity undefined")
	}

	lhs := 1/cosine - cosine
	rhs := math.Sin(x) * math.Tan(x)
	if math.Abs(lhs-rhs) > identityTolerance {
		return 0, errors.New("identity failed")
	}
	return rhs, nil
}

func DeriveModulus(x float64) (int, error) {
	witness, err := IdentityWitness(x)
	if err != nil {
		return 0, err
	}
	modulus := int64(math.Trunc(math.Abs(witness)*1000))%251 + 2
	if modulus < 2 {
		modulus = 2
	}
	return int(modulus), nil
}

func DeriveShift(x float64) (int, error) {
	modulus, err := DeriveModulus(x)
	if err != nil {
		return 0, err
	}
	return modulus % 26, nil
}

func DeriveSeed(x float64) (int, error) {
	witness, err := IdentityWitness(x)
	if err != nil {
		return 0, err
	}
	scaledWitness := int64(math.Trunc(math.Abs(witness * 1_000_000)))
	return int(scaledWitness % 10_000), nil
}

func EncodeLock33Pipe(text string, x float64) (string, error) {
	if !utf8.ValidString(text) {
		return "", errors.New("text must be valid UTF-8")
	}
	seed, err := DeriveSeed(x)
	if err != nil {
		return "", err
	}

	shift := int64(seed%33 + 1)
	var encoded strings.Builder
	first := true
	for _, character := range text {
		if !first {
			encoded.WriteByte('|')
		}
		first = false
		encoded.WriteString(strconv.FormatInt(int64(character)+shift, 16))
	}
	return encoded.String(), nil
}

func DecodeLock33Pipe(payload string, x float64) (string, error) {
	if payload == "" {
		return "", nil
	}
	seed, err := DeriveSeed(x)
	if err != nil {
		return "", err
	}

	shift := int64(seed%33 + 1)
	var decoded strings.Builder
	for _, token := range strings.Split(payload, "|") {
		if token == "" {
			return "", errors.New("invalid lock33-pipe payload")
		}
		for _, character := range token {
			if !isHexDigit(character) {
				return "", errors.New("invalid lock33-pipe payload")
			}
		}

		encodedCodePoint, err := strconv.ParseUint(token, 16, 32)
		if err != nil {
			return "", errors.New("invalid lock33-pipe payload")
		}
		codePoint := int64(encodedCodePoint) - shift
		if codePoint < 0 || codePoint > utf8.MaxRune || !utf8.ValidRune(rune(codePoint)) {
			return "", errors.New("invalid lock33-pipe code point")
		}
		decoded.WriteRune(rune(codePoint))
	}
	return decoded.String(), nil
}

func isHexDigit(character rune) bool {
	return character >= '0' && character <= '9' ||
		character >= 'a' && character <= 'f' ||
		character >= 'A' && character <= 'F'
}
