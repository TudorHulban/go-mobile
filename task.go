package main

type Task struct {
	Name   string
	Owner  string
	Status string
}

var _Tasks = []Task{
	{
		Name:   "Operation 1",
		Status: "init",
	},
	{
		Name:   "Revision 1",
		Status: "not started",
	},
	{
		Name:   "Operation 2",
		Status: "assigned",
		Owner:  "John Smith",
	},
	{
		Name:   "Operation 3",
		Status: "in work",
		Owner:  "John Smith",
	},
	{
		Name:   "Operation 7",
		Status: "work done",
		Owner:  "Tom Sawyer",
	},
	{
		Name:   "Operation 8",
		Status: "to bill",
		Owner:  "Tom Sawyer",
	},
	{
		Name:   "Operation 9",
		Status: "invoiced",
		Owner:  "Mary Black",
	},
	{
		Name:   "Operation 10",
		Status: "closed",
		Owner:  "Tom Sawyer",
	},
}
