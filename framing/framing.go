package framing

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"hash/crc32"
	"io"
)

func join(first []byte, other ...[]byte) []byte {
	for _, o := range other {
		first = append(first, o...)
	}
	return first
}

const (
	flagByte   byte = 0x7E
	escapeByte byte = 0x7d
	xorPattern byte = 0x20
)

func escape(data []byte) []byte {
	length := len(data)

	length += bytes.Count(data, []byte{flagByte})
	length += bytes.Count(data, []byte{escapeByte})

	result := make([]byte, length)
	pos := 0

	for _, b := range data {
		switch b {
		case flagByte, escapeByte:
			result[pos] = escapeByte
			pos++
			result[pos] = b ^ xorPattern
			pos++
		default:
			result[pos] = b
			pos++
		}
	}

	return result
}

func Encode(frame []byte) []byte {
	hash32 := crc32.New(crc32.IEEETable)
	_, err := hash32.Write(frame)
	if err != nil {
		panic(err)
	}
	sum := hash32.Sum(nil)

	// TODO: is it better to checksum after escaping the body, or before?

	body := join(frame, sum)
	escaped := escape(body)

	return join(
		[]byte{flagByte},
		escaped,
		[]byte{flagByte},
	)
}

func split() bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, splitErr error) {
		if len(data) == 0 && atEOF {
			return 0, nil, io.EOF
		}

		if data[0] != flagByte {
			// advance to first flag
			idx0 := bytes.IndexByte(data[1:], flagByte)
			if idx0 == -1 {
				return len(data), nil, nil
			} else {
				return idx0 - 1, nil, nil
			}
		}

		idx1 := bytes.IndexByte(data[1:], flagByte)
		if idx1 == -1 {
			return 0, nil, nil
		}
		idx1++ // account for starting at 1
		body := data[1:idx1]
		var err error
		body, err = unescape(body)
		if err != nil {
			return 1, nil, nil
		}

		if len(body) < 4 {
			return 1, nil, nil
		}

		payload := body[0 : len(body)-4]
		checksum := body[len(body)-4:]

		hash32 := crc32.New(crc32.IEEETable)
		_, _ = hash32.Write(payload)
		digest := hash32.Sum(nil)

		if !bytes.Equal(digest, checksum) {
			return 1, nil, nil
		}

		return idx1 + 1, payload, nil
	}
}

func unescape(body []byte) ([]byte, error) {
	result := make([]byte, 0, len(body))

	escaped := false

	for _, b := range body {
		if escaped {
			switch b {
			case flagByte ^ xorPattern:
				result = append(result, flagByte)
			case escapeByte ^ xorPattern:
				result = append(result, escapeByte)
			default:
				return nil, errors.New("unknown escaped byte")
			}
			escaped = false
		} else {
			switch b {
			case escapeByte:
				escaped = true
			default:
				result = append(result, b)
			}
		}
	}

	if escaped {
		return nil, errors.New("escape in last position")
	}

	return result, nil
}

func Decode(reader io.Reader, callback func([]byte)) error {
	// TODO: limit frame size, and enforce it

	scanner := bufio.NewScanner(reader)
	scanner.Split(split())

	for scanner.Scan() {
		token := scanner.Bytes()
		frame := make([]byte, len(token))
		copy(frame, token)
		callback(frame)
	}

	return scanner.Err()
}
