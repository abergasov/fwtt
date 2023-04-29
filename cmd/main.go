package main

import (
	"flag"
	"fmt"
	"fwtt/internal/config"
	"fwtt/internal/logger"
	"fwtt/internal/repository/quotes"
	validatorRepo "fwtt/internal/repository/validator"
	"fwtt/internal/routes"
	"fwtt/internal/service/hasher"
	"fwtt/internal/service/quoter"
	validatorService "fwtt/internal/service/validator"
	"fwtt/internal/storage/database"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	confFile = flag.String("config", "configs/app_conf.yml", "Configs file path")
	appHash  = os.Getenv("GIT_HASH")
)

func main() {
	flag.Parse()
	appLog, err := logger.NewAppLogger(appHash)
	if err != nil {
		log.Fatalf("unable to create logger: %s", err)
	}
	appLog.Info("app starting", zap.String("conf", *confFile))
	appConf, err := config.InitConf(*confFile)
	if err != nil {
		appLog.Fatal("unable to init config", err, zap.String("config", *confFile))
	}

	appLog.Info("create storage connections")
	dbConn, err := getDBConnect(appLog, &appConf.ConfigDB, appConf.MigratesFolder)
	if err != nil {
		appLog.Fatal("unable to connect to db", err, zap.String("host", appConf.ConfigDB.Address))
	}
	defer func() {
		if err = dbConn.Close(); err != nil {
			appLog.Fatal("unable to close db connection", err)
		}
	}()

	appLog.Info("init repositories")
	repoValidator := validatorRepo.NewRepository(dbConn)
	repoQuotes := quotes.NewRepository(dbConn)

	appLog.Info("init services")
	serviceQuotes := quoter.NewService(appLog, repoQuotes)
	serviceHasher := hasher.NewService()
	serviceValidator := validatorService.NewService(
		appLog,
		appConf.ConfigValidator.ChallengeDifficulty,
		appConf.ConfigValidator.ChallengeMaxAllowed,
		appConf.ConfigValidator.ChallengeTTL,
		serviceHasher,
		repoValidator,
	)
	defer serviceValidator.Stop()

	appLog.Info("init http service")
	appHTTPServer := routes.InitAppRouter(appLog, serviceQuotes, serviceValidator, fmt.Sprintf(":%d", appConf.AppPort))
	defer func() {
		if err = appHTTPServer.Stop(); err != nil {
			appLog.Fatal("unable to stop http service", err)
		}
	}()
	go func() {
		if err = appHTTPServer.Run(); err != nil {
			appLog.Fatal("unable to start http service", err)
		}
	}()

	// register app shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // This blocks the main thread until an interrupt is received
}

func getDBConnect(appLog logger.AppLogger, cnf *config.DBConf, migratesFolder string) (*database.DBConnect, error) {
	for i := 0; i < 5; i++ {
		dbConnect, err := database.InitDBConnect(cnf, migratesFolder)
		if err == nil {
			return dbConnect, nil
		}
		appLog.Error("can't connect to db", err, zap.Int("attempt", i))
		time.Sleep(time.Duration(i) * time.Second * 5)
	}
	return nil, fmt.Errorf("can't connect to db")
}
