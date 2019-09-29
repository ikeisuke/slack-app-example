package command

import (
	"errors"
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/repository/command"
	"reflect"
	"testing"
)

type ChannelRepositoryMock struct {
	list entity.SlackAPIChannelListResponse
	err  error
}

func (c *ChannelRepositoryMock) List() (entity.SlackAPIChannelListResponse, error) {
	return c.list, c.err
}

func TestNewChannelInteraction(t *testing.T) {
	type args struct {
		repo command.IChannelRepository
	}
	tests := []struct {
		name string
		args args
		want *ChannelInteraction
	}{
		{
			args: args{
				repo: &ChannelRepositoryMock{},
			},
			want: &ChannelInteraction{
				repository: &ChannelRepositoryMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChannelInteraction(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelInteraction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannelInteraction_List(t *testing.T) {
	type fields struct {
		repository command.IChannelRepository
	}
	type args struct {
		input ChannelListInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.SlackMessage
		wantErr bool
	}{
		{
			name: "empty channels response",
			fields: fields{
				repository: &ChannelRepositoryMock{
					list: make(entity.SlackAPIChannelListResponse, 0),
					err:  nil,
				},
			},
			args: args{
				input: ChannelListInput{},
			},
			want: &entity.SlackMessage{
				ResponseType:    "ephemeral",
				Text:            "No channels found.",
				Attachments:     nil,
				Blocks:          nil,
				ThreadTS:        "",
				ReplaceOriginal: false,
			},
			wantErr: false,
		},
		{
			name: "1 channel response",
			fields: fields{
				repository: &ChannelRepositoryMock{
					list: entity.SlackAPIChannelListResponse{
						{
							ID:   "1234",
							Name: "channel_name_1",
						},
					},
					err: nil,
				},
			},
			args: args{
				input: ChannelListInput{},
			},
			want: &entity.SlackMessage{
				ResponseType: "in_channel",
				Text:         "",
				Attachments:  nil,
				Blocks: []entity.SlackMessageBlock{
					{
						Type: "actions",
						Text: nil,
						Elements: []entity.SlackMessageBlockActionElement{
							{
								Type:     "button",
								ActionID: "channel_detail_1234",
								Text: &entity.SlackMessageTextObject{
									Type: "mrkdwn",
									Text: "channel_name_1",
								},
							},
						},
					},
				},
				ThreadTS:        "",
				ReplaceOriginal: false,
			},
			wantErr: false,
		},
		{
			name: "multiple channel response",
			fields: fields{
				repository: &ChannelRepositoryMock{
					list: entity.SlackAPIChannelListResponse{
						{
							ID:   "1234",
							Name: "channel_name_1",
						},
						{
							ID:   "r564",
							Name: "channel_name_2",
						},
					},
					err: nil,
				},
			},
			args: args{
				input: ChannelListInput{},
			},
			want: &entity.SlackMessage{
				ResponseType: "in_channel",
				Text:         "",
				Attachments:  nil,
				Blocks: []entity.SlackMessageBlock{
					{
						Type: "actions",
						Text: nil,
						Elements: []entity.SlackMessageBlockActionElement{
							{
								Type:     "button",
								ActionID: "channel_detail_1234",
								Text: &entity.SlackMessageTextObject{
									Type: "mrkdwn",
									Text: "channel_name_1",
								},
							},
							{
								Type:     "button",
								ActionID: "channel_detail_r564",
								Text: &entity.SlackMessageTextObject{
									Type: "mrkdwn",
									Text: "channel_name_2",
								},
							},
						},
					},
				},
				ThreadTS:        "",
				ReplaceOriginal: false,
			},
			wantErr: false,
		},
		{
			name: "error channel response",
			fields: fields{
				repository: &ChannelRepositoryMock{
					list: entity.SlackAPIChannelListResponse{},
					err:  errors.New("error response"),
				},
			},
			args: args{
				input: ChannelListInput{},
			},
			want: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChannelInteraction{
				repository: tt.fields.repository,
			}
			got, err := c.List(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}
