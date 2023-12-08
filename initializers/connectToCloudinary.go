package initializers

import (
	"errors"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var CLD *cloudinary.Cloudinary

func ConnectToCloudinary() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLD_URI"))

	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to intialize Cloudinary")
	}
	CLD = cld

	return cld, nil
}
