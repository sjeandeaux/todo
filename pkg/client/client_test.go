package client_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sjeandeaux/todo/pkg/client"
	pb "github.com/sjeandeaux/todo/todo-grpc/v1"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

var _ = Describe("Todo", func() {
	var (
		manager *client.ToDoManager
		mock    *mockToDoServiceClient
	)

	BeforeEach(func() {
		mock = &mockToDoServiceClient{}
		manager = &client.ToDoManager{
			Client: mock,
		}
	})

	Describe("Create", func() {
		Context("With a good todo", func() {
			It("should create one", func() {
				mock.expectedRequest = &pb.CreateRequest{
					ToDo: &pb.ToDo{
						Id:          "id",
						Title:       "title",
						Description: "descrption",
						Tags:        []string{"", ""},
						Reminder:    &timestamp.Timestamp{Seconds: int64(666)},
						State:       pb.ToDo_IN_PROGRESS,
					},
				}
				mock.response = &pb.CreateResponse{Id: "id-636"}

				Ω(manager.Create(context.TODO(), client.ToDo{
					ID:          "id",
					Title:       "title",
					Description: "descrption",
					Tags:        []string{"", ""},
					Reminder:    666,
					State:       "IN_PROGRESS",
				})).Should(Equal("id-636"))
			})
		})

		Context("With an issue", func() {
			It("should fail", func() {
				mock.expectedRequest = &pb.CreateRequest{
					ToDo: &pb.ToDo{
						Id:          "id",
						Title:       "title",
						Description: "descrption",
						Tags:        []string{"go", "rust"},
						Reminder:    &timestamp.Timestamp{Seconds: int64(666)},
						State:       pb.ToDo_IN_PROGRESS,
					},
				}
				mock.err = errors.New("error")
				response, err := manager.Create(context.TODO(), client.ToDo{
					ID:          "id",
					Title:       "title",
					Description: "descrption",
					Tags:        []string{"go", "rust"},
					Reminder:    666,
					State:       "IN_PROGRESS",
				})

				Ω(response).Should(Equal(""))
				Ω(err).Should(Equal(mock.err))
			})
		})
	})

	Describe("Read", func() {
		Context("With a todo which exists", func() {
			It("should get the todo", func() {
				mock.expectedRequest = &pb.ReadRequest{
					Id: "id",
				}
				mock.response = &pb.ReadResponse{
					ToDo: &pb.ToDo{
						Id:          "id",
						Title:       "title",
						Description: "descrption",
						Tags:        []string{"go", "rust"},
						Reminder:    &timestamp.Timestamp{Seconds: int64(666)},
						State:       pb.ToDo_IN_PROGRESS,
					},
				}

				Ω(manager.Read(context.TODO(), "id")).Should(Equal(&client.ToDo{
					ID:          "id",
					Title:       "title",
					Description: "descrption",
					Tags:        []string{"go", "rust"},
					Reminder:    666,
					State:       "IN_PROGRESS",
				}))
			})
		})

		Context("With a todo which doesn't exist", func() {
			It("should get nothing", func() {
				mock.expectedRequest = &pb.ReadRequest{
					Id: "id",
				}
				mock.response = &pb.ReadResponse{
					ToDo: nil,
				}

				Ω(manager.Read(context.TODO(), "id")).Should(BeNil())
			})
		})

		Context("With an issue", func() {
			It("should fail", func() {
				mock.expectedRequest = &pb.ReadRequest{
					Id: "id",
				}
				mock.err = errors.New("error")
				response, err := manager.Read(context.TODO(), "id")
				Ω(response).Should(BeNil())
				Ω(err).Should(Equal(mock.err))
			})
		})
	})

	Describe("Update", func() {
		Context("With a good todo which is updated", func() {
			It("should be true", func() {
				mock.expectedRequest = &pb.UpdateRequest{
					ToDo: &pb.ToDo{
						Id:          "id",
						Title:       "title",
						Description: "descrption",
						Tags:        []string{"", ""},
						Reminder:    &timestamp.Timestamp{Seconds: int64(666)},
						State:       pb.ToDo_IN_PROGRESS,
					},
				}
				mock.response = &pb.UpdateResponse{Updated: 6}

				Ω(manager.Update(context.TODO(), client.ToDo{
					ID:          "id",
					Title:       "title",
					Description: "descrption",
					Tags:        []string{"", ""},
					Reminder:    666,
					State:       "IN_PROGRESS",
				})).Should(Equal(true))
			})
		})

		Context("With a good todo which is updated", func() {
			It("should be false", func() {
				mock.expectedRequest = &pb.UpdateRequest{
					ToDo: &pb.ToDo{
						Id:          "id",
						Title:       "title",
						Description: "descrption",
						Tags:        []string{"", ""},
						Reminder:    &timestamp.Timestamp{Seconds: int64(666)},
						State:       pb.ToDo_IN_PROGRESS,
					},
				}
				mock.response = &pb.UpdateResponse{Updated: 0}

				Ω(manager.Update(context.TODO(), client.ToDo{
					ID:          "id",
					Title:       "title",
					Description: "descrption",
					Tags:        []string{"", ""},
					Reminder:    666,
					State:       "IN_PROGRESS",
				})).Should(Equal(false))
			})
		})

		Context("With an issue", func() {
			It("should fail", func() {
				mock.expectedRequest = &pb.UpdateRequest{
					ToDo: &pb.ToDo{
						Id:          "id",
						Title:       "title",
						Description: "descrption",
						Tags:        []string{"", ""},
						Reminder:    &timestamp.Timestamp{Seconds: int64(666)},
						State:       pb.ToDo_IN_PROGRESS,
					},
				}
				mock.err = errors.New("error")
				response, err := manager.Update(context.TODO(), client.ToDo{
					ID:          "id",
					Title:       "title",
					Description: "descrption",
					Tags:        []string{"", ""},
					Reminder:    666,
					State:       "IN_PROGRESS",
				})

				Ω(response).Should(Equal(false))
				Ω(err).Should(Equal(mock.err))
			})
		})
	})
	Describe("Delete", func() {
		Context("With a good todo which is deleted", func() {
			It("should be true", func() {
				mock.expectedRequest = &pb.DeleteRequest{
					Id: "id",
				}
				mock.response = &pb.DeleteResponse{Deleted: 6}

				Ω(manager.Delete(context.TODO(), "id")).Should(Equal(true))
			})
		})

		Context("With a good todo which is updated", func() {
			It("should be false", func() {
				mock.expectedRequest = &pb.DeleteRequest{
					Id: "id",
				}
				mock.response = &pb.DeleteResponse{Deleted: 0}

				Ω(manager.Delete(context.TODO(), "id")).Should(Equal(false))
			})
		})

		Context("With an issue", func() {
			It("should fail", func() {
				mock.expectedRequest = &pb.DeleteRequest{
					Id: "id",
				}
				mock.err = errors.New("error")
				response, err := manager.Delete(context.TODO(), "id")

				Ω(response).Should(Equal(false))
				Ω(err).Should(Equal(mock.err))
			})
		})
	})

	Describe("Search", func() {
		Context("With todos", func() {
			It("should get the todos", func() {
				mock.expectedRequest = &pb.SearchRequest{
					Pattern: "pattern",
					Tags:    []string{"t1", "t2"},
					States:  []pb.ToDo_State{pb.ToDo_DONE},
				}
				mock.response = &pb.SearchResponse{
					ToDos: []*pb.ToDo{{
						Id:          "id",
						Title:       "title",
						Description: "descrption",
						Tags:        []string{"go", "rust"},
						Reminder:    &timestamp.Timestamp{Seconds: int64(666)},
						State:       pb.ToDo_IN_PROGRESS,
					}},
				}

				Ω(manager.Search(context.TODO(), "pattern", []string{"t1", "t2"}, []string{"DONE", "ignore"})).Should(Equal([]client.ToDo{{
					ID:          "id",
					Title:       "title",
					Description: "descrption",
					Tags:        []string{"go", "rust"},
					Reminder:    666,
					State:       "IN_PROGRESS",
				}}))
			})
		})

		Context("With an issue", func() {
			It("should fail", func() {
				mock.expectedRequest = &pb.SearchRequest{
					States: []pb.ToDo_State{},
				}
				mock.err = errors.New("error")
				response, err := manager.Search(context.TODO(), "", nil, nil)
				Ω(response).Should(BeNil())
				Ω(err).Should(Equal(mock.err))
			})
		})
	})
})
