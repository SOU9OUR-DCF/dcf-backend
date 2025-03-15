package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	baseConfigPath := fmt.Sprintf("%s/config.yaml", path)
	baseConfig, err := processConfigFile(baseConfigPath)
	if err != nil {
		return nil, err
	}

	if err := v.ReadConfig(strings.NewReader(baseConfig)); err != nil {
		return nil, fmt.Errorf("error reading base config: %w", err)
	}

	env := os.Getenv("APP_ENV")
	if env != "" {
		envConfigPath := fmt.Sprintf("%s/config.%s.yaml", path, env)
		if _, err := os.Stat(envConfigPath); err == nil {
			envConfig, err := processConfigFile(envConfigPath)
			if err != nil {
				return nil, err
			}

			envViper := viper.New()
			envViper.SetConfigType("yaml")
			if err := envViper.ReadConfig(strings.NewReader(envConfig)); err != nil {
				return nil, fmt.Errorf("error reading env config: %w", err)
			}

			if err := v.MergeConfigMap(envViper.AllSettings()); err != nil {
				return nil, fmt.Errorf("error merging configs: %w", err)
			}

			fmt.Println("Merged environment config file:", envConfigPath)
		}
	}

	v.SetDefault("server.port", "8080")
	v.SetDefault("server.readTimeout", time.Second*10)
	v.SetDefault("server.writeTimeout", time.Second*10)
	v.SetDefault("server.shutdownTimeout", time.Second*30)
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("redis.address", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("jwt.expiresIn", time.Hour*24)

	if !v.IsSet("jwt.secret") {
		return nil, fmt.Errorf("jwt secret is required")
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
}

func processConfigFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	processed := processEnvVarSyntax(string(data))

	return processed, nil
}

func processEnvVarSyntax(content string) string {
	result := content

	for {
		start := strings.Index(result, "${")
		if start == -1 {
			break
		}

		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		end += start

		varExpr := result[start+2 : end]

		var varName, defaultVal string
		if strings.Contains(varExpr, ":-") {
			parts := strings.SplitN(varExpr, ":-", 2)
			varName = parts[0]
			defaultVal = parts[1]
		} else {
			varName = varExpr
		}

		value := os.Getenv(varName)
		if value == "" {
			value = defaultVal
		}

		result = result[:start] + value + result[end+1:]
	}

	return result
}
