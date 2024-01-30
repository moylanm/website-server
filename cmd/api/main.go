package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"mylesmoylan.net/internal/data"
)

type config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Db   struct {
		Dsn          string        `yaml:"dsn"`
		MaxOpenConns int           `yaml:"maxOpenConns"`
		MaxIdleConns int           `yaml:"maxIdleConns"`
		MaxIdleTime  time.Duration `yaml:"maxIdleTime"`
	} `yaml:"db"`
	Limiter struct {
		Rps     float64 `yaml:"rps"`
		Burst   int     `yaml:"burst"`
		Enabled bool    `yaml:"enabled"`
	} `yaml:"limiter"`
	Admin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"admin"`
}

type application struct {
	config        config
	logger        *slog.Logger
	models        data.Models
	templateCache map[string]*template.Template
}

func readConfig(path string) (config, error) {
	var cfg config

	data, err := os.ReadFile(path)
	if err != nil {
		return config{}, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return config{}, err
	}

	overrideConfigWithEnv(&cfg)

	return cfg, nil
}

func overrideConfigWithEnv(cfg *config) {
	if dsn := os.Getenv("WEBSITE_DB_DSN"); dsn != "" {
		cfg.Db.Dsn = dsn
	}

	if username := os.Getenv("WEBSITE_USER"); username != "" {
		cfg.Admin.Username = username
	}

	if password := os.Getenv("WEBSITE_PASS"); password != "" {
		cfg.Admin.Password = password
	}
}

func validateConfig(cfg *config) error {
	// Validate DB configuration
	if cfg.Db.Dsn == "" {
		return fmt.Errorf("database DSN is required")
	}
	if cfg.Db.MaxOpenConns <= 0 {
		return fmt.Errorf("max open connections must be positive")
	}
	if cfg.Db.MaxIdleConns <= 0 {
		return fmt.Errorf("max idle connections must be positive")
	}
	if cfg.Db.MaxIdleTime <= 0 {
		return fmt.Errorf("max idle time must be positive")
	}

	// Validate server configuration
	if cfg.Host != "" && cfg.Host != "localhost" {
		return fmt.Errorf("invalid hostname")
	}
	if cfg.Port < 1024 || cfg.Port > 65535 {
		return fmt.Errorf("port must be between 1024 and 65535")
	}

	// Validate admin credentials
	if cfg.Admin.Username == "" || cfg.Admin.Password == "" {
		return fmt.Errorf("admin username and password are required")
	}

	// Validate limiter configuration if enabled
	if cfg.Limiter.Enabled {
		if cfg.Limiter.Rps <= 0 {
			return fmt.Errorf("rate limiter RPS must be positive")
		}
		if cfg.Limiter.Burst <= 0 {
			return fmt.Errorf("rate limiter burst must be positive")
		}
	}

	return nil
}

func main() {
	var cfgPath string

	flag.StringVar(&cfgPath, "config-path", "config.yaml", "Configuration file path")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := readConfig(cfgPath)
	if err != nil {
		logger.Error(err.Error())
	}

	err = validateConfig(&cfg)
	if err != nil {
		logger.Error("Invalid configuration: " + err.Error())
		os.Exit(1)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		config:        cfg,
		logger:        logger,
		models:        data.NewModels(db),
		templateCache: templateCache,
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Db.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
