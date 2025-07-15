import "@typespec/http";
import "@typespec/openapi3";
import "./common.tsp";

using Http;
using Common;

namespace Errors {
  @error
  @doc("Bad request error response")
  model BadRequest {
    @statusCode code: 400;
    ...ApiResponse<never>;
  }

  @error
  @doc("Unauthorized error response")
  model Unauthorized {
    @statusCode code: 401;
    ...ApiResponse<never>;
  }

  @error
  @doc("Forbidden error response")
  model Forbidden {
    @statusCode code: 403;
    ...ApiResponse<never>;
  }

  @error
  @doc("Not found error response")
  model NotFound {
    @statusCode code: 404;
    ...ApiResponse<never>;
  }

  @error
  @doc("Conflict error response")
  model Conflict {
    @statusCode code: 409;
    ...ApiResponse<never>;
  }

  @error
  @doc("Internal server error response")
  model InternalServerError {
    @statusCode code: 500;
    ...ApiResponse<never>;
  }

  @error
  @doc("Resource created successfully")
  model Created {
    @statusCode code: 201;
    ...ApiResponse<unknown>;
  }
}