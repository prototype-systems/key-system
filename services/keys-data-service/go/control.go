package main

import (
	"keys-data-service/channels"
	"keys-data-service/libraries"
	"keys-data-service/routes"
	"keys-data-service/settings"
	"keys-data-service/utilities"
	"net"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func RunControl() {
	RunIpc()

	RunLog()

	server := fiber.New(fiber.Config{
		AppName: settings.ServiceName,

		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 5,

		ErrorHandler: func(fiberContext fiber.Ctx, issue error) error {
			statusCode := fiber.StatusInternalServerError
			if fiberError, castSucceeded := issue.(*fiber.Error); castSucceeded {
				statusCode = fiberError.Code
			}

			libraries.LogIssue("keys/control", "service request with issue", map[string]any{
				"status": statusCode,
				"issue":  issue,
			})

			return utilities.IssueResponse(fiberContext, statusCode, issue.Error())
		},
	})

	server.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	server.Use(logger.New())
	server.Use(cors.New())

	serviceGroup := server.Group("/service")
	dataGroup := serviceGroup.Group("/data")
	keysGroup := dataGroup.Group("/keys")

	routes.RegisterControls(keysGroup)

	server.Use(func(fiberContext fiber.Ctx) error {
		return utilities.IssueResponse(fiberContext, fiber.StatusNotFound, "request undetermined")
	})

	go func() {
		if utilities.CheckServiceFileDescriptor() {
			listenerFile := os.NewFile(utilities.GetServiceFileDescriptor(), "listener")
			if listenerFile != nil {
				networkListener, issue := net.FileListener(listenerFile)

				if issue != nil {
					_ = listenerFile.Close()

					libraries.LogIssue("keys/control", "service file listener issue", map[string]any{"issue": issue})

					channels.SignalKill()
				} else {
					if issue := server.Listener(networkListener); issue != nil {
						_ = networkListener.Close()
						_ = listenerFile.Close()

						libraries.LogIssue("keys/control", "service network listener issue", map[string]any{"issue": issue})

						channels.SignalKill()
					}
				}
			}

			return
		}

		if issue := server.Listen(":" + utilities.GetPort()); issue != nil {
			libraries.LogIssue("keys/control", "service listen with issue", map[string]any{"issue": issue})

			channels.SignalKill()
		}
	}()

	select {
	case <-channels.GetStartChannel():
		libraries.LogData("keys/control", "received start signal")
	case <-channels.GetKillChannel():
		libraries.LogData("keys/control", "received kill signal")
	}

	if issue := server.Shutdown(); issue != nil {
		libraries.LogIssue("keys/control", "service done with issue", map[string]any{"issue": issue})
	}
}
