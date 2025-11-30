package usecase

import (
	"context"
	"errors"
	"reflect"
	"repo-guardian/internal/domain"
	"testing"
	"time"
)

type mockUserRepository struct {
	createFunc  func(ctx context.Context, user *domain.User) error
	getByIDFunc func(ctx context.Context, id int64) (*domain.User, error)
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockUserRepository) Delete(ctx context.Context, id int64) error {
	return m.deleteFunc(ctx, id)
}

func TestNewUserUsecase(t *testing.T) {
	repo := &mockUserRepository{}
	timeout := 5 * time.Second
	usecase := NewUserUsecase(repo, timeout)

	if usecase == nil {
		t.Errorf("NewUserUsecase returned nil")
	}

	expectedType := "*usecase.userUsecase"
	actualType := reflect.TypeOf(usecase).String()

	if actualType != expectedType {
		t.Errorf("NewUserUsecase returned wrong type. Expected %s, got %s", expectedType, actualType)
	}

	u, ok := usecase.(*userUsecase)
	if !ok {
		t.Errorf("NewUserUsecase did not return *userUsecase")
	}

	if u.userRepo != repo {
		t.Errorf("NewUserUsecase did not set userRepo correctly")
	}

	if u.contextTimeout != timeout {
		t.Errorf("NewUserUsecase did not set contextTimeout correctly")
	}
}

func TestUserUsecase_Register(t *testing.T) {
	type args struct {
		c    context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		mockRepo func() *mockUserRepository
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "success",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					createFunc: func(ctx context.Context, user *domain.User) error {
						return nil
					},
				}
			},
			args: args{
				c:    context.Background(),
				user: &domain.User{Name: "test"},
			},
			wantErr: false,
		},
		{
			name: "failure",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					createFunc: func(ctx context.Context, user *domain.User) error {
						return errors.New("create error")
					},
				}
			},
			args: args{
				c:    context.Background(),
				user: &domain.User{Name: "test"},
			},
			wantErr: true,
			err:     errors.New("create error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &userUsecase{
				userRepo:       tt.mockRepo(),
				contextTimeout: time.Second,
			}
			err := a.Register(tt.args.c, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("Expected error %v, got %v", tt.err, err)
			}

		})
	}
}

func TestUserUsecase_GetUser(t *testing.T) {
	type args struct {
		c  context.Context
		id int64
	}
	tests := []struct {
		name    string
		mockRepo func() *mockUserRepository
		args    args
		want    *domain.User
		wantErr bool
		err     error
	}{
		{
			name: "success",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					getByIDFunc: func(ctx context.Context, id int64) (*domain.User, error) {
						return &domain.User{ID: id, Name: "test"}, nil
					},
				}
			},
			args: args{
				c:  context.Background(),
				id: 1,
			},
			want:    &domain.User{ID: 1, Name: "test"},
			wantErr: false,
		},
		{
			name: "not found",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					getByIDFunc: func(ctx context.Context, id int64) (*domain.User, error) {
						return nil, errors.New("not found")
					},
				}
			},
			args: args{
				c:  context.Background(),
				id: 1,
			},
			want:    nil,
			wantErr: true,
			err:     errors.New("not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &userUsecase{
				userRepo:       tt.mockRepo(),
				contextTimeout: time.Second,
			}
			got, err := a.GetUser(tt.args.c, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("Expected error %v, got %v", tt.err, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userUsecase.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	type args struct {
		c  context.Context
		id int64
	}
	tests := []struct {
		name    string
		mockRepo func() *mockUserRepository
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "success",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					deleteFunc: func(ctx context.Context, id int64) error {
						return nil
					},
				}
			},
			args: args{
				c:  context.Background(),
				id: 1,
			},
			wantErr: false,
		},
		{
			name: "failure",
			mockRepo: func() *mockUserRepository {
				return &mockUserRepository{
					deleteFunc: func(ctx context.Context, id int64) error {
						return errors.New("delete error")
					},
				}
			},
			args: args{
				c:  context.Background(),
				id: 1,
			},
			wantErr: true,
			err:     errors.New("delete error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &userUsecase{
				userRepo:       tt.mockRepo(),
				contextTimeout: time.Second,
			}
			err := a.DeleteUser(tt.args.c, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("Expected error %v, got %v", tt.err, err)
			}
		})
	}
}
