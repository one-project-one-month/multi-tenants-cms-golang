import "@typespec/http";
import "@typespec/openapi3";
import "../lib/common.tsp";

import "./user.tsp";

using Common;
using Models.User;

namespace Models.Auth {
  /** Login request payload */
  model LoginRequest {
    @doc("User's email address")
    email: email;
    
    @doc("User's password")
    @minLength(6)
    password: string;
  }

  /** User registration request */
  model RegisterRequest {
    @doc("User's full name")
    @minLength(2)
    name: string;
    
    @doc("User's email address")
    email: email;
    
    @doc("User's password")
    @minLength(6)
    password: string;
    
    @doc("User's role (optional)")
    role?: UserRole;
  }

  /** Email verification request */
  model EmailVerificationRequest {
    @doc("Email address to verify")
    email: email;
    
    @doc("6-digit verification code")
    @minLength(6)
    @maxLength(6)
    @pattern("^[0-9]{6}$")
    code: string;
  }

  /** Refresh token request */
  model RefreshTokenRequest {
    @doc("Refresh token")
    refresh_token: string;
  }

  /** Authentication response */
  model AuthResponse {
    @doc("User information")
    user: UserResponse;
    
    @doc("JWT access token")
    access_token?: string;
    
    @doc("Refresh token")
    refresh_token?: string;
    
    @doc("Token expiration time")
    expires_at?: utcDateTime;
  }

  /** Token response */
  model TokenResponse {
    @doc("JWT access token")
    access_token: string;
    
    @doc("Refresh token")
    refresh_token: string;
    
    @doc("Token expiration time")
    expires_at: utcDateTime;
  }
}