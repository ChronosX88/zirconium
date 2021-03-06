package core

import (
	"context"
	"fmt"
	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type AppContext struct {
	router          *Router
	authManager     *AuthManager
	sessionManager  *SessionManager
	websocketServer *WebsocketServer
	cfg             *Config
	database        *mongo.Database
	userManager     *UserManager
}

func NewAppContext(cfg *Config) *AppContext {
	var err error
	appContext := &AppContext{}
	appContext.cfg = cfg

	appContext.connectToDatabase()

	router, err := NewRouter(appContext)
	if err != nil {
		logger.Fatalf("Unable to initialize router: %s", err.Error())
	}
	appContext.router = router

	um, err := NewUserManager(appContext.database)
	if err != nil {
		logger.Fatalf("Unable to initialize user manager: %s", err.Error())
	}
	appContext.userManager = um

	authManager, err := NewAuthManager(um, cfg.ServerID, router)
	if err != nil {
		logger.Fatalf("Unable to initialize authentication manager: %s", err.Error())
	}
	appContext.authManager = authManager

	sessionManager := NewSessionManager(router, cfg.ServerDomains)
	appContext.sessionManager = sessionManager

	wss := NewWebsocketServer(cfg, sessionManager)
	appContext.websocketServer = wss

	return appContext
}

func (ac *AppContext) connectToDatabase() {
	ctx := context.TODO()
	dbUri := fmt.Sprintf("mongodb://%s:%s@%s:%d", ac.cfg.Mongo.User, ac.cfg.Mongo.Password, ac.cfg.Mongo.Host, ac.cfg.Mongo.Port)
	opts := options.Client().ApplyURI(dbUri)
	err := opts.Validate()
	if err != nil {
		logger.Fatalf("invalid database config: %s", err)
	}
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Fatalf("cannot connect to mongo database: %s", err)
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatalf("cannot connect to mongo database: %s", err)
	}

	db := mongoClient.Database(ac.cfg.Mongo.Database)
	ac.database = db
}

func (ac *AppContext) Run() error {
	return ac.websocketServer.Run()
}
