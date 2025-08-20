package db

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var dsn string
	environment := os.Getenv("ENV")

	if environment == "Production" {
		dsn = getProductionDSN()
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}

	var err error
	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return nil
		}
		fmt.Println("Waiting for MySQL to be ready...")
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("failed to connect to database after retries: %w", err)
}

func getProductionDSN() string {
	caCert := os.Getenv("DB_CA_CERT")
	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM([]byte(caCert)); !ok {
		panic("Failed to append CA cert")
	}

	tlsConfig := &tls.Config{
		RootCAs: rootCertPool,
	}

	// Register TLS config
	if err := mysql.RegisterTLSConfig("aiven", tlsConfig); err != nil {
		panic(err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=aiven",
		user, password, host, port, dbname,
	)
}
