package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func longestPalindrome(s string) string {
	if s == "" || len(s) < 1 {
		return ""
	}
	var start, end int
	for i := 0; i < len(s); i++ {
		len1 := expandFromMiddle(s, i, i)
		len2 := expandFromMiddle(s, i, i+1)
		len := Max(len1, len2)
		if len > (end - start) {
			start = i - (len-1)/2
			end = i + len/2

		}
	}
	return s[start : end+1]
}
func expandFromMiddle(s string, l, r int) int {
	if s == "" || l > r {
		return 0
	}
	for l >= 0 && r < len(s) && string(s[l]) == string(s[r]) {
		l--
		r++
	}
	return r - l - 1
}

func writeToDatabase(db *sql.DB, str, pal string) error {
	stmtIns, err := db.Prepare("INSERT INTO palindromes VALUES( ?, ? )")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(str, pal)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}
func initDatabase(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( %s VARCHAR(255) NOT NULL, %s VARCHAR(255) NOT NULL)", "palindromes", "string", "longest_palindrome")
	stmtIns, err := db.Prepare(stmt)
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec()
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}
func main() {
	var inputStr *string
	inputStr = flag.String("inputString", "", "Input string to find longest palindromic substring")
	flag.Parse()
	if inputStr != nil && *inputStr != "" {
		user := os.Getenv("MYSQL_USER")
		pass := os.Getenv("MYSQL_PASS")
		host := os.Getenv("MYSQL_HOST")
		port := os.Getenv("MYSQL_PORT")
		database := os.Getenv("MYSQL_DB")
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, database)
		log.Println("Opeing connection to database...")
		db, err := sql.Open("mysql", connStr)
		if err != nil {
			log.Fatalf("Error connecting to mysql database. Error: %s", err.Error())
		}
		log.Println("Initializing database...")
		err = initDatabase(db)
		if err != nil {
			log.Fatalf("Error initializing database. Error: %s", err.Error())
		}
		log.Printf("Calculating longest palindrome for %s...", *inputStr)
		pal := longestPalindrome(*inputStr)
		log.Printf("Longest palindrome for string %s is %s", *inputStr, pal)
		log.Println("writing to database...")
		err = writeToDatabase(db, *inputStr, pal)
		if err != nil {
			log.Fatalf("Error writing to database. Error: %s", err.Error())
		}
		defer db.Close()
	} else {
		log.Fatalf("No inputString specified. Exiting...")
	}
}
