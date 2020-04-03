package config

import "os"

//Config ...
type Config struct {
	AuthServerAddr  string
	AuthPort        string
	HelloServerAddr string
	HelloPort       string
	JWTKey          string
	JWTPubKey       string
	JWTToken        string
	TLSCert         string
	TLSKey          string
	CACrt           string
	Username        string
	Password        string
	UseTLS          string
}

//Value ...
func Value() *Config {
	return &Config{
		AuthServerAddr:  getEnv("AUTH_SERVER_ADDR", "127.0.0.1"),
		AuthPort:        getEnv("AUTH_PORT", "50051"),
		HelloServerAddr: getEnv("HELLO_SERVER_ADDR", "127.0.0.1"),
		HelloPort:       getEnv("HELLO_PORT", "50052"),
		JWTKey:          getEnv("JWT_KEY", "ssl/jwt-key.pem"),
		JWTPubKey:       getEnv("JWT_PUBLIC_KEY", "ssl/jwt.pem"),
		JWTToken:        getEnv("JWT_TOKEN", "ssl/token.pem"),
		TLSCert:         getEnv("TLS_CERT", "ssl/auth.pem"),
		TLSKey:          getEnv("TLS_KEY", "ssl/auth-key.pem"),
		CACrt:           getEnv("CA_CERT", "ssl/ca.pem"),
		Username:        getEnv("USERNAME", "ec2ainun"),
		Password:        getEnv("PASSWORD", "isPartOfHumanRace"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
