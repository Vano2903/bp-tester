package model

type Event struct {
	Action      string
	ContainerID string
	ExitCode    string //only for simplicity atm
}
