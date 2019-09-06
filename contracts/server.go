package contracts

type ServerConfig struct {
	Timeout    int64  `json:"timeout "yaml:"timeout" mapstructure:"timeout"`
	BindAddr   string `json:"bind_addr" yaml:"bind_addr" mapstructure:"bind_addr"`
	PoolSize   int    `json:"pool_size" yaml:"pool_size" mapstructure:"pool_size"`
	CryptoType string `json:"crypto_type" yaml:"crypto_type" mapstructure:"crypto_type"`
	Key        string `json:"key" yaml:"key"`
}

type PoolStatus struct {
	WorkerStausMap map[int]int
	FreeNum        int
}

type Server interface {
	Start()
	Stop()
	Status() PoolStatus
	GetConfig() ServerConfig
}
