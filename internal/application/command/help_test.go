package command

import (
	"errors"
	"fmt"
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/repository/command"
	"reflect"
	"testing"
)

type HelpRepositoryMock struct {
	message string
	err     error
}

func (h *HelpRepositoryMock) HelpMessage(input command.HelpRepositoryInput) (string, error) {
	return fmt.Sprintf(h.message, input.SlashCommand), h.err
}

func TestNewHelpInteraction(t *testing.T) {
	type args struct {
		repository command.IHelpRepository
	}
	tests := []struct {
		name string
		args args
		want *HelpInteraction
	}{
		{
			args: args{
				repository: &HelpRepositoryMock{},
			},
			want: &HelpInteraction{repository: &HelpRepositoryMock{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHelpInteraction(tt.args.repository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHelpInteraction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHelpInteraction_Run(t *testing.T) {
	type fields struct {
		repository command.IHelpRepository
	}
	type args struct {
		input HelpInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.SlackMessage
		wantErr bool
	}{
		{
			fields: fields{
				repository: &HelpRepositoryMock{
					message: "help message: %s",
					err:     nil,
				},
			},
			args: args{
				input: HelpInput{SlashCommand: "command_name"},
			},
			want: &entity.SlackMessage{
				ResponseType: "ephemeral",
				Text:         "",
				Attachments:  nil,
				Blocks: []entity.SlackMessageBlock{
					{
						Type: "section",
						Text: &entity.SlackMessageTextObject{
							Type: "mrkdwn",
							Text: "help message: command_name",
						},
					},
				},
				ThreadTS:        "",
				ReplaceOriginal: false,
			},
			wantErr: false,
		},
		{
			fields: fields{
				repository: &HelpRepositoryMock{
					message: "help message: %s",
					err:     errors.New("error"),
				},
			},
			args: args{
				input: HelpInput{SlashCommand: "command_name"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HelpInteraction{
				repository: tt.fields.repository,
			}
			got, err := h.Run(tt.args.input)
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
