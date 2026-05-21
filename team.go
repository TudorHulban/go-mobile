package main

import (
	"fmt"
	"slices"
)

type Team struct {
	TeamMembers map[uint8]string
	TeamTasks   map[uint8]*Task

	lastTaskID uint8
}

func NewTeam(members ...string) *Team {
	result := Team{
		TeamMembers: make(map[uint8]string),
		TeamTasks:   map[uint8]*Task{},
	}

	for ix, member := range members {
		result.TeamMembers[uint8(ix)] = member
	}

	return &result
}

func (t *Team) GetTeamMembers() []string {
	result := make([]string, len(t.TeamMembers))

	for ix, teamMember := range t.TeamMembers {
		result[ix] = teamMember
	}

	return result
}

func (t *Team) UpsertTask(task *Task) error {
	if _, exists := t.TeamMembers[task.OwnerID]; !exists && task.OwnerID != 0 {
		return fmt.Errorf(
			"unknown owner ID: %d",
			task.OwnerID,
		)
	}

	// update
	if task.ID != 0 {
		currentTask, errGet := t.GetTaskBy(task.ID)
		if errGet != nil {
			return errGet
		}

		// currentTask.Name = task.Name
		currentTask.OwnerID = task.OwnerID
		currentTask.Status = task.Status

		return nil
	}

	// create
	t.lastTaskID++

	t.TeamTasks[t.lastTaskID] = task

	return nil
}

func (t *Team) GetTaskBy(id uint8) (*Task, error) {
	if result, exists := t.TeamTasks[id]; exists {
		return result, nil
	}

	return nil,
		fmt.Errorf(
			"task with ID %d not found",
			id,
		)
}

func (t *Team) GetTeamMemberBy(id uint8) (string, error) {
	if result, exists := t.TeamMembers[id]; exists {
		return result, nil
	}

	return "",
		fmt.Errorf(
			"task with ID %d not found",
			id,
		)
}

func (t *Team) GetTeamMemberIDBy(name string) (uint8, error) {
	for result, teamMemberName := range t.TeamMembers {
		if teamMemberName == name {
			return result, nil

		}
	}

	return 0,
		fmt.Errorf(
			"team member with name %s not found",
			name,
		)
}

func (t *Team) GetTasksIDs() []uint8 {
	result := make([]uint8, 0, len(t.TeamTasks))

	for id := range t.TeamTasks {
		result = append(result, id)
	}

	slices.Sort(result)

	return result
}
