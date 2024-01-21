package httpServer

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hkm15022001/Supply-Chain-Event-Management/api/middleware"
	"github.com/hkm15022001/Supply-Chain-Event-Management/api/router"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

// RunServer will start 2 server for app and web
func RunServer() {
	// gin.SetMode(gin.ReleaseMode)
	// export GIN_MODE=debug

	webServer := &http.Server{
		Addr:         os.Getenv("WEB_PORT"),
		Handler:      webRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	appServer := &http.Server{
		Addr:         os.Getenv("APP_PORT"),
		Handler:      appRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		err := webServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	g.Go(func() error {
		err := appServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func webRouter() http.Handler {
	e := gin.Default()

	e.Static("/scem-ship/api/images", os.Getenv("IMAGE_FILE_PATH"))
	e.Static("/scem-ship//api/qrcode", os.Getenv("QR_CODE_FILE_PATH"))
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	e.MaxMultipartMemory = 8 << 20 // 8 MiB

	api := e.Group("/scem-ship/api")

	// Active web auth
	if os.Getenv("RUN_WEB_AUTH") == "yes" {
		api.Use(middleware.ValidateWebSession())
	}

	router.WebOrderShipRoutes(api)
	return e
}

func appRouter() http.Handler {
	e := gin.Default()
	e.Static("/scem-ship//api/images", os.Getenv("IMAGE_FILE_PATH"))
	api := e.Group("/scem-ship//api")
	router.AppOrderShipRoutes(api)
	return e
}
