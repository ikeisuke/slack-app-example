package application

import (
	"encoding/json"
	"errors"
	"github.com/ikeisuke/slack-app-example/internal/entity"
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

type CommandRepositoryMock struct {
	err error
}

func NewCommandRepository(err error) *CommandRepositoryMock {
	return &CommandRepositoryMock{
		err: err,
	}
}
func (s *CommandRepositoryMock) Run(input repository.CommandRepositoryInput) (*entity.SlackMessage, error) {
	return &entity.SlackMessage{
		Text: input.UserName,
	}, s.err
}

type PresenterMock struct {
}

func (p *PresenterMock) Output(data interface{}) (string, error) {
	buf, err := json.Marshal(data)
	return string(buf), err
}

func TestNewSlashCommandInteraction(t *testing.T) {
	type args struct {
		r repository.ISignatureRepository
		s repository.ICommandRepository
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
				s: NewCommandRepository(nil),
				p: &PresenterMock{},
			},
			want: NewSlashCommandInteraction(
				NewSignatureRepositoryMock(true),
				NewCommandRepository(nil),
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
		signature repository.ISignatureRepository
		command   repository.ICommandRepository
		presenter presenter.IPresenter
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
				signature: NewSignatureRepositoryMock(false),
				command:   NewCommandRepository(nil),
				presenter: &PresenterMock{},
			},
			args: args{
				input: &SlashCommandInput{},
			},
			want: "{\"response_type\":\"ephemeral\",\"text\":\"Sorry, that didn't work. Please try again. (invalid)\"}",
		},
		{
			name: "valid signature, empty request body",
			fields: fields{
				signature: NewSignatureRepositoryMock(true),
				command:   NewCommandRepository(nil),
				presenter: &PresenterMock{},
			},
			args: args{
				input: &SlashCommandInput{},
			},
			want: "{}",
		},
		{
			name: "valid signature, full request body",
			fields: fields{
				signature: NewSignatureRepositoryMock(true),
				command:   NewCommandRepository(nil),
				presenter: &PresenterMock{},
			},
			args: args{
				input: &SlashCommandInput{
					Body: "channel_id=a&channel_name=b&command=c&response_url=d&team_domain=e&team_id=f&text=h&token=i&trigger_id=j&user_id=k&user_name=l",
				},
			},
			want: "{\"text\":\"l\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SlackCommandInteraction{
				signature:  tt.fields.signature,
				subCommand: tt.fields.command,
				presenter:  tt.fields.presenter,
			}
			if got, _ := s.Run(tt.args.input); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
