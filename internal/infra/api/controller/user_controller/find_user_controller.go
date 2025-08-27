package user_controller

import (
	"l03/configuration/rest_err"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (uc *UserController) FindUserById(c *gin.Context) {
	userId := c.Param("user-id")

	if err := uuid.Validate(userId); err != nil {
		errRest := rest_err.NewBadRequestError("invalid user id", rest_err.Causes{
			Field:   "user-id",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	data, err := uc.userUseCase.FindById(c, userId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (uc *UserController) FindUsers(c *gin.Context) {
	data, err := uc.userUseCase.FindUsers(c)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, data)
}
