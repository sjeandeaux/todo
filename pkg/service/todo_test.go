package service_test

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	. "github.com/sjeandeaux/todo/pkg/service"
	pb "github.com/sjeandeaux/todo/todo-grpc/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type todoInMongo struct {
	Title       string
	Description string
	Tags        []string
	Reminder    int64
	State       string
}

var _ = Describe("Todo", func() {

	const (
		mongoURI     = "mongodb://devroot:devroot@localhost:27017/?authSource=admin"
		databaseName = "challenge"
		collection   = "todo"
	)

	var (
		server *ToDoServiceServer
		client *mongo.Client
	)

	//Before each test it cleans the database
	BeforeEach(func() {
		if testing.Short() {
			Skip("Skip in short mode (need database access)")
		}
		By("Clean the database, create a client mongo and a server")

		var err error
		client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
		Ω(err).NotTo(HaveOccurred())
		client.Database(databaseName).Collection(collection).Drop(context.TODO())

		//create the server
		server, err = NewToDoServiceServer(context.TODO(), "mongodb://devroot:devroot@localhost:27017/?authSource=admin")
		Ω(err).NotTo(HaveOccurred())
		Ω(server).ShouldNot(BeNil())
	})

	AfterEach(func() {
		server.Close()
	})

	//helper find a todo
	find := func(id string) (*todoInMongo, error) {
		objectID, _ := primitive.ObjectIDFromHex(id)
		result := client.Database(databaseName).Collection(collection).FindOne(context.TODO(), bson.M{"_id": objectID})
		actual := &todoInMongo{}
		return actual, result.Decode(actual)
	}

	//helper insert a todo
	insert := func(todo *todoInMongo) string {
		result, err := client.Database(databaseName).Collection(collection).InsertOne(context.TODO(), todo)
		Ω(err).NotTo(HaveOccurred())
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			return oid.Hex()
		}
		return ""
	}

	Describe("Create", func() {
		Context("With a todo payload", func() {
			It("should create one todo", func() {
				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}

				request := &pb.CreateRequest{
					ToDo: &pb.ToDo{
						Title:       "Challenge - todo",
						Description: "Should create a micro service with 12factor",
						Tags:        []string{"golang", "12factor", "k8s"},
						Reminder:    &timestamp.Timestamp{Seconds: 1573046180},
						State:       pb.ToDo_IN_PROGRESS,
					},
				}
				response, err := server.Create(context.TODO(), request)
				Ω(err).NotTo(HaveOccurred())
				Ω(response.GetId()).ShouldNot(BeEmpty())

				expected := &todoInMongo{
					Title:       "Challenge - todo",
					Description: "Should create a micro service with 12factor",
					Tags:        []string{"golang", "12factor", "k8s"},
					Reminder:    1573046180,
					State:       pb.ToDo_IN_PROGRESS.String(),
				}

				Ω(find(response.GetId())).Should(Equal(expected))
			})
		})
	})

	Describe("Read", func() {
		Context("An existing todo", func() {
			It("should return it", func() {
				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}
				id := insert(&todoInMongo{
					Title:       "Read - Challenge - todo",
					Description: "Read - Should create a micro service with 12factor",
					Tags:        []string{"golang", "12factor", "k8s"},
					Reminder:    1573046180,
					State:       pb.ToDo_DONE.String(),
				})

				request := &pb.ReadRequest{
					Id: id,
				}
				response, err := server.Read(context.TODO(), request)
				Ω(err).NotTo(HaveOccurred())

				expected := &pb.ToDo{
					Id:          id,
					Title:       "Read - Challenge - todo",
					Description: "Read - Should create a micro service with 12factor",
					Tags:        []string{"golang", "12factor", "k8s"},
					Reminder:    &timestamp.Timestamp{Seconds: 1573046180},
					State:       pb.ToDo_DONE,
				}

				Ω(response.GetToDo()).Should(Equal(expected))

			})
		})

		Context("With an inexistant id", func() {
			It("should fail", func() {
				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}

				request := &pb.ReadRequest{
					Id: "nope",
				}
				_, err := server.Read(context.TODO(), request)
				Ω(err).ShouldNot(BeNil())
			})
		})

		Context("An non existing todo", func() {
			It("should return nothing", func() {
				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}
				request := &pb.ReadRequest{
					Id: "5dc2d3d4aba443c197307ea2",
				}
				response, err := server.Read(context.TODO(), request)
				Ω(err).NotTo(HaveOccurred())
				Ω(response.GetToDo()).Should(BeNil())

			})
		})
	})

	Describe("Update", func() {
		Context("An existing todo", func() {
			It("should update it", func() {
				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}
				id := insert(&todoInMongo{
					Title:       "Should be updated - Challenge - todo",
					Description: "Should be updated - New Should create a micro service with 12factor",
					Tags:        []string{"golang", "12factor", "k8s"},
					Reminder:    1573046240,
					State:       pb.ToDo_NOT_STARTED.String(),
				})

				request := &pb.UpdateRequest{
					ToDo: &pb.ToDo{
						Id:          id,
						Title:       "New Challenge - todo",
						Description: "New Should create a micro service with 12factor",
						Tags:        []string{"golang", "12factor", "k8s", "ci/cd"},
						Reminder:    &timestamp.Timestamp{Seconds: 1573046666},
						State:       pb.ToDo_IN_PROGRESS,
					},
				}
				response, err := server.Update(context.TODO(), request)
				Ω(err).NotTo(HaveOccurred())
				Ω(response.GetUpdated()).Should(Equal(int64(1)))

				expected := &todoInMongo{
					Title:       "New Challenge - todo",
					Description: "New Should create a micro service with 12factor",
					Tags:        []string{"golang", "12factor", "k8s", "ci/cd"},
					Reminder:    1573046666,
					State:       pb.ToDo_IN_PROGRESS.String(),
				}

				Ω(find(id)).Should(Equal(expected))
			})
		})
		Context("With a bad id", func() {
			It("should fail", func() {

				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}

				request := &pb.UpdateRequest{
					ToDo: &pb.ToDo{
						Id:          "bad id",
						Title:       "New Challenge - todo",
						Description: "New Should create a micro service with 12factor",
						Tags:        []string{"golang", "12factor", "k8s", "ci/cd"},
						Reminder:    &timestamp.Timestamp{Seconds: 1573046666},
					},
				}

				_, err := server.Update(context.TODO(), request)
				Ω(err).ShouldNot(BeNil())
			})
		})

		Context("With a not existing id", func() {
			It("should update nothing", func() {

				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}

				request := &pb.UpdateRequest{
					ToDo: &pb.ToDo{
						Id:          "5dc2d3d4aba443c197307ea2",
						Title:       "New Challenge - todo",
						Description: "New Should create a micro service with 12factor",
						Tags:        []string{"golang", "12factor", "k8s", "ci/cd"},
						Reminder:    &timestamp.Timestamp{Seconds: 1573046666},
					},
				}

				response, err := server.Update(context.TODO(), request)
				Ω(err).NotTo(HaveOccurred())
				Ω(response.GetUpdated()).Should(Equal(int64(0)))
			})
		})
	})

	Describe("Delete", func() {
		Context("An existing todo", func() {
			It("should delete it", func() {
				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}
				id := insert(&todoInMongo{
					Title:       "Should be deleted - Challenge - todo",
					Description: "Should be deleted - New Should create a micro service with 12factor",
					Tags:        []string{"golang", "12factor", "k8s"},
					Reminder:    1573046240,
					State:       pb.ToDo_NOT_STARTED.String(),
				})

				request := &pb.DeleteRequest{
					Id: id,
				}
				response, err := server.Delete(context.TODO(), request)
				Ω(err).NotTo(HaveOccurred())
				Ω(response.GetDeleted()).Should(Equal(int64(1)))

				_, err = find(id)
				Ω(err).Should(Equal(mongo.ErrNoDocuments))
			})
		})
		Context("With a bad id", func() {
			It("should fail", func() {

				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}

				request := &pb.DeleteRequest{
					Id: "bad id",
				}

				_, err := server.Delete(context.TODO(), request)
				Ω(err).ShouldNot(BeNil())
			})
		})

		Context("With a not existing id", func() {
			It("should delete nothing", func() {

				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}

				request := &pb.DeleteRequest{
					Id: "5dc2d3d4aba443c197307ea2",
				}

				response, err := server.Delete(context.TODO(), request)
				Ω(err).NotTo(HaveOccurred())
				Ω(response.GetDeleted()).Should(Equal(int64(0)))
			})
		})
	})

	Describe("Search", func() {
		Context("With a good pattern, good tags and good state", func() {
			It("should return a todo", func() {
				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}
				id := insert(&todoInMongo{
					Title:       "Read - Challenge - todo",
					Description: "Read - Should create a micro service with 12factor",
					Tags:        []string{"golang", "12factor", "k8s"},
					Reminder:    1573046180,
					State:       pb.ToDo_DONE.String(),
				})

				request := &pb.SearchRequest{
					Pattern: "^Read.*",
					States:  []pb.ToDo_State{pb.ToDo_DONE},
					Tags:    []string{"golang", "12factor"},
				}
				response, err := server.Search(context.TODO(), request)
				Ω(err).NotTo(HaveOccurred())

				expected := &pb.ToDo{
					Id:          id,
					Title:       "Read - Challenge - todo",
					Description: "Read - Should create a micro service with 12factor",
					Tags:        []string{"golang", "12factor", "k8s"},
					Reminder:    &timestamp.Timestamp{Seconds: 1573046180},
					State:       pb.ToDo_DONE,
				}

				Ω(response.GetToDos()).Should(HaveLen(1))
				Ω(response.GetToDos()[0]).Should(Equal(expected))

			})
		})

		Context("With a pattern, tags and state which doesn't match", func() {
			It("should find it", func() {
				if testing.Short() {
					Skip("Skip in short mode (need database access)")
				}
				insert(&todoInMongo{
					Title:       "Read - Challenge - todo",
					Description: "Read - Should create a micro service with 12factor",
					Tags:        []string{"golang", "12factor", "k8s"},
					Reminder:    1573046180,
					State:       pb.ToDo_DONE.String(),
				})

				request := &pb.SearchRequest{
					Pattern: "^Write.*",
					States:  []pb.ToDo_State{pb.ToDo_DONE},
					Tags:    []string{"golang", "12factor"},
				}
				response, err := server.Search(context.TODO(), request)
				Ω(err).NotTo(HaveOccurred())
				Ω(response.GetToDos()).Should(HaveLen(0))
			})
		})
	})

})
