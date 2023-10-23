package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/kriskelly/dating-app-example/internal/graph/generated"
	"github.com/kriskelly/dating-app-example/internal/model"
)

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.User, error) {
	uid, err := model.LoginUser(ctx, r.dgraphClient, email, password)
	if err != nil {
		return nil, err
	}
	r.sessionManager.Put(ctx, "userID", uid)
	fields := getFieldStr(ctx)
	options := &model.FindUserOptions{
		Fields: fields,
	}
	user, err := model.FindUser(ctx, r.dgraphClient, *uid, options)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *mutationResolver) Signup(ctx context.Context, input model.NewUser) (*model.User, error) {
	newUser := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}
	err := model.CreateUser(ctx, r.dgraphClient, newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (r *mutationResolver) LikeUser(ctx context.Context, userID string) (*model.LikedResponse, error) {
	currentUser, err := r.getCurrentUser(ctx, nil)
	if err != nil {
		return nil, err
	}

	return model.LikeUser(ctx, r.dgraphClient, currentUser, userID)
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	fields := getFieldStr(ctx)
	options := &model.FindUserOptions{Fields: fields}
	return r.getCurrentUser(ctx, options)
}

func (r *queryResolver) Matches(ctx context.Context) ([]*model.User, error) {
	fields := getFieldStr(ctx)
	options := &model.FindUserOptions{
		Fields: fields,
	}
	// Fetch current user with default values
	currentUser, err := r.getCurrentUser(ctx, nil)
	if err != nil {
		return nil, err
	}
	return model.FindMatches(ctx, r.dgraphClient, currentUser, options)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
