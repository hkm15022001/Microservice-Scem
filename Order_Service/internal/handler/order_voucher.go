package handler

import (
	"net/http"
	"time"

	"github.com/hkm15022001/Supply-Chain-Event-Management/internal/model"
	"github.com/gin-gonic/gin"
)

// -------------------- DELIVERY LOCATION HANDLER FUNTION --------------------

// GetOrderVoucherListHandler in database
func GetOrderVoucherListHandler(c *gin.Context) {
	queryStrings := c.Request.URL.Query()
	sortByCondition := queryStrings.Get("sortByCondition")
	orderVoucherList := []model.OrderVoucher{}
	if sortByCondition == "comming" {
		db.Order("id asc").Find(&orderVoucherList, "start_date > ?", time.Now().Unix())
	} else if sortByCondition == "happening" {
		db.Order("id asc").Find(&orderVoucherList, "start_date < ? AND end_date > ?", time.Now().Unix(), time.Now().Unix())
	} else {
		db.Order("id asc").Find(&orderVoucherList)
	}
	c.JSON(http.StatusOK, gin.H{"order_voucher_list": &orderVoucherList})
	return
}
