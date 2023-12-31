package config

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gorilla/websocket"
	"github.com/raihan2bd/vidverse/initializers"
	"github.com/raihan2bd/vidverse/repository"
	dbrepo "github.com/raihan2bd/vidverse/repository/dbRepo"
	"gorm.io/gorm"
)

type Application struct {
	DB               *gorm.DB
	CLD              *cloudinary.Cloudinary
	DBMethods        repository.DatabaseRepo
	NotificationChan chan *NotificationEvent
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

	return &Application{
		DB:               db,
		DBMethods:        dbrepo.NewPostgresRepo(initializers.DB, initializers.CLD),
		CLD:              cld,
		NotificationChan: make(chan *NotificationEvent),
	}, nil
}
