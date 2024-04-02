package utils

type Config struct {
    AuthenticatedPrefix   string `yaml:"authenticated-prefix"`
    UnauthenticatedPrefix string `yaml:"unauthenticated-prefix"` // 1. CHANGE THESE TO SNAKE CASE
	JwtRSAPublicKey       string `yaml:"jwt-rsa-public-key"`
    JwtRSAPrivateKey      string `yaml:"jwt-rsa-private-key"`
}