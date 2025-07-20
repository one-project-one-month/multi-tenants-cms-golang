package modifier

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

func ResponseModifier(
	ctx context.Context,
	w http.ResponseWriter,
	resp proto.Message,
) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if cookies := md.HeaderMD.Get("Set-Cookie"); len(cookies) > 0 {
		for _, cookie := range cookies {
			w.Header().Add("Set-Cookie", cookie)
		}
		delete(md.HeaderMD, "Set-Cookie")
	}

	cookieHeaders := []string{
		"access-token-cookie",
		"refresh-token-cookie",
		"user-info-cookie",
		"csrf-token-cookie",
		"session-id-cookie",
		"mfa-code",
	}

	for _, header := range cookieHeaders {
		if values := md.HeaderMD.Get(header); len(values) > 0 {
			for _, value := range values {
				w.Header().Add("Set-Cookie", value)
			}
			delete(md.HeaderMD, header)
		}
	}

	return nil
}
