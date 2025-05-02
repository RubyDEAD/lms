package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Request struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	EmailConfirm bool   `json:"email_confirm"`
}

type Response struct {
	ID string `json:"id"`
}

func CreateSupabaseAuthUser(email string, password string) (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	baseURL := os.Getenv("SUPABASE_URL_PROJECT")
	if baseURL == "" {
		log.Fatal("SUPABASE_URL_PROJECT environment variable is not set")
	}

	url := baseURL + "/auth/v1/admin/users"

	anonKey := os.Getenv("SUPABASE_ANONKEY")
	if anonKey == "" {
		log.Fatal("SUPABASE_ANONKEY environment variable is not set")
	}

	reqBody := Request{
		Email:        email,
		Password:     password,
		EmailConfirm: true,
	}

	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+anonKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", anonKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("supabase auth user creation failed: %v", resp.Status)
	}

	var supabaseResponse Response
	err = json.NewDecoder(resp.Body).Decode(&supabaseResponse)
	if err != nil {
		return "", err
	}

	return supabaseResponse.ID, nil
}
