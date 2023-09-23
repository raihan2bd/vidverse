package initializers

import (
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var CLD *cloudinary.Cloudinary

func ConnectToCloudinary() {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLD_URI"))

	if err != nil {
		panic("failed to intialize Cloudinary")
	}
	CLD = cld
}
