package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

var Default *Daemon

func init() {
	Default = New()
}

type Daemon struct {
	app    *fiber.App
	config DaemonConfig
}

func New() (d *Daemon) {
	d = &Daemon{app: fiber.New()}
	d.config.Populate()
	return
}

func (d *Daemon) App() *fiber.App {
	return d.app
}

func (d *Daemon) Config() *DaemonConfig {
	return &d.config
}

func (d *Daemon) Start() error {
	ch := make(chan error)
	go func() {
		ch <- d.Listen()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-ch:
		return err
	case sig := <-sigs:
		if err := d.app.Shutdown(); err != nil {
			return fmt.Errorf("error while shutting down due to %v signal: %v", sig, err)
		} else {
			return nil
		}
	}
}
