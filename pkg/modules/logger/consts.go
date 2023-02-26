package logger

const (
	KeyRequestID  = "request_id"
	KeyUserClient = "user_client"

	KeyServiceNamespace = "service.namespace"
	KeyServiceDomain    = "service.domain"
	KeyServiceName      = "service.name"

	KeyGRPCService          = "grpc.service"
	KeyGRPCMethod           = "grpc.method"
	KeyGRPCRequestBody      = "grpc.request.body"
	KeyGRPCRequestStartTime = "grpc.request.start_time"
	KeyGRPCRequestDeadline  = "grpc.request.deadline"
	KeyGRPCRequestResponse  = "grpc.request.response"
)
