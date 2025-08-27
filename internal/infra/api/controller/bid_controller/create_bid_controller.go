package bid_controller

import (
	"l03/configuration/rest_err"
	"l03/internal/infra/api/web/validation"
	"l03/internal/usecase/bid_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (bc *BidController) CreateBid(c *gin.Context) {
	var bidInputDTO bid_usecase.BidInputDTO

	if err := c.ShouldBindJSON(&bidInputDTO); err != nil {
		restErr := validation.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	outputDTO, err := bc.BidUseCase.Create(c, bidInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusCreated, outputDTO)
}
