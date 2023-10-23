package model

import (
	"context"
	"testing"

	"github.com/kriskelly/dating-app-example/internal/dgraph"
)

func TestCreateUser(t *testing.T) {
	dg := dgraph.NewClient()
	dg.Connect()
	defer dg.Close()

	ctx := context.Background()

	tests := map[string]struct {
		input *User
	}{
		"simple": {input: &User{Name: "Foo", Password: "foobar"}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if err := CreateUser(ctx, dg, test.input); err != nil {
				t.Fatal(err)
			}
			user, err := FindUser(ctx, dg, test.input.ID, nil)
			if err != nil {
				t.Fatal(err)
			}
			if user.Name != test.input.Name {
				t.Fatal("Could not find the user", test.input.Name)
			}
		})
	}
}

func TestFindUser(t *testing.T) {
	dg := dgraph.NewClient()
	dg.Connect()
	defer dg.Close()
	ctx := context.Background()

	subject := &User{Name: "Foo", Password: "Barfoo"}
	CreateUser(ctx, dg, subject)

	tests := map[string]struct {
		setup  func() *User
		expect func(t *testing.T, u *User)
	}{
		"happy path": {
			setup: func() *User {
				user, err := FindUser(ctx, dg, subject.ID, &FindUserOptions{
					Fields: "name",
				})
				if err != nil {
					t.Fatal(err)
				}
				return user
			},
			expect: func(t *testing.T, u *User) {
				if u.Name != "Foo" {
					t.Fatal("Did not return the correct user", u)
				}
			},
		},
		"with_fields": {
			setup: func() *User {
				user, err := FindUser(ctx, dg, subject.ID, &FindUserOptions{
					Fields: "id name",
				})
				if err != nil {
					t.Fatal(err)
				}
				return user
			},
			expect: func(t *testing.T, u *User) {
				if u.ID == "" {
					t.Fatal("Did not fetch the ID")
				}
				if u.Name == "" {
					t.Fatal("Did not fetch the name")
				}
				if u.UID != "" {
					t.Fatal("Should not have fetched the UID")
				}
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			u := test.setup()
			test.expect(t, u)
		})
	}
}

func TestLikeUser(t *testing.T) {
	// Tests:
	// 1) One user liking another user
	// 2) One user liking a user who has already liked them

	dg := dgraph.NewClient()
	dg.Connect()
	defer dg.Close()

	subject := &User{Name: "Foo", Password: "Barfoo"}
	object := &User{Name: "Bar", Password: "Bazfoo"}
	ctx := context.Background()
	CreateUser(ctx, dg, subject)
	CreateUser(ctx, dg, object)

	tests := map[string]struct {
		setup  func() *LikedResponse
		expect func(t *testing.T, resp *LikedResponse)
	}{
		"happy path": {
			setup: func() *LikedResponse {
				resp, err := LikeUser(ctx, dg, subject, object.ID)
				if err != nil {
					t.Fatal(err)
				}
				return resp
			},
			expect: func(t *testing.T, resp *LikedResponse) {
				if resp.Success != true {
					t.Fatal("Did not successfully like the user")
				}
			},
		},
		"matching": {
			setup: func() *LikedResponse {
				// Was liked already
				resp, err := LikeUser(ctx, dg, object, subject.ID)
				if err != nil {
					t.Fatal(err)
				}
				if resp.Success != true {
					t.Fatal("Did not successfully like the user in reverse (object -> subject)")
				}
				return resp
			},
			expect: func(t *testing.T, resp *LikedResponse) {
				if !resp.Matched {
					t.Fatal("Did not match with the other user")
				}
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.setup()
			resp, err := LikeUser(ctx, dg, subject, object.ID)
			if err != nil {
				t.Fatal(err)
			}
			test.expect(t, resp)
		})
	}
}

func TestFindMatches(t *testing.T) {
	dg := dgraph.NewClient()
	dg.Connect()
	defer dg.Close()

	subject := &User{Name: "Foo", Password: "Barfoo"}
	object := &User{Name: "Bar", Password: "Bazfoo"}
	ctx := context.Background()
	CreateUser(ctx, dg, subject)
	CreateUser(ctx, dg, object)
	LikeUser(ctx, dg, subject, object.ID)
	LikeUser(ctx, dg, object, subject.ID)

	tests := map[string]struct {
		setup  func() []*User
		expect func(t *testing.T, matches []*User)
	}{
		"happy path": {
			setup: func() []*User {
				matches, err := FindMatches(ctx, dg, subject, nil)
				if err != nil {
					t.Fatal(err)
				}
				return matches
			},
			expect: func(t *testing.T, matches []*User) {
				if len(matches) != 1 {
					t.Fatal("Did not find matches")
				}
			},
		},
		"with field context": {
			setup: func() []*User {
				matches, err := FindMatches(ctx, dg, subject, &FindUserOptions{
					Fields: "name",
				})
				if err != nil {
					t.Fatal(err)
				}
				return matches
			},
			expect: func(t *testing.T, matches []*User) {
				match := matches[0]
				if match.Name == "" {
					t.Fatal("Did not fetch name")
				}
				if match.UID != "" {
					t.Fatal("Did fetch UID")
				}
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			matches := test.setup()
			test.expect(t, matches)
		})
	}
}
