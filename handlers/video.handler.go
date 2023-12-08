package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/initializers"
	"github.com/raihan2bd/vidverse/models"
	validator "github.com/raihan2bd/vidverse/validators"
)

func (m *Repo) HandleGetAllVideos(c *gin.Context) {
	// search query
	searchQuery := c.DefaultQuery("search", "")

	// pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid page number",
		})
		return
	}

	// limit
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "24"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit number",
		})
		return
	}

	// Get all videos
	videos, count, err := m.App.DBMethods.GetAllVideos(page, limit, searchQuery)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// has next page
	hasNextPage := false
	if count > int64(page*limit) {
		hasNextPage = true
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"videos":        videos,
		"page":          page,
		"limit":         limit,
		"has_next_page": hasNextPage,
	})
}

func (m *Repo) HandleCreateVideo(c *gin.Context) {
	videoFile, fileInfo, err := c.Request.FormFile("video")
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(400, gin.H{"error": "File is required."})
		return
	}
	defer videoFile.Close()

	validator := validator.New()

	// validate video
	validator.IsVideo(fileInfo.Header.Get("Content-Type"), "video")
	validator.IsVideoSize(fileInfo.Size, 100*1024*1024, "video")

	// Todo: add image upload system later

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
	resp, err := initializers.CLD.Upload.Upload(ctx, videoFile, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	video := models.Video{Title: title, Description: description, PublicID: resp.PublicID, SecureURL: resp.SecureURL, ChannelID: 1}

	// user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&video)

	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to upload the video",
		})

		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully", "video_id": video.ID})
}

func (m *Repo) HandleGetSingleVideo(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("videoID"))
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{
			"error": "Invalid ID",
		})
		return
	}
	var video *models.Video
	if id > 0 {
		video, err = m.App.DBMethods.GetVideoByID(id)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	var channel = &models.ChannelPayload{
		ID:          video.Channel.ID,
		Title:       video.Channel.Title,
		Description: video.Channel.Description,
		Logo:        video.Channel.Logo,
		UserID:      video.Channel.UserID,
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"title":       video.Title,
		"description": video.Description,
		"id":          video.ID,
		"vid_src":     video.SecureURL,
		"channel":     channel,
		"likes":       len(video.Likes),
		"views":       video.Views,
	})

}

// Get related videos
func (m *Repo) HandleGetRelatedVideos(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("channelID"))
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	var videos []models.VideoDTO
	videos, _, err = m.App.DBMethods.GetVideosByChannelID(id, 1, 24)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	if len(videos) == 0 {
		videos, _, err = m.App.DBMethods.GetAllVideos(1, 24, "")
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"videos": videos,
	})
}

func (m *Repo) StreamVideoBuff(c *gin.Context) {
	filename := filepath.Join("./uploads/videos", "./test.mp4")
	videoFile, err := os.Open(filename)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening video file")
		return
	}
	defer videoFile.Close()

	fileInfo, err := videoFile.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting video file information")
		return
	}

	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		c.String(http.StatusBadRequest, "Range header not provided")
		return
	}

	parts := strings.SplitN(rangeHeader, "=", 2)
	if len(parts) != 2 || parts[0] != "bytes" {
		c.String(http.StatusBadRequest, "Invalid Range header format")
		return
	}

	byteRange := parts[1]
	byteRanges := strings.SplitN(byteRange, "-", 2)
	start, end := byteRanges[0], byteRanges[1]

	startPos := 0
	endPos := int(fileInfo.Size()) - 1

	if start != "" {
		startPos, _ = strconv.Atoi(start)
	}
	if end != "" {
		endPos, _ = strconv.Atoi(end)
	}

	contentLength := endPos - startPos + 1

	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startPos, endPos, fileInfo.Size()))
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))
	c.Header("Content-Type", "video/mp4")

	c.Status(http.StatusPartialContent)

	videoFile.Seek(int64(startPos), 0)
	bufSize := 258 * 1024
	buf := make([]byte, bufSize)
	readSize := 0

	for readSize < contentLength {
		n := contentLength - readSize
		if n > bufSize {
			n = bufSize
		}
		n, err := videoFile.Read(buf[:n])
		if err != nil {
			break
		}
		c.Writer.Write(buf[:n])
		readSize += n
	}
}

// func (m *Repo) StreamVideo(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Params.ByName("videoID"))
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "404 video not found!",
// 		})
// 		return
// 	}

// 	var SecureURL string
// 	result := initializers.DB.Table("videos").Select("secure_url").Where("id = ?", id).Scan(&SecureURL)
// 	fmt.Println(SecureURL)

// 	if result.Error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "404 video not found!",
// 		})
// 		return
// 	}

// 	retriveVideo, err := http.Get(SecureURL)

// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "404 video not found!"})
// 		return
// 	}

// 	defer retriveVideo.Body.Close()
// 	buffSize := 128 * 1024
// 	buff := make([]byte, buffSize)

// 	for {
// 		n, err := retriveVideo.Body.Read(buff)
// 		if err != nil && err != io.EOF {
// 			// Handle the error
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream video"})
// 			return
// 		}
// 		if n == 0 {
// 			break
// 		}
// 		c.Writer.Write(buff[:n])
// 		c.Writer.Flush()
// 	}
// }

func (m *Repo) StreamVideo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("videoID")) // Use c.Param instead of c.Params.ByName
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	var SecureURL string
	result := initializers.DB.Table("videos").Select("secure_url").Where("id = ?", id).Scan(&SecureURL)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	if SecureURL == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video URL not found!",
		})
		return
	}

	// Log the SecureURL for debugging purposes
	fmt.Println("SecureURL:", SecureURL)

	http.ServeFile(c.Writer, c.Request, SecureURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		// You should also log the error for debugging
		fmt.Println("ServeFile error:", err)
		return
	}
}

// Delete video
func (m *Repo) HandleDeleteVidoe(c *gin.Context) {
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

	_, err = initializers.CLD.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID, ResourceType: "video"})

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete the video.",
		})
		return
	}

	// delete the video
	err = m.App.DBMethods.DeleteVideoByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "something went wrong. failed to delete the video",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Successfully deleted the video",
	})

}

// func (m *Repo) UploadVideo(c *gin.Context) {
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

// Get comments
func (m *Repo) HandleGetComments(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("videoID"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "404 video not found!",
			"errors": err,
		})
		return
	}

	var page, limit int
	page, err = strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	limit, err = strconv.Atoi(c.DefaultQuery("limit", "24"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	var comments []models.CommentDTO
	var count int64
	comments, count, err = m.App.DBMethods.GetCommentsByVideoID(id, page, limit)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	var hasNextPage bool

	if count > int64(page*limit) {
		hasNextPage = true
	} else {
		hasNextPage = false
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"page":          page,
		"comments":      comments,
		"has_next_page": hasNextPage,
	})
}

// get videos by channelID with pagination
func (m *Repo) HandleGetVideosByChannelID(c *gin.Context) {
	chanID, err := strconv.Atoi(c.Params.ByName("channelID"))

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "404 Channel not found!",
		})
		return
	}

	var page, limit int
	page, err = strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	var videos []models.VideoDTO
	var count int64
	videos, count, err = m.App.DBMethods.GetVideosByChannelID(chanID, page, limit)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	var hasNextPage bool

	if count > int64(page*limit) {

		hasNextPage = true
	} else {
		hasNextPage = false
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"page":          page,
		"videos":        videos,
		"has_next_page": hasNextPage,
	})
}
