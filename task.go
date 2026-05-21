package main

type Task struct {
	Name   string
	Status string

	OwnerID uint8
}

type Tasks []*Task

func (t Tasks) GetNoTasksPerOwner() map[uint8]uint8 {
	allTasks := make(map[uint8]uint8, 0)

	for _, task := range t {
		_, exists := allTasks[task.OwnerID]
		if exists {
			allTasks[task.OwnerID]++

			continue
		}

		allTasks[task.OwnerID] = 1
	}

	return allTasks
}
