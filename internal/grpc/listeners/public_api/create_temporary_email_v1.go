package public_api

import (
	"context"
	"fmt"
	"math/rand"

	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
)

func (m *manager) CreateTemporaryEMailV1(
	ctx context.Context,
	_ *pigeomail_api_pb.PublicAPICreateTemporaryEMailV1Request,
) (*pigeomail_api_pb.PublicAPICreateTemporaryEMailV1Response, error) {
	l := getLogger(ctx)

	l.Info("CreateTemporaryEMailV1")
	email := generateRandomEmail(m.smtpDomain)
	l.Info("Generated email: " + email)

	m.clients.sse.CreateStream(email)

	l.Info("Created stream for: " + email)

	return &pigeomail_api_pb.PublicAPICreateTemporaryEMailV1Response{
		Email: email,
	}, nil
}

// generateRandomEmail generates random email with domain from config
func generateRandomEmail(smtpDomain string) string {
	return fmt.Sprintf(
		"%s@%s",
		generateRandomString(32),
		smtpDomain,
	)
}

// generateRandomString generates random string with length maxNum
func generateRandomString(maxNum int) string {
	var (
		chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		b     = make([]rune, maxNum)
	)

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
