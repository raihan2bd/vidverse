package helpers

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raihan2bd/vidverse/config"
	"github.com/raihan2bd/vidverse/models"
)

// Decode the token
func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// Validate the token
func ValidateToken(claims jwt.MapClaims) bool {
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return false
	}

	// check the user_id from the token as well
	return claims["sub"] != 0
}

// Validate user_id and convert to uint
func ValidateAndGetUserByID(app *config.Application, id any) (*models.User, error) {
	userID, ok := id.(uint)
	if !ok {
		return nil, errors.New("invalid User")
	}

	user, err := app.DBMethods.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("invalid User")
	}

	if user.ID <= 0 {
		return nil, errors.New("invalid User")
	}
	user.Password = ""

	return user, nil
}

// Upload image to cloudinary
func UploadImageToCloudinary(ctx context.Context, CLD *cloudinary.Cloudinary, image multipart.File, uploadPath ...string) (string, string, error) {
	folder := "vidverse/uploads/images"
	if len(uploadPath) > 0 {
		folder = uploadPath[0]
	}

	resp, err := CLD.Upload.Upload(ctx, image, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", "", err
	}
	return resp.SecureURL, resp.PublicID, nil
}

// delete image from cloudinary
func DeleteImageFromCloudinary(ctx context.Context, CLD *cloudinary.Cloudinary, publicID string) error {
	result, err := CLD.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID, ResourceType: "image"})
	if err != nil {
		return errors.New("failed to delete image")
	}

	if result.Result != "ok" {
		return errors.New("failed to delete image")
	}

	return nil
}

// Upload video to cloudinary
func UploadVideoToCloudinary(ctx context.Context, CLD *cloudinary.Cloudinary, video multipart.File, uploadPath ...string) (string, string, error) {
	folder := "vidverse/uploads/videos"
	if len(uploadPath) > 0 {
		folder = uploadPath[0]
	}
	resp, err := CLD.Upload.Upload(ctx, video, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", "", err
	}
	return resp.SecureURL, resp.PublicID, nil
}

// delete video from cloudinary
func DeleteVideoFromCloudinary(ctx context.Context, CLD *cloudinary.Cloudinary, publicID string) error {
	_, err := CLD.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID, ResourceType: "video"})
	if err != nil {
		return err
	}
	return nil
}
