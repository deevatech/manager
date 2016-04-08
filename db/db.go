package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

type Config struct {
	Host string
	User string
	Pass string
	Name string
}

var config *Config
var DB *gorm.DB

func init() {
	config = &Config{
		Host: dbGetEnv("DEEVA_DB_HOST", true),
		User: dbGetEnv("DEEVA_DB_USER", true),
		Pass: dbGetEnv("DEEVA_DB_PASS", false),
		Name: dbGetEnv("DEEVA_DB_NAME", true),
	}

	log.Printf("DB Connection: %s", config)
	db, err := gorm.Open("mysql", config.String())
	if err != nil {
		log.Fatalf("Error opening database connection: %s", err)
	}

	DB = db
}

func dbGetEnv(name string, required bool) string {
	value := os.Getenv(name)
	if len(value) == 0 && required == true {
		log.Fatalf("Missing %s from the environment!", name)
	}

	return value
}

func (c Config) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", c.User, c.Pass, c.Host, c.Name)
}
