// Package chisel helps creating chisels: programs that process kafcat's cat output. For example, a chisel can interpret the key/value bytes as Avro.
// TODO: this can and should be so much cleaner.
package chisel

import (
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

var hexBlobSLineRe = regexp.MustCompile("^[0-9a-fA-F]{8}\\s+[0-9a-fA-F\\s]+\\s+\\|.+\\|")
var hexBlobLineRe = regexp.MustCompile("([0-9a-fA-F]{8})\\s+([0-9a-fA-F\\s]+).*")

// ParseHexDump reverses hex.Dump(...). Suitable for <1MB inputs only for performance reasons
func ParseHexDump(s string) ([]byte, error) {
	ss := hexBlobLineRe.FindAllStringSubmatch(s, -1)

	var bs []byte

	for _, l := range ss {
		if len(l) == 3 {
			offset, hexs := l[1], strings.Replace(l[2], " ", "", -1)
			_ = offset // TODO
			subss, err := hex.DecodeString(hexs)
			if err != nil {
				return nil, err
			}
			bs = append(bs, subss...)
		}
	}
	return bs, nil
}

// putInMapSlice sets the key to the value
func putInMapSlice(ms yaml.MapSlice, key interface{}, value interface{}) yaml.MapSlice {
	for i, m := range ms {
		if m.Key == key {
			ms[i].Value = value
		}
	}
	return ms
}

// converts a mapSliceToMap into a (non-ordered) map
func mapSliceToMap(ms yaml.MapSlice) map[interface{}]interface{} {
	ma := map[interface{}]interface{}{}
	for _, m := range ms {
		switch m.Value.(type) {
		case yaml.MapSlice:
			ma[m.Key] = mapSliceToMap(m.Value.(yaml.MapSlice))
		default:
			ma[m.Key] = m.Value
		}
	}
	return ma
}

type SerializeFunc func(interface{}) ([]byte, error)
type DeserializeFunc func([]byte) (interface{}, error)

type scannerState int

const (
	searching scannerState = iota
	swallowKeyHex
	swallowValueHex
)

func parseAsBytes(i interface{}) ([]byte, error) {
	k, ok := i.(string)
	if !ok {
		return nil, fmt.Errorf("bad")
	}
	if hexBlobSLineRe.MatchString(k) { //TODO: reconsider. do we really want this?!
		return ParseHexDump(k)
	}
	return []byte(k), nil
}

func parseAndDeserialize(i interface{}, t DeserializeFunc) (interface{}, error) {
	bs, err := parseAsBytes(i)
	if err != nil {
		return nil, err
	}
	return t(bs)
}

// bytes -> interface{}
func Deserialize(in io.Reader, out io.Writer, deserializeKey DeserializeFunc, deserializeValue DeserializeFunc) error {
	dec := yaml.NewDecoder(in)

	for {
		var msEntry yaml.MapSlice
		if err := dec.Decode(&msEntry); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		entry := mapSliceToMap(msEntry)

		key, err := parseAndDeserialize(entry["key"], deserializeKey)
		if err != nil {
			return err
		}

		value, err := parseAndDeserialize(entry["value"], deserializeValue)
		if err != nil {
			return err
		}

		putInMapSlice(msEntry, "key", key)
		putInMapSlice(msEntry, "value", value)

		ybs, err := yaml.Marshal(msEntry)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, "---")
		fmt.Fprint(out, string(ybs))
	}
}

// interface{} -> []byte
func Serialize(in io.Reader, out io.Writer, serializeKey SerializeFunc, serializeValue SerializeFunc, fmtHex bool) error {
	dec := yaml.NewDecoder(in)

	for {
		var msEntry yaml.MapSlice
		if err := dec.Decode(&msEntry); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		entry := mapSliceToMap(msEntry)

		key, err := serializeKey(entry["key"])
		if err != nil {
			return err
		}

		value, err := serializeValue(entry["value"])
		if err != nil {
			return err
		}

		if fmtHex {
			putInMapSlice(msEntry, "key", hex.Dump(key))
			putInMapSlice(msEntry, "value", hex.Dump(value))
		} else {
			putInMapSlice(msEntry, "key", string(key))
			putInMapSlice(msEntry, "value", string(value))
		}

		ybs, err := yaml.Marshal(msEntry)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, "---")
		fmt.Fprint(out, string(ybs))
	}
}
