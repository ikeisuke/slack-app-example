package presenter

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewResponsePresenter(t *testing.T) {
	tests := []struct {
		name string
		want *ResponsePresenter
	}{
		{
			want: &ResponsePresenter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResponsePresenter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponsePresenter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponsePresenter_Output(t *testing.T) {
	type args struct {
		data interface{}
		err  error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "output with error",
			args: args{
				data: "",
				err:  errors.New("error message"),
			},
			want: "{\"response_type\":\"ephemeral\",\"text\":\"Sorry, that didn't work. Please try again. (error message)\"}",
		},
		{
			name: "output without error",
			args: args{
				data: map[string]string{"test": "data"},
				err:  nil,
			},
			want: "{\"test\":\"data\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResponsePresenter{}
			if got := r.Output(tt.args.data, tt.args.err); got != tt.want {
				t.Errorf("Output() = %v, want %v", got, tt.want)
			}
		})
	}
}
