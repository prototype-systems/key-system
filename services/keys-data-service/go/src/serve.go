package main

import (
	"context"
	"keys-data-service/channels"
	"keys-data-service/data"
	"keys-data-service/libraries"
	"keys-data-service/routes"
	"keys-data-service/settings"
	"keys-data-service/utilities"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func RunService() {
	RunIpc()

	RunLog()

	cache, issue := data.Open()
	if issue != nil {
		libraries.LogIssue("keys/serve", "cannot open cache", map[string]any{"issue": issue})

		utilities.KillService()
	}

	data.SetCache(cache)

	defer cache.Close()

	server := fiber.New(fiber.Config{
		AppName: settings.ServiceName,

		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 5,
		BodyLimit:    1 * 1024 * 1024,

		ErrorHandler: func(fiberContext fiber.Ctx, issue error) error {
			statusCode := fiber.StatusInternalServerError
			if fiberError, castSucceeded := issue.(*fiber.Error); castSucceeded {
				statusCode = fiberError.Code
			}

			libraries.LogIssue("keys/serve", "service request with issue", map[string]any{
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

	routes.RegisterService(keysGroup)
	routes.RegisterData(keysGroup, cache)

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

					libraries.LogIssue("keys/serve", "service file listener with issue", map[string]any{"issue": issue})

					channels.SignalStop()
				} else {
					if issue := server.Listener(networkListener); issue != nil {
						_ = networkListener.Close()
						_ = listenerFile.Close()

						libraries.LogIssue("keys/serve", "service network listener with issue", map[string]any{"issue": issue})

						channels.SignalStop()
					}
				}
			}

			return
		}

		if issue := server.Listen(":" + utilities.GetPort()); issue != nil {
			libraries.LogIssue("keys/serve", "service listen with issue", map[string]any{"issue": issue})

			channels.SignalStop()
		}
	}()

	terminateContext, cancelTerminateContext := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancelTerminateContext()

	select {
	case <-terminateContext.Done():
		libraries.LogData("keys/serve", "received terminate signal")
	case <-channels.GetStopChannel():
		libraries.LogData("keys/serve", "received stop signal")
	case <-channels.GetAbortChannel():
		libraries.LogData("keys/serve", "received abort signal")
	}

	if issue := server.Shutdown(); issue != nil {
		libraries.LogIssue("keys/serve", "service done with issue", map[string]any{"issue": issue})
	}
}
