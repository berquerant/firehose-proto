package grpcx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func metadataFromContext(ctx context.Context) metadata.MD {
	if md, exist := metadata.FromIncomingContext(ctx); exist {
		return md
	}
	return map[string][]string{}
}

func extractFromMetadata(md metadata.MD, tag string) (string, bool) {
	if len(md.Get(tag)) > 0 {
		return md.Get(tag)[0], true
	}
	return "", false
}

const (
	ClientRequestIDTag = "grpcx.client_request_id"
	ServerRequestIDTag = "grpcx.server_request_id"
)

func ExtractClientRequestIDFromMetadata(md metadata.MD) (string, bool) {
	return extractFromMetadata(md, ClientRequestIDTag)
}

func ExtractServerRequestIDFromMetadata(md metadata.MD) (string, bool) {
	return extractFromMetadata(md, ServerRequestIDTag)
}

func AppendClientRequestIDIntoMetadata(md metadata.MD, id string) {
	md.Append(ClientRequestIDTag, id)
}

func AppendServerRequestIDIntoMetadata(md metadata.MD, id string) {
	md.Append(ServerRequestIDTag, id)
}
