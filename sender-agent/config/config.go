package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Frequency int
	Role      string
	Network   string
}

var conf Config

func GetConfig() Config {
	return conf
}

func init() {
	//Load environment variables from file
	err := godotenv.Load("sender.env")
	if err != nil {
		log.Println("Error loading sender.env file")
	}

	//set frequency of recorded metrics from .env file, if an error occurs set to 60 seconds
	conf.Frequency, err = strconv.Atoi(os.Getenv("SEND_FREQUENCY"))
	if err != nil {
		conf.Frequency = 60
	}
	conf.Role = os.Getenv("ROLE_HOST")
	conf.Network = os.Getenv("NETWORK_NAME")
}
