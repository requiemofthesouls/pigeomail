package private_api

import (
	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
)

func New() pigeomail_api_pb.PrivateAPIServer {
	return &manager{}
}

type (
	manager struct{}
)
