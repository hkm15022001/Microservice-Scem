package common

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hkm15022001/Supply-Chain-Event-Management/internal/model"
	"github.com/skip2/go-qrcode"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var db *gorm.DB

// MappingGormDBConnection to open connect with database
func MappingGormDBConnection(_db *gorm.DB) {
	db = _db
}

// ------------------------------ CALL FROM STATE SERVICE ------------------------------

// GetOrderLongShipList function
func GetOrderLongShipList(longShipID uint) ([]model.OrderLongShip, error) {
	orderLongShips := []model.OrderLongShip{}
	if err := db.Where("long_ship_id = ?", longShipID).Order("id asc").Find(&orderLongShips).Error; err != nil {
		return nil, err
	}
	return orderLongShips, nil
}

// ++++++++++++++++++++ Order Long Ship Worker ++++++++++++++++++++

func getOrderInfoOrNotFoundForShipment(orderID uint) (*model.OrderInfoForShipment, error) {

	orderInfoForShipment := &model.OrderInfoForShipment{}
	err := db.Model(&model.OrderInfo{}).Order("id asc").First(&orderInfoForShipment, orderID).Error
	if err != nil {
		return orderInfoForShipment, err
	}
	return orderInfoForShipment, nil
}

func getOrderWorkflowDataByOrderID(orderID uint) (*model.OrderWorkflowData, error) {

	orderWorkflowData := &model.OrderWorkflowData{}
	err := db.Model(orderWorkflowData).Order("id asc").First(orderWorkflowData, "order_id = ?", orderID).Error
	if err != nil {
		return orderWorkflowData, err
	}
	return orderWorkflowData, nil
}

// CreateOrderLongShip in database
func CreateOrderLongShip(orderID uint) (uint, error) {

	orderInfoForShipment, err := getOrderInfoOrNotFoundForShipment(orderID)
	if err != nil {
		return uint(0), err
	}

	orderLongShip := &model.OrderLongShip{}
	orderLongShip.OrderID = orderID
	orderLongShip.LongShipID = orderInfoForShipment.LongShipID
	orderLongShip.CustomerSendID = orderInfoForShipment.CustomerSendID
	orderLongShip.CustomerReceiveID = orderInfoForShipment.CustomerReceiveID

	if err := db.Create(orderLongShip).Error; err != nil {
		return uint(0), err
	}

	orderInfo := &model.OrderInfo{}
	orderInfo.ID = orderID
	orderInfo.OrderLongShipID = orderLongShip.ID
	if err = db.Model(&orderInfo).Updates(&orderInfo).Error; err != nil {
		return uint(0), err
	}
	return orderLongShip.ID, nil
}

// ++++++++++++++++++++ Order Short Ship Worker ++++++++++++++++++++

