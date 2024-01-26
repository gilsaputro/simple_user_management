// This file contains the repository implementation layer.
package repository

import (
	"testing"

	_ "github.com/lib/pq"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()
	type args struct {
		opts NewRepositoryOptions
	}
	tests := []struct {
		name string
		args args
		want RepositoryInterface
	}{
		{
			name: "TestNewRepository",
			args: args{
				opts: NewRepositoryOptions{
					Dsn: "fake-dsn",
				},
			},
			want: &Repository{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(tt.args.opts)
			r.Db = db

			err := r.Db.Ping()
			if err != nil {
				t.Errorf("NewRepository() = %v", err)
			}
		})
	}
}
