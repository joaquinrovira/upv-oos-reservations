package agent

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"time"

	"github.com/joaquinrovira/upv-oos-reservations/lib/model/config"
	"github.com/joaquinrovira/upv-oos-reservations/lib/util"
	"golang.org/x/time/rate"
)

type Agent struct {
	cfg    Config
	target config.TargetValue
	client util.RLHTTPClient
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
	client := util.NewClient(
		&http.Client{
			// Enable cookie jar for user authentication
			Jar: jar,
			// Do not follow redirects
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		rate.NewLimiter(rate.Every(1*time.Second), 3)) 

	// Build agent
	ctx, cancel := context.WithCancel(context.Background())
	a = &Agent{cfg: c, target: target, client: client, ctx: ctx, cancel: cancel}

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
