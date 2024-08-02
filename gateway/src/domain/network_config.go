package domain

type NetworkConfig struct {
	ListenIp   string
	ListenPort int
	Timeout    int
	HMACKey    string
}
