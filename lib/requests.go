package lib

import (
	"github.com/joaquinrovira/upv-oos-reservations/internal/model"
	"github.com/joaquinrovira/upv-oos-reservations/internal/requests"
)

func (agent Agent) Login() error {
	return requests.Login(agent.Client, agent.Cfg.User, agent.Cfg.Pass)
}

func (agent Agent) GetReservationsData() (*model.ReservationsWeek, error) {
	res, _ := requests.GetReservationsData(agent.Client)
	selection, _ := model.FindTable(res)
	data, _ := model.ParseHTMLTable(selection)
	reservations, _ := model.MarshalTable(&data)
	return reservations, nil
}
