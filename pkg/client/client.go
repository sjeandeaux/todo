package client

import (
	"context"

	"github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/sjeandeaux/todo/todo-grpc/v1"
)

//ToDoManager manages your todo
type ToDoManager struct {
	Client pb.ToDoServiceClient
}

// ToDo a todo
type ToDo struct {
	//ID of todo
	ID string
	//Title of todo
	Title string
	//Description of todo
	Description string
	//State of todo
	State string
	//Tags of todo
	Tags []string
	//Reminder of todo
	Reminder int64
}

// Create create a toDo
func (m *ToDoManager) Create(cxt context.Context, toDo ToDo) (string, error) {
	state, ok := pb.ToDo_State_value[toDo.State]
	if !ok {
		state = int32(pb.ToDo_NOT_STARTED)
	}
	request := &pb.CreateRequest{
		ToDo: &pb.ToDo{
			Id:          toDo.ID,
			Title:       toDo.Title,
			Description: toDo.Description,
			Tags:        toDo.Tags,
			Reminder:    &timestamp.Timestamp{Seconds: int64(toDo.Reminder)},
			State:       pb.ToDo_State(state),
		},
	}

	response, err := m.Client.Create(cxt, request)
	if err != nil {
		return "", err
	}
	return response.GetId(), nil
}

// Update update a todo
func (m *ToDoManager) Update(cxt context.Context, toDo ToDo) (bool, error) {
	state, ok := pb.ToDo_State_value[toDo.State]
	if !ok {
		state = int32(pb.ToDo_NOT_STARTED)
	}
	request := &pb.UpdateRequest{
		ToDo: &pb.ToDo{
			Id:          toDo.ID,
			Title:       toDo.Title,
			Description: toDo.Description,
			Tags:        toDo.Tags,
			Reminder:    &timestamp.Timestamp{Seconds: int64(toDo.Reminder)},
			State:       pb.ToDo_State(state),
		},
	}

	response, err := m.Client.Update(cxt, request)
	if err != nil {
		return false, err
	}
	return response.Updated > 0, nil
}

// Delete a todo
func (m *ToDoManager) Delete(cxt context.Context, id string) (bool, error) {
	request := &pb.DeleteRequest{
		Id: id,
	}
	response, err := m.Client.Delete(cxt, request)
	if err != nil {
		return false, err
	}
	return response.Deleted > 0, nil
}

// Search todos
func (m *ToDoManager) Read(cxt context.Context, id string) (*ToDo, error) {
	request := &pb.ReadRequest{
		Id: id,
	}
	response, err := m.Client.Read(cxt, request)
	if err != nil {
		return nil, err
	}

	if todo := response.GetToDo(); todo != nil {
		return &ToDo{
			ID:          todo.GetId(),
			Title:       todo.GetTitle(),
			Description: todo.GetDescription(),
			Tags:        todo.GetTags(),
			State:       todo.GetState().String(),
			Reminder:    todo.GetReminder().GetSeconds(),
		}, nil
	}

	return nil, nil

}

// Search todos
func (m *ToDoManager) Search(cxt context.Context, pattern string, tags []string, states []string) ([]ToDo, error) {
	request := &pb.SearchRequest{
		Pattern: pattern,
		Tags:    tags,
		States:  toState(states),
	}
	response, err := m.Client.Search(cxt, request)
	if err != nil {
		return nil, err
	}

	result := []ToDo{}
	for _, todo := range response.GetToDos() {
		result = append(result, ToDo{
			ID:          todo.GetId(),
			Title:       todo.GetTitle(),
			Description: todo.GetDescription(),
			Tags:        todo.GetTags(),
			State:       todo.GetState().String(),
			Reminder:    todo.GetReminder().GetSeconds(),
		})
	}
	return result, nil
}

func toState(states []string) []pb.ToDo_State {
	result := []pb.ToDo_State{}
	for _, s := range states {
		if state, ok := pb.ToDo_State_value[s]; ok {
			result = append(result, pb.ToDo_State(state))
		}
	}
	return result
}
