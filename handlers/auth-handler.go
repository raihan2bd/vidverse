package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raihan2bd/vidverse/helpers"
	"github.com/raihan2bd/vidverse/internal/mail"
	"github.com/raihan2bd/vidverse/models"
	validator "github.com/raihan2bd/vidverse/validators"
	"golang.org/x/crypto/bcrypt"
)

func (m *Repo) LoginHandler(c *gin.Context) {
	// Get user credentials from req body
	type UserCreds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var payload UserCreds

	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	user, _ := m.App.DBMethods.GetUserByEmail(payload.Email)

	if user == nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	if user.ID <= 0 {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare the password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	exp := time.Now().Add(time.Hour * 24 * 7).Unix()

	// Generate Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.ID,
		"user_role": user.UserRole,
		"user_name": user.Name,
		"exp":       exp,
	})

	jwtSecret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create token",
		})

		return
	}

	var userResponse struct {
		ID       uint   `json:"id"`
		Username string `json:"user_name"`
		UserRole string `json:"user_role"`
		Avatar   string `json:"avatar"`
	}

	userResponse.ID = user.ID
	userResponse.UserRole = user.UserRole
	userResponse.Avatar = user.Avatar
	userResponse.Username = user.Name

	// send it as a response
	c.JSON(http.StatusOK, gin.H{
		"user":       userResponse,
		"token":      tokenString,
		"expires_at": exp,
	})

}

func (m *Repo) HandleMyAuthInfo(c *gin.Context) {
}

func (m *Repo) SignupHandler(c *gin.Context) {
	fmt.Println("SignupHandler")
	var payload models.UserPayload

	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payload inputs",
		})
		log.Println(err.Error())
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)

	// initialize the validator
	v := validator.New()

	// validate user name
	v.Required(payload.Name, "name", "Name is Required")
	v.IsLength(payload.Name, "name", 3, 100)
	v.IsValidFullName(payload.Name, "name")

	// validate email
	v.Required(payload.Email, "email", "Email is required")
	v.IsEmail(payload.Email, "email", "Invalid email address")

	// validate password
	v.IsLength(payload.Password, "password", 6, 255)
	v.IsValidPassword(payload.Password, "password")

	if !v.Valid() {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"errors": v.Errors,
		})
		return
	}

	// // check username is exist
	// user, _ := m.App.DBMethods.GetUserByUsername(payload.UserName)
	// if user != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{
	// 		"error": "username is already taken. please try another one",
	// 	})
	// 	return
	// }

	// check email is exist
	user, _ := m.App.DBMethods.GetUserByEmail(payload.Email)
	if user != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "email address is already exist. please try another one",
		})
		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong. please try again later.",
		})
		return
	}

	// save the info into the database
	newUser := models.User{
		Name:     payload.Name,
		Password: string(hash),
		Email:    payload.Email,
	}

	id, err := m.App.DBMethods.CreateNewUser(&newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "You account successfully created!. Login into your account now!",
		"id":      id,
	})
}

