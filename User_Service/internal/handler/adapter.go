package handler

import (
	"log"
	"os"

	"github.com/hkm15022001/Supply-Chain-Event-Management/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// GetGormInstance function
func GetGormInstance() *gorm.DB {
	return db
}

// ConnectPostgres to open connect with database
func ConnectPostgres() (err error) {
	// https://github.com/go-gorm/postgres
	// https://stackoverflow.com/questions/50085286/postgresql-fatal-ident-authentication-failed-for-user
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("POSTGRES_DSN"),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	return err
}

// ConnectMySQL to open connect with database
func ConnectMySQL() (err error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := os.Getenv("MYSQL_DSN")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

// ConnectSQLite to open connect with database
func ConnectSQLite() (err error) {
	// github.com/mattn/go-sqlite3
	dsn := os.Getenv("SQLITE_DSN")
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	return err
}

// MigrationDatabase when update database
func MigrationDatabase() (err error) {
	return db.AutoMigrate(
		&model.Employee{},
		&model.Customer{},
		&model.CustomerCredit{},
		&model.EmployeeType{},
		&model.DeliveryLocation{},
		&model.OrderInfo{},
		&model.TransportType{},
		&model.UserAuthenticate{},
		&model.UserFCMToken{},
		&model.OrderPay{},
		&model.LongShip{},
		&model.OrderLongShip{},
		&model.OrderShortShip{},
		&model.OrderWorkflowData{},
		&model.LongShipWorkflowData{},
		&model.OrderVoucher{},
		&model.CustomerNotification{},
	)
}

// RefreshDatabase remove all table and create new data
func RefreshDatabase() (err error) {
	if err := deleteDatabase(); err != nil {
		return err
	}
	if err := MigrationDatabase(); err != nil {
		return err
	}
	if err := createDeliveryLocation(); err != nil {
		return err
	}
	if err := createEmployeeType(); err != nil {
		return err
	}
	if err := createDefaultEmployee(); err != nil {
		return err
	}
	if err := createDefaultCustomer(); err != nil {
		return err
	}
	if err := createCustomerCredit(); err != nil {
		return err
	}
	if err := createTransportType(); err != nil {
		return err
	}
	if err := createLongShip(); err != nil {
		return err
	}
	if err := createExampleOrder(); err != nil {
		return err
	}
	if err := createExampleOrderPay(); err != nil {
		return err
	}
	if err := createExampleOrderShortShip(); err != nil {
		return err
	}
	if err := createExampleOrderLongShip(); err != nil {
		return err
	}
	if err := createExampleOrderWorkflowData(); err != nil {
		return err
	}
	if err := createExampleOrder2(); err != nil {
		return err
	}
	if err := createOrderVoucher(); err != nil {
		return err
	}
	if err := createCustomerNotification(); err != nil {
		return err
	}
	return
}

func deleteDatabase() (err error) {
	return db.Migrator().DropTable(
		&model.Employee{},
		&model.Customer{},
		&model.CustomerCredit{},
		&model.EmployeeType{},
		&model.DeliveryLocation{},
		&model.OrderInfo{},
		&model.TransportType{},
		&model.UserAuthenticate{},
		&model.UserFCMToken{},
		&model.OrderPay{},
		&model.LongShip{},
		&model.OrderLongShip{},
		&model.OrderShortShip{},
		&model.OrderWorkflowData{},
		&model.LongShipWorkflowData{},
		&model.OrderVoucher{},
		&model.CustomerNotification{},
	)
}

func createEmployeeType() error {
	employeeType := &model.EmployeeType{Name: "Admin"}
	if err := db.Create(employeeType).Error; err != nil {
		return err
	}
	employeeType = &model.EmployeeType{Name: "Input staff"}
	if err := db.Create(employeeType).Error; err != nil {
		return err
	}
	employeeType = &model.EmployeeType{Name: "Delivery staff"}
	if err := db.Create(employeeType).Error; err != nil {
		return err
	}
	employeeType = &model.EmployeeType{Name: "Load package staff"}
	if err := db.Create(employeeType).Error; err != nil {
		return err
	}
	employeeType = &model.EmployeeType{Name: "Load package staff1"}
	if err := db.Create(employeeType).Error; err != nil {
		return err
	}
	return nil
}

func createDefaultEmployee() error {
	userAuth := &model.UserAuthenticate{Email: "admin@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee := &model.Employee{UserAuthID: userAuth.ID, Name: "Binh", Age: 40, Phone: 223334444, Gender: "male", Address: "12 Tran Hung Dao", IdentityCard: "17687t562765786", EmployeeTypeID: 1, Avatar: "image1.jpg"}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "inputstaff@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Huan", Age: 35, Phone: 448883333, Gender: "male", Address: "21 Huynh Thuc Khang", IdentityCard: "17687t562765786", EmployeeTypeID: 2, Avatar: "image2.jpg"}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "deliverystaff@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Tuan", Age: 37, Phone: 776664993, Gender: "male", Address: "21 Nhat Tao", IdentityCard: "17687t562765786", EmployeeTypeID: 3, Avatar: "image3.jpg", DeliveryLocationID: 6}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "loadpackagestaff@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Hung", Age: 47, Phone: 776334958, Gender: "male", Address: "84 Nguyen Trau", IdentityCard: "17687t562765786", EmployeeTypeID: 4, Avatar: "image3.jpg"}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "deliverystaff2@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Hieu", Age: 37, Phone: 776664993, Gender: "male", Address: "21 Trung Son", IdentityCard: "17687t562765786", EmployeeTypeID: 3, Avatar: "image3.jpg", DeliveryLocationID: 11}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "deliverystaff3@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Thao", Age: 37, Phone: 776662293, Gender: "female", Address: "32 Xuan Son", IdentityCard: "17687t562774286", EmployeeTypeID: 3, Avatar: "image3.jpg", DeliveryLocationID: 1}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}
	return nil
}

func createDefaultCustomer() error {
	userAuth := &model.UserAuthenticate{Email: "customer@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	log.Print("create auth!")
	customer := &model.Customer{UserAuthID: userAuth.ID, Name: "Customer One", Age: 18, Phone: 223334444, Gender: "male", Address: "12 Tran Hung Dao, Phuong 1, Quan 5"}
	if err := db.Create(customer).Error; err != nil {
		return err
	}
	userAuth.CustomerID = customer.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "customer3@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	customer = &model.Customer{UserAuthID: userAuth.ID, Name: "Customer Three", Age: 18, Phone: 223334444, Gender: "female", Address: "13 Tran Hung Dao"}
	if err := db.Create(customer).Error; err != nil {
		return err
	}
	userAuth.CustomerID = customer.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "customer2@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	customer = &model.Customer{UserAuthID: userAuth.ID, Name: "Customer Two", Age: 18, Phone: 223334444, Gender: "male", Address: "14 Tran Hung Dao"}
	if err := db.Create(customer).Error; err != nil {
		return err
	}
	userAuth.CustomerID = customer.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	return nil
}

func createCustomerCredit() error {
	customerCredit := &model.CustomerCredit{CustomerID: 1, Phone: 223334444, ValidatePhone: true, AccountBalance: 30000000}
	if err := db.Create(customerCredit).Error; err != nil {
		return err
	}
	customerCredit = &model.CustomerCredit{CustomerID: 2, Phone: 223334444, ValidatePhone: true, AccountBalance: 30000000}
	if err := db.Create(customerCredit).Error; err != nil {
		return err
	}
	customerCredit = &model.CustomerCredit{CustomerID: 3, Phone: 223334444, ValidatePhone: true, AccountBalance: 30000000}
	if err := db.Create(customerCredit).Error; err != nil {
		return err
	}
	return nil
}

func createTransportType() error {
	transportType := &model.TransportType{SameCity: true, LocationOne: "Hà Nội", ShortShipPricePerKm: 20000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{SameCity: true, LocationOne: "Bắc Ninh", ShortShipPricePerKm: 26000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "Hà Nội", LocationTwo: "Bắc Bộ", ServiceType: "Normal", LongShipPricePerKm: 1150, ShortShipPricePerKm: 16000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "Hà Nội", LocationTwo: "Bắc Bộ", ServiceType: "Express", LongShipPricePerKm: 2550, ShortShipPricePerKm: 20000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "Đà Nẵng", LocationTwo: "Trung Bộ", ServiceType: "Express", LongShipPricePerKm: 2220, ShortShipPricePerKm: 16000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "Đà Nẵng", LocationTwo: "Trung Bộ", ServiceType: "Normal", LongShipPricePerKm: 1050, ShortShipPricePerKm: 16000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "Đà Nẵng", LocationTwo: "Tây Nguyên", ServiceType: "Normal", LongShipPricePerKm: 1000, ShortShipPricePerKm: 16000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "Hồ Chí Minh", LocationTwo: "Nam Bộ", ServiceType: "Normal", LongShipPricePerKm: 1200, ShortShipPricePerKm: 16000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "Hồ Chí Minh", LocationTwo: "Nam Bộ", ServiceType: "Express", LongShipPricePerKm: 2600, ShortShipPricePerKm: 16000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "Hồ Chí Minh", LocationTwo: "Đồng Bằng Sông Cửu Long", ServiceType: "Normal", LongShipPricePerKm: 1270, ShortShipPricePerKm: 17000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "Hồ Chí Minh", LocationTwo: "Đồng Bằng Sông Cửu Longộ", ServiceType: "Express", LongShipPricePerKm: 2650, ShortShipPricePerKm: 16000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	return nil
}

func createLongShip() error {
	longShip := &model.LongShip{LSQrCode: "1611463545_4f163f5f0f9a621d.jpg", TransportTypeID: 3, LicensePlate: "51A 435.22", EstimatedTimeOfDeparture: 1610599301, EstimatedTimeOfArrival: 1610999301, Finished: true}
	if err := db.Create(longShip).Error; err != nil {
		return err
	}
	longShip = &model.LongShip{LSQrCode: "1611463464_d90bf90bbaa6cd39.jpg", TransportTypeID: 4, LicensePlate: "51B 425.82", EstimatedTimeOfDeparture: 1610099301, EstimatedTimeOfArrival: 1610399301, Finished: true}
	if err := db.Create(longShip).Error; err != nil {
		return err
	}
	return nil
}

func createDeliveryLocation() error {
	vietnamCities := []model.DeliveryLocation{
		{City: "Hà Giang", AreaCode: "0219", Latitude: 22.8333, Longitude: 104.9833, Region: "Bắc Bộ"},
		{City: "Cao Bằng", AreaCode: "0206", Latitude: 22.6667, Longitude: 106.25, Region: "Bắc Bộ"},
		{City: "Lào Cai", AreaCode: "0214", Latitude: 22.4833, Longitude: 103.9667, Region: "Bắc Bộ"},
		{City: "Sơn La", AreaCode: "0212", Latitude: 21.3333, Longitude: 103.9, Region: "Bắc Bộ"},
		{City: "Lai Châu", AreaCode: "0213", Latitude: 22.3964, Longitude: 103.4431, Region: "Bắc Bộ"},
		{City: "Bắc Kạn", AreaCode: "0209", Latitude: 22.1333, Longitude: 105.8333, Region: "Bắc Bộ"},
		{City: "Lạng Sơn", AreaCode: "0205", Latitude: 21.8333, Longitude: 106.7333, Region: "Bắc Bộ"},
		{City: "Tuyên Quang", AreaCode: "0207", Latitude: 21.8333, Longitude: 105.25, Region: "Bắc Bộ"},
		{City: "Yên Bái", AreaCode: "0216", Latitude: 21.7167, Longitude: 104.9, Region: "Bắc Bộ"},
		{City: "Thái Nguyên", AreaCode: "0208", Latitude: 21.6, Longitude: 105.8333, Region: "Bắc Bộ"},
		{City: "Điện Biên", AreaCode: "0215", Latitude: 21.3833, Longitude: 103, Region: "Bắc Bộ"},
		{City: "Phú Thọ", AreaCode: "0210", Latitude: 21.3333, Longitude: 105.1667, Region: "Bắc Bộ"},
		{City: "Vĩnh Phúc", AreaCode: "0211", Latitude: 21.3, Longitude: 105.6, Region: "Bắc Bộ"},
		{City: "Bắc Giang", AreaCode: "0204", Latitude: 21.2667, Longitude: 106.2, Region: "Bắc Bộ"},
		{City: "Bắc Ninh", AreaCode: "0222", Latitude: 21.1861, Longitude: 106.075, Region: "Bắc Bộ"},
		{City: "Hà Nội", AreaCode: "024", Latitude: 21.0285, Longitude: 105.8542, Region: "Bắc Bộ"},
		{City: "Quảng Ninh", AreaCode: "0203", Latitude: 20.9531, Longitude: 107.0656, Region: "Bắc Bộ"},
		{City: "Hải Dương", AreaCode: "0220", Latitude: 20.9333, Longitude: 106.3333, Region: "Bắc Bộ"},
		{City: "Hải Phòng", AreaCode: "0225", Latitude: 20.8481, Longitude: 106.6883, Region: "Bắc Bộ"},
		{City: "Hòa Bình", AreaCode: "0218", Latitude: 20.8133, Longitude: 105.3383, Region: "Bắc Bộ"},
		{City: "Hưng Yên", AreaCode: "0221", Latitude: 20.6467, Longitude: 106.0514, Region: "Bắc Bộ"},
		{City: "Hà Nam", AreaCode: "0226", Latitude: 20.5417, Longitude: 105.9222, Region: "Bắc Bộ"},
		{City: "Thái Bình", AreaCode: "0227", Latitude: 20.4556, Longitude: 106.3319, Region: "Bắc Bộ"},
		{City: "Nam Định", AreaCode: "0228", Latitude: 20.4267, Longitude: 106.165, Region: "Bắc Bộ"},
		{City: "Ninh Bình", AreaCode: "0229", Latitude: 20.2544, Longitude: 105.975, Region: "Bắc Bộ"},
		{City: "Thanh Hóa", AreaCode: "0237", Latitude: 19.8066, Longitude: 105.7667, Region: "Trung Bộ"},
		{City: "Nghệ An", AreaCode: "0238", Latitude: 18.8236, Longitude: 105.6328, Region: "Trung Bộ"},
		{City: "Hà Tĩnh", AreaCode: "0239", Latitude: 18.3333, Longitude: 105.9, Region: "Trung Bộ"},
		{City: "Quảng Bình", AreaCode: "0232", Latitude: 17.5364, Longitude: 106.5174, Region: "Trung Bộ"},
		{City: "Quảng Trị", AreaCode: "0233", Latitude: 16.75, Longitude: 107.2, Region: "Trung Bộ"},
		{City: "Thừa Thiên Huế", AreaCode: "0234", Latitude: 16.4633, Longitude: 107.5956, Region: "Trung Bộ"},
		{City: "Đà Nẵng", AreaCode: "0236", Latitude: 16.0544, Longitude: 108.2022, Region: "Trung Bộ"},
		{City: "Quảng Nam", AreaCode: "0235", Latitude: 15.8801, Longitude: 108.3385, Region: "Trung Bộ"},
		{City: "Quảng Ngãi", AreaCode: "0255", Latitude: 15.1201, Longitude: 108.798, Region: "Trung Bộ"},
		{City: "Kon Tum", AreaCode: "0260", Latitude: 14.3544, Longitude: 107.9833, Region: "Tây Nguyên"},
		{City: "Gia Lai", AreaCode: "0269", Latitude: 13.9833, Longitude: 108.2167, Region: "Tây Nguyên"},
		{City: "Bình Định", AreaCode: "0256", Latitude: 14.25, Longitude: 109.2, Region: "Trung Bộ"},
		{City: "Phú Yên", AreaCode: "0257", Latitude: 13.1667, Longitude: 109.1667, Region: "Trung Bộ"},
		{City: "Đắk Lắk", AreaCode: "0262", Latitude: 12.8333, Longitude: 108.05, Region: "Tây Nguyên"},
		{City: "Khánh Hòa", AreaCode: "0258", Latitude: 12.25, Longitude: 109.2, Region: "Trung Bộ"},
		{City: "Đắk Nông", AreaCode: "0261", Latitude: 12.25, Longitude: 107.8, Region: "Tây Nguyên"},
		{City: "Lâm Đồng", AreaCode: "0263", Latitude: 11.5, Longitude: 108.8333, Region: "Tây Nguyên"},
		{City: "Ninh Thuận", AreaCode: "0259", Latitude: 11.75, Longitude: 108.3667, Region: "Trung Bộ"},
		{City: "Bình Phước", AreaCode: "0271", Latitude: 11.75, Longitude: 106.75, Region: "Nam Bộ"},
		{City: "Tây Ninh", AreaCode: "0276", Latitude: 11.3333, Longitude: 106.2, Region: "Nam Bộ"},
		{City: "Bình Dương", AreaCode: "0274", Latitude: 11.1667, Longitude: 106.6667, Region: "Nam Bộ"},
		{City: "Đồng Nai", AreaCode: "0251", Latitude: 11.1667, Longitude: 107.1667, Region: "Nam Bộ"},
		{City: "Bình Thuận", AreaCode: "0252", Latitude: 10.9333, Longitude: 108.1, Region: "Trung Bộ"},
		{City: "Hồ Chí Minh", AreaCode: "028", Latitude: 10.7769, Longitude: 106.7009, Region: "Nam Bộ"},
		{City: "Long An", AreaCode: "0272", Latitude: 10.6667, Longitude: 106, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Bà Rịa – Vũng Tàu", AreaCode: "0254", Latitude: 10.4091, Longitude: 107.1386, Region: "Nam Bộ"},
		{City: "Đồng Tháp", AreaCode: "0277", Latitude: 10.6667, Longitude: 105.6667, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "An Giang", AreaCode: "0296", Latitude: 10.5, Longitude: 105.1667, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Tiền Giang", AreaCode: "0273", Latitude: 10.3547, Longitude: 106.3406, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Vĩnh Long", AreaCode: "0270", Latitude: 10.2556, Longitude: 105.9722, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Bến Tre", AreaCode: "0275", Latitude: 10.2417, Longitude: 106.375, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Cần Thơ", AreaCode: "0292", Latitude: 10.0333, Longitude: 105.7833, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Kiên Giang", AreaCode: "0297", Latitude: 10, Longitude: 105.1667, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Trà Vinh", AreaCode: "0294", Latitude: 9.935, Longitude: 106.3456, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Hậu Giang", AreaCode: "0293", Latitude: 9.7875, Longitude: 105.4678, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Sóc Trăng", AreaCode: "0299", Latitude: 9.6039, Longitude: 105.9811, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Bạc Liêu", AreaCode: "0291", Latitude: 9.1667, Longitude: 105.15, Region: "Đồng Bằng Sông Cửu Long"},
		{City: "Cà Mau", AreaCode: "0290", Latitude: 9.1769, Longitude: 105.152, Region: "Đồng Bằng Sông Cửu Long"},
	}

	for _, cityInfor := range vietnamCities {
		location := &cityInfor
		if err := db.Create(location).Error; err != nil {
			return err
		}
	}
	// location := &{City: "Hà Nội", Latitude: 21.0285, Longitude: 105.8542}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "Hà Nội", District: "2"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "Hà Nội", District: "3"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "Hà Nội", District: "Tan Binh"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "Hà Nội", District: "Phu Nhuan"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "DL", District: "1"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "DL", District: "2"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "DL", District: "3"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "DL", District: "4"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "DL", District: "5"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "VT", District: "1"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	// location = &{City: "VT", District: "2"}
	// if err := db.Create(location).Error; err != nil {
	// 	return err
	// }
	return nil
}

func createExampleOrder() error {
	orderInfo := &model.OrderInfo{
		Weight: 2, Volume: 10, Type: "Normal", Image: "box.jpg",
		CustomerSendID: 1, EmplCreateID: 2,
		Sender:           "Customer One - 269 Ngo Quyen, Hà Nội - 5676765678",
		Receiver:         "Mai Thi Cuc - 38 Tran Hung Dao, Hà Nội - 6765677867",
		Detail:           "May vi tinh ca nhan va ban phim may tinh",
		ReceiversAddress: "Hà Nội",
		OrderShortShipID: 1, TransportTypeID: 1, LongShipDistance: 20,
		TotalPrice: 200000, Note: "Giao hang vao buoi sang",
	}
	if err := db.Create(orderInfo).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrderPay() error {
	orderPay := &model.OrderPay{
		OrderID: 1, PayStatus: true, TotalPrice: 200000, PayMethod: "cash",
	}
	if err := db.Create(orderPay).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrderWorkflowData() error {
	orderPay := &model.OrderWorkflowData{
		OrderID: 1, WorkflowKey: "123abc", WorkflowInstanceKey: 321123321123, OrderPayID: 1, CustomerSendID: 1,
	}
	if err := db.Create(orderPay).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrderLongShip() error {
	orderLongShip := &model.OrderLongShip{
		OrderID: 1, LongShipID: 1, CustomerSendID: 1,
	}
	if err := db.Create(orderLongShip).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrderShortShip() error {
	orderShortShip := &model.OrderShortShip{
		OrderID: 1, ShipperID: 3, CustomerSendID: 1, OSSQrCode: "1611465837_1e708c0f8e9a0213.jpg",
		Sender:   "Customer One - 269 Ngo Quyen, Quan 5, Hà Nội - 5676765678",
		Receiver: "Mai Thi Cuc - 38 Tran Hung Dao, Quan 1, Hà Nội - 6765677867",
	}
	if err := db.Create(orderShortShip).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrder2() error {
	orderInfo := &model.OrderInfo{
		Weight: 3, Volume: 50, Type: "Special", Image: "box.jpg",
		CustomerSendID: 1,
		Sender:         "Customer one - 231-233 Le Hong Phong - 6578678678",
		Receiver:       "Mac Thi Buoi - 74 Phan Chau Trinh, Quan 3, DL - 567865676",
		Detail:         "May xay thit", UseLongShip: true, LongShipID: 1,
		OrderShortShipID: 1, TransportTypeID: 3, LongShipDistance: 20,
		TotalPrice: 200000, Note: "Giao hang vao buoi trua",
	}
	if err := db.Create(orderInfo).Error; err != nil {
		return err
	}
	return nil
}

func createOrderVoucher() error {
	orderVoucher := &model.OrderVoucher{
		Title:     "Khuyen mai ngay he 1",
		Content:   "Tan huong nhung ngay he sang khoang voi Move nice VN nha. Giam gia 30.000 tu hom nay den cuoi thang.",
		StartDate: 1611390576, EndDate: 1619390576, Discount: 30000,
	}
	if err := db.Create(orderVoucher).Error; err != nil {
		return err
	}
	orderVoucher = &model.OrderVoucher{
		Title:     "Khuyen mai ngay xuan 2",
		Content:   "Tan huong nhung ngay xuan mat me voi Move nice VN nha. Giam gia 50.000 tu hom nay den cuoi thang.",
		StartDate: 1615690576, EndDate: 1619390576, Discount: 50000,
	}
	if err := db.Create(orderVoucher).Error; err != nil {
		return err
	}
	orderVoucher = &model.OrderVoucher{
		Title:     "Khuyen mai ngay dong 3",
		Content:   "Tan huong nhung ngay xuan mat me voi Move nice VN nha. Giam gia 70.000 tu hom nay den cuoi thang.",
		StartDate: 1612390576, EndDate: 1620390576, Discount: 70000,
	}
	if err := db.Create(orderVoucher).Error; err != nil {
		return err
	}
	orderVoucher = &model.OrderVoucher{
		Title:     "Khuyen mai ngay phu nu 5",
		Content:   "Chuc mung ngay quoc te phu nu voi Move nice VN nha. Giam gia 100.000 tu hom nay den cuoi thang.",
		StartDate: 1612390576, EndDate: 1620390576, Discount: 100000,
	}
	if err := db.Create(orderVoucher).Error; err != nil {
		return err
	}
	return nil
}

func createCustomerNotification() error {
	customerNotification := &model.CustomerNotification{
		CustomerID: 1, Title: "Your order has started long ship trip",
		Content: "Order id: 1349014 Long ship id: 1268070",
	}
	if err := db.Create(customerNotification).Error; err != nil {
		return err
	}
	customerNotification = &model.CustomerNotification{
		CustomerID: 1, Title: "Your order has finished long ship trip",
		Content: "Order id: 1349014 Long ship id: 1268070",
	}
	if err := db.Create(customerNotification).Error; err != nil {
		return err
	}
	customerNotification = &model.CustomerNotification{
		CustomerID: 1, Title: "Shipper has called you",
		Content: "Your order id: 1349014 has been verified",
	}
	if err := db.Create(customerNotification).Error; err != nil {
		return err
	}
	customerNotification = &model.CustomerNotification{
		CustomerID: 1, Title: "Shipper has confirmed your package",
		Content: "Thanks for using our service. Finished order id: 1349014",
	}
	if err := db.Create(customerNotification).Error; err != nil {
		return err
	}
	return nil
}
