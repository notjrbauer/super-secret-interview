package api

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// AuthzInterceptor intercepts calls in order to auth a user.
func AuthzInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := authz(ctx, info.FullMethod); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	return handler(ctx, req)
}

// StreamAuthInterceptor intercepts calls in order to auth a user.

func StreamAuthInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := authz(stream.Context(), info.FullMethod); err != nil {
		return status.Error(codes.PermissionDenied, err.Error())
	}
	return handler(srv, stream)
}

// authz verifies a user can perform specified operations.
// The verification looks at the SANS / DNS.
func authz(ctx context.Context, method string) error {
	// reads the peer information from context
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return errors.New("error to read peer information")
	}
	// reads user tls information.
	tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return errors.New("error to get auth information")
	}

	certs := tlsInfo.State.VerifiedChains
	if len(certs) == 0 || len(certs[0]) == 0 {
		return errors.New("missing certificate chain")
	}
	for _, role := range certs[0][0].DNSNames {
		if Can(method, role) {
			return nil
		}
	}
	return errors.New("unauthorized")
}
