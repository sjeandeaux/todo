package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	pb "github.com/sjeandeaux/todo/todo-grpc/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	database   = "challenge"
	collection = "todo"
)

func getMongoClient(url string) (*mongo.Client, error) {
	return mongo.Connect(context.Background(), options.Client().ApplyURI(url))
}

// ToDoServiceServer manages the todos list
type ToDoServiceServer struct {
	client         *mongo.Client
	todoCollection *mongo.Collection
	ctx            context.Context
}

const (
	keyID          = "_id"
	keyTitle       = "title"
	keyDescription = "description"
	keyTags        = "tags"
	keyState       = "state"
	keyReminder    = "reminder"
)

//Use in create
type todoInMongo struct {
	Title       string
	Description string
	Tags        []string
	Reminder    int64
	State       string
}

//Use in read
type todoInMongoWithID struct {
	ID   primitive.ObjectID `bson:"_id"`
	Todo todoInMongo        `bson:"inline"`
}

func newTodo(r *pb.ToDo) (t todoInMongo) {
	t.Title = r.GetTitle()
	t.Description = r.GetDescription()
	t.Tags = r.GetTags()
	if reminder := r.GetReminder(); reminder != nil {
		t.Reminder = reminder.GetSeconds()
	}
	t.State = r.GetState().String()
	return
}

func (t *todoInMongoWithID) todo() *pb.ToDo {
	return &pb.ToDo{
		Id:          t.ID.Hex(),
		Title:       t.Todo.Title,
		Description: t.Todo.Description,
		Tags:        t.Todo.Tags,
		State:       pb.ToDo_State(pb.ToDo_State_value[t.Todo.State]),
		Reminder:    &timestamp.Timestamp{Seconds: t.Todo.Reminder},
	}
}

// validate the implementation
var _ pb.ToDoServiceServer = &ToDoServiceServer{}

// NewToDoServiceServer TODO
func NewToDoServiceServer(ctx context.Context, url string) (*ToDoServiceServer, error) {
	client, err := getMongoClient(url)
	if err != nil {
		return nil, err
	}
	database := client.Database(database)
	return &ToDoServiceServer{
		client:         client,
		todoCollection: database.Collection(collection),
	}, nil
}

// Close TODO
func (s *ToDoServiceServer) Close() {
	if s.client != nil {
		s.client.Disconnect(s.ctx)
	}
}

// HealthChecher the checker
func (s *ToDoServiceServer) HealthChecher() func() error {
	return func() error {
		if s.client == nil {
			return errors.New("mongo client is nil")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := s.client.Ping(ctx, nil); err != nil {
			return fmt.Errorf("ping failure %v", err)
		}
		return nil
	}
}

// Create a todo
func (s *ToDoServiceServer) Create(ctx context.Context, r *pb.CreateRequest) (*pb.CreateResponse, error) {
	todo := newTodo(r.ToDo)
	result, err := s.todoCollection.InsertOne(ctx, todo)
	if err != nil {
		return nil, err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return &pb.CreateResponse{
			Id: oid.Hex(),
		}, nil
	}
	return nil, errors.New("unexpected error")
}

// Read a todo
func (s *ToDoServiceServer) Read(ctx context.Context, r *pb.ReadRequest) (*pb.ReadResponse, error) {
	id, err := primitive.ObjectIDFromHex(r.GetId())
	if err != nil {
		return nil, err
	}
	result := s.todoCollection.FindOne(ctx, bson.M{keyID: id})

	objMongo := &todoInMongoWithID{}

	if err := result.Decode(objMongo); err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.ReadResponse{}, nil
		}
		return nil, err
	}
	return &pb.ReadResponse{
		ToDo: objMongo.todo(),
	}, nil
}

// Update a todo
func (s *ToDoServiceServer) Update(ctx context.Context, r *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id, err := primitive.ObjectIDFromHex(r.ToDo.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.M{keyID: bson.M{"$eq": id}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: keyTags, Value: r.GetToDo().GetTags()}}},
		primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: keyState, Value: r.GetToDo().GetState().String()}}},
	}

	if title := r.GetToDo().GetTitle(); title != "" {
		update = append(update, primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: keyTitle, Value: r.GetToDo().GetTitle()}}})
	}

	if description := r.GetToDo().GetDescription(); description != "" {
		update = append(update, primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: keyDescription, Value: description}}})
	}

	if reminder := r.GetToDo().GetReminder(); reminder != nil {
		update = append(update, primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: keyReminder, Value: reminder.GetSeconds()}}})
	}

	result, err := s.todoCollection.UpdateOne(
		ctx,
		filter,
		update,
	)

	if err != nil {
		return nil, err
	}
	return &pb.UpdateResponse{Updated: result.ModifiedCount}, nil

}

// Delete a todo
func (s *ToDoServiceServer) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	id, err := primitive.ObjectIDFromHex(r.GetId())
	if err != nil {
		return nil, err
	}
	result, err := s.todoCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{
		Deleted: result.DeletedCount,
	}, nil
}

// Search todos
func (s *ToDoServiceServer) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {

	filter := bson.D{}

	if pattern := r.GetPattern(); pattern != "" {
		filter = append(filter, bson.E{Key: keyDescription, Value: primitive.Regex{Pattern: pattern}})
	}

	if states := r.GetStates(); len(states) > 0 {
		filter = append(filter, bson.E{Key: keyState, Value: bson.D{{Key: "$in", Value: toString(states)}}})
	}

	if tags := r.GetTags(); len(tags) > 0 {
		filter = append(filter, bson.E{Key: keyTags, Value: bson.D{{Key: "$all", Value: tags}}})
	}
	log.Info(filter)
	cur, err := s.todoCollection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	elements := []*pb.ToDo{}
	for cur.Next(ctx) {
		var result todoInMongoWithID
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		log.Debug(result)
		elements = append(elements, result.todo())
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &pb.SearchResponse{ToDos: elements}, nil

}

func toString(values []pb.ToDo_State) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = value.String()
	}
	return result
}
