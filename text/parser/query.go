package query

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"time"
)

type requiredArgumentNotSetError struct {
	key string
}

func requiredArgumentNotSet(key string) *requiredArgumentNotSetError {
	return &requiredArgumentNotSetError{key}
}

func (err *requiredArgumentNotSetError) Error() string {
	return "required argument " + err.key + " not set"
}

type parseError struct {
	key string
	typ string
	err error
}

func parseIntError(key string, err error) *parseError {
	return &parseError{
		key: key,
		typ: "int",
		err: err,
	}
}

func parseUintError(key string, err error) *parseError {
	return &parseError{
		key: key,
		typ: "uint",
		err: err,
	}
}

func parseFloatError(key string, err error) *parseError {
	return &parseError{
		key: key,
		typ: "float",
		err: err,
	}
}

func parseBoolError(key string, err error) *parseError {
	return &parseError{
		key: key,
		typ: "bool",
		err: err,
	}
}

func (err *parseError) Error() string {
	return "parse " + err.key + " failed: " + err.err.Error()
}

func (err *parseError) Unwrap() error {
	return err.err
}

// Maybe rawurl is of the form scheme:path.
// (Scheme must be [a-zA-Z][a-zA-Z0-9+-.]*)
// If so, return scheme, path; else return "", rawurl.
func getscheme(rawurl string) (scheme, path string, err error) {
	for i := 0; i < len(rawurl); i++ {
		c := rawurl[i]
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
		// do nothing
		case '0' <= c && c <= '9' || c == '+' || c == '-' || c == '.':
			if i == 0 {
				return "", rawurl, nil
			}
		case c == ':':
			if i == 0 {
				return "", "", errors.New("missing protocol scheme")
			}
			return rawurl[:i], rawurl[i+1:], nil
		default:
			// we have encountered an invalid character,
			// so there is no valid scheme
			return "", rawurl, nil
		}
	}
	return "", rawurl, nil
}

// ParseURL parses rawurl as an url.URL with defaultScheme
// if scheme(like xxx://) not found in rawurl
func ParseURL(rawurl string, defaultScheme string) (*url.URL, error) {
	scheme, rest, err := getscheme(rawurl)
	if err != nil || len(scheme) == 0 || len(rest) < 2 || rest[0] != '/' || rest[1] != '/' {
		if defaultScheme != "" {
			rawurl = defaultScheme + "://" + rawurl
		} else {
			return nil, errors.New("url shoud be start with <scheme>://, e.g. http://, tcp://, unix://")
		}
	}
	return url.Parse(rawurl)
}

// Query alias map[string][]string
type Query = map[string][]string

func getArgument(q Query, key string, required bool) (value string, err error) {
	if vs := q[key]; len(vs) > 0 {
		value = vs[0]
	} else if required {
		err = requiredArgumentNotSet(key)
	}
	return
}

func parseInt64(q Query, key string, required bool, dft int64) (int64, error) {
	if value, err := getArgument(q, key, required); err != nil {
		return dft, err
	} else {
		if value == "" {
			return dft, nil
		}
		x, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return dft, parseIntError(key, err)
		}
		return x, nil
	}
}

func parseUint64(q Query, key string, required bool, dft uint64) (uint64, error) {
	if value, err := getArgument(q, key, required); err != nil {
		return dft, err
	} else {
		if value == "" {
			return dft, nil
		}
		x, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return dft, parseIntError(key, err)
		}
		return x, nil
	}
}

func parseFloat64(q Query, key string, required bool, dft float64) (float64, error) {
	if value, err := getArgument(q, key, required); err != nil {
		return dft, err
	} else {
		if value == "" {
			return dft, nil
		}
		x, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return dft, parseFloatError(key, err)
		}
		return x, nil
	}
}

