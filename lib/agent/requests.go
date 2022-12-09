package agent

import (
	"github.com/joaquinrovira/upv-oos-reservations/lib/model"
	"github.com/joaquinrovira/upv-oos-reservations/lib/requests"
)

func (agent Agent) Login() error {
	return requests.Login(agent.client, agent.cfg.User, agent.cfg.Pass)
}

func (agent Agent) GetReservationsData() (*model.ReservationsWeek, error) {
	res, _ := requests.GetReservationsData(agent.client)
	selection, _ := model.FindTable(res)
	data, _ := model.ParseHTMLTable(selection)
	reservations, _ := model.MarshalTable(&data)
	return reservations, nil
}

// Agent
// Read config
// Register fsnotify to update on config changes
// Ready to run
