// models/common.tsp
import "@typespec/http";

using TypeSpec.Http;

namespace CmsApi.Models.Common;

// Common response models
model SuccessResponse {
  success: true;
  message: string;
  data?: unknown;
}

model ErrorResponse {
  success: false;
  message: string;
  error?: string;
  details?: unknown;
}

// Authentication
model BearerAuth {
  @header("Authorization")
  authorization: string;
}

// Pagination
model PaginateRequest {
  @query page?: int32;
  @query limit?: int32;
}

// Status enums
enum RequestStatus {
  pending: "pending",
  approved: "approved",
  rejected: "rejected",
  in_review: "in_review",
}