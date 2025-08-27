package auction_controller

import (
	"l03/configuration/rest_err"
	"l03/internal/usecase/auction_usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ac *AuctionController) FindById(c *gin.Context) {
	id := c.Param("id")

	if err := uuid.Validate(id); err != nil {
		errRest := rest_err.NewBadRequestError("invalid auction id", rest_err.Causes{
			Field:   "auctionId",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	data, err := ac.auctionUseCase.FindById(c, id)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, data)
}

func (ac *AuctionController) FindAuctions(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	productName := c.Query("product-name")

	statusInt, errConv := strconv.Atoi(status)
	if errConv != nil {
		errRest := rest_err.NewBadRequestError("invalid status")
		c.JSON(errRest.Code, errRest)
		return
	}

	auctions, err := ac.auctionUseCase.FindActions(c, auction_usecase.AuctionStatus(statusInt), category, productName)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (ac *AuctionController) FindWinningBidbyAuctionId(c *gin.Context) {
	auctionId := c.Param("auction-id")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := rest_err.NewBadRequestError("invalid auction id", rest_err.Causes{
			Field:   "auction-id",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	auctionData, err := ac.auctionUseCase.FindWinnigBidByAuctionId(c, auctionId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctionData)
}
