package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvOptions struct {
	MongoURI          string
	DB                string
	Out               string
	Sample            int
	OptionalThreshold float64
	DateAs            string
	ObjectIDAs        string
}

func LoadDotEnv(path string) error {
	if path == "" {
		if _, err := os.Stat(".env"); err == nil {
			return godotenv.Load(".env")
		}
		return nil
	}
	return godotenv.Load(path)
}

func ReadEnv() EnvOptions {
	get := func(keys ...string) string {
		for _, k := range keys {
			if v := os.Getenv(k); v != "" {
				return v
			}
		}
		return ""
	}

	parseInt := func(s string, def int) int {
		if s == "" {
			return def
		}
		if v, err := strconv.Atoi(s); err == nil {
			return v
		}
		return def
	}

	parseFloat := func(s string, def float64) float64 {
		if s == "" {
			return def
		}
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return v
		}
		return def
	}

	envOptions := EnvOptions{
		MongoURI:          get("MONGOTS_MONGO_URI", "MONGO_URI"),
		Out:               get("MONGOTS_OUT"),
		Sample:            parseInt(get("MONGOTS_SAMPLE"), 20),
		OptionalThreshold: parseFloat(get("MONGOTS_OPTIONAL_THRESHOLD"), 0.98),
		DateAs:            get("MONGOTS_DATE_AS"),
		ObjectIDAs:        get("MONGOTS_OBJECTID_AS"),
	}

	fmt.Println("Env Options")
	fmt.Println(envOptions)

	return envOptions
}
