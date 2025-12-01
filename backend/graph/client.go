package graph

import (
	"context"
	"todo-app/graph/model"
	"todo-app/models"
)

type Client struct {
	resolver *Resolver
}

func NewClient(resolver *Resolver) *Client {
	return &Client{resolver: resolver}
}

func (c *Client) CreateUser(ctx context.Context, email, password string) (*models.User, error) {
	mutation := &mutationResolver{c.resolver}
	return mutation.CreateUser(ctx, model.CreateUserInput{
		Email:    email,
		Password: password,
	})
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := &queryResolver{c.resolver}
	return query.UserByEmail(ctx, email)
}

func (c *Client) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := &queryResolver{c.resolver}
	return query.User(ctx, id)
}

func (c *Client) GetUsers(ctx context.Context) ([]*models.User, error) {
	query := &queryResolver{c.resolver}
	return query.Users(ctx)
}

func (c *Client) DeleteUser(ctx context.Context, id string) (bool, error) {
	mutation := &mutationResolver{c.resolver}
	return mutation.DeleteUser(ctx, id)
}

func (c *Client) UpdateUserAdmin(ctx context.Context, id string, isAdmin bool) (*models.User, error) {
	mutation := &mutationResolver{c.resolver}
	return mutation.UpdateUserAdmin(ctx, id, model.UpdateUserAdminInput{
		IsAdmin: isAdmin,
	})
}

func (c *Client) GetUserCount(ctx context.Context) (int, error) {
	query := &queryResolver{c.resolver}
	return query.UserCount(ctx)
}

func (c *Client) CreateTodo(ctx context.Context, userID, title, description string) (*models.Todo, error) {
	mutation := &mutationResolver{c.resolver}
	return mutation.CreateTodo(ctx, userID, model.CreateTodoInput{
		Title:       title,
		Description: description,
	})
}

func (c *Client) GetTodos(ctx context.Context, userID string) ([]*models.Todo, error) {
	query := &queryResolver{c.resolver}
	return query.Todos(ctx, userID)
}

func (c *Client) GetTodo(ctx context.Context, id, userID string) (*models.Todo, error) {
	query := &queryResolver{c.resolver}
	return query.Todo(ctx, id, userID)
}

func (c *Client) UpdateTodo(ctx context.Context, id, userID string, title, description *string, completed *bool) (*models.Todo, error) {
	mutation := &mutationResolver{c.resolver}
	return mutation.UpdateTodo(ctx, id, userID, model.UpdateTodoInput{
		Title:       title,
		Description: description,
		Completed:   completed,
	})
}

func (c *Client) DeleteTodo(ctx context.Context, id, userID string) (bool, error) {
	mutation := &mutationResolver{c.resolver}
	return mutation.DeleteTodo(ctx, id, userID)
}
