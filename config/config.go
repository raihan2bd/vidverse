package config

import (
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gorilla/websocket"
	"github.com/raihan2bd/vidverse/initializers"
	"github.com/raihan2bd/vidverse/internal/mail"
	"github.com/raihan2bd/vidverse/repository"
	dbrepo "github.com/raihan2bd/vidverse/repository/dbRepo"
	"gorm.io/gorm"
)

type Application struct {
	DB               *gorm.DB
	CLD              *cloudinary.Cloudinary
	DBMethods        repository.DatabaseRepo
	NotificationChan chan *NotificationEvent
	Mailer           mail.Mail
}

type NotificationEvent struct {
	BroadcasterID uint
	Action        string
	Data          interface{}
	Conn          *websocket.Conn
}

func LoadConfig() (*Application, error) {
	var (
		cld *cloudinary.Cloudinary
		db  *gorm.DB
		err error
	)

	db, err = initializers.ConnectToDB()
	if err != nil {
		return nil, err
	}

	cld, err = initializers.ConnectToCloudinary()
	if err != nil {
		return nil, err
	}

	err = initializers.SyncDatabase()
	if err != nil {
		return nil, err
	}

	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := mail.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        mailPort,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRIPTION"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
	}

	return &Application{
		DB:               db,
		DBMethods:        dbrepo.NewPostgresRepo(initializers.DB, initializers.CLD),
		CLD:              cld,
		NotificationChan: make(chan *NotificationEvent),
		Mailer:           m,
	}, nil
}
