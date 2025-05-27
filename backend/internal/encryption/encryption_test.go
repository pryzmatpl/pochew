package encryption

import (
	"testing"
)

func TestEncryptionService(t *testing.T) {
	service := NewService(1000) // Use fewer iterations for testing

	t.Run("GenerateUserKey", func(t *testing.T) {
		key, err := service.GenerateUserKey()
		if err != nil {
			t.Fatalf("Failed to generate user key: %v", err)
		}
		
		if len(key) == 0 {
			t.Error("Generated key is empty")
		}
		
		// Generate another key and ensure they're different
		key2, err := service.GenerateUserKey()
		if err != nil {
			t.Fatalf("Failed to generate second user key: %v", err)
		}
		
		if key == key2 {
			t.Error("Generated keys should be unique")
		}
	})

	t.Run("EncryptDecrypt", func(t *testing.T) {
		userKey, err := service.GenerateUserKey()
		if err != nil {
			t.Fatalf("Failed to generate user key: %v", err)
		}

		testData := "This is a test message for encryption"
		
		// Encrypt
		encrypted, err := service.Encrypt(testData, userKey)
		if err != nil {
			t.Fatalf("Failed to encrypt data: %v", err)
		}
		
		if encrypted == testData {
			t.Error("Encrypted data should not match original data")
		}
		
		// Decrypt
		decrypted, err := service.Decrypt(encrypted, userKey)
		if err != nil {
			t.Fatalf("Failed to decrypt data: %v", err)
		}
		
		if decrypted != testData {
			t.Errorf("Decrypted data doesn't match original. Expected: %s, Got: %s", testData, decrypted)
		}
	})

	t.Run("EncryptDecryptFile", func(t *testing.T) {
		userKey, err := service.GenerateUserKey()
		if err != nil {
			t.Fatalf("Failed to generate user key: %v", err)
		}

		testContent := []byte("This is test file content with some binary data: \x00\x01\x02\x03")
		
		// Encrypt file
		encrypted, err := service.EncryptFile(testContent, userKey)
		if err != nil {
			t.Fatalf("Failed to encrypt file: %v", err)
		}
		
		// Decrypt file
		decrypted, err := service.DecryptFile(encrypted, userKey)
		if err != nil {
			t.Fatalf("Failed to decrypt file: %v", err)
		}
		
		if string(decrypted) != string(testContent) {
			t.Errorf("Decrypted file content doesn't match original")
		}
	})

	t.Run("DecryptWithWrongKey", func(t *testing.T) {
		userKey1, err := service.GenerateUserKey()
		if err != nil {
			t.Fatalf("Failed to generate user key 1: %v", err)
		}
		
		userKey2, err := service.GenerateUserKey()
		if err != nil {
			t.Fatalf("Failed to generate user key 2: %v", err)
		}

		testData := "This is a test message"
		
		// Encrypt with key 1
		encrypted, err := service.Encrypt(testData, userKey1)
		if err != nil {
			t.Fatalf("Failed to encrypt data: %v", err)
		}
		
		// Try to decrypt with key 2 (should fail)
		_, err = service.Decrypt(encrypted, userKey2)
		if err == nil {
			t.Error("Decryption should fail with wrong key")
		}
	})

	t.Run("EmptyData", func(t *testing.T) {
		userKey, err := service.GenerateUserKey()
		if err != nil {
			t.Fatalf("Failed to generate user key: %v", err)
		}

		// Test empty string
		encrypted, err := service.Encrypt("", userKey)
		if err != nil {
			t.Fatalf("Failed to encrypt empty data: %v", err)
		}
		
		decrypted, err := service.Decrypt(encrypted, userKey)
		if err != nil {
			t.Fatalf("Failed to decrypt empty data: %v", err)
		}
		
		if decrypted != "" {
			t.Errorf("Decrypted empty data should be empty, got: %s", decrypted)
		}
	})

	t.Run("LargeData", func(t *testing.T) {
		userKey, err := service.GenerateUserKey()
		if err != nil {
			t.Fatalf("Failed to generate user key: %v", err)
		}

		// Create large test data (1MB)
		largeData := make([]byte, 1024*1024)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}
		
		testData := string(largeData)
		
		// Encrypt
		encrypted, err := service.Encrypt(testData, userKey)
		if err != nil {
			t.Fatalf("Failed to encrypt large data: %v", err)
		}
		
		// Decrypt
		decrypted, err := service.Decrypt(encrypted, userKey)
		if err != nil {
			t.Fatalf("Failed to decrypt large data: %v", err)
		}
		
		if decrypted != testData {
			t.Error("Decrypted large data doesn't match original")
		}
	})
}

func BenchmarkEncryption(b *testing.B) {
	service := NewService(10000)
	userKey, _ := service.GenerateUserKey()
	testData := "This is a benchmark test message for encryption performance"

	b.Run("Encrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := service.Encrypt(testData, userKey)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	encrypted, _ := service.Encrypt(testData, userKey)
	
	b.Run("Decrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := service.Decrypt(encrypted, userKey)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
} 