package repository

import (
	"agent-platform/internal/model/ent"
	"agent-platform/internal/model/ent/user"
	"context"
	"fmt"
	"time"
)

// UserRepository handles user data access
type UserRepository struct {
	client *ent.Client
}

// NewUserRepository creates a new user repository
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, u *ent.User) (*ent.User, error) {
	builder := r.client.User.
		Create().
		SetID(u.ID).
		SetUsername(u.Username).
		SetEmail(u.Email).
		SetPasswordHash(u.PasswordHash).
		SetRole(u.Role).
		SetStatus(u.Status)

	if u.Metadata != nil {
		builder = builder.SetMetadata(u.Metadata)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	return created, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*ent.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("failed querying user: %w", err)
	}

	return u, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.Email(email)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, fmt.Errorf("failed querying user: %w", err)
	}

	return u, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*ent.User, error) {
	u, err := r.client.User.
		Query().
		Where(user.Username(username)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found with username: %s", username)
		}
		return nil, fmt.Errorf("failed querying user: %w", err)
	}

	return u, nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, page, pageSize int32, role, status string) ([]*ent.User, int, error) {
	query := r.client.User.Query()

	// Apply filters
	if role != "" {
		query = query.Where(user.Role(role))
	}
	if status != "" {
		query = query.Where(user.Status(status))
	}

	// Get total count
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting users: %w", err)
	}

	// Apply pagination
	offset := int((page - 1) * pageSize)
	users, err := query.
		Order(ent.Desc(user.FieldCreatedAt)).
		Offset(offset).
		Limit(int(pageSize)).
		All(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed listing users: %w", err)
	}

	return users, total, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) (*ent.User, error) {
	updateQuery := r.client.User.UpdateOneID(id)

	for key, value := range updates {
		switch key {
		case "username":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetUsername(v)
			}
		case "email":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetEmail(v)
			}
		case "password_hash":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetPasswordHash(v)
			}
		case "role":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetRole(v)
			}
		case "status":
			if v, ok := value.(string); ok {
				updateQuery = updateQuery.SetStatus(v)
			}
		case "metadata":
			if v, ok := value.(map[string]interface{}); ok {
				updateQuery = updateQuery.SetMetadata(v)
			}
		}
	}

	updated, err := updateQuery.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("failed updating user: %w", err)
	}

	return updated, nil
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.client.User.
		UpdateOneID(id).
		SetLastLoginAt(now).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed updating last login: %w", err)
	}

	return nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	err := r.client.User.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user not found: %s", id)
		}
		return fmt.Errorf("failed deleting user: %w", err)
	}

	return nil
}

// EmailExists checks if an email is already registered
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	count, err := r.client.User.
		Query().
		Where(user.Email(email)).
		Count(ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UsernameExists checks if a username is already taken
func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	count, err := r.client.User.
		Query().
		Where(user.Username(username)).
		Count(ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
