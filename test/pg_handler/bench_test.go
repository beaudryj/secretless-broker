package main

import (
	"os"
	"fmt"
	"testing"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	socketPath = "/sock/.s.PGSQL.5432"
	select1Query = "SELECT 1"
	select10Query = "SELECT generate_series(0, 9)"
	select100Query = "SELECT generate_series(0, 99)"
	select1000Query = "SELECT generate_series(0, 999)"
	select10000Query = "SELECT generate_series(0, 9999)"
	envPostgresPort = "POSTGRES_PORT"
	envSecretlessPort = "SECRETLESS_PORT"
	dockerPortPostgres = "5432"
	dockerPortSecretless = "15432"
)

func getPostgresConnection() (*sql.DB, error) {
	port := os.Getenv(envPostgresPort)
	if port == "" {
		port = dockerPortPostgres
	}
	connStr := fmt.Sprintf("port=%s sslmode=disable database=postgres user=test", port)
	return configureDB(sql.Open("postgres", connStr))
}

func getSecretlessConnection() (*sql.DB, error) {
	port := os.Getenv(envSecretlessPort)
	if port == "" {
		port = dockerPortSecretless
	}
	connStr := fmt.Sprintf("port=%s sslmode=disable database=postgres", port)
	return configureDB(sql.Open("postgres", connStr))
}

func configureDB(db *sql.DB, err error) (*sql.DB, error) {
	if err != nil {
		return nil, err
	}

	// Make sure pooling does not occur
	db.SetMaxOpenConns(1)

	return db, nil
}

// runQuery executes a query. Expects the timer to already have been stopped.
func runQuery(b *testing.B, db *sql.DB, query string) {
	b.StartTimer()
	rows, err := db.Query(query)
	if err != nil {
		b.Fatal(err)
	}
	b.StopTimer()
	rows.Close()
}

func benchmarkBaselineQuery(b *testing.B, query string) {
	b.StopTimer()
	db, err := getPostgresConnection()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < b.N; i++ {
		runQuery(b, db, query)
	}
}

func benchmarkSecretlessQuery(b *testing.B, query string) {
	b.StopTimer()
	db, err := getSecretlessConnection()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < b.N; i++ {
		runQuery(b, db, query)
	}
}

func BenchmarkBaseline_Select1(b *testing.B) { 
	benchmarkBaselineQuery(b, select1Query)
}

func BenchmarkBaseline_Select10(b *testing.B) { 
	benchmarkBaselineQuery(b, select10Query)
}

func BenchmarkBaseline_Select100(b *testing.B) { 
	benchmarkBaselineQuery(b, select100Query)
}

func BenchmarkBaseline_Select1000(b *testing.B) { 
	benchmarkBaselineQuery(b, select1000Query)
}

func BenchmarkBaseline_Select10000(b *testing.B) { 
	benchmarkBaselineQuery(b, select10000Query)
}

func BenchmarkSecretless_Select1(b *testing.B) {
	benchmarkSecretlessQuery(b, select1Query)
}

func BenchmarkSecretless_Select10(b *testing.B) {
	benchmarkSecretlessQuery(b, select10Query)
}

func BenchmarkSecretless_Select100(b *testing.B) {
	benchmarkSecretlessQuery(b, select100Query)
}

func BenchmarkSecretless_Select1000(b *testing.B) {
	benchmarkSecretlessQuery(b, select1000Query)
}

func BenchmarkSecretless_Select10000(b *testing.B) {
	benchmarkSecretlessQuery(b, select10000Query)
}