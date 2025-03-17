package config

import "github.com/ilyakaznacheev/cleanenv"

func LoadConfig(path string, cfg interface{}) error {
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return err
	}
	return nil
}
