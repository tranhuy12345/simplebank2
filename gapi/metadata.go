package gapi

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedFor              = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	fmt.Println(md)
	if userAgent := md.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
		mtdt.UserAgent = userAgent[0]
	}
	if clientIp := md.Get(xForwardedFor); len(clientIp) > 0 {
		mtdt.ClientIp = clientIp[0]
	}
	return mtdt
}
