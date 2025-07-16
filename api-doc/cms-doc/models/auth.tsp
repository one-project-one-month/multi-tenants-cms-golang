// models/auth.tsp
import "@typespec/http";
import "./common.tsp";

using TypeSpec.Http;
using CmsApi.Models.Common;

namespace CmsApi.Models;

// Request Models
model LoginRequest {
  @pattern("^[^@]+@[^@]+\\.[^@]+$")
  email: string;
  
  @minLength(6)
  password: string;
}

model RegisterRequest {
  @minLength(2)
  name: string;
  
  @pattern("^[^@]+@[^@]+\\.[^@]+$")
  email: string;
  
  @minLength(6)
  password: string;
  
  role?: string;
}

model MFAVerificationRequest {
  token_id: uint32;
  
  @pattern("^[0-9]{6}$")
  code: string;
}

model MFALoginRequest {
  @pattern("^[^@]+@[^@]+\\.[^@]+$")
  email: string;
  
  password: string;
  
  @pattern("^[0-9]{6}$")
  mfa_code?: string;
}

model RefreshTokenRequest {
  refresh_token: string;
}

model UserUpdateRequest {
  @minLength(2)
  name: string;
}

model EmailVerificationRequest {
  @pattern("^[^@]+@[^@]+\\.[^@]+$")
  email: string;
  
  @pattern("^[0-9]{6}$")
  code: string;
}

model LogoutRequest {
  refresh_token?: string;
}

// Response Models
model AuthResponse {
  user: UserResponse;
  access_token?: string;
  refresh_token?: string;
  expires_at?: utcDateTime;
}

model MFASetupResponse {
  secret: string;
  qr_code_url: string;
  qr_code_image: string; // Base64 encoded PNG
  token_id: uint32;
  manual_entry: string;
}

model UserResponse {
  id: string; // UUID
  name: string;
  email: string;
  role: string;
  verified: boolean;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model TokenResponse {
  access_token: string;
  refresh_token: string;
  expires_at: utcDateTime;
}

model RegistrationResponse {
  user: UserResponse;
  next_step: "email_verification";
  message: string;
  expires_in: string;
}

model MFASetupSuccessResponse {
  setup_data: MFASetupResponse;
  next_step: "verify_mfa_setup";
  message: string;
}