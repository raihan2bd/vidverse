package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raihan2bd/vidverse/config"
	"github.com/raihan2bd/vidverse/helpers"
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

func (m *Repo) generateThumbURL(publicID string) string {
	return fmt.Sprintf("https://res.cloudinary.com/%s/video/upload/%s.jpeg", initializers.CLD.Config.Cloud.CloudName, publicID)
}

func (m *Repo) HandleCreateVideo(c *gin.Context) {
	// authorization
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access denied! Please login first",
		})
		return
	}

	userID := uint(user_id.(float64))
	var user *models.User
	user, err := m.App.DBMethods.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access denied! Please login first",
		})
		return
	}

	// check the user role
	if user.UserRole != "admin" {
		if user.UserRole != "author" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied! You are not allowed to upload video",
			})
			return
		}
	}

	videoFile, fileInfo, err := c.Request.FormFile("video")
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": "File is required."})
		return
	} else if fileInfo == nil {
		defer videoFile.Close()
		c.IndentedJSON(400, gin.H{"error": "video is required."})
		return
	}

	var thumbSecureURL, thumbPublicID string
	thumbFile, thumbFileInfo, err := c.Request.FormFile("thumb")
	if err != nil {
	} else if thumbFileInfo != nil {
		defer thumbFile.Close()
	}

	validator := validator.New()

	// validate video
	validator.IsVideo(fileInfo.Header.Get("Content-Type"), "video")
	validator.IsVideoSize(fileInfo.Size, 100*1024*1024, "video")

	// Todo: add image upload system later

	title := c.PostForm("title")
	description := c.PostForm("description")
	channel_id := c.PostForm("channel_id")

	channelID, err := strconv.Atoi(channel_id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid channel id",
		})
		return
	}

	validator.Required(title, "title", "title is required.")
	validator.IsLength(title, "title", 5, 255)
	validator.Required(description, "description", "description is required")
	validator.IsLength(description, "description", 25, 500)
	validator.Required(channel_id, "channel_id", "channel_id is required")

	if thumbFileInfo != nil && thumbFile != nil {
		validator.IsImage(thumbFileInfo.Header.Get("Content-Type"), "thumb")
		validator.IsImageSize(thumbFileInfo.Size, 5*1024*1024, "thumb")
	}

	if !validator.Valid() {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": validator.GetErrMsg(),
		})
		return
	}

	// check the channel is available or not
	var channel *models.CustomChannelDTO
	channel, err = m.App.DBMethods.GetChannelByID(channelID)
	if err != nil || channel.ID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "The channel you want to upload video is not found!",
		})
		return
	}

	// check if the channel user is the same or not
	if channel.UserID != userID {
		if user.UserRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied! You are not allowed to upload video to this channel",
			})
			return
		}
	}

	// upload video to cloudinary
	ctx := context.Background()
	var secureURL, videoPublicID string
	secureURL, videoPublicID, err = helpers.UploadVideoToCloudinary(ctx, m.App.CLD, videoFile)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload the video",
		})
		return
	}

	if thumbFileInfo != nil && thumbFile != nil {
		// upload thumb to cloudinary
		thumbSecureURL, thumbPublicID, err = helpers.UploadImageToCloudinary(ctx, m.App.CLD, thumbFile, "vidverse/uploads/thumbs")
		if err != nil {
			thumbSecureURL = m.generateThumbURL(videoPublicID)
		}
	} else {
		// generate thumb url
		thumbSecureURL = m.generateThumbURL(videoPublicID)
	}

	video := models.Video{Title: title, Description: description, PublicID: videoPublicID, SecureURL: secureURL, ChannelID: channel.ID, Thumb: thumbSecureURL, ThumbPublicID: thumbPublicID}

	videoID, err := m.App.DBMethods.CreateVideo(&video)
	if err != nil {
		// delete thumbnail from cloudinary
		_ = helpers.DeleteImageFromCloudinary(ctx, m.App.CLD, thumbPublicID)
		// delete video from cloudinary
		_ = helpers.DeleteImageFromCloudinary(ctx, m.App.CLD, videoPublicID)

		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create the video",
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"message":  "Successfully created the video",
		"video_id": videoID,
	})
}

