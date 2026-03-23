package lpstream

import (
	"bytes"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	var b bytes.Buffer
	w := NewEncoder(&b)

	data1 := []byte("hello")
	data2 := []byte("world!")

	// write two frames
	if err := w.WriteFrame(data1); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteFrame(data2); err != nil {
		t.Fatal(err)
	}

	// read back
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
