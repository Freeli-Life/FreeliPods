package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Whitelist []string
	Blacklist []string
	Port      int
	WebPort   int
	TLSKey    string
	TLSCert   string
	Domain    string
}

const defaultConfigPath = "freelipods.conf"

func LoadConfig() (*Config, error) {
	if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
		log.Println("Configuration file not found, creating default freelipods.conf")
		return createDefaultConfig()
	}

	return readConfigFromFile()
}

func createDefaultConfig() (*Config, error) {
	cfg := &Config{
		Whitelist: []string{},
		Blacklist: []string{},
		Port:      50051,
		WebPort:   8080,
		TLSKey:    "server.key",
		TLSCert:   "server.crt",
		Domain:    "localhost",
	}

	file, err := os.Create(defaultConfigPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("warning: failed to close config file: %v", err)
		}
	}()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(
		"domain=localhost\n" +
			"port=50051\n" +
			"web_port=8080\n" +
			"tls_key=server.key\n" +
			"tls_cert=server.crt\n" +
			"whitelist=\n" +
			"blacklist=\n",
	)
	if err != nil {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func readConfigFromFile() (*Config, error) {
	file, err := os.Open(defaultConfigPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("warning: failed to close config file: %v", err)
		}
	}()

	cfg := &Config{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "domain":
			cfg.Domain = value
		case "port":
			cfg.Port, _ = strconv.Atoi(value)
		case "web_port":
			cfg.WebPort, _ = strconv.Atoi(value)
		case "tls_key":
			cfg.TLSKey = value
		case "tls_cert":
			cfg.TLSCert = value
		case "whitelist":
			cfg.Whitelist = strings.Split(value, ",")
		case "blacklist":
			cfg.Blacklist = strings.Split(value, ",")
		}
	}
	return cfg, scanner.Err()
}