// handle update video
func (m *Repo) HandleUpdateVideo(c *gin.Context) {
	// authorization
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access denied! Please login first",
		})
		return
	}

	userID := uint(user_id.(float64))
	var user *models.User
	user, err := m.App.DBMethods.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access denied! Please login first",
		})
		return
	}

	// check the user role
	if user.UserRole != "admin" {
		if user.UserRole != "author" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied! You are not allowed to upload video",
			})
			return
		}
	}

	videoID, err := strconv.Atoi(c.Params.ByName("videoID"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "404 video not found! Invalid video id",
		})
		return
	}

	// check the video is available or not
	var video *models.Video
	video, err = m.App.DBMethods.GetVideoByID(videoID)
	if err != nil || video.ID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	videoFile, fileInfo, _ := c.Request.FormFile("video")
	if fileInfo != nil && videoFile != nil {
		defer videoFile.Close()
	}

	thumbFile, thumbFileInfo, _ := c.Request.FormFile("thumb")
	if thumbFileInfo != nil && thumbFile != nil {
		defer thumbFile.Close()
	}

	var thumbUrl, thumbPublicID, videoUrl, videoPublicID string

	// validate form data
	validator := validator.New()
	title := c.PostForm("title")
	description := c.PostForm("description")

	if title != "" {
		validator.IsLength(title, "title", 5, 255)
	} else {
		title = video.Title
	}
	if description != "" {
		validator.IsLength(description, "description", 25, 500)
	} else {
		description = video.Description
	}

	if thumbFileInfo != nil && thumbFile != nil {
		validator.IsImage(thumbFileInfo.Header.Get("Content-Type"), "thumb")
		validator.IsImageSize(thumbFileInfo.Size, 5*1024*1024, "thumb")
	}

	if videoFile != nil && fileInfo != nil {
		validator.IsVideo(fileInfo.Header.Get("Content-Type"), "video")
		validator.IsVideoSize(fileInfo.Size, 100*1024*1024, "video")
	}

	if !validator.Valid() {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": validator.GetErrMsg(),
		})
		return
	}

	// upload video to cloudinary if video file is available
	ctx := context.Background()
	if videoFile != nil && fileInfo != nil {
		videoUrl, videoPublicID, err = helpers.UploadVideoToCloudinary(ctx, m.App.CLD, videoFile)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to upload the video",
			})
			return
		}
	} else {
		videoUrl = video.SecureURL
		videoPublicID = video.PublicID
	}

	// upload thumb to cloudinary if thumb file is available
	if thumbFileInfo != nil && thumbFile != nil {
		thumbUrl, thumbPublicID, err = helpers.UploadImageToCloudinary(ctx, m.App.CLD, thumbFile, "vidverse/uploads/thumbs")
		if err != nil {
			if thumbUrl == "" && thumbPublicID == "" && videoPublicID != "" && videoUrl != "" {
				thumbUrl = m.generateThumbURL(videoPublicID)
			} else {
				thumbUrl = video.Thumb
			}
		}
	} else {
		thumbUrl = video.Thumb
	}

	var oldVideoPublicID, oldThumbPublicID string = video.PublicID, video.ThumbPublicID

	// update video
	video.Title = title
	video.Description = description
	video.SecureURL = videoUrl
	video.Thumb = thumbUrl
	video.ThumbPublicID = thumbPublicID
	video.PublicID = videoPublicID

	err = m.App.DBMethods.UpdateVideo(video)
	if err != nil {
		// delete thumbnail from cloudinary
		if thumbPublicID != video.ThumbPublicID {
			_ = helpers.DeleteImageFromCloudinary(ctx, m.App.CLD, thumbPublicID)
		}

		if videoPublicID != video.PublicID {
			_ = helpers.DeleteVideoFromCloudinary(ctx, m.App.CLD, videoPublicID)
		}

		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update the video",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":  "Successfully updated the video",
		"video_id": video.ID,
	})

	// delete old thumbnail from cloudinary
	if oldThumbPublicID != "" && (thumbPublicID != oldThumbPublicID) {
		_ = helpers.DeleteImageFromCloudinary(ctx, m.App.CLD, oldThumbPublicID)
	}

	if oldVideoPublicID != "" && (videoPublicID != oldVideoPublicID) {
		_ = helpers.DeleteVideoFromCloudinary(ctx, m.App.CLD, oldVideoPublicID)
	}

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
		ID:            video.Channel.ID,
		Title:         video.Channel.Title,
		Description:   video.Channel.Description,
		Logo:          video.Channel.Logo,
		UserID:        video.Channel.UserID,
		Subscriptions: video.Channel.Subscriptions,
	}
	// check the user is logged in or not
	var isLiked bool
	userID, ok := c.Get("user_id")
	if !ok {
		channel.IsSubscribed = false
		isLiked = false
	} else {
		userIDUint := uint(userID.(float64))
		channel.IsSubscribed = m.App.DBMethods.IsSubscribed(userIDUint, video.Channel.ID)
		_, err := m.App.DBMethods.GetLikeByVideoIDAndUserID(video.ID, userIDUint)
		if err != nil {
			isLiked = false
		} else {
			isLiked = true
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"title":       video.Title,
		"description": video.Description,
		"id":          video.ID,
		"vid_src":     video.SecureURL,
		"channel":     channel,
		"likes":       len(video.Likes),
		"views":       video.Views,
		"thumb":       video.Thumb,
		"is_liked":    isLiked,
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

	http.ServeFile(c.Writer, c.Request, SecureURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
}

// Delete video
func (m *Repo) HandleDeleteVideo(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("videoID"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "404 page not found!",
		})
		return
	}

	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access denied! Please login first",
		})
		return
	}

	userID := uint(user_id.(float64))
	var user *models.User
	user, err = m.App.DBMethods.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access denied! Please login first",
		})
		return
	}

	// check the video is available or not
	var video *models.Video
	video, err = m.App.DBMethods.GetVideoByID(id)
	if err != nil || video.ID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	// check the user role
	if user.UserRole != "admin" {
		if user.UserRole != "author" || video.Channel.UserID != user.ID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied! You are not allowed to delete video",
			})
			return
		}
	}

	// delete video from database
	err = m.App.DBMethods.DeleteVideoModel(video)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete the video",
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

