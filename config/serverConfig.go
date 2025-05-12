package config

type ServerConfig struct {
	Port               int    `env:"PORT" envDefault:"80"`
	ListenAddress      string `env:"ADDRESS" envDefault:"0.0.0.0"`
	ConfigPath         string `env:"CONFIG_PATH"`
	ImageCacheDir      string `env:"IMAGE_CACHE_DIR" envDefault:"/tmp"`
	TalosFactoryServer string `env:"TALOS_FACTORY_SERVER" envDefault:"https://pxe.factory.talos.dev/image"`
	ServerAddress      string `env:"SERVER_ADDRESS" envDefault:"http://172.16.100.51"`
}

var ServerConf ServerConfig
