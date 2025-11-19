package protocol

import (
	"bytes"
	"testing"
)

func Test_EncryptAndDecrypt(t *testing.T) {
	inputData := []byte("Lorem Ipsum is simply dummy text of the printing and typesetting industry.")

	encryptedData, err := EncryptRSA(&keyForClientCommunication.PublicKey, inputData)
	if err != nil {
		t.Fatal(err)
	}
	decryptedData, err := DecryptRSA(encryptedData)
	if err != nil {
		t.Fatal(err)
	}
	if string(decryptedData) != string(inputData) {
		t.Fatalf("Decrypted data does not match original. Got '%s', expected '%s'", string(decryptedData), string(inputData))
	}
}

func TestDecryptRSA_LeadingZero(t *testing.T) {
	// Arrange: Create a plaintext that is guaranteed to start with a zero byte.
	// This mimics our real protocol's check byte.
	originalPlaintext := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}

	t.Logf("Original plaintext (len=%d): %x", len(originalPlaintext), originalPlaintext)

	// Act: Encrypt the data. We assume EncryptRSA is correct.
	ciphertext, err := EncryptRSA(&keyForClientCommunication.PublicKey, originalPlaintext)
	if err != nil {
		t.Fatalf("Encryption failed unexpectedly: %v", err)
	}

	// Decrypt the data using the function we are testing.
	decryptedBlock, err := DecryptRSA(ciphertext)
	if err != nil {
		t.Fatalf("Decryption failed unexpectedly: %v", err)
	}

	t.Logf("Decrypted block (len=%d):   %x", len(decryptedBlock), decryptedBlock)

	// Assert: Check if the decrypted block has the expected length.
	// The key size is the expected length for raw RSA plaintext.
	expectedLength := keyForClientCommunication.Size()
	if len(decryptedBlock) != expectedLength {
		// This assertion will fail if DecryptRSA uses m.Bytes()
		t.Errorf("Decrypted block has incorrect length. Got %d, want %d", len(decryptedBlock), expectedLength)
		t.Log("This failure indicates the decryption function is likely stripping leading zero bytes.")
	}

	// As a secondary check, verify the content.
	// We need to compare the end of the decrypted block with our original plaintext.
	if !bytes.HasSuffix(decryptedBlock, originalPlaintext) {
		t.Errorf("Decrypted block does not contain the original plaintext at the end.")
	}
}
