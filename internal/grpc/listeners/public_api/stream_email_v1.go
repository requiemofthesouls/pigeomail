package public_api

import (
	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
)

func (m *manager) StreamEMailV1(
	req *pigeomail_api_pb.PublicAPIStreamEMailV1Request,
	// stream pigeomail_api_pb.PublicAPI_StreamEMailV1Server,
) error {
	//TODO implement me

	//for {
	//	if err := stream.Send(&pigeomail_api_pb.EMail{
	//		Id:        "1",
	//		Email:     "test",
	//		Subject:   "test",
	//		Body:      "test",
	//		Sender:    "aaa@test.com",
	//		Recipient: "keepo@pigeomail.ddns.net",
	//		CreatedAt: time.Now().UnixNano(),
	//	}); err != nil {
	//		return err
	//	}
	//
	//	time.Sleep(1 * time.Second)
	//}

	return nil
}
