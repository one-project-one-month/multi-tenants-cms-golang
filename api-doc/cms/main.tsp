import "@typespec/http";
import "@typespec/openapi3";
import "./service/auth-service.tsp";
import "./service/owner-service.tsp";
import "./lib/common.tsp";

using TypeSpec.Http;
using Services.Auth;
using Services.Owner;

@service(#{title: "CMS Platform API"})
@server("https://api.cms.example.com", "Production server")
@server("https://staging-api.cms.example.com", "Staging server")
namespace CMSAuthAPI;

@doc("This API provides authentication, user management, and owner management functionalities for the CMS platform.")
model BearerAuth {
  @doc("Bearer token for authentication")
  type: "http";
  scheme: "bearer";
  bearerFormat: "JWT";
}

interface API extends AuthOperations, OwnerOperations {}