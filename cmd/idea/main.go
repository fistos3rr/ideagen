package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fistos3rr/ideagen/internal/ai"
	"github.com/fistos3rr/ideagen/internal/data"
	"github.com/fistos3rr/ideagen/internal/jsonlog"
	"github.com/fistos3rr/ideagen/internal/prompt"
	redis_repo "github.com/fistos3rr/ideagen/internal/redis"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type config struct {
	port           int
	env            string
	aiProviderType string
	db             struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	jwt struct {
		secret          []byte
		accessTokenTTL  time.Duration
		refreshTokenTTL time.Duration
	}
	redis struct {
		addr          string
		password      string
		ideaBufferTTL time.Duration
	}
}

type application struct {
	config          config
	logger          *jsonlog.Logger
	aiProvider      ai.Provider
	aiConfig        ai.Config
	models          data.Models
	redisRepository redis_repo.Repository
	wg              sync.WaitGroup
	promptManager   *prompt.PromptManager
}

func (cfg *config) parseEnv() {
	// JWT
	// --------------------------------------------------------
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		panic("no jwt secret provided")
	}
	jwtSecret := []byte(secret)

	strVal := os.Getenv("ACCESS_TOKEN_TTL_MINUTES")
	accessTTL, err := strconv.Atoi(strVal)
	if err != nil {
		accessTTL = 15
	}

	strVal = os.Getenv("REFRESH_TOKEN_TTL_DAYS")
	refreshTTL, err := strconv.Atoi(strVal)
	if err != nil {
		refreshTTL = 30
	}

	cfg.jwt.secret = jwtSecret
	cfg.jwt.accessTokenTTL = time.Duration(accessTTL) * time.Minute
	cfg.jwt.refreshTokenTTL = time.Duration(refreshTTL) * time.Hour * 24
	// --------------------------------------------------------

	// APP
	// --------------------------------------------------------
	strVal = os.Getenv("PORT")
	port, err := strconv.Atoi(strVal)
	if err != nil {
		port = 8080
	}

	strVal = os.Getenv("AI_PROVIDER")
	aiProviderType := strVal
	if aiProviderType == "" {
		aiProviderType = "groq"
	}

	cfg.port = port
	cfg.aiProviderType = aiProviderType
	// --------------------------------------------------------

	// DATABASE
	// --------------------------------------------------------
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	strVal = os.Getenv("DB_MAX_OPEN_CONNS")
	maxOpenConns, err := strconv.Atoi(strVal)
	if err != nil {
		maxOpenConns = 25
	}

	strVal = os.Getenv("DB_MAX_IDLE_CONNS")
	maxIdleConns, err := strconv.Atoi(strVal)
	if err != nil {
		maxIdleConns = 25
	}

	strVal = os.Getenv("DB_MAX_IDLE_TIME")
	maxIdleTime := strVal
	if maxIdleTime == "" {
		maxIdleTime = "15m"
	}

	cfg.db.dsn = dsn
	cfg.db.maxOpenConns = maxOpenConns
	cfg.db.maxIdleConns = maxIdleConns
	cfg.db.maxIdleTime = maxIdleTime
	// --------------------------------------------------------

	// REDIS
	// --------------------------------------------------------
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	strVal = os.Getenv("REDIS_IDEA_TTL_MINUTE")
	ideaBufferTTL, err := strconv.Atoi(strVal)
	if err != nil {
		ideaBufferTTL = 15
	}

	cfg.redis.addr = redisHost + ":" + redisPort
	cfg.redis.password = redisPassword
	cfg.redis.ideaBufferTTL = time.Duration(ideaBufferTTL) * time.Minute
	// --------------------------------------------------------
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	var cfg config
	cfg.parseEnv()
	cfg.env = "development"

	logger.PrintInfo("jwt token config", map[string]string{
		"access_token_ttl":  string(int(cfg.jwt.accessTokenTTL / time.Minute)),
		"refresh_token_ttl": string(int(cfg.jwt.refreshTokenTTL / (24 * time.Hour))),
	})

	logger.PrintInfo("database config", map[string]string{
		"dsn": cfg.db.dsn,
	})

	aicfg := ai.Config{
		APIKey: os.Getenv("AI_API_KEY"),
		Model:  os.Getenv("AI_MODEL"),
		APIURL: os.Getenv("AI_API_URL"),
	}

	aiLogData := map[string]string{
		"api_url": aicfg.APIURL,
		"model":   aicfg.Model,
	}

	if cfg.env == "development" {
		aiLogData["api_key"] = aicfg.APIKey
	}

	var provider ai.Provider
	switch cfg.aiProviderType {
	case "groq":
		logger.PrintInfo("running groq ai provider", aiLogData)
		provider = ai.NewGroqClient(aicfg)
	default:
		logger.PrintFatal(errors.New("unknown provider"), map[string]string{
			"provider": cfg.aiProviderType,
		})
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	rdb, err := openRedis(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer rdb.Close()

	logger.PrintInfo("redis connection established", nil)

	redisConfig := redis_repo.Config{
		IdeaTTL: cfg.redis.ideaBufferTTL,
	}

	promptManager, err := prompt.NewPromptManager(".")
	ok := prompt.IsDefaultErr(err)
	if ok {
		logger.PrintInfo("initializing prompt manager", map[string]string{
			"files": err.Error(),
		})
	} else if err != nil {
		panic(err)
	}

	app := &application{
		config:          cfg,
		logger:          logger,
		aiProvider:      provider,
		aiConfig:        aicfg,
		models:          data.NewModels(db),
		redisRepository: redis_repo.NewRepository(rdb, redisConfig),
		promptManager:   promptManager,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func openRedis(cfg config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.redis.addr,
		Password: cfg.redis.password,
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
