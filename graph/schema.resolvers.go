package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"webmane_go/graph/model"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	// panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
	randNumber, _ := rand.Int(rand.Reader, big.NewInt(100))

	user := &model.User{
		ID: input.UserID,
	}

	todo := &model.Todo{
		Text: input.Text,
		ID:   fmt.Sprintf("T%d", randNumber),
		User: user,
	}
	r.todos = append(r.todos, todo)
	return todo, nil
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	// panic(fmt.Errorf("not implemented: Todos - todos"))
	return r.todos, nil
}

// User is the resolver for the user field.
func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	// panic(fmt.Errorf("not implemented: User - user"))
	return &model.User{ID: obj.UserID, Name: "user " + obj.User.ID}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Todo returns TodoResolver implementation.
func (r *Resolver) Todo() TodoResolver { return &todoResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
