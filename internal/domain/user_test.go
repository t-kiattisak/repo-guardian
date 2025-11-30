package domain

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type mockUserRepository struct {
	createFunc  func(ctx context.Context, user *User) error
	getByIDFunc func(ctx context.Context, id int64) (*User, error)
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockUserRepository) Create(ctx context.Context, user *User) error {
	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockUserRepository) Delete(ctx context.Context, id int64) error {
	return m.deleteFunc(ctx, id)
}

type userUsecase struct {
	userRepo UserRepository
}

func NewUserUsecase(userRepo UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) Register(ctx context.Context, user *User) error {
	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) GetUser(ctx context.Context, id int64) (*User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id int64) error {
	return u.userRepo.Delete(ctx, id)
}

func TestUserUsecase_Register(t *testing.T) {
	type fields struct {
		userRepo UserRepository
	}
	type args struct {
		ctx  context.Context
		user *User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		repoErr error
	}{
		{
			name: "success",
			fields: fields{
				userRepo: &mockUserRepository{
					createFunc: func(ctx context.Context, user *User) error {
						return nil
					},
				},
			},
			args: args{
				ctx:  context.Background(),
				user: &User{Username: "test", Email: "test@example.com"},
			},
			wantErr: false,
		},
		{
			name: "failure - repo error",
			fields: fields{
				userRepo: &mockUserRepository{
					createFunc: func(ctx context.Context, user *User) error {
						return errors.New("repo error")
					},
				},
			},
			args: args{
				ctx:  context.Background(),
				user: &User{Username: "test", Email: "test@example.com"},
			},
			wantErr: true,
			repoErr: errors.New("repo error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userUsecase{
				userRepo: tt.fields.userRepo,
			}
			err := u.Register(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserUsecase.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.repoErr != nil {
				if err.Error() != tt.repoErr.Error() {
					t.Errorf("UserUsecase.Register() error = %v, want %v", err, tt.repoErr)
				}
			}
		})
	}
}

func TestUserUsecase_GetUser(t *testing.T) {
	type fields struct {
		userRepo UserRepository
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
		repoErr error
	}{
		{
			name: "success",
			fields: fields{
				userRepo: &mockUserRepository{
					getByIDFunc: func(ctx context.Context, id int64) (*User, error) {
						return &User{ID: id, Username: "test", Email: "test@example.com"}, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    &User{ID: 1, Username: "test", Email: "test@example.com"},
			wantErr: false,
		},
		{
			name: "failure - repo error",
			fields: fields{
				userRepo: &mockUserRepository{
					getByIDFunc: func(ctx context.Context, id int64) (*User, error) {
						return nil, errors.New("repo error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: true,
			repoErr: errors.New("repo error"),
		},
		{
			name: "failure - user not found",
			fields: fields{
				userRepo: &mockUserRepository{
					getByIDFunc: func(ctx context.Context, id int64) (*User, error) {
						return nil, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userUsecase{
				userRepo: tt.fields.userRepo,
			}
			got, err := u.GetUser(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserUsecase.GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserUsecase.GetUser() = %v, want %v", got, tt.want)
			}
			if tt.wantErr && err != nil && tt.repoErr != nil {
				if err.Error() != tt.repoErr.Error() {
					t.Errorf("UserUsecase.GetUser() error = %v, want %v", err, tt.repoErr)
				}
			}
		})
	}
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	type fields struct {
		userRepo UserRepository
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		repoErr error
	}{
		{
			name: "success",
			fields: fields{
				userRepo: &mockUserRepository{
					deleteFunc: func(ctx context.Context, id int64) error {
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: false,
		},
		{
			name: "failure - repo error",
			fields: fields{
				userRepo: &mockUserRepository{
					deleteFunc: func(ctx context.Context, id int64) error {
						return errors.New("repo error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: true,
			repoErr: errors.New("repo error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userUsecase{
				userRepo: tt.fields.userRepo,
			}
			err := u.DeleteUser(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserUsecase.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.repoErr != nil {
				if err.Error() != tt.repoErr.Error() {
					t.Errorf("UserUsecase.DeleteUser() error = %v, want %v", err, tt.repoErr)
				}
			}
		})
	}
}
