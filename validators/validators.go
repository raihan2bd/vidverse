package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Validator struct {
	Errors map[string]string `json:"errors"`
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) GetErrMsg() string {
	for _, errMsg := range v.Errors {
		return errMsg
	}
	return ""
}

func (v *Validator) Required(data, key, message string) {
	if len(strings.Trim(data, "")) <= 0 {
		v.AddError("key", message)
	}
}

func (v *Validator) IsLength(data, key string, minLength, maxLength int, message ...string) {
	trimData := len(strings.Trim(data, ""))

	if trimData < minLength || trimData > maxLength {
		msg := fmt.Sprintf("%s must be between %d and %d characters ", key, minLength, maxLength)
		if len(message) > 0 {
			msg = message[0]
		}
		v.AddError(key, msg)
	}
}

func (v *Validator) IsEmail(email, key, message string) {
	// A simple regex pattern to validate email
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(email) {
		v.AddError(key, message)
	}
}

// isValidPassword validates the password with the given key
func (v *Validator) IsValidPassword(password, key string) {
	// Set initial flags to false
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
		hasSpace   = false
	)

	// Loop through each character in the password
	for _, char := range password {
		switch {
		// Check if the character is uppercase
		case unicode.IsUpper(char):
			hasUpper = true
		// Check if the character is lowercase
		case unicode.IsLower(char):
			hasLower = true
		// Check if the character is a number
		case unicode.IsNumber(char):
			hasNumber = true
		// Check if the character is a special character
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		// Check if the character is whitespace or tab
		case unicode.IsSpace(char):
			hasSpace = true
		}
	}

	// Check if the password has at least one uppercase letter
	if !hasUpper {
		v.AddError(key, "Password must have at least one uppercase letter")
	}

	// Check if the password has at least one lowercase letter
	if !hasLower {
		v.AddError(key, "Password must have at least one lowercase letter")
	}

	// Check if the password has at least one number
	if !hasNumber {
		v.AddError(key, "Password must have at least one number")
	}

	// Check if the password has at least one special character
	if !hasSpecial {
		v.AddError(key, "Password must have at least one special character")
	}

	// Check if the password contains whitespace or tab characters
	if hasSpace {
		v.AddError(key, "Password cannot contain whitespace or tab characters")
	}
}

// IsValidFullName validates the FullName with the given key
func (v *Validator) IsValidFullName(fullName, key string) {
	// Regular expression pattern for full name
	pattern := `^([A-Za-z])(\s?[A-Za-z])+$`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Check if the name matches the regular expression pattern
	if !regex.MatchString(fullName) {
		v.AddError(key, "Invalid full name. Please enter a valid full name that contains only letters, spaces, and the following characters: ', . -")
	}
}

// Validate Video Type
func (v *Validator) IsVideo(videoType, key string) {
	if videoType != "video/mp4" && videoType != "video/ogg" && videoType != "video/webm" && videoType != "video/3gp" && videoType != "video/mov" && videoType != "video/avi" && videoType != "video/wmv" && videoType != "video/flv" && videoType != "video/mkv" {
		v.AddError(key, "Invalid video type. Please upload a valid video type")
	}
}

// Validate Video Size
func (v *Validator) IsVideoSize(videoSize, max int64, key string) {
	if videoSize > max {
		v.AddError(key, "Invalid video size. Please upload a valid video size. Maximum size is 10MB")
	}
}

// Validate Image Type
func (v *Validator) IsImage(imageType, key string) {
	if imageType != "image/png" && imageType != "image/jpg" && imageType != "image/jpeg" {
		v.AddError(key, "Invalid image type. Please upload a valid image type")
	}
}

// Validate Image Size
func (v *Validator) IsImageSize(imageSize, max int64, key string) {
	if imageSize > max {
		v.AddError(key, "Invalid image size. Please upload a valid image size. Maximum size is 5MB")
	}
}
