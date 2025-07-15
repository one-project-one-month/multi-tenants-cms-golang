// models/extended.tsp
import "@typespec/http";
import "./common.tsp";

using CmsApi.Models.Common;

namespace CmsApi.Models;

// Owner related models
model OwnerResponse {
  id: string; // UUID
  name: string;
  email: string;
  role: string;
  name_space: string;
  verified: boolean;
}

model OwnerCreateRequest {
  @minLength(2)
  name: string;
  
  @pattern("^[^@]+@[^@]+\\.[^@]+$")
  email: string;
  
  @minLength(3)
  name_space: string;  // Changed from 'namespace' to 'name_space'
  
  @minLength(6)
  password: string;
}

model OwnerUpdateRequest {
  @minLength(2)
  name: string;
  
  @minLength(3)
  name_space: string;  // Changed from 'namespace' to 'name_space'
}

// Page request models
model PageRequestResponse {
  id: string; // UUID
  ownerId: string; // UUID
  requestType: string;
  title: string;
  description: string;
  status: RequestStatus;
  pageUrl?: string;
  logoUrl?: string;
  createdAt: utcDateTime;
}

model CreatePageRequest {
  ownerId: string; // UUID
  requestType: string;
  title: string;
  description: string;
  pageUrl?: string;
  // Note: multipart.FileHeader would be handled differently in TypeSpec
  // You'd typically use @body contentType: "multipart/form-data"
}

model ChangeStatusPageRequest {
  requestId: string; // UUID
  status: RequestStatus;
}

model OwnerDeleteRequest {
  ids: string[]; // UUID array
  forceDelete: boolean;
}