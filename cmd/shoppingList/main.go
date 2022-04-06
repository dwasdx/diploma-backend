package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"shopingList/api"
	auth2 "shopingList/api/auth"
	"shopingList/api/controllers"
	"shopingList/api/controllers/sync"
	"shopingList/api/controllers/users"
	"shopingList/cmd/shoppingList/pkg"
	"shopingList/pkg/events"
	"shopingList/pkg/listeners"
	"shopingList/pkg/models"
	"shopingList/pkg/readModels"
	"shopingList/pkg/repositories"
	"shopingList/pkg/services"
	"shopingList/pkg/services/login_limiter"
	"shopingList/pkg/services/sms"
	"shopingList/store/mysql"
)

var (
	configPath = flag.String("c", "./config.json", "path to config file")
	config     models.AppConfig
)

func main() {
	flag.Parse()
	loadConfig()
	logInit(config.LogLevel)

	applicationStopped := make(chan struct{})
	authenticator := auth2.NewService(auth2.Opts{
		SigningKey: []byte("cc77ec27b21f31f7de3bd0dc2652b5f99g34ac28e38ecfa7c06f4737d7c809a1"),
	})

	db, err := openDb(config.Database)
	if err != nil {
		log.Fatalln("[ERROR]: can't open mysql db: ", err)
		return
	}

	dataService := mysql.NewDataStore(db)

	var smsService sms.SmsService
	var codeGenerator models.UserAuthCodeGenerator

	if config.SmsEnabled {
		if _, err := config.SmsAeroConfig.Validate(); err != nil {
			log.Fatalln("SmsAeroConfig is wrong", err)
		}

		smsService = createSmsService(config.SmsAeroConfig)
		codeGenerator = services.AuthCodeGenerator{}
	} else {
		smsService = nil
		codeGenerator = services.FixedCodeGenerator{}
	}

	tokenStorage := mysql.NewFCMTokenStorage(db)
	pushChannel := make(chan services.PushNotificationMessage)

	if config.HasFirebaseCredentials() {
		pushService, err := services.NewPushService(config.FirebaseCredentialsFile, tokenStorage, false)
		if err != nil {
			log.Fatalln("[ERROR]: can't start pushService: ", err)
		}

		go pushService.ListenChannel(pushChannel)
	}

	restServer := api.New(authenticator, dataService, config.Server.Port, applicationStopped)

	// Репозиторий для уведомлений
	notificationRepository := repositories.NewNotificationsRepository(db)
	notificationReadRepository := readModels.NewNotificationsReadRepository(db)

	// Channels для listeners
	chanGoodsChange := make(chan events.GoodsChangeEvent)
	goodChangeListener := listeners.GoodChangeListener{Repository: notificationRepository, PushChannel: pushChannel}
	go goodChangeListener.Run(chanGoodsChange)

	chanShareChange := make(chan events.ShareListEvent)
	shareChangeListener := listeners.ShareListChangeListener{Repository: notificationRepository, PushChannel: pushChannel}
	go shareChangeListener.Run(chanShareChange)

	// Публичные контроллеры
	publicController := controllers.NewPublic()

	loginLimiter := getRedisLoginLimiter(config.LoginLimiterConfig, config.RedisConfig)

	userController := users.NewUserController(authenticator, dataService, smsService, loginLimiter, codeGenerator)
	userController.SetDebugPhones(config.DebugPhones)
	restServer.AddPublicRoutes(publicController.Routes()...)
	restServer.AddPublicRoutes(userController.Routes()...)

	// Контроллеры под авторизацией
	privateController := controllers.NewPrivate(dataService)
	syncController := sync.NewSyncController(authenticator, dataService, chanGoodsChange, chanShareChange)
	tokenController := controllers.NewFCMTokenController(authenticator, tokenStorage)
	sharedListController := controllers.NewSharedListsController(authenticator, dataService)
	sharedListController.ChanShareChange = chanShareChange
	refbookController := controllers.NewRefbookController(
		repositories.NewRefbookCategoriesRepository(db),
		repositories.NewRefbookProductsRepository(db))

	notificationController := controllers.NewNotificationController(
		authenticator, notificationRepository, notificationReadRepository)

	if config.HasFirebaseCredentials() {
		notificationController.PushChannel = pushChannel
	}

	restServer.AddPrivateRoutes(privateController.Routes()...)
	restServer.AddPrivateRoutes(syncController.Routes()...)
	restServer.AddPrivateRoutes(notificationController.Routes()...)
	restServer.AddPrivateRoutes(tokenController.Routes()...)
	restServer.AddPrivateRoutes(refbookController.Routes()...)
	restServer.AddPrivateRoutes(sharedListController.Routes()...)

	tgListener, err := pkg.CreateTgListener(config.TelegramBotToken, db)
	if err != nil {
		log.Errorln(errors.Wrap(err, "Error create telegram listener"))
	} else {
		go tgListener.Run()
	}

	go restServer.Run()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		restServer.Stop()
		close(applicationStopped)
		log.Info("Server stopped")
	}()
	<-applicationStopped
}

func loadConfig() {
	config = LoadConfigFromPath(*configPath)
}

func LoadConfigFromPath(configPath string) models.AppConfig {
	config := models.AppConfig{}

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalln("[ERROR]: can't open config file: ", err)
	}
	defer file.Close() // nolint: errcheck, gosec - not critic here
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatalln("[ERROR]: can't parse config file: ", err) // nolint: gocritic
	}

	return config
}

func openDb(config models.DatabaseConfig) (*sql.DB, error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.User, config.Password, config.Address, config.DbName)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatalln("[ERROR]: can't open mysql db: ", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func createSmsService(c models.SmsAeroConfig) sms.SmsService {
	log.Info("init sms service")
	service, err := sms.NewSmsService(c.Email, c.ApiKey, c.Sign, c.Channel, c.TestMode)

	if err != nil {
		log.Fatalln("[ERROR]: can't init sms service: ", err)
		return nil
	}

	return service
}

func logInit(logLevel string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	if logLevel == "" {
		logLevel = "info"
		log.Warnf("Log level not specified. Setting = %v.", logLevel)
	}

	// parse string, this is built-in feature of logrus
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		lvl = log.DebugLevel
		log.Errorf("Bad format logLevel: %v. Set %v", logLevel, lvl)
	}

	// set global log level
	log.SetLevel(lvl)
}

func getRedisLoginLimiter(limCfg models.LoginLimiterConfig, redisConfig models.RedisConfig) users.LoginLimiter {
	redisStorage := login_limiter.NewRedisStorage(redisConfig.Address, redisConfig.Password, redisConfig.DB)

	return login_limiter.NewLoginLimiter(redisStorage, limCfg.SeqLimitSeconds, limCfg.DailyLimitCount)
}
