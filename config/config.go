package config

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	env "github.com/Netflix/go-env"
	"github.com/go-playground/validator/v10"

	"github.com/rangidev/rangi/database"
)

var (
	validate             = validator.New()
	contentPathParts     = []string{"content"}
	blueprintsPathParts  = append(contentPathParts, "blueprints")
	sqlite3FilePathParts = append(contentPathParts, "sqlite3", "rangi.db")
)

type Config struct {
	// Server
	HostAndPort string `env:"RANGI_HOST_AND_PORT,default=:6532"`
	// Log
	LogLevel  string `env:"RANGI_LOG_LEVEL,default=info" validate:"oneof=debug info warn error"`
	LogFormat string `env:"RANGI_LOG_FROMAT,default=text" validate:"oneof=text json"`
	// Templates
	EnableTemplateDevelopment bool `env:"RANGI_ENABLE_TEMPLATE_DEVELOPMENT,default=false"`
	// Blueprints
	BlueprintsPath string `env:"RANGI_BLUEPRINTS_PATH"`
	// Database
	DatabaseType string `env:"RANGI_DATABASE_TYPE,default=sqlite3" validate:"oneof=sqlite3 postgresql"`
	// Sqlite3
	// Used to override the default file path
	Sqlite3DatabaseFile string `env:"RANGI_SQLITE3_DATABASE_FILE"`
	// Admin interface
	AdminItemsLimit int `env:"RANGI_ADMIN_ITEMS_LIMIT,default=50" validate:"gte=1,lte=200"`

	// Used to reference assets that are stored relative to the binary. TODO: Do we want to store everything relative to the binary as default behavior?
	ExecutableDir    string
	Logger           *slog.Logger
	DatabaseInstance *database.DB
	Validate         *validator.Validate
}

func New() *Config {
	// Read env vars
	var config Config
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatalln("Error while unmarshaling config from environment:", err)
	}
	// Validate
	err = validate.Struct(&config)
	if err != nil {
		log.Fatalln("Error while validating config:", err)
	}

	// Logger
	var logLevel slog.Level
	err = logLevel.UnmarshalText([]byte(config.LogLevel))
	if err != nil {
		log.Fatalln("Error while determining log level", err)
	}
	switch config.LogFormat {
	case "text":
		config.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		}))
	case "json":
		config.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		}))
	default:
		log.Fatalln("Unrecognized log format", config.LogFormat)
	}

	// Executabe path
	executablePath, err := os.Executable()
	if err != nil {
		config.Logger.Error("Could not get executable path", "error", err)
		os.Exit(1)
	}
	config.ExecutableDir = filepath.Dir(executablePath)

	// Blueprints
	blueprintsPath := config.BlueprintsPath
	if blueprintsPath == "" {
		// Use default blueprints path
		blueprintsPath = filepath.Join(append([]string{config.ExecutableDir}, blueprintsPathParts...)...)
		err = os.MkdirAll(blueprintsPath, os.ModePerm)
		if err != nil {
			config.Logger.Error("Could not make directory for blueprints", "error", err)
			os.Exit(1)
		}
	}
	config.BlueprintsPath = blueprintsPath

	// Database
	switch config.DatabaseType {
	case "sqlite3":
		databaseFilename := config.Sqlite3DatabaseFile
		if databaseFilename == "" {
			// Use default sqlite3 database file
			databaseFilename = filepath.Join(append([]string{config.ExecutableDir}, sqlite3FilePathParts...)...)
			err = os.MkdirAll(filepath.Dir(databaseFilename), os.ModePerm)
			if err != nil {
				config.Logger.Error("Could not make directory for sqlite3 database", "error", err)
				os.Exit(1)
			}
		}
		dbInstance, err := database.NewSqlite3Instance(databaseFilename)
		if err != nil {
			config.Logger.Error("Could not connect to sqlite3 database", "error", err)
			os.Exit(1)
		}
		config.DatabaseInstance = dbInstance
	default:
		config.Logger.Error("Invalid database type provided", "type", config.DatabaseType)
		os.Exit(1)
	}

	// Validator
	config.Validate = validate

	return &config
}
