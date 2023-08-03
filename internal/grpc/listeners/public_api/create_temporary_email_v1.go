package public_api

import (
	"context"

	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
)

func (m *manager) CreateTemporaryEMailV1(
	ctx context.Context,
	req *pigeomail_api_pb.PublicAPICreateTemporaryEMailV1Request,
) (*pigeomail_api_pb.PublicAPICreateTemporaryEMailV1Response, error) {
	//TODO implement me
	return &pigeomail_api_pb.PublicAPICreateTemporaryEMailV1Response{
		Email: "keepo@pigeomail.ddns.net",
	}, nil
}
