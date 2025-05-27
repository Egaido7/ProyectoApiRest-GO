package main

import (
	"fmt"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	//framework
	"parte3/api"
)

// Custom validation function for regexp
func regexpValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	regex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return regex.MatchString(value)
}
func main() {
	r := gin.Default()

	// Register custom validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("regexp", regexpValidation)
	}

	api.InitRoutes(r)

	if err := r.Run(":8080"); err != nil {
		panic(fmt.Errorf("error trying to start server: %v", err))
	}
}
