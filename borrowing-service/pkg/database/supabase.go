package database

import (
	"fmt"
	"log"
	"os"

	"github.com/supabase-community/supabase-go"
)

var Client *supabase.Client

func Init() error {
	// Get DATABASE_URL from environment variable
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL is missing")
	}

	// Initialize the Supabase client with the DATABASE_URL
	var err error
	Client, err = supabase.NewClient(databaseURL, "", nil) // Empty string for the key (supabase-go should handle it)
	if err != nil {
		log.Printf("Failed to connect to Supabase: %v", err)
		return err
	}

	log.Println("✅ Successfully connected to Supabase")
	return nil
}
