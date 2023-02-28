package main

import (
	"fmt"

	"github.com/joaquinrovira/upv-oos-reservations/backend/daemon"
)

func main() {
	if err := daemon.Default.Start(); err != nil {
		fmt.Println(err)
	}
}
