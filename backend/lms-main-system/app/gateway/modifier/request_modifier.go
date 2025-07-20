package modifier

import (
	"context"
	"google.golang.org/grpc/metadata"
	"net/http"
)

func RequestModifier(ctx context.Context, req *http.Request) metadata.MD {
	md := metadata.MD{}

	if auth := req.Header.Get("Authorization"); auth != "" {
		md.Set("authorization", auth)
	}
	if org := req.Header.Get("X-Organisation"); org != "" {
		md.Set("x-organisation", org)
	}

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
		case "organisation":
			md.Set("organisation-cookie", cookie.Value)
		}
	}

	return md
}
