package repository

import (
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/infrastructure"
	"reflect"
	"testing"
)

type SlackInfrastructureMock struct {
	infrastructure.ISlack
}

func TestNewSubCommandRepository(t *testing.T) {
	type args struct {
		infra infrastructure.ISlack
	}
	tests := []struct {
		name string
		args args
		want *CommandRepository
	}{
		{
			args: args{
				infra: SlackInfrastructureMock{},
			},
			want: &CommandRepository{
				infrastructure: SlackInfrastructureMock{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommandRepository(tt.args.infra); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommandRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandRepository_Run(t *testing.T) {
	type fields struct {
		infrastructure infrastructure.ISlack
	}
	type args struct {
		input CommandRepositoryInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.SlackMessage
		wantErr bool
	}{
		{
			name: "unknown command",
			fields: fields{
				infrastructure: &SlackInfrastructureMock{},
			},
			args: args{
				input: CommandRepositoryInput{
					ChannelID:   "",
					ChannelName: "",
					Command:     "",
					ResponseUrl: "",
					TeamDomain:  "",
					TeamID:      "",
					Text:        "unknown_command unknown_subcommand",
					Token:       "",
					TriggerID:   "",
					UserID:      "",
					UserName:    "",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CommandRepository{
				infrastructure: tt.fields.infrastructure,
			}
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

func Test_parseArgs(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "1 arg",
			args: args{
				input: "arg",
			},
			want: []string{
				"arg",
			},
		},
		{
			name: "2 args",
			args: args{
				input: "args1 args2",
			},
			want: []string{
				"args1",
				"args2",
			},
		},
		{
			name: "3 args",
			args: args{
				input: "args1 args2 args3",
			},
			want: []string{
				"args1",
				"args2",
				"args3",
			},
		},
		{
			name: "contiguous separators",
			args: args{
				input: "args1  args2   args3",
			},
			want: []string{
				"args1",
				"args2",
				"args3",
			},
		},
		{
			name: "quoted value",
			args: args{
				input: "args1  \"args2 args3\"",
			},
			want: []string{
				"args1",
				"args2 args3",
			},
		},
		{
			name: "contiguous quoted value",
			args: args{
				input: "args1  \"args2 args3\"\"args4\"args5",
			},
			want: []string{
				"args1",
				"args2 args3",
				"args4",
				"args5",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseArgs(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
