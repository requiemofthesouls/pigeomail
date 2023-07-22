package logger

const (
	KeyRequestID  = "request_id"
	KeyUserClient = "user_client"

	KeyServiceName = "service.name"

	KeyGRPCService          = "grpc.service"
	KeyGRPCMethod           = "grpc.method"
	KeyGRPCRequestBody      = "grpc.request.body"
	KeyGRPCRequestStartTime = "grpc.request.start_time"
	KeyGRPCRequestDeadline  = "grpc.request.deadline"
	KeyGRPCRequestResponse  = "grpc.request.response"

	KeyKafkaClientID         = "kafka.client_id"
	KeyKafkaService          = "kafka.service"
	KeyKafkaConsumerGroupID  = "kafka.consumer_group_id"
	KeyKafkaMsgTopic         = "kafka.msg.topic"
	KeyKafkaMsgValue         = "kafka.msg.value"
	KeyKafkaHandlerStartTime = "kafka.handler.start_time"

	KeyWatcherName   = "watcher.name"
	KeyWatcherParams = "watcher.params"

	KeyRMQConnectionName   = "rmq.connection_name"
	KeyRMQServerName       = "rmq.server_name"
	KeyRMQExchange         = "rmq.exchange"
	KeyRMQRoutingKey       = "rmq.routing_key"
	KeyRMQMsgBody          = "rmq.msg.body"
	KeyRMQHandlerStartTime = "rmq.handler.start_time"
	KeyRMQHandlerPanicMsg  = "rmq.handler.panic_msg"
)
