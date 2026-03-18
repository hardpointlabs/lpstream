package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hardpointlabs/lpstream"
)

func generateSecureAlphanumeric(length int) ([]byte, error) {
	// Define the character set to use for the random string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a byte slice to hold the random characters
	b := make([]byte, length)

	// Use crypto/rand to read random bytes
	// The number of bytes needed is slightly more than the desired length to ensure enough randomness after filtering
	randomBytes := make([]byte, length+100) // Read extra bytes to be safe
	if _, err := rand.Read(randomBytes); err != nil {
		return make([]byte, 0), err // Return error if reading fails
	}

	// Populate the byte slice with characters from the charset
	for i := 0; i < length; i++ {
		// Use modulo to select a character from the charset using a random byte
		// This approach is fast and avoids bias when using enough random source bytes
		b[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return b, nil
}

func main() {
	log.Println("Starting client")
	conn, err := net.DialTimeout("tcp", "test-server:8124", 10*time.Second)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		return
	}
	defer conn.Close()

	msg, err := generateSecureAlphanumeric(1 << 16)
	if err != nil {
		log.Fatal(err)
	}
	writer := lpstream.NewWriter(conn)
	reader := lpstream.NewReader(conn)

	err = writer.WriteFrame(msg)
	if err != nil {
		log.Fatal(err)
	}

	read, err := reader.ReadFrame()
	if err != nil {
		log.Fatal(err)
	}

	if !bytes.Equal(msg, read) {
		log.Fatal("Written & received messages are not equal")
	}

	log.Println("Received OK!")
}
