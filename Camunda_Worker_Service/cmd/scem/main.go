package main

import (
	"log"
	"os"
	"sync"

	httpServer "github.com/hkm15022001/Supply-Chain-Event-Management/api/server"
	"github.com/hkm15022001/Supply-Chain-Event-Management/internal/handler"
	CommonService "github.com/hkm15022001/Supply-Chain-Event-Management/internal/service/common"
	CommonMessage "github.com/hkm15022001/Supply-Chain-Event-Management/internal/service/common_message"
	ZBMessage "github.com/hkm15022001/Supply-Chain-Event-Management/internal/service/zeebe/message"
	ZBWorker "github.com/hkm15022001/Supply-Chain-Event-Management/internal/service/zeebe/worker"
	ZBWorkflow "github.com/hkm15022001/Supply-Chain-Event-Management/internal/service/zeebe/workflow"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("RUNENV") != "docker" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Select database
	if os.Getenv("SELECT_DATABASE") == "1" {
		connectPostgress()
	} else if os.Getenv("SELECT_DATABASE") == "2" {
		connectMySQL()
	} else if os.Getenv("SELECT_DATABASE") == "3" {
		connectSQLite()
	} else {
		log.Println("No database selected!")
		os.Exit(1)
	}

	gormDB := handler.GetGormInstance()
	CommonService.MappingGormDBConnection(gormDB)
	CommonMessage.MappingGormDBConnection(gormDB)

	if os.Getenv("STATE_SERVICE") == "1" {
		connectZeebeClient()
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.RunServer()
	}()
	wg.Wait()
}

func connectPostgress() {
	if err := handler.ConnectPostgres(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Connected with posgres database!")
}

func connectMySQL() {
	if err := handler.ConnectMySQL(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Println("Connected with posgres database!")
}

func connectSQLite() {
	if err := handler.ConnectSQLite(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Println("Connected with sqlite database!")
}

func connectZeebeClient() {
	if err := ZBWorkflow.ConnectZeebeEngine(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Println("Zeebe workflow package connected with zeebe!")
	if err := ZBMessage.ConnectZeebeEngine(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	log.Println("Zeebe message package connected with zeebe!")
	// Run Zebee service
	ZBWorker.RunOrderLongShip()
	ZBWorker.RunOrderShortShip()
	ZBWorker.RunLongShipFinish()
}
