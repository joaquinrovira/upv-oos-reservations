package daemon

import (
	"fmt"
	"net"
	"os"
)

// Listen serves HTTP requests based on the value of config.Listen
func (d *Daemon) Listen() error {
	Listen := d.config.Listen
	switch Listen {
	case Net:
		return d.listenNet()
	case Socket:
		return d.listenSocket()
	default:
		return fmt.Errorf("daemon listen kind '%s' unknown or not implemented", Listen)
	}
}

// Listen serves HTTP requests from the given interface and port
func (d *Daemon) listenNet() error {
	Interface := d.config.Net.Interface
	Port := d.config.Net.Port
	return d.app.Listen(fmt.Sprintf("%s:%d", Interface, Port))
}

// Listen serves HTTP requests from the given socket file
func (d *Daemon) listenSocket() error {
	Filename := d.config.Socket.Filename
	unixListener, err := net.Listen("unix", Filename)
	if err != nil {
		return fmt.Errorf("error requesting socket - %v", err)
	}
	defer os.Remove(Filename)
	return d.app.Listener(unixListener)
}
