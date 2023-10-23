package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/google/uuid"

	"github.com/kriskelly/dating-app-example/internal/dgraph"
)

var (
	defaultUserFields = "id uid name"
)

// User model
type User struct {
	ID       string   `json:"id"`
	UID      string   `json:"uid,omitempty"`
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	DType    []string `json:"dgraph.type,omitempty"`
}

type loginResponse struct {
	Login []struct {
		ID      string
		Success bool
	}
}

type userResponse struct {
	User []*User
}

type matchUpsertResponse struct {
	Res []struct {
		UID string
	}
}

// FindUserOptions Options for FindUser()
type FindUserOptions struct {
	Fields string
}

// FindUser find the user based on their ID.
func FindUser(ctx context.Context, client *dgraph.Client, id string, options *FindUserOptions) (*User, error) {
	var fields string
	if options == nil {
		fields = defaultUserFields
	} else {
		fields = options.Fields
	}
	str := fmt.Sprintf(`{
		var(func: eq(id, %s)) {
			u as uid
		}		
		user(func: uid(u)) {
			%s
		}
	}`, id, fields)

	data, err := client.NewTxn().Query(ctx, str)
	if err != nil {
		return nil, err
	}

	var decode userResponse
	if err := json.Unmarshal(data.Json, &decode); err != nil {
		return nil, err
	}
	if len(decode.User) == 0 {
		return nil, errors.New("User not found")
	}
	return decode.User[0], nil
}

// LoginUser verifies that the email and password are valid and returns the UID of the user.
func LoginUser(ctx context.Context, client *dgraph.Client, email string, password string) (*string, error) {
	data, err := client.NewTxn().QueryWithVars(ctx, `
		query Login($email: string, $password: string) {
			login(func: eq(email, $email)) {
				id
				success: checkpwd(password, $password)
		   }   
		}`, map[string]string{"$email": email, "$password": password})

	if err != nil {
		return nil, err
	}

	var decode loginResponse
	if err := json.Unmarshal(data.Json, &decode); err != nil {
		return nil, err
	}
	if len(decode.Login) > 0 && decode.Login[0].Success {
		return &decode.Login[0].ID, nil
	} else {
		return nil, errors.New("Login failed")
	}
}

// CreateUser will persist the user in the Dgraph database
func CreateUser(ctx context.Context, client *dgraph.Client, user *User) error {
	user.UID = "_:newuser"
	user.DType = []string{"User"}

	id := uuid.New().String()
	user.ID = id

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(user)
	if err != nil {
		return err
	}

	mu.SetJson = pb
	response, err := client.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return err
	}
	uid := response.Uids["newuser"]
	user.UID = uid
	return nil
}

// LikeUser Like a user
func LikeUser(ctx context.Context, client *dgraph.Client, currentUser *User, targetUserID string) (*LikedResponse, error) {
	targetUser, err := FindUser(ctx, client, targetUserID, nil)
	if err != nil {
		return nil, err
	}
	mu := &api.Mutation{
		CommitNow: true,
		SetNquads: []byte(fmt.Sprintf("<%s> <liked> <%s> . ", currentUser.UID, targetUser.UID)),
	}
	if _, err := client.NewTxn().Mutate(ctx, mu); err != nil {
		return nil, err
	}

	// If the target user already liked the current user, make it a match
	query := fmt.Sprintf(`{
		is_matched as var(func:uid(<%s>)) @filter(uid_in(liked, <%s>)) {
			uid
		}
		res(func: uid(is_matched)) {
			uid
		}
	}`, targetUser.UID, currentUser.UID)
	req := &api.Request{
		Query: query,
		Mutations: []*api.Mutation{
			&api.Mutation{
				CommitNow: true,
				Cond:      `@if(eq(len(is_matched), 1))`,
				SetNquads: []byte(fmt.Sprintf(`
					<%s> <matched> <%s> .
					<%s> <matched> <%s> .
				`, currentUser.UID, targetUser.UID,
					targetUser.UID, currentUser.UID)),
			},
		},
	}
	txn := client.NewTxn()
	resp, err := txn.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	var decode matchUpsertResponse
	if err := json.Unmarshal(resp.Json, &decode); err != nil {
		return nil, err
	}
	isMatched := len(decode.Res) > 0

	if err := txn.Commit(ctx); err != nil {
		return nil, err
	}
	return &LikedResponse{Success: true, Matched: isMatched}, nil
}

// FindMatches Find the matches for the given user
func FindMatches(ctx context.Context, client *dgraph.Client, currentUser *User, options *FindUserOptions) ([]*User, error) {
	var fields string
	if options == nil {
		fields = defaultUserFields
	} else {
		fields = options.Fields
	}

	str := fmt.Sprintf(`
	query matches($me: string) {
		user(func: has(matched)) @filter(uid_in(~matched, $me)) {
			%s
		}
	}`, fields)
	data, err := client.NewTxn().QueryWithVars(ctx, str,
		map[string]string{"$me": currentUser.UID})
	if err != nil {
		return nil, err
	}
	var decode userResponse
	if err := json.Unmarshal(data.Json, &decode); err != nil {
		return nil, err
	}
	return decode.User, nil
}
