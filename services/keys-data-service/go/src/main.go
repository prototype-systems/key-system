package main

import (
	"keys-data-service/utilities"
)

func main() {
	utilities.AllowProcess()

	if utilities.CheckChild() {
		RunChild()

		return
	}

	RunParent()
}
