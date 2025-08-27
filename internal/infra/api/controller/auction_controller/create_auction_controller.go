package auction_controller

import (
	"l03/configuration/rest_err"
	"l03/internal/infra/api/web/validation"
	"l03/internal/usecase/auction_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ac *AuctionController) CreateAuction(c *gin.Context) {
	var auctionInputDTO auction_usecase.AuctionInputDTO

	if err := c.ShouldBindJSON(&auctionInputDTO); err != nil {
		restErr := validation.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	outputDTO, err := ac.auctionUseCase.Create(c, auctionInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusCreated, outputDTO)
}
