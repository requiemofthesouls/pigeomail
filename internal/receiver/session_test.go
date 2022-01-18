package receiver

import (
	"io"
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	"pigeomail/internal/repository"
	"pigeomail/rabbitmq"
)

func TestSession_parseMail(t *testing.T) {
	type fields struct {
		publisher rabbitmq.IRMQEmailPublisher
		repo      repository.IEmailRepository
		logger    logr.Logger
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantM   *rabbitmq.ParsedEmail
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				publisher: tt.fields.publisher,
				repo:      tt.fields.repo,
				logger:    tt.fields.logger,
			}
			gotM, err := s.parseMail(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("parseMail() gotM = %v, want %v", gotM, tt.wantM)
			}
		})
	}
}
