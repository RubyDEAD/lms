package database

import (
	"os"

	"github.com/supabase-community/supabase-go"
)

var Client *supabase.Client

func Init() error {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	var err error
	Client, err = supabase.NewClient(url, key, nil)
	return err
}
