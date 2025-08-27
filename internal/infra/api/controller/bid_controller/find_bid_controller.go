package bid_controller

import (
	"l03/configuration/rest_err"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (bc *BidController) FindBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auction-id")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := rest_err.NewBadRequestError("invalid auction id", rest_err.Causes{
			Field:   "auction-id",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	data, err := bc.BidUseCase.FindByAuctionId(c, auctionId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, data)
}
