package main

import (
	"github.com/Meromen/mfc-backend/controller"
	"github.com/Meromen/mfc-backend/db"
	"github.com/Meromen/mfc-backend/httpserver"
	"github.com/Meromen/mfc-backend/preferences"
	"github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
)

type Exit struct{ Code int }

func handleExit(logger *logrus.Logger) {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true {
			os.Exit(exit.Code)
		}
		logger.Errorf("%v Stask trace: %s", e, debug.Stack())
		os.Exit(1)
	}
}

func main() {
	logger := logrus.New()

	p, err := preferences.Get()
	if err != nil {
		logger.Fatalf("Failed to set preferences: %v\n", err)
	}

	if p.LogAsJSON {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	logger.Level = logrus.Level(p.LogLevel)

	logger.Formatter = &logrus.JSONFormatter{}
	logrus.SetOutput(logger.Writer())

	defer handleExit(logger)

	dbConn, err := db.Connect(&p.PostgresUrl)
	if err != nil {
		logger.Errorf("Failed to open postgres connection: %v\n", err)
		panic(Exit{1})
	}
	defer dbConn.Close()
	err = dbConn.Ping()
	if err != nil {
		logger.Errorf("Failed to check postgres connection: %v\n", err)
		panic(Exit{1})
	}

	mfcStorage, err := db.NewMfcStorage(dbConn)
	if err != nil {
		logger.Errorf("Failed to create mfc storage: %v\n", err)
		panic(Exit{1})
	}

	c := controller.NewController(mfcStorage, logger)

	server := httpserver.NewServer(10, 10, p.ServerAddress, c, logger)

	server.Start()
}
