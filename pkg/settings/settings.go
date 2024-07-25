package settings

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
	Redis struct {
		Host     string
		Port     int
		Password string
	}
	KEY struct {
		TG string
	}
	API struct {
		Port string
	}
	Logging struct {
		Level string
	}
}

func (c *Config) ReadConfig() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found or error loading .env file")
	}

	c.Database.Host = os.Getenv("POSTGRES_HOST")
	c.Database.Port, err = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return fmt.Errorf("invalid POSTGRES_PORT: %v", err)
	}
	c.Database.User = os.Getenv("POSTGRES_USER")
	c.Database.Password = os.Getenv("POSTGRES_PASSWORD")
	c.Database.DBName = os.Getenv("POSTGRES_DB")

	c.Redis.Host = os.Getenv("REDIS_HOST")
	c.Redis.Port, err = strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return fmt.Errorf("invalid REDIS_PORT: %v", err)
	}

	c.Redis.Password = os.Getenv("REDIS_PASSWORD")

	c.KEY.TG = os.Getenv("TG_KEY")

	c.API.Port = ":" + os.Getenv("API_PORT")

	c.Logging.Level = os.Getenv("LOG_LEVEL")

	return nil
}
