package agent

type Config struct {
	APIServerAddress string `yaml:"apiServerAddress"`
	APIServerPort    int32  `yaml:"apiServerPort"`

	HumstackAPIServerAddress string `yaml:"humstackAPIServerAddress"`
	HumstackAPIServerPort    int32  `yaml:"humstackAPIServerPort"`
}
