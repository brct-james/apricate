package auth

import (
	"crypto/rand"
	"math/big"
	"os"

	"apricate/filemngr"
	"apricate/log"

	"github.com/joho/godotenv"
)

// Creates or updates the APRICATE_ACCESS_SECRET value in secrets.env
func CreateOrUpdateAuthSecretInFile() {
	// Ensure exists
	filemngr.Touch("data/secrets.env")
	// Read file to lines array splitting by newline
	lines, readErr := filemngr.ReadFileToLineSlice("data/secrets.env")
	if readErr != nil {
		// Auth is mission-critical, using Fatal
		log.Error.Fatalf("Could not read lines from secrets.env. Err: %v", readErr)
	}

	// Securely generate new 64 character secret
	newSecret, generationErr := GenerateRandomSecureString(64)
	if generationErr != nil {
		log.Error.Fatalf("Could not generate secure string: %v", generationErr)
	}
	secretString :=  "APRICATE_ACCESS_SECRET=" + string(newSecret)
	log.Debug.Printf("New Secret Generated: %s", secretString)
	
	// Search existing file for secret identifier
	found, i := filemngr.KeyInSliceOfLines("APRICATE_ACCESS_SECRET=", lines)
	if found {
		// Update existing secret
		lines [i] = secretString
	} else {
		// Create secret in env file since could not find one to update
		// If empty file then replace 1st line else append to end
		log.Debug.Printf("Creating new secret in env file. secrets.env[0] == ''? %v", lines[0] == "")
		if lines[0] == "" {
			log.Debug.Printf("Blank secrets.env, replacing line 0")
			lines[0] = secretString
		} else {
			log.Debug.Printf("Not blank secrets.env, appending to end")
			lines = append(lines, secretString)
		}
	}
	
	// Join and write out
	writeErr := filemngr.WriteLinesToFile("data/secrets.env", lines)
	if writeErr != nil {
		log.Error.Fatalf("Could not write secrets.env: %v", writeErr)
	}
	log.Info.Println("Wrote auth secret to secrets.env")
}

// Generate random string of n characters
func GenerateRandomSecureString(n int) (string, error) {
	const allowed = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowed))))
		if err != nil {
			return "", err
		}
		ret[i] = allowed[num.Int64()]
	}
	return string(ret), nil
}


// Load secrets.env file to environment
func LoadSecretsToEnv() {
	godotenvErr := godotenv.Load("data/secrets.env")
	if godotenvErr != nil {
		// Loading secrets is mission-critical, fatal
		log.Error.Fatalf("Error loading secrets.env file. %v", godotenvErr)
	} else {
		log.Info.Println("Loaded secrets.env file successfully")
		log.Debug.Printf("APRICATE_ACCESS_SECRET: %s", os.Getenv("APRICATE_ACCESS_SECRET"))
	}
}