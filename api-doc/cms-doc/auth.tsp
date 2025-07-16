import "@typespec/http";
import "@typespec/rest";
import "./models/auth.tsp";
import "./models/common.tsp";
import "./models/extends.tsp";
using TypeSpec.Http;
using TypeSpec.Rest;
using CmsApi.Models;

namespace CmsApi.Auth;

model ErrorResponse {
    Success: boolean, 
    message: string, 
    Error: string,
}
model SuccessResponse {
    Success: boolean, 
    message: string, 
    Error: string,
}
@route("/auth")
interface AuthService {
  @post
  @route("/login")
  login(@body request: LoginRequest): AuthResponse | ErrorResponse;

  @post
  @route("/register")
  register(@body request: RegisterRequest): {
    @statusCode _: 201;
    @body response: RegistrationResponse;
  } | ErrorResponse;

  @post
  @route("/logout")
  logout(
    @header Authorization?: string,
    @body request?: LogoutRequest
  ): SuccessResponse | ErrorResponse;

  @post
  @route("/refresh")
  refresh(@body request: RefreshTokenRequest): TokenResponse | ErrorResponse;

  @post
  @route("/me")
  @useAuth(BearerAuth)
  getMe(): UserResponse | ErrorResponse;

  @put
  @route("/profile/{id}")
  @useAuth(BearerAuth)
  updateProfile(
    @path id: string,
    @body request: UserUpdateRequest
  ): UserResponse | ErrorResponse;

  @post
  @route("/mfa/setup")
  @useAuth(BearerAuth)
  setupMFA(): MFASetupSuccessResponse | ErrorResponse;

  @post
  @route("/mfa/verify")
  @useAuth(BearerAuth)
  verifyMFASetup(@body request: MFAVerificationRequest): SuccessResponse | ErrorResponse;

  @post
  @route("/mfa/login")
  loginWithMFA(@body request: MFALoginRequest): AuthResponse | ErrorResponse;

  @post
  @route("/email/verify")
  verifyEmail(@body request: EmailVerificationRequest): SuccessResponse | ErrorResponse;
}