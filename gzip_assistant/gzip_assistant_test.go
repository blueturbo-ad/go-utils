package gzipassistant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecrypt(t *testing.T) {
	t.Run("DecryptEmpty", func(t *testing.T) {
		dst, err := Decrypt([]byte(""))
		assert.NotNil(t, err)
		assert.Nil(t, dst)

	})
}

func TestEncrypt(t *testing.T) {
	t.Run("EncryptEmpty", func(t *testing.T) {
		src, err := Encrypt([]byte(""))
		assert.Nil(t, err)
		assert.NotNil(t, src)
	})
}

func TestDeEncrypt(t *testing.T) {

	t.Run("DeEncryptEmpty", func(t *testing.T) {
		src := "1234567890abcdefghijklmnopqrstuvwxyz!@#$%^&*()_+"
		dst, err := Encrypt([]byte(src))
		assert.Nil(t, err)
		assert.NotNil(t, dst)
		ori, err := Decrypt(dst)
		assert.Nil(t, err)
		assert.NotNil(t, dst)
		assert.Equal(t, src, string(ori))
	})
}
