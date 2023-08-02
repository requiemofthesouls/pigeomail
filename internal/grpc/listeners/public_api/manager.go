package public_api

import (
	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
)

func New() pigeomail_api_pb.PublicAPIServer {
	return &manager{}
}

type (
	manager struct{}
)