// Handle Video Like
func (m *Repo) HandleVideoLike(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access denied!",
		})
		return
	}

	userID := uint(user_id.(float64))

	// check the user is available or not
	user, err := m.App.DBMethods.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access denied!",
		})
		return
	}

	// Get Video ID
	videoID, err := strconv.Atoi(c.Params.ByName("videoID"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	// Check the video is available or not
	video, err := m.App.DBMethods.GetVideoByID(videoID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
		return
	}

	if video.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "404 video not found!",
		})
	}

	// Check the video is already liked by the user or not
	var like *models.Like
	like, err = m.App.DBMethods.GetLikeByVideoIDAndUserID(uint(videoID), userID)

	if err != nil {
		var newLike = models.Like{
			UserID:  userID,
			VideoID: uint(videoID),
		}

		// Create a new like
		id, err := m.App.DBMethods.CreateLike(&newLike)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to like the video",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully liked the video",
			"like_id": id,
		})

		if video.Channel.UserID == user.ID {
			return
		}

		// send notification to the video owner
		ownerID := video.Channel.UserID
		notification := models.Notification{LikeID: id, Type: "like", VideoID: uint(videoID), SenderID: user.ID, ReceiverID: ownerID, SenderName: user.Name, IsRead: false}

		nID, err := m.App.DBMethods.CreateNotification(&notification)

		if err == nil {
			// send notification to the user
			notification.SenderAvatar = user.Avatar
			notification.Thumb = video.Thumb
			notification.ID = nID
			m.App.NotificationChan <- &config.NotificationEvent{BroadcasterID: ownerID, Action: "a_new_notification", Data: &notification}
		}

		return
	}

	// Delete the like
	err = m.App.DBMethods.DeleteLikeByID(like.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to unlike the video",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully unliked the video",
	})

}
