package application

import (
	"encoding/json"
	"errors"
	"github.com/ikeisuke/slack-app-example/internal/presenter"
	"github.com/ikeisuke/slack-app-example/internal/repository"
	"reflect"
	"testing"
)

type SignatureRepositoryMock struct {
	valid bool
}

func NewSignatureRepositoryMock(valid bool) *SignatureRepositoryMock {
	return &SignatureRepositoryMock{valid: valid}
}
func (s *SignatureRepositoryMock) Verify(input *repository.SignatureInput) error {
	if s.valid {
		return nil
	}
	return errors.New("invalid")
}

type SubCommandRepositoryMock struct {
	err error
}

func NewSubCommandRepository(err error) *SubCommandRepositoryMock {
	return &SubCommandRepositoryMock{
		err: err,
	}
}
func (s *SubCommandRepositoryMock) Run(input repository.SubCommandRepositoryInput) (interface{}, error) {
	return input, s.err
}

type PresenterMock struct {
}

func (p *PresenterMock) Output(data interface{}, err error) string {
	message := ""
	if err != nil {
		message = err.Error()
	}
	t, _ := json.Marshal(map[string]interface{}{
		"data": data,
		"err":  message,
	})
	return string(t)
}

func TestNewSlashCommandInteraction(t *testing.T) {
	type args struct {
		r repository.ISignatureRepository
		s repository.ISubCommandRepository
		p presenter.IPresenter
	}
	tests := []struct {
		name string
		args args
		want *SlackCommandInteraction
	}{
		{
			args: args{
				r: NewSignatureRepositoryMock(true),
				s: NewSubCommandRepository(nil),
				p: &PresenterMock{},
			},
			want: NewSlashCommandInteraction(
				NewSignatureRepositoryMock(true),
				NewSubCommandRepository(nil),
				&PresenterMock{},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlashCommandInteraction(tt.args.r, tt.args.s, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSlashCommandInteraction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackCommandInteraction_Run(t *testing.T) {
	type fields struct {
		signature  repository.ISignatureRepository
		subCommand repository.ISubCommandRepository
		presenter  presenter.IPresenter
	}
	type args struct {
		input *SlashCommandInput
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "invalid signature",
			fields: fields{
				signature:  NewSignatureRepositoryMock(false),
				subCommand: NewSubCommandRepository(nil),
				presenter:  &PresenterMock{},
			},
			args: args{
				input: &SlashCommandInput{},
			},
			want: "{\"data\":null,\"err\":\"invalid\"}",
		},
		{
			name: "valid signature, invalid request body",
			fields: fields{
				signature:  NewSignatureRepositoryMock(true),
				subCommand: NewSubCommandRepository(nil),
				presenter:  &PresenterMock{},
			},
			args: args{
				input: &SlashCommandInput{
					Body: "a=b&c=d=f",
				},
			},
			want: "{\"data\":null,\"err\":\"invalid request body detected\"}",
		},
		{
			name: "valid signature, empty request body",
			fields: fields{
				signature:  NewSignatureRepositoryMock(true),
				subCommand: NewSubCommandRepository(nil),
				presenter:  &PresenterMock{},
			},
			args: args{
				input: &SlashCommandInput{},
			},
			want: "{\"data\":{},\"err\":\"\"}",
		},
		{
			name: "valid signature, full request body",
			fields: fields{
				signature:  NewSignatureRepositoryMock(true),
				subCommand: NewSubCommandRepository(nil),
				presenter:  &PresenterMock{},
			},
			args: args{
				input: &SlashCommandInput{
					Body: "channel_id=a&channel_name=b&command=c&response_url=d&team_domain=e&team_id=f&text=h&token=i&trigger_id=j&user_id=k&user_name=l",
				},
			},
			want: "{\"data\":{\"channel_id\":\"a\",\"channel_name\":\"b\",\"command\":\"c\",\"response_url\":\"d\",\"team_domain\":\"e\",\"team_id\":\"f\",\"text\":\"h\",\"token\":\"i\",\"trigger_id\":\"j\",\"user_id\":\"k\",\"user_name\":\"l\"},\"err\":\"\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SlackCommandInteraction{
				signature:  tt.fields.signature,
				subCommand: tt.fields.subCommand,
				presenter:  tt.fields.presenter,
			}
			if got := s.Run(tt.args.input); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
