package contracts

type ClientConfig struct {
	Timeout    int    `json:"timeout "yaml:"timeout" mapstructure:"timeout"`
	ServerAddr string `json:"server_addr" yaml:"server_addr" mapstructure:"server_addr"`
	ClientNum  int    `json:"client_num" yaml:"client_num" mapstructure:"client_num"`
	CryptoType string `json:"crypto_type" yaml:"crypto_type" mapstructure:"crypto_type"`
	Key        string `json:"key" yaml:"key"`
}

type Client interface {
	Send(msg *Msg) error
	Read() (*Msg, error)
	Connect() error
	Close() error
}
