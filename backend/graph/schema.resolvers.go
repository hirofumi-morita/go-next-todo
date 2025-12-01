package graph

import (
        "context"
        "fmt"
        "strconv"
        "todo-app/graph/model"
        "todo-app/models"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*models.User, error) {
        user := &models.User{
                Email: input.Email,
        }

        if err := user.HashPassword(input.Password); err != nil {
                return nil, fmt.Errorf("failed to hash password: %w", err)
        }

        var count int64
        r.DB.Model(&models.User{}).Count(&count)
        if count == 0 {
                user.IsAdmin = true
        }

        if err := r.DB.Create(user).Error; err != nil {
                return nil, fmt.Errorf("failed to create user: %w", err)
        }

        return user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
        userID, err := strconv.ParseUint(id, 10, 64)
        if err != nil {
                return false, fmt.Errorf("invalid user ID: %w", err)
        }

        var user models.User
        if err := r.DB.First(&user, userID).Error; err != nil {
                return false, fmt.Errorf("user not found: %w", err)
        }

        if err := r.DB.Where("user_id = ?", userID).Delete(&models.Todo{}).Error; err != nil {
                return false, fmt.Errorf("failed to delete user's todos: %w", err)
        }

        result := r.DB.Delete(&user)
        if result.Error != nil {
                return false, fmt.Errorf("failed to delete user: %w", result.Error)
        }

        return result.RowsAffected > 0, nil
}

func (r *mutationResolver) UpdateUserAdmin(ctx context.Context, id string, input model.UpdateUserAdminInput) (*models.User, error) {
        userID, err := strconv.ParseUint(id, 10, 64)
        if err != nil {
                return nil, fmt.Errorf("invalid user ID: %w", err)
        }

        var user models.User
        if err := r.DB.First(&user, userID).Error; err != nil {
                return nil, fmt.Errorf("user not found: %w", err)
        }

        user.IsAdmin = input.IsAdmin
        if err := r.DB.Save(&user).Error; err != nil {
                return nil, fmt.Errorf("failed to update user: %w", err)
        }

        return &user, nil
}

func (r *mutationResolver) CreateTodo(ctx context.Context, userID string, input model.CreateTodoInput) (*models.Todo, error) {
        uid, err := strconv.ParseUint(userID, 10, 64)
        if err != nil {
                return nil, fmt.Errorf("invalid user ID: %w", err)
        }

        todo := &models.Todo{
                Title:       input.Title,
                Description: input.Description,
                UserID:      uint(uid),
        }

        if err := r.DB.Create(todo).Error; err != nil {
                return nil, fmt.Errorf("failed to create todo: %w", err)
        }

        return todo, nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, userID string, input model.UpdateTodoInput) (*models.Todo, error) {
        todoID, err := strconv.ParseUint(id, 10, 64)
        if err != nil {
                return nil, fmt.Errorf("invalid todo ID: %w", err)
        }

        uid, err := strconv.ParseUint(userID, 10, 64)
        if err != nil {
                return nil, fmt.Errorf("invalid user ID: %w", err)
        }

        var todo models.Todo
        if err := r.DB.Where("id = ? AND user_id = ?", todoID, uid).First(&todo).Error; err != nil {
                return nil, fmt.Errorf("todo not found: %w", err)
        }

        if input.Title != nil {
                todo.Title = *input.Title
        }
        if input.Description != nil {
                todo.Description = *input.Description
        }
        if input.Completed != nil {
                todo.Completed = *input.Completed
        }

        if err := r.DB.Save(&todo).Error; err != nil {
                return nil, fmt.Errorf("failed to update todo: %w", err)
        }

        return &todo, nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, id string, userID string) (bool, error) {
        todoID, err := strconv.ParseUint(id, 10, 64)
        if err != nil {
                return false, fmt.Errorf("invalid todo ID: %w", err)
        }

        uid, err := strconv.ParseUint(userID, 10, 64)
        if err != nil {
                return false, fmt.Errorf("invalid user ID: %w", err)
        }

        result := r.DB.Where("id = ? AND user_id = ?", todoID, uid).Delete(&models.Todo{})
        if result.Error != nil {
                return false, fmt.Errorf("failed to delete todo: %w", result.Error)
        }

        return result.RowsAffected > 0, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
        userID, err := strconv.ParseUint(id, 10, 64)
        if err != nil {
                return nil, fmt.Errorf("invalid user ID: %w", err)
        }

        var user models.User
        if err := r.DB.Preload("Todos").First(&user, userID).Error; err != nil {
                return nil, fmt.Errorf("user not found: %w", err)
        }

        return &user, nil
}

func (r *queryResolver) UserByEmail(ctx context.Context, email string) (*models.User, error) {
        var user models.User
        if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
                return nil, fmt.Errorf("user not found: %w", err)
        }

        return &user, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
        var users []*models.User
        if err := r.DB.Find(&users).Error; err != nil {
                return nil, fmt.Errorf("failed to fetch users: %w", err)
        }

        return users, nil
}

func (r *queryResolver) UserCount(ctx context.Context) (int, error) {
        var count int64
        if err := r.DB.Model(&models.User{}).Count(&count).Error; err != nil {
                return 0, fmt.Errorf("failed to count users: %w", err)
        }

        return int(count), nil
}

func (r *queryResolver) Todo(ctx context.Context, id string, userID string) (*models.Todo, error) {
        todoID, err := strconv.ParseUint(id, 10, 64)
        if err != nil {
                return nil, fmt.Errorf("invalid todo ID: %w", err)
        }

        uid, err := strconv.ParseUint(userID, 10, 64)
        if err != nil {
                return nil, fmt.Errorf("invalid user ID: %w", err)
        }

        var todo models.Todo
        if err := r.DB.Where("id = ? AND user_id = ?", todoID, uid).First(&todo).Error; err != nil {
                return nil, fmt.Errorf("todo not found: %w", err)
        }

        return &todo, nil
}

func (r *queryResolver) Todos(ctx context.Context, userID string) ([]*models.Todo, error) {
        uid, err := strconv.ParseUint(userID, 10, 64)
        if err != nil {
                return nil, fmt.Errorf("invalid user ID: %w", err)
        }

        var todos []*models.Todo
        if err := r.DB.Where("user_id = ?", uid).Find(&todos).Error; err != nil {
                return nil, fmt.Errorf("failed to fetch todos: %w", err)
        }

        return todos, nil
}

func (r *queryResolver) TodosByUser(ctx context.Context, userID string) ([]*models.Todo, error) {
        return r.Todos(ctx, userID)
}

func (r *todoResolver) ID(ctx context.Context, obj *models.Todo) (string, error) {
        return strconv.FormatUint(uint64(obj.ID), 10), nil
}

func (r *todoResolver) UserID(ctx context.Context, obj *models.Todo) (string, error) {
        return strconv.FormatUint(uint64(obj.UserID), 10), nil
}

func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
        return strconv.FormatUint(uint64(obj.ID), 10), nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

func (r *Resolver) Todo() TodoResolver { return &todoResolver{r} }

func (r *Resolver) User() UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
