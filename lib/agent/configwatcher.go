package agent

import (
	"github.com/fsnotify/fsnotify"
	"github.com/joaquinrovira/upv-oos-reservations/lib/logging"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/config"
)

func (a *Agent) ConfigWatch() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logging.Out().Fatal().Err(err)
	}
	defer watcher.Close()

	// Add a path.
	err = watcher.Add(a.cfg.ConfigFile)
	if err != nil {
		logging.Out().Fatal().Err(err)
	}

	for {
		select {
		case <-a.ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				logging.Out().Info().Msg("config file change detected, loading")
				target, err := config.LoadConfig(a.cfg.ConfigFile)
				if err != nil {
					logging.Out().Error().Err(err).Msg("error loading new config")
				} else {
					a.target = target
					logging.Out().Info().Interface("target", target).Msg("new config loaded successfully")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logging.Out().Error().Err(err).Msg("error with file-watcher")
		}
	}
}