// CreateOrderShortShip function
func CreateOrderShortShip(orderID uint) (uint, error) {

	orderShortShip := &model.OrderShortShip{}
	orderShortShip.OrderID = orderID
	orderInfoForShipment, err := getOrderInfoOrNotFoundForShipment(orderID)
	if err != nil {
		return uint(0), err
	}
	orderWorkflowData, err := getOrderWorkflowDataByOrderID(orderID)
	if err != nil {
		return uint(0), err
	}

	transportType := &model.TransportType{}
	if err := db.First(transportType, orderInfoForShipment.TransportTypeID).Error; err != nil {
		return uint(0), err
	}
	deliveryLocation := &model.DeliveryLocation{}

	if orderInfoForShipment.UseLongShip == true {
		if err := db.Where("city = ?", orderInfoForShipment.ReceiversAddress).First(deliveryLocation).Error; err != nil {
			return uint(0), err
		}
	} else {
		if err := db.Where("city = ?", transportType.LocationOne).First(deliveryLocation).Error; err != nil {
			return uint(0), err
		}
	}
	// Pick random employee type driver base on location
	employeeList := []model.EmployeeInfoForShortShip{}
	selectPart := "e.id, e.employee_type_id, e.delivery_location_id "
	err = db.Table("employees as e").Select(selectPart).
		Where("e.employee_type_id = ? AND e.delivery_location_id = ?", 3, deliveryLocation.ID).Find(&employeeList).Error
	if err != nil {
		return uint(0), err
	}

	length := len(employeeList)
	if length > 1 {
		length = length - 1
	}
	log.Print("length of delivery staff short ship in location:", length)
	index := rand.Intn(length)
	orderShortShip.ShipperID = employeeList[index].ID
	orderShortShip.CustomerSendID = orderInfoForShipment.CustomerSendID
	orderShortShip.CustomerReceiveID = orderInfoForShipment.CustomerReceiveID
	orderShortShip.Sender = orderInfoForShipment.Sender
	orderShortShip.Receiver = orderInfoForShipment.Receiver
	orderShortShip.ShipperReceiveMoney = orderWorkflowData.ShipperReceiveMoney

	// Create QR code
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return uint(0), err
	}
	newName := fmt.Sprintf("%x", b)
	createTime := fmt.Sprintf("%d", time.Now().Unix())
	newName = createTime + "_" + newName + ".jpg"
	filepath := os.Getenv("QR_CODE_FILE_PATH") + newName
	if err := qrcode.WriteFile(newName, qrcode.Medium, 256, filepath); err != nil {
		return uint(0), err
	}
	orderShortShip.OSSQrCode = newName

	if err := db.Create(orderShortShip).Error; err != nil {
		return uint(0), err
	}
	orderInfo := &model.OrderInfo{}
	orderInfo.ID = orderID
	orderInfo.OrderShortShipID = orderShortShip.ID
	if err = db.Model(&orderInfo).Updates(&orderInfo).Error; err != nil {
		return uint(0), err
	}
	return orderShortShip.ID, nil
}
func CreateLongShipHandlerGRPC() (uint, error) {
	currentTime := time.Now()
	// Tính toán thời gian cụ thể 6 giờ sáng hôm sau
	nextDay := currentTime.AddDate(0, 0, 1) // Thêm 1 ngày
	DepartureTime := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 6, 0, 0, 0, nextDay.Location())
	// Create long ship id base on Time
	DepartureTimeTimestamp := DepartureTime.Unix()
	log.Print(DepartureTimeTimestamp)
	longShip := &model.LongShip{}
	longShip.TransportTypeID = 3
	longShip.EstimatedTimeOfDeparture = DepartureTimeTimestamp

	current := uuid.New().Time()
	currentString := fmt.Sprintf("%d", current)
	rawUint, _ := strconv.ParseUint(currentString, 10, 64)
	longShip.ID = uint(rawUint / 100000000000)
	log.Print("1", longShip)
	// Run concurrency
	var g errgroup.Group

	// Create QR code
	g.Go(func() error {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			return err
		}
		newQrCode := fmt.Sprintf("%x", b)
		createTime := fmt.Sprintf("%d", time.Now().Unix())
		newQrCode = createTime + "_" + newQrCode + ".jpg"
		filepath := os.Getenv("QR_CODE_FILE_PATH") + newQrCode
		if err := qrcode.WriteFile(newQrCode, qrcode.Medium, 256, filepath); err != nil {
			log.Print("error when create qrcode", err)
			return err
		}
		longShip.LSQrCode = newQrCode
		log.Print("2", longShip)

		if err := db.Create(longShip).Error; err != nil {
			log.Print("Error when create LongShip in db", err)
			return err
		}
		return nil
	})

	// Create workflow instance in zeebe in state machine and save to long ship workflow data
	g.Go(func() error {
		WorkflowKey, WorkflowInstanceKey, err := CreateWorkflowLongShipInstanceHandler(longShip.ID)
		if err != nil {
			log.Print("Error when create LongShipInstance", err)
			return err
		}

		longShipWorkflowData := &model.LongShipWorkflowData{}
		longShipWorkflowData.LongShipID = longShip.ID
		longShipWorkflowData.WorkflowKey = WorkflowKey
		longShipWorkflowData.WorkflowInstanceKey = WorkflowInstanceKey
		if err := db.Create(longShipWorkflowData).Error; err != nil {
			log.Print("Error when create longshipworkflow in db", err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Print("Error when wait group", err)
		return uint(0), err
	}
	return longShip.ID, nil
}

func UpdateOrderLongShipGRPC(orderId uint, longShipId uint) (uint, error) {
	orderInfo, err := getOrderInfoOrNotFoundGRPC(orderId)
	if err != nil {
		return uint(0), err
	}
	orderInfo.LongShipID = longShipId
	if err = db.Model(&orderInfo).Updates(&orderInfo).Error; err != nil {
		return uint(0), err
	}
	orderPay, err := getOrderPayOrNotFoundByOrderID(orderId)
	if err != nil {
		log.Print("This OrderId have not been pay")
		return uint(0), err
	}
	// shipper_receive_money will be sent when using app
	orderWorkflowData := &model.OrderWorkflowData{
		OrderID:             orderPay.OrderID,
		ShipperReceiveMoney: orderPay.ShipperReceiveMoney,
		UseLongShip:         orderInfo.UseLongShip,
		CustomerSendID:      orderInfo.CustomerSendID,
		CustomerReceiveID:   orderInfo.CustomerReceiveID,
	}
	// Create workflow instance in zeebe
	WorkflowKey, WorkflowInstanceKey, err := CreateWorkflowFullShipInstanceHandler(orderWorkflowData)
	if err != nil {
		log.Print("Error when create fullship workflow instance from plan service:", err)
		return uint(0), err
	}
	orderWorkflowData.OrderPayID = orderPay.ID
	orderWorkflowData.LongShipID = longShipId
	orderWorkflowData.WorkflowKey = WorkflowKey
	orderWorkflowData.WorkflowInstanceKey = WorkflowInstanceKey
	// Create workflow data in database
	if err := db.Create(orderWorkflowData).Error; err != nil {
		return uint(0), err
	}
	return orderId, nil
}
func getOrderInfoOrNotFoundGRPC(id uint) (*model.OrderInfo, error) {
	orderInfo := &model.OrderInfo{}
	if err := db.First(orderInfo, id).Error; err != nil {
		return orderInfo, err
	}
	return orderInfo, nil
}
func getOrderPayOrNotFoundByOrderID(orderID uint) (*model.OrderPay, error) {
	orderPay := &model.OrderPay{}
	if err := db.Where("order_id = ?", orderID).First(orderPay).Error; err != nil {
		return orderPay, err
	}
	return orderPay, nil
}
