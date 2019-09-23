package repository

import (
	"reflect"
	"testing"
)

func TestNewSubCommandRepository(t *testing.T) {
	tests := []struct {
		name string
		want *SubCommandRepository
	}{
		{
			want: &SubCommandRepository{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSubCommandRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSubCommandRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubCommandRepository_Run(t *testing.T) {
	type args struct {
		input SubCommandRepositoryInput
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			args: args{
				input: SubCommandRepositoryInput{},
			},
			want: map[string]string{
				"text": "{}",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SubCommandRepository{}
			got, err := s.Run(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() got = %v, want %v", got, tt.want)
			}
		})
	}
}
