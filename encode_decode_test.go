package lpstream

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	var b bytes.Buffer
	w := NewEncoder(&b)

	data1 := []byte("hello")
	data2 := []byte("world!")

	if err := w.WriteFrame(data1); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteFrame(data2); err != nil {
		t.Fatal(err)
	}

	r := NewDecoder(&b)
	got1, err := r.ReadFrame()
	if err != nil {
		t.Fatal(err)
	}
	if string(got1) != string(data1) {
		t.Errorf("expected %q, got %q", data1, got1)
	}

	got2, err := r.ReadFrame()
	if err != nil {
		t.Fatal(err)
	}
	if string(got2) != string(data2) {
		t.Errorf("expected %q, got %q", data2, got2)
	}
}

func TestECDHKeyExchange(t *testing.T) {
	prv1, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	prv2, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	pub1 := prv1.PublicKey()
	pub2 := prv2.PublicKey()

	var bufSide1 bytes.Buffer
	codecSide1 := NewFrameCodec(&bufSide1)

	if err := codecSide1.WriteFrame(pub1.Bytes()); err != nil {
		t.Fatal(err)
	}
	if err := codecSide1.WriteFrame(pub2.Bytes()); err != nil {
		t.Fatal(err)
	}

	var bufSide2 bytes.Buffer
	codecSide2 := NewFrameCodec(&bufSide2)

	if err := codecSide2.WriteFrame(pub2.Bytes()); err != nil {
		t.Fatal(err)
	}
	if err := codecSide2.WriteFrame(pub1.Bytes()); err != nil {
		t.Fatal(err)
	}

	decoderSide1 := NewDecoder(&bufSide2)
	_, err = decoderSide1.ReadFrame()
	if err != nil {
		t.Fatal(err)
	}
	_, err = decoderSide1.ReadFrame()
	if err != nil {
		t.Fatal(err)
	}

	decoderSide2 := NewDecoder(&bufSide1)
	side2SentPub1, err := decoderSide2.ReadFrame()
	if err != nil {
		t.Fatal(err)
	}
	side2SentPub2, err := decoderSide2.ReadFrame()
	if err != nil {
		t.Fatal(err)
	}

	secret1, err := ecdh256(prv1, side2SentPub2)
	if err != nil {
		t.Fatal(err)
	}
	secret2, err := ecdh256(prv2, side2SentPub1)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(secret1, secret2) {
		t.Fatalf("shared secrets do not match: %x vs %x", secret1, secret2)
	}

	plaintext := []byte("hello world")
	ciphertext, err := encryptAESGCM(secret1, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	var bufCT bytes.Buffer
	codecCT := NewFrameCodec(&bufCT)
	if err := codecCT.WriteFrame(ciphertext); err != nil {
		t.Fatal(err)
	}

	decoderCT := NewDecoder(&bufCT)
	decryptedFrame, err := decoderCT.ReadFrame()
	if err != nil {
		t.Fatal(err)
	}

	plaintextDecrypted, err := decryptAESGCM(secret2, decryptedFrame)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(plaintext, plaintextDecrypted) {
		t.Errorf("expected %q, got %q", plaintext, plaintextDecrypted)
	}
}