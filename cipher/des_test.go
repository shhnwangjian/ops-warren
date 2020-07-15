package cipher

import (
	"encoding/base64"
	"testing"
)

func TestDesDecrypt(t *testing.T) {
	val, _ := base64.StdEncoding.DecodeString("VuU0fC1+C+UVAx83vWczd4xMjaEcQqiY")
	res, err := DesDecrypt(val, []byte("sfe023f_"))
	if err != nil {
		t.Fatalf("test DesDecrypt failed. %s", err)
	}
	t.Logf("test DesDecrypt succ, %s", string(res))
}

func TestDesEncrypt(t *testing.T) {
	res, err := DesEncrypt([]byte("7GDWY4SUE05LUQXV7R44"), []byte("sfe023f_"))
	if err != nil {
		t.Fatalf("test DesEncrypt failed. %s", err)
	}
	t.Logf("test DesEncrypt succ, %s", base64.StdEncoding.EncodeToString(res))
}

func TestTripleDesDecrypt(t *testing.T) {
	key := []byte("sfe023f_sefiel#fi32lf3e!")
	val, _ := base64.StdEncoding.DecodeString("RmSetwALzx3hQUUO/ZwjqJEaMyxqRy0T")
	res, err := TripleDesDecrypt(val, key)
	if err != nil {
		t.Fatalf("test TripleDesDecrypt failed. %s", err)
	}
	t.Logf("test TripleDesDecrypt succ, %s", string(res))
}

func TestTripleDesEncrypt(t *testing.T) {
	key := []byte("sfe023f_sefiel#fi32lf3e!")
	res, err := TripleDesEncrypt([]byte("7GDWY4SUE05LUQXV7R44"), key)
	if err != nil {
		t.Fatalf("test TripleDesEncrypt failed. %s", err)
	}
	t.Logf("test TripleDesEncrypt succ, %s", base64.StdEncoding.EncodeToString(res))
}
