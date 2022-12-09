package agent

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"

	"github.com/joaquinrovira/upv-oos-reservations/lib/model/config"
)

type Agent struct {
	cfg    Config
	target config.TargetValue
	client *http.Client
	ctx    context.Context
	cancel context.CancelFunc
}

type Config struct {
	User       string
	Pass       string
	ConfigFile string
}

func New(c Config) (a *Agent, err error) {
	// Initial config load
	target, err := config.LoadConfig(c.ConfigFile)
	if err != nil {
		panic(err)
	}

	// Initialize HTTP agent with cookie Jar
	jar, err := cookiejar.New(nil)
	client := http.Client{
		// Enable cookie jar for user authentication
		Jar: jar,
		// Do not follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	// Build agent
	ctx, cancel := context.WithCancel(context.Background())
	a = &Agent{cfg: c, target: target, client: &client, ctx: ctx, cancel: cancel}

	// Register interrupt listener
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		cancel()
	}()

	// Startup thread to check for config changes
	go a.ConfigWatch()

	return
}

func (a *Agent) Stop() {
	a.cancel()
}
