package settings

import (
	"path/filepath"
	"time"
)

const TargetPort = "7375"

const TargetAddress = "http://127.0.0.1:" + TargetPort

const TargetNamespace = "keys-data-service"

var TargetTestDirectory = filepath.Join("services", "general-systems", "keys-data-service", "go", "src")

const TestTimeout = 10 * time.Second

const TestCommandTimeout = 3 * time.Second

const TestRequestTimeout = 3 * time.Second

const RootTestDirectory = "test"
