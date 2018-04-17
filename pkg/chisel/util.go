// Package chisel helps creating chisels: programs that process kafcat's cat output. For example, a chisel can interpret the key/value bytes as Avro.
// TODO: this can and should be so much cleaner.
package chisel

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strings"
)

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

func formatKV(out io.Writer, swallow []string, tr ByteSliceTransformer) error {
	if len(swallow) <= 1 {
		return fmt.Errorf("too short")
	}
	bs, err := ParseHexDump(strings.Join(swallow[1:], "\n"))
	if err != nil {
		return err
	}
	return tr(bs, out)
}

func mustFormatKV(out io.Writer, swallow []string, tr ByteSliceTransformer) {
	err := formatKV(out, swallow, tr)
	if err != nil {
		fmt.Printf("| # error: %s\n", err)
		if len(swallow) == 1 {
			return
		}
		fmt.Printf(strings.Join(swallow[1:], "\n") + "\n")
		return
	}
}

type ByteSliceTransformer func([]byte, io.Writer) error

type scannerState int

const (
	searching scannerState = iota
	swallowKeyHex
	swallowValueHex
)

func Transform(in io.Reader, out io.Writer, key ByteSliceTransformer, value ByteSliceTransformer) error {
	scn := bufio.NewScanner(in)

	st := searching
	var swallow []string
	for scn.Scan() {
		t := scn.Text()
	again:
		switch {
		case st == searching && t == "key: |":
			st = swallowKeyHex
			swallow = []string{t}

		case st == searching && t == "value: |":
			st = swallowValueHex
			swallow = []string{t}

		// searching && not key / value? directly dump
		case st == searching:
			fmt.Fprintln(out, t)

		// not searching, swallowing a hex blob and input looks hexish? extend
		case hexBlobLineRe.MatchString(t):
			swallow = append(swallow, t)

		// not a string and hexishness ended? process
		default:
			if st == swallowKeyHex {
				fmt.Fprintf(out, "key: ")
				mustFormatKV(out, swallow, key)
			} else {
				fmt.Fprintf(out, "value: ")
				mustFormatKV(out, swallow, value)
			}
			st = searching
			goto again
		}
	}
	if st != searching && len(swallow) > 0 {
		if st == swallowKeyHex {
			mustFormatKV(out, swallow, key)
		} else {
			mustFormatKV(out, swallow, value)
		}
	}
	return scn.Err()
}
