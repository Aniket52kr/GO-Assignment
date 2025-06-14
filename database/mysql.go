package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func init() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Build DSN from individual env vars
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("MySQL open error:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("MySQL connection error:", err)
	}

	log.Println("Connected to MySQL")

	// Initialize DB schema from SQL file
	if err := runInitSQL("database/init.sql"); err != nil {
		log.Fatal("Init SQL error:", err)
	}
}

func runInitSQL(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	statements := splitSQLStatements(string(data))
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			log.Printf("Error executing statement: %s\nErr: %v\n", stmt, err)
			return err
		}
	}
	return nil
}

func splitSQLStatements(script string) []string {
	var stmts []string
	current := ""
	for _, line := range strings.Split(script, "\n") {
		line = trimComment(line)
		current += line + " "
		if endsWithSemicolon(line) {
			stmts = append(stmts, strings.TrimSpace(current))
			current = ""
		}
	}
	if current != "" {
		stmts = append(stmts, strings.TrimSpace(current))
	}
	return stmts
}

func trimComment(line string) string {
	if i := strings.Index(line, "--"); i != -1 {
		return strings.TrimSpace(line[:i])
	}
	return strings.TrimSpace(line)
}

func endsWithSemicolon(line string) bool {
	return strings.HasSuffix(strings.TrimSpace(line), ";")
}

// Expose db globally
func GetDB() *sql.DB {
	return db
}
