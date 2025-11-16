package grpc

import (
	"context"
	"fmt"

	pb "agent-platform/gen/go"
	"agent-platform/internal/auth"
	"agent-platform/internal/model/ent"
	entuser "agent-platform/internal/model/ent/user"
	"agent-platform/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserServer gRPC User service implementation
type UserServer struct {
	pb.UnimplementedUserServiceServer
	client     *ent.Client
	repo       *repository.UserRepository
	jwtService *auth.JWTService
}

// NewUserServer creates a new User service instance
func NewUserServer(client *ent.Client, jwtService *auth.JWTService) *UserServer {
	return &UserServer{
		client:     client,
		repo:       repository.NewUserRepository(client),
		jwtService: jwtService,
	}
}

// Register registers a new user
func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Validate input
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Check if email already exists
	emailExists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check email existence")
	}
	if emailExists {
		return nil, status.Error(codes.AlreadyExists, "email already registered")
	}

	// Check if username already exists
	usernameExists, err := s.repo.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check username existence")
	}
	if usernameExists {
		return nil, status.Error(codes.AlreadyExists, "username already taken")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// Create user
	user := &ent.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "user", // default role
		Status:       "active",
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create user: %v", err))
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(created.ID, created.Username, created.Email, created.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	// Convert to protobuf
	pbUser := &pb.User{
		Id:       created.ID,
		Username: created.Username,
		Email:    created.Email,
		Role:     created.Role,
		Status:   created.Status,
	}

	if !created.LastLoginAt.IsZero() {
		pbUser.LastLoginAt = timestamppb.New(created.LastLoginAt)
	}
	if !created.CreatedAt.IsZero() {
		pbUser.CreatedAt = timestamppb.New(created.CreatedAt)
	}
	if !created.UpdatedAt.IsZero() {
		pbUser.UpdatedAt = timestamppb.New(created.UpdatedAt)
	}

	return &pb.RegisterResponse{
		User:  pbUser,
		Token: token,
	}, nil
}

// Login authenticates a user
func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Validate input
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	// Verify password
	if !auth.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, status.Error(codes.PermissionDenied, "user account is not active")
	}

	// Update last login time
	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log but don't fail the login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	// Convert to protobuf
	pbUser := &pb.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
	}

	if !user.LastLoginAt.IsZero() {
		pbUser.LastLoginAt = timestamppb.New(user.LastLoginAt)
	}
	if !user.CreatedAt.IsZero() {
		pbUser.CreatedAt = timestamppb.New(user.CreatedAt)
	}
	if !user.UpdatedAt.IsZero() {
		pbUser.UpdatedAt = timestamppb.New(user.UpdatedAt)
	}

	return &pb.LoginResponse{
		User:  pbUser,
		Token: token,
	}, nil
}

// GetProfile gets the current user's profile
func (s *UserServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.User, error) {
	// Get user ID from context (set by auth middleware)
	userID := auth.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Get user from database
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	// Convert to protobuf
	pbUser := &pb.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
	}

	if !user.LastLoginAt.IsZero() {
		pbUser.LastLoginAt = timestamppb.New(user.LastLoginAt)
	}
	if !user.CreatedAt.IsZero() {
		pbUser.CreatedAt = timestamppb.New(user.CreatedAt)
	}
	if !user.UpdatedAt.IsZero() {
		pbUser.UpdatedAt = timestamppb.New(user.UpdatedAt)
	}

	return pbUser, nil
}

// UpdateProfile updates the current user's profile
func (s *UserServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.User, error) {
	// Get user ID from context
	userID := auth.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if req.Username != "" {
		// Check if username is already taken by another user
		existingUser, err := s.client.User.Query().
			Where(entuser.Username(req.Username)).
			Where(entuser.IDNEQ(userID)).
			First(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return nil, status.Error(codes.Internal, "failed to check username")
		}
		if existingUser != nil {
			return nil, status.Error(codes.AlreadyExists, "username already taken")
		}
		updates["username"] = req.Username
	}

	if req.Email != "" {
		// Check if email is already taken by another user
		existingUser, err := s.client.User.Query().
			Where(entuser.Email(req.Email)).
			Where(entuser.IDNEQ(userID)).
			First(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return nil, status.Error(codes.Internal, "failed to check email")
		}
		if existingUser != nil {
			return nil, status.Error(codes.AlreadyExists, "email already registered")
		}
		updates["email"] = req.Email
	}

	// Update user
	updated, err := s.repo.Update(ctx, userID, updates)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	// Convert to protobuf
	pbUser := &pb.User{
		Id:       updated.ID,
		Username: updated.Username,
		Email:    updated.Email,
		Role:     updated.Role,
		Status:   updated.Status,
	}

	if !updated.LastLoginAt.IsZero() {
		pbUser.LastLoginAt = timestamppb.New(updated.LastLoginAt)
	}
	if !updated.CreatedAt.IsZero() {
		pbUser.CreatedAt = timestamppb.New(updated.CreatedAt)
	}
	if !updated.UpdatedAt.IsZero() {
		pbUser.UpdatedAt = timestamppb.New(updated.UpdatedAt)
	}

	return pbUser, nil
}

// ChangePassword changes the current user's password
func (s *UserServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*emptypb.Empty, error) {
	// Get user ID from context
	userID := auth.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Validate input
	if req.OldPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "old password is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new password is required")
	}

	// Get user
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	// Verify old password
	if !auth.VerifyPassword(req.OldPassword, user.PasswordHash) {
		return nil, status.Error(codes.PermissionDenied, "old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// Update password
	updates := map[string]interface{}{
		"password_hash": hashedPassword,
	}
	_, err = s.repo.Update(ctx, userID, updates)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update password")
	}

	return &emptypb.Empty{}, nil
}
