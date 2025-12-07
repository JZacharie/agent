package runtime

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

type Config struct {
	Core struct {
		Log struct {
			EnableFile bool   `json:"enable_file" mapstructure:"enable_file"`
			Path       string `json:"path" mapstructure:"path"`
		} `json:"log" mapstructure:"log"`
	} `json:"core" mapstructure:"core"`
	Server struct {
		Port int `json:"port" mapstructure:"port"`
	} `json:"server" mapstructure:"server"`
	Library struct {
		Path           string   `json:"path" mapstructure:"path"`
		Blacklist      []string `json:"blacklist" mapstructure:"blacklist"`
		IgnoreDotFiles bool     `json:"ignore_dot_files" mapstructure:"ignore_dot_files"`
	} `json:"library" mapstructure:"library"`
	Render struct {
		MaxWorkers      int    `json:"max_workers" mapstructure:"max_workers"`
		ModelColor      string `json:"model_color" mapstructure:"model_color"`
		BackgroundColor string `json:"background_color" mapstructure:"background_color"`
	} `json:"render" mapstructure:"render"`
	Integrations struct {
		Thingiverse struct {
			Token string `json:"token" mapstructure:"token"`
		} `json:"thingiverse" mapstructure:"thingiverse"`
	} `json:"integrations" mapstructure:"integrations"`
	Database struct {
		Type     string `json:"type" mapstructure:"type"`
		Postgres struct {
			Host     string `json:"host" mapstructure:"host"`
			Port     int    `json:"port" mapstructure:"port"`
			User     string `json:"user" mapstructure:"user"`
			Password string `json:"password" mapstructure:"password"`
			Database string `json:"database" mapstructure:"database"`
			SSLMode  string `json:"sslmode" mapstructure:"sslmode"`
		} `json:"postgres" mapstructure:"postgres"`
	} `json:"database" mapstructure:"database"`
}

var Cfg *Config

var dataPath = "/data"

func init() {
	viper.BindEnv("DATA_PATH")
	if viper.GetString("DATA_PATH") != "" {
		dataPath = viper.GetString("DATA_PATH")
	}
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		err := os.MkdirAll(dataPath, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
	}

	bindEnv()

	if v := viper.GetInt("PORT"); v == 0 {
		viper.SetDefault("server.port", 8000)
	} else {
		viper.SetDefault("server.port", v)
	}
	if v := viper.GetString("LIBRARY_PATH"); v == "" {
		viper.SetDefault("library.path", "/library")
	} else {
		viper.SetDefault("library.path", v)
	}
	if v := viper.GetString("MODEL_RENDER_COLOR"); v == "" {
		viper.SetDefault("render.model_color", "#167DF0")
	} else {
		viper.SetDefault("render.model_color", v)
	}
	if v := viper.GetString("MODEL_BACKGROUND_COLOR"); v == "" {
		viper.SetDefault("render.background_color", "#FFFFFF")
	} else {
		viper.SetDefault("render.background_color", v)
	}

	viper.SetDefault("library.blacklist", []string{})
	viper.SetDefault("library.ignore_dot_files", true)
	viper.SetDefault("render.max_workers", 5)
	viper.SetDefault("core.log.enable_file", false)

	viper.SetDefault("server.hostname", "localhost")

	// Database defaults
	if v := viper.GetString("DATABASE_TYPE"); v == "" {
		viper.SetDefault("database.type", "sqlite")
	} else {
		viper.SetDefault("database.type", v)
	}
	if v := viper.GetString("POSTGRES_HOST"); v == "" {
		viper.SetDefault("database.postgres.host", "localhost")
	} else {
		viper.SetDefault("database.postgres.host", v)
	}
	if v := viper.GetInt("POSTGRES_PORT"); v == 0 {
		viper.SetDefault("database.postgres.port", 5432)
	} else {
		viper.SetDefault("database.postgres.port", v)
	}
	if v := viper.GetString("POSTGRES_USER"); v != "" {
		viper.SetDefault("database.postgres.user", v)
	}
	if v := viper.GetString("POSTGRES_PASSWORD"); v != "" {
		viper.SetDefault("database.postgres.password", v)
	}
	if v := viper.GetString("POSTGRES_DATABASE"); v != "" {
		viper.SetDefault("database.postgres.database", v)
	}
	if v := viper.GetString("POSTGRES_SSLMODE"); v == "" {
		viper.SetDefault("database.postgres.sslmode", "disable")
	} else {
		viper.SetDefault("database.postgres.sslmode", v)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(dataPath)
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	cfg := &Config{}
	viper.Unmarshal(cfg)

	cfg.Library.Blacklist = append(cfg.Library.Blacklist, ".project.stlib", ".thumb.png", ".render.png")

	if _, err := os.Stat(path.Join(dataPath, "config.toml")); os.IsNotExist(err) {
		log.Println("config.toml not found, creating...")
		SaveConfig(cfg)
	}

	Cfg = cfg
}

func bindEnv() {
	viper.BindEnv("PORT")
	viper.BindEnv("LIBRARY_PATH")
	viper.BindEnv("MAX_RENDER_WORKERS")
	viper.BindEnv("MODEL_RENDER_COLOR")
	viper.BindEnv("MODEL_BACKGROUND_COLOR")
	viper.BindEnv("LOG_PATH")
	viper.BindEnv("THINGIVERSE_TOKEN")
	viper.BindEnv("DATABASE_TYPE")
	viper.BindEnv("POSTGRES_HOST")
	viper.BindEnv("POSTGRES_PORT")
	viper.BindEnv("POSTGRES_USER")
	viper.BindEnv("POSTGRES_PASSWORD")
	viper.BindEnv("POSTGRES_DATABASE")
	viper.BindEnv("POSTGRES_SSLMODE")
}

func GetDataPath() string {
	return dataPath
}

func SaveConfig(cfg *Config) error {
	f, err := os.OpenFile(filepath.Join(GetDataPath(), "config.toml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	if err := toml.NewEncoder(f).Encode(cfg); err != nil {
		log.Println(err)
	}
	if err := f.Close(); err != nil {
		log.Println(err)
	}
	Cfg = cfg
	return err
}
