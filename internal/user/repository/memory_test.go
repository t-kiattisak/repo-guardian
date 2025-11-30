package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"repo-guardian/internal/domain"
)

func TestMemoryUserRepository_Create(t *testing.T) {
	type fields struct {
		users map[int64]*domain.User
	}
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(f *fields, a *args)
	}{
		{
			name: "Successfully create a user",
			fields: fields{
				users: make(map[int64]*domain.User),
			},
			args: args{
				ctx: context.Background(),
				user: &domain.User{
					ID:    1,
					Name:  "John Doe",
					Email: "john.doe@example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail to create a user - user already exists",
			fields: fields{
				users: map[int64]*domain.User{
					1: {
						ID:    1,
						Name:  "John Doe",
						Email: "john.doe@example.com",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				user: &domain.User{
					ID:    1,
					Name:  "John Doe",
					Email: "john.doe@example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "Successfully create a user after deleting the same id",
			fields: fields{
				users: map[int64]*domain.User{},
			},
			args: args{
				ctx: context.Background(),
				user: &domain.User{
					ID:    1,
					Name:  "John Doe",
					Email: "john.doe@example.com",
				},
			},
			wantErr: false,
			setup: func(f *fields, a *args) {
				f.users[1] = &domain.User{
					ID: 1,
				}
				delete(f.users, 1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &memoryUserRepository{
				users: tt.fields.users,
			}
			if tt.setup != nil {
				tt.setup(&tt.fields, &tt.args)
			}
			if err := r.Create(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryUserRepository_GetByID(t *testing.T) {
	type fields struct {
		users map[int64]*domain.User
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "Successfully get user by ID",
			fields: fields{
				users: map[int64]*domain.User{
					1: {
						ID:    1,
						Name:  "John Doe",
						Email: "john.doe@example.com",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &domain.User{
				ID:    1,
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name: "Fail to get user by ID - user not found",
			fields: fields{
				users: make(map[int64]*domain.User),
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Fail to get user by ID - different user",
			fields: fields{
				users: map[int64]*domain.User{
					2: {
						ID:    2,
						Name:  "John Doe",
						Email: "john.doe@example.com",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &memoryUserRepository{
				users: tt.fields.users,
			}
			got, err := r.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryUserRepository_Delete(t *testing.T) {
	type fields struct {
		users map[int64]*domain.User
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
	}{
		{
			name: "Successfully delete user by ID",
			fields: fields{
				users: map[int64]*domain.User{
					1: {
						ID:    1,
						Name:  "John Doe",
						Email: "john.doe@example.com",
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
			name: "Fail to delete user by ID - user not found",
			fields: fields{
				users: make(map[int64]*domain.User),
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: true,
		},
		{
			name: "Fail to delete user by ID - different user",
			fields: fields{
				users: map[int64]*domain.User{
					2: {
						ID:    2,
						Name:  "John Doe",
						Email: "john.doe@example.com",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &memoryUserRepository{
				users: tt.fields.users,
			}
			err := r.Delete(tt.args.ctx, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if _, exists := r.users[tt.args.id]; exists {
					t.Errorf("Delete() user not deleted")
				}
			}

			if tt.wantErr && !errors.Is(err, errors.New("user not found")) {
				t.Errorf("Delete() error = %v, wantErr %v", err, errors.New("user not found"))
			}
		})
	}
}
