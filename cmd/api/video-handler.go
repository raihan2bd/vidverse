package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/initializers"
	"github.com/raihan2bd/vidverse/models"
	validator "github.com/raihan2bd/vidverse/validators"
)

func UploadVideo(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	validator := validator.New()

	title := c.PostForm("title")
	description := c.PostForm("description")

	validator.Required(title, "title", "title is required.")
	validator.IsLength(title, "title", 5, 255)
	validator.Required(description, "description", "description is required")
	validator.IsLength(description, "description", 25, 500)

	if !validator.Valid() {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": validator.GetErrMsg(),
		})
		return
	}

	// Define the upload directory
	folder := "vidverse/uploads/videos"

	ctx := context.Background()

	// Upload the file to Cloudinary with specified folder and transformations
	resp, err := initializers.CLD.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	video := models.Video{Title: title, Description: description, PublicID: resp.PublicID, SecureURL: resp.SecureURL}

	// user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&video)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to upload the video",
		})

		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully", "video_id": video.ID})
}

func GetSingleVideo(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("videoID"))
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	var video models.Video
	if id > 0 {
		initializers.DB.First(&video, "id = ?", id)

		if video.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "404 video not found!",
			})
			return
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"title":       video.Title,
		"description": video.Description,
		"id":          video.ID,
		"vid_src":     fmt.Sprintf("/api/v1/file/video/%d", video.ID),
	})
}

func StreamVideo(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("videoID"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	var SecureURL string
	result := initializers.DB.Table("videos").Select("secure_url").Where("id = ?", id).Scan(&SecureURL)
	fmt.Println(SecureURL)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	retriveVideo, err := http.Get(SecureURL)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 video not found!"})
		return
	}

	defer retriveVideo.Body.Close()
	c.Header("Content-Type", "video/mp4")
	buffSize := 1024 * 1024
	buff := make([]byte, buffSize)

	for {
		n, err := retriveVideo.Body.Read(buff)
		if err != nil && err != io.EOF {
			// Handle the error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream video"})
			return
		}
		if n == 0 {
			break
		}
		c.Writer.Write(buff[:n])
		c.Writer.Flush()
	}
}

// Delete video
func DeleteVidoe(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("videoID"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "404 page not found!",
		})
		return
	}

	// Check the videoID is available or not
	var publicID string
	result := initializers.DB.Table("videos").Select("public_id").Where("id = ?", id).Scan(&publicID)

	if result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "404 the video you want to delete is not found!",
		})
		return
	}

	fmt.Println(publicID)

	// delete the video

	result = initializers.DB.Unscoped().Delete(&models.Video{}, id)
	if result.Error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "something went wrong. failed to delete the video",
		})
		return
	}

	adminResp, err := initializers.CLD.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID, ResourceType: "video"})

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete the video.",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":  "Successfully deleted the video",
		"cld_resp": adminResp,
	})

}

// func UploadVideo(c *gin.Context) {
// 	file, header, err := c.Request.FormFile("file")
// 	if err != nil {
// 		c.IndentedJSON(400, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer file.Close()

// 	// Define the upload directory
// 	uploadDir := "./uploads/videos"

// 	// Ensure the upload directory is exits
// 	err = os.MkdirAll(uploadDir, os.ModePerm)
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	filename := filepath.Join(uploadDir, header.Filename)

// 	// Create a buffer to store the file data in chunks
// 	buffer := make([]byte, 1024)

// 	// Create a destination file to write the uploaded data
// 	// destinationFile, err := os.Create(uploadDir + header.Filename)
// 	destinationFile, err := os.Create(filename)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer destinationFile.Close()

// 	// Loop through and write the file data in chunks
// 	for {
// 		bytesRead, err := file.Read(buffer)
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}
// 		destinationFile.Write(buffer[:bytesRead])
// 	}

// 	c.JSON(200, gin.H{"message": "File uploaded successfully", "video_id": 1, "video_url": "/video"})
// }

// Upload file with cloudinary
