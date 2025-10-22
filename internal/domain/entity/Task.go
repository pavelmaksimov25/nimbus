package entity

type Task struct {
	ID string
	Items []TaskItem
}

type TaskItem struct {
	ID string
}

