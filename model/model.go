package model

import "gorm.io/gorm"

type Todomvc struct {
	gorm.Model
	Item   string
	Status uint
}

type ToDoMvcDel struct {
	Id uint
}

type ToDoMvcAdd struct {
	Item string
}

type ToDoMvcUpdate struct {
	Id     uint
	Item   string
	Status uint
}

type ToDoMvcFind struct {
	Item   string
	Status int
}
