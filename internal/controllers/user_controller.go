package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/services"
	"money-tracker-backend/internal/utils"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (userController *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload"))
		return
	}

	createdUser, err := userController.userService.CreateUser(&user)
	if err != nil {
		if utils.IsUniqueViolation(err) {
			response := utils.ErrorResponse(err.Error())
			c.JSON(http.StatusConflict, response)
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("User created successfully", createdUser))
}

func (userController *UserController) GetUser(c *gin.Context) {
	email, _ := c.Get("email")
	users, err := userController.userService.GetUser(email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Error fetching users"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Users retrieved successfully", users))
}

func (userController *UserController) SearchUser(c *gin.Context) {
	var payload struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload"))
		return
	}

	user, err := userController.userService.SearchByEmail(payload.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("User found", user))
}
