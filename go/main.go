package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type DatabaseConfig struct {
	Host   string
	Port   uint16
	DBName string
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("User is ?")
	ok := scanner.Scan()
	if !ok {
		fmt.Fprintf(os.Stderr, "Unable to read username\n")
		os.Exit(1)
	}
	username := scanner.Text()
	fmt.Println("Password is ?")
	ok = scanner.Scan()
	if !ok {
		fmt.Fprintf(os.Stderr, "Unable to read password\n")
		os.Exit(1)
	}
	password := scanner.Text()
	fmt.Println("User is " + username)
	fmt.Println("Password is " + password)

	// Вызываем ParseConn для присвоения данных из конфинурации
	connConfig, err := ParseConn("connect.conf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading connect.conf: %v\n", err)
		os.Exit(1)
	}

	urlExample := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, connConfig.Host, connConfig.Port, connConfig.DBName)
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var version string
	err = conn.QueryRow(context.Background(), "SELECT VERSION()").Scan(&version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(version)
}

func ParseConn(filePath string) (DatabaseConfig, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	uri := strings.TrimSpace(string(fileContent))

	connConfig, err := pgx.ParseConfig(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing URI: %v\n", err)
		os.Exit(1)
	}

	if connConfig.Host == "" || connConfig.Port == 0 || connConfig.Database == "" {
		err := errors.New("Missing required fields in the configuration")
		fmt.Printf("Error validating configuration: %v\n", err)
		os.Exit(1)
	}
	// Создание структуры DatabaseConfig и присвоение значений
	dbConfig := DatabaseConfig{
		Host:   connConfig.Host,
		Port:   connConfig.Port,
		DBName: connConfig.Database,
	}

	return dbConfig, err
}
