package framing

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestEncode(t *testing.T) {
	msg := []byte("abc123")
	encode := Encode(msg)
	fmt.Println(hex.Dump(msg))
	fmt.Println(hex.Dump(encode))
}

func TestDecode(t *testing.T) {
	msgs := join(
		Encode([]byte("abc123")),
		Encode([]byte("xyx456")),
		Encode([]byte("a")),
	)

	err := Decode(bytes.NewBuffer(msgs), func(frame []byte) {
		fmt.Println("cb", hex.Dump(frame))
	})

	require.NoError(t, err)
}

func TestCombinations(t *testing.T) {
	r, w := io.Pipe()
	ch := make(chan []byte, 0)

	go func() {
		err := Decode(r, func(frame []byte) {
			go func(frame []byte) {
				ch <- frame
			}(frame)
		})
		if err != nil {
			panic(err)
		}
	}()

	t.Run("empty message", func(t *testing.T) {
		_, err := w.Write(Encode([]byte{}))
		require.NoError(t, err)
		got := <-ch
		assert.Equal(t, []byte{}, got)
	})

	t.Run("single byte messages", func(t *testing.T) {
		for i := 0; i < 255; i++ {
			b := byte(i)
			fmt.Printf(" ==> 0x%02x\n", b)
			msg := []byte{b}
			_, err := w.Write(Encode(msg))
			require.NoError(t, err)
			got := <-ch
			assert.Equal(t, msg, got)
		}
	})

	t.Run("double byte messages", func(t *testing.T) {
		for i := 0; i < 255; i++ {
			for j := 0; j < 255; j++ {
				msg := []byte{byte(i), byte(j)}
				_, err := w.Write(Encode(msg))
				require.NoError(t, err)
				got := <-ch
				assert.Equal(t, msg, got)
			}
		}
	})
}

func TestBadInput(t *testing.T) {
	t.Fail() // need to implement
}
