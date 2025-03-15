package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"money-tracker-backend/internal/interfaces"
	"money-tracker-backend/internal/models"
	"money-tracker-backend/internal/utils"
)

type UserController struct {
	userService interfaces.UserServiceInterface
}

func NewUserController(userService interfaces.UserServiceInterface) interfaces.UserControllerInterface {
	return &UserController{userService: userService}
}

func (userController *UserController) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload"))
		return
	}

	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
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
