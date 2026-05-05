package config

type Config struct {
	InputDir  string `yaml:"input"`
	OutputDir string `yaml:"output"`
	Recursive bool   `yaml:"recursive"`
	Overwrite bool   `yaml:"overwrite"`

	Template struct {
		Layout string `yaml:"layout"`
		Style  string `yaml:"style"`
	} `yaml:"template"`

	Extensions string  `yaml:"extensions"`
}

