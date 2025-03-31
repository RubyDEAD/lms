package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func Connect_to_DB() (*pgx.Conn, error) {
	var err error
	DB, err = pgx.Connect(context.Background(), "postgresql://postgres.roifumfsdyhyegikikpd:Cih4Q9lIpPH0MGMD@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected to PostgresSQL")
	return DB, nil
}
