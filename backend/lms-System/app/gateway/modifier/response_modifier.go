package modifier

import (
	"context"
	"google.golang.org/grpc/metadata"
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
	}

	for _, headerName := range cookieHeaders {
		if values := md.HeaderMD.Get(headerName); len(values) > 0 {
			w.Header().Add("Set-Cookie", values[0])
			delete(md.HeaderMD, headerName)
		}
	}
	return nil

}

func RequestModifier(ctx context.Context, req *http.Request) metadata.MD {
	md := metadata.MD{}
	
	for _, cookie := range req.Cookies() {
		switch cookie.Name {
		case "access-token":
			md.Set("access-token-cookie", cookie.Value)
		case "refresh-token":
			md.Set("refresh-token-cookie", cookie.Value)
		case "user-info":
			md.Set("user-info-cookie", cookie.Value)
		case "csrf-token":
			md.Set("csrf-token-cookie", cookie.Value)
		case "session-id":
			md.Set("session-id-cookie", cookie.Value)
		}
	}

	return md
}
