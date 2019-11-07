package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	pb "github.com/sjeandeaux/ori/todo-grpc/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Config to connect to database
type Config struct {
	Host     string
	Port     string
	Login    string
	Password string

	Database   string
	Collection string
}

func (c *Config) getMongoURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", c.Login, c.Password, c.Host, c.Port)
}

func (c *Config) getMongoClient() (*mongo.Client, error) {
	return mongo.Connect(context.Background(), options.Client().ApplyURI(c.getMongoURI()))
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

type todoInMongo struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string
	Description string
	Tags        []string
	Reminder    int64
	State       string
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

func (t *todoInMongo) todo() *pb.ToDo {
	return &pb.ToDo{
		Id:          t.ID.Hex(),
		Title:       t.Title,
		Description: t.Description,
		Tags:        t.Tags,
		State:       pb.ToDo_State(pb.ToDo_State_value[t.State]),
		Reminder:    &timestamp.Timestamp{Seconds: t.Reminder},
	}
}

// validate the implementation
var _ pb.ToDoServiceServer = &ToDoServiceServer{}

// NewToDoServiceServer TODO
func NewToDoServiceServer(ctx context.Context, c Config) (*ToDoServiceServer, error) {
	client, err := c.getMongoClient()
	if err != nil {
		return nil, err
	}
	database := client.Database(c.Database)
	return &ToDoServiceServer{
		client:         client,
		todoCollection: database.Collection(c.Collection),
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
		return s.client.Ping(ctx, nil)
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

	objMongo := &todoInMongo{}

	if err := result.Decode(objMongo); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
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
		primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: keyTitle, Value: r.GetToDo().GetTitle()}}},
		primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: keyDescription, Value: r.GetToDo().GetDescription()}}},
		primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: keyTags, Value: r.GetToDo().GetTags()}}},
		primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: keyState, Value: r.GetToDo().GetState().String()}}},
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
		var result todoInMongo
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