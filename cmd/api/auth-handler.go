package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raihan2bd/vidverse/models"
	validator "github.com/raihan2bd/vidverse/validators"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) LoginHandler(c *gin.Context) {
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

	user, _ := app.DB.GetUserByEmail(payload.Email)

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

	// Generate Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.ID,
		"user_role": user.UserRole,
		"user_name": user.UserName,
		"exp":       time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte("My-Secret"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create token",
		})

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

	// send it as a response
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged in",
		"token":   tokenString,
	})

}

func (app *application) HandleMyAuthInfo(c *gin.Context) {

}

func (app *application) SignupHandler(c *gin.Context) {
	var payload models.UserPayload

	if err := c.BindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payload inputs",
		})
		log.Println(err.Error())
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	payload.UserName = strings.TrimSpace(strings.ToLower(payload.UserName))

	// initialize the validator
	v := validator.New()

	// validate user name
	v.Required(payload.Name, "name", "Name is Required")
	v.IsLength(payload.Name, "name", 3, 100)
	v.IsValidFullName(payload.Name, "name")

	// validate username
	v.IsLength(payload.UserName, "username", 5, 100)
	regex := regexp.MustCompile(`^[a-z][a-z0-9]*$`)
	if !regex.MatchString(payload.UserName) {
		v.AddError("username", "Username must start with a letter and can contain letters or numbers only.")
	}

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

	// check username is exist
	user, _ := app.DB.GetUserByUsername(payload.UserName)
	if user != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "username is already taken. please try another one",
		})
		return
	}

	// check email is exist
	user, _ = app.DB.GetUserByEmail(payload.Email)
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
		UserName: payload.UserName,
		Password: string(hash),
		Email:    payload.Email,
	}

	id, err := app.DB.CreateNewUser(&newUser)
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
