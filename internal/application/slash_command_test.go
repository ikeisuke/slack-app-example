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

func (s *SignatureRepositoryMock) Verify(input repository.SignatureInput) error {
	if s.valid {
		return nil
	}
	return errors.New("invalid")
}

type CommandRepositoryMock struct {
	err error
}

func (s *CommandRepositoryMock) Run(input repository.CommandRepositoryInput) (*entity.SlackMessage, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &entity.SlackMessage{
		Text: input.UserName,
	}, nil
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
				r: &SignatureRepositoryMock{valid: true},
				s: &CommandRepositoryMock{},
				p: &PresenterMock{},
			},
			want: NewSlashCommandInteraction(
				&SignatureRepositoryMock{valid: true},
				&CommandRepositoryMock{},
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
				signature: &SignatureRepositoryMock{valid: false},
				command:   &CommandRepositoryMock{},
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
				signature: &SignatureRepositoryMock{valid: true},
				command:   &CommandRepositoryMock{},
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
				signature: &SignatureRepositoryMock{valid: true},
				command:   &CommandRepositoryMock{},
				presenter: &PresenterMock{},
			},
			args: args{
				input: &SlashCommandInput{
					Body: "channel_id=a&channel_name=b&command=c&response_url=d&team_domain=e&team_id=f&text=h&token=i&trigger_id=j&user_id=k&user_name=l",
				},
			},
			want: "{\"text\":\"l\"}",
		},
		{
			name: "valid signature, repository error",
			fields: fields{
				signature: &SignatureRepositoryMock{valid: true},
				command:   &CommandRepositoryMock{err: errors.New("err")},
				presenter: &PresenterMock{},
			},
			args: args{
				input: &SlashCommandInput{
					Body: "channel_id=a&channel_name=b&command=c&response_url=d&team_domain=e&team_id=f&text=h&token=i&trigger_id=j&user_id=k&user_name=l",
				},
			},
			want: "{\"response_type\":\"ephemeral\",\"text\":\"Sorry, that didn't work. Please try again. (err)\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SlackCommandInteraction{
				signature: tt.fields.signature,
				command:   tt.fields.command,
				presenter: tt.fields.presenter,
			}
			if got, _ := s.Run(tt.args.input); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
