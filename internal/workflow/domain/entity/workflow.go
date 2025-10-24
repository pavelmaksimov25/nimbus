package entity

type Workflow struct {
	ID      string
	Payload string
	Tasks   []Task
}