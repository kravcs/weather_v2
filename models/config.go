package models

type Config struct {
	API struct {
		Apikey   string `yaml:"apikey"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"api"`
	Cache struct {
		Duration string `yaml:"duration"`
	} `yaml:"cache"`
	Database struct {
		Driver string `yaml:"driver"`
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
	} `yaml:"database"`
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
}
