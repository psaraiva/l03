package user_controller

import (
	"l03/configuration/rest_err"
	"l03/internal/infra/api/web/validation"
	"l03/internal/usecase/user_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (uc *UserController) Create(c *gin.Context) {
	var userInputDTO user_usecase.UserInputDTO

	if err := c.ShouldBindJSON(&userInputDTO); err != nil {
		restErr := validation.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	outputDTO, err := uc.userUseCase.Create(c, userInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusCreated, outputDTO)
}
