package repository

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNewSignatureRepository(t *testing.T) {
	tests := []struct {
		name string
		want *SignatureRepository
	}{
		{
			want: NewSignatureRepository(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSignatureRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSignatureRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignatureRepository_Verify(t *testing.T) {
	type args struct {
		input *SignatureInput
	}
	timestamp := int(time.Now().Unix())
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{
			name: "more than 5 minutes ago",
			args: args{
				input: &SignatureInput{
					Timestamp:        timestamp - (5*60 + 1),
					Signature:        "",
					SigningSecret:    "",
					Body:             "",
					SignatureVersion: "",
				},
			},
			wantErr: true,
			err:     "input timestamp is newer or older than 5 minutes",
		},
		{
			name: "more than 5 minutes in the future",
			args: args{
				input: &SignatureInput{
					Timestamp:        timestamp + (5*60 + 1),
					Signature:        "",
					SigningSecret:    "",
					Body:             "",
					SignatureVersion: "",
				},
			},
			wantErr: true,
			err:     "input timestamp is newer or older than 5 minutes",
		},
		{
			name: "5 minutes ago and invalid signature",
			args: args{
				input: &SignatureInput{
					Timestamp:        timestamp - (5 * 60),
					Signature:        "",
					SigningSecret:    "",
					Body:             "",
					SignatureVersion: "",
				},
			},
			wantErr: true,
			err:     "invalid signature detected",
		},
		{
			name: "5 minutes in the future and invalid signature",
			args: args{
				input: &SignatureInput{
					Timestamp:        timestamp + (5 * 60),
					Signature:        "",
					SigningSecret:    "",
					Body:             "",
					SignatureVersion: "",
				},
			},
			wantErr: true,
			err:     "invalid signature detected",
		},
		{
			name: "valid signature",
			args: args{
				input: &SignatureInput{
					Timestamp:        int(timestamp),
					Signature:        "v0=" + sign("v0:"+strconv.Itoa(timestamp)+":a=b&c=d&e=f", "signature_secret_txt"),
					SigningSecret:    "signature_secret_txt",
					Body:             "a=b&c=d&e=f",
					SignatureVersion: "v0",
				},
			},
			wantErr: false,
			err:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := &SignatureRepository{}
			if err := si.Verify(tt.args.input); err != nil {
				if tt.wantErr {
					if tt.err != err.Error() {
						t.Errorf("Verify() error = %v, wantErr %v", err.Error(), tt.err)
					}
				} else {
					t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if tt.wantErr {
					t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
