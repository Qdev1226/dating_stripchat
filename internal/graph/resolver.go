package graph

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/alexedwards/scs/v2"
	"github.com/kriskelly/dating-app-example/internal/dgraph"
	"github.com/kriskelly/dating-app-example/internal/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	users          map[string]*model.User
	dgraphClient   *dgraph.Client
	sessionManager *scs.SessionManager
}

func NewResolver(sessionManager *scs.SessionManager, client *dgraph.Client) *Resolver {
	return &Resolver{
		dgraphClient:   client,
		users:          make(map[string]*model.User),
		sessionManager: sessionManager,
	}
}

func (r *Resolver) getCurrentUser(ctx context.Context, options *model.FindUserOptions) (*model.User, error) {
	rawCurrentUserID := r.sessionManager.Get(ctx, "userID")
	if rawCurrentUserID == nil {
		return nil, errors.New("No current user")
	}
	userID := rawCurrentUserID.(string)
	user, err := model.FindUser(ctx, r.dgraphClient, userID, options)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("No current user")
	}
	return user, nil
}

func getFieldStr(ctx context.Context) string {
	return getFieldStrRecursive(
		graphql.GetRequestContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
	)
}

func getFieldStrRecursive(ctx *graphql.RequestContext, fields []graphql.CollectedField) string {
	fieldStr := ""
	for _, field := range fields {
		fieldStr += field.Name
		if len(field.SelectionSet) > 0 {
			fieldStr += "{\n " + getFieldStrRecursive(ctx, graphql.CollectFields(ctx, field.SelectionSet, nil)) + " }"
		}
		fieldStr += "\n"
	}
	return fieldStr
}