func (m *Repo) RequestForgotPassword(c *gin.Context) {
	// Get user credentials from req body
	type UserEmail struct {
		Email string `json:"email"`
	}
	var payload UserEmail

	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email",
		})
		return
	}

	v := validator.New()
	v.IsEmail(payload.Email, "email", "Invalid email. Please provide a valid email")

	if !v.Valid() {
		c.IndentedJSON(400, gin.H{
			"error": "Invalid email address. Please provide a valid email.",
		})
		return
	}

	user, err := m.App.DBMethods.GetUserByEmail(payload.Email)
	if err != nil {
		c.IndentedJSON(500, gin.H{
			"error": "Internal server error. Please try again",
		})
		return
	}

	if user == nil || (user.Email != payload.Email) {
		c.IndentedJSON(404, gin.H{
			"error": "The account you are trying to reset is not found!",
		})
		return
	}

	tokenString, err := helpers.GenerateRandomToken(60)
	if err != nil {
		c.IndentedJSON(500, gin.H{
			"error": "Internal server error. Please try again",
		})
		return
	}

	var userToken = models.Token{
		UserID: user.ID,
		Token:  tokenString,
	}

	// add token to the database
	err = m.App.DBMethods.AddForgotPasswordToken(&userToken)
	if err != nil {
		c.IndentedJSON(500, gin.H{
			"error": "Internal server error. Please try again",
		})
		return
	}

	data := map[string]any{
		"UserName":   user.Name,
		"VerifyLink": fmt.Sprintf("%s/verify/%v?token=%s", os.Getenv("APP_DOMAIN"), user.ID, tokenString),
	}

	// send email
	msg := mail.Message{
		From:        m.App.Mailer.FromAddress,
		To:          user.Email,
		Subject:     "Password Reset Request",
		IsResetPass: true,
		DataMap:     data,
	}

	err = m.App.Mailer.SendSmtpMessage(msg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send the email. Please make sure your email is correct."})
		return
	}

	c.IndentedJSON(201, gin.H{
		"message": "An email has been sent to you. Please check your inbox and verify yourself",
	})
}

func (m *Repo) ForgotPassword(c *gin.Context) {
	// type UserPayload struct {
	// 	Token    string `json:"token"`
	// 	Password string `json:"password"`
	// }
	// var payload UserPayload

	// if err := c.BindJSON(&payload); err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{
	// 		"error": "Invalid payload",
	// 	})
	// 	return
	// }

	// // validate token string
	// v := validator.New()
	// v.IsLength(payload.Token, "token", 8, 255)
	// v.IsValidPassword(payload.Password, "password")
	// if !v.Valid() {
	// 	c.IndentedJSON(500, gin.H{
	// 		"error": v.GetErrMsg(),
	// 	})
	// 	return
	// }

	// // check if the token is exist
	// userToken, isValid := m.App.DBMethods.ValidateForgotPasswordToken(payload.Token)
	// if !isValid {
	// 	c.IndentedJSON(403, gin.H{
	// 		"error": "The link is already expired",
	// 	})
	// 	return
	// }

	// // fetch the user
	// user, err := m.App.DBMethods.GetUserByID(userToken.ID)
	// if err != nil || user.ID == 0 {
	// 	c.IndentedJSON(404, gin.H{
	// 		"error": "The user account password you want to change is not found.",
	// 	})
	// 	return
	// }

	// // Update password
	// // hash the password
	// hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	// if err != nil {
	// 	c.IndentedJSON(http.StatusInternalServerError, gin.H{
	// 		"error": "something went wrong. please try again later.",
	// 	})
	// 	return
	// }

	// user.Password = string(hash)
	// err = m.App.DBMethods.UpdateUserPassword(user)
	// if err != nil {
	// 	c.IndentedJSON(http.StatusInternalServerError, gin.H{
	// 		"error": "something went wrong. please try again later.",
	// 	})
	// 	return
	// }

	// // delete token
	// _ = m.App.DBMethods.DeleteUserForgotToken(userToken.ID)

	// c.JSON(200, gin.H{
	// 	"message": "User is verified",
	// 	"user_id": userToken.UserID,
	// })

}

func (m *Repo) SendMail(c *gin.Context) {
	// type mailMessage struct {
	// 	From    string `json:"from"`
	// 	To      string `json:"to"`
	// 	Subject string `json:"subject"`
	// 	Message string `json:"message"`
	// }

	// var payload mailMessage

	// if err := c.BindJSON(&payload); err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{
	// 		"error": "Invalid email",
	// 	})
	// 	return
	// }

	msg := mail.Message{
		From:    m.App.Mailer.FromAddress,
		To:      "raihan2bd.official@gmail.com",
		Subject: "sending mail from vidverse",
		Data:    "This is a dummy message",
	}

	err := m.App.Mailer.SendSmtpMessage(msg)
	if err != nil {
		log.Println("failed to send the mail", err)
	}

	c.JSON(200, gin.H{
		"message": "Email has been send successfully",
	})
}