func Int(q Query, key string, dft int) (int, error) {
	i, err := parseInt64(q, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int(i), nil
}

func RequiredInt(q Query, key string) (int, error) {
	i, err := parseInt64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func Int8(q Query, key string, dft int8) (int8, error) {
	i, err := parseInt64(q, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int8(i), nil
}

func RequiredInt8(q Query, key string) (int8, error) {
	i, err := parseInt64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int8(i), nil
}

func Int16(q Query, key string, dft int16) (int16, error) {
	i, err := parseInt64(q, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int16(i), nil
}

func RequiredInt16(q Query, key string) (int16, error) {
	i, err := parseInt64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int16(i), nil
}

func Int32(q Query, key string, dft int32) (int32, error) {
	i, err := parseInt64(q, key, false, int64(dft))
	if err != nil {
		return dft, err
	}
	return int32(i), nil
}

func RequiredInt32(q Query, key string) (int32, error) {
	i, err := parseInt64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func Int64(q Query, key string, dft int64) (int64, error) {
	return parseInt64(q, key, false, dft)
}

func RequiredInt64(q Query, key string) (int64, error) {
	return parseInt64(q, key, true, 0)
}

func Uint(q Query, key string, dft uint) (uint, error) {
	i, err := parseUint64(q, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint(i), nil
}

func RequiredUint(q Query, key string) (uint, error) {
	i, err := parseUint64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

func Uint8(q Query, key string, dft uint8) (uint8, error) {
	i, err := parseUint64(q, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint8(i), nil
}

func RequiredUint8(q Query, key string) (uint8, error) {
	i, err := parseUint64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint8(i), nil
}

func Uint16(q Query, key string, dft uint16) (uint16, error) {
	i, err := parseUint64(q, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint16(i), nil
}

func RequiredUint16(q Query, key string) (uint16, error) {
	i, err := parseUint64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}

func Uint32(q Query, key string, dft uint32) (uint32, error) {
	i, err := parseUint64(q, key, false, uint64(dft))
	if err != nil {
		return dft, err
	}
	return uint32(i), nil
}

func RequiredUint32(q Query, key string) (uint32, error) {
	i, err := parseUint64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

func Uint64(q Query, key string, dft uint64) (uint64, error) {
	return parseUint64(q, key, false, dft)
}

func RequiredUint64(q Query, key string) (uint64, error) {
	return parseUint64(q, key, true, 0)
}

func Float32(q Query, key string, dft float32) (float32, error) {
	f, err := parseFloat64(q, key, false, float64(dft))
	if err != nil {
		return dft, err
	}
	return float32(f), err
}

func RequiredFloat32(q Query, key string) (float32, error) {
	f, err := parseFloat64(q, key, true, 0)
	if err != nil {
		return 0, err
	}
	return float32(f), err
}

func Float64(q Query, key string, dft float64) (float64, error) {
	return parseFloat64(q, key, false, dft)
}

func RequiredFloat64(q Query, key string) (float64, error) {
	return parseFloat64(q, key, true, 0)
}

func Bool(q Query, key string, dft bool) (bool, error) {
	if value, err := getArgument(q, key, false); err != nil {
		return dft, err
	} else {
		if value == "" {
			return dft, nil
		}
		x, err := strconv.ParseBool(value)
		if err != nil {
			return dft, parseBoolError(key, err)
		}
		return x, nil
	}
}

func RequiredBool(q Query, key string) (bool, error) {
	if value, err := getArgument(q, key, true); err != nil {
		return false, err
	} else {
		x, err := strconv.ParseBool(value)
		if err != nil {
			return false, parseBoolError(key, err)
		}
		return x, nil
	}
}

func String(q Query, key string, dft string) string {
	if value, _ := getArgument(q, key, false); value == "" {
		return dft
	} else {
		return value
	}
}

func RequiredString(q Query, key string) (string, error) {
	if value, err := getArgument(q, key, true); err != nil {
		return "", err
	} else {
		return value, nil
	}
}

func JSON(q Query, key string, ptr any) error {
	if value, _ := getArgument(q, key, false); value == "" {
		return nil
	} else {
		return json.Unmarshal([]byte(value), ptr)
	}
}

func RequiredDuration(q Query, key string) (time.Duration, error) {
	if value, err := getArgument(q, key, true); err != nil {
		return 0, err
	} else {
		return time.ParseDuration(value)
	}
}

func Duration(q Query, key string, dft time.Duration) (time.Duration, error) {
	if value, _ := getArgument(q, key, false); value == "" {
		return dft, nil
	} else {
		return time.ParseDuration(value)
	}
}
