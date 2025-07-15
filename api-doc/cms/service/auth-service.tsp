import "@typespec/http";
import "@typespec/openapi3";
import "../lib/common.tsp";
import "../lib/error.tsp";
import "../model/auth.tsp";
import "../model/user.tsp";
import "../model/mfa.tsp";

using TypeSpec.Http;
using Common;
using Errors;
using Models.Auth;
using Models.User;
using Models.MFA;

namespace Services.Auth {
  @route("/cms/auth")
  @tag("Authentication")
  interface AuthOperations {
    @doc("User login")
    @post
    @route("/login")
    authLogin(
      @body @doc("Login credentials") request: LoginRequest
    ): ApiResponse<AuthResponse> | BadRequest | Unauthorized;

    @doc("User registration")
    @post
    @route("/register")
    authRegister(
      @body @doc("Registration details") request: RegisterRequest
    ): ApiResponse<{user: UserResponse, next_step: "email_verification", message: string, expires_in: string}> | BadRequest | Conflict | InternalServerError;

    @doc("Email verification")
    @post
    @route("/verify-email")
    authVerifyEmail(
      @body @doc("Email verification details") request: EmailVerificationRequest
    ): ApiResponse<never> | BadRequest | InternalServerError;

    @doc("User logout")
    @post
    @route("/logout")
    authLogout(
      @header("Authorization") @doc("Bearer token") authorization?: string,
      @body @doc("Refresh token (optional)") request?: {refresh_token?: string}
    ): ApiResponse<never> | BadRequest | InternalServerError;

    @doc("Refresh access token")
    @post
    @route("/refresh")
    authRefresh(
      @body @doc("Refresh token") request: RefreshTokenRequest
    ): ApiResponse<TokenResponse> | BadRequest | Unauthorized | InternalServerError;

    @doc("Get current user profile")
    @post
    @route("/me")
    @useAuth(BearerAuth)
    authGetMe(): ApiResponse<UserResponse> | Unauthorized | NotFound | InternalServerError;

    @doc("Update user profile")
    @put
    @route("/profile/{id}")
    @useAuth(BearerAuth)
    authUpdateProfile(
      @path @doc("User ID") id: uuid,
      @body @doc("Profile update data") request: UserUpdateRequest
    ): ApiResponse<UserResponse> | BadRequest | Unauthorized | Forbidden | NotFound | InternalServerError;

    @doc("Setup MFA for user")
    @post
    @route("/mfa/setup/{userid}")
    @useAuth(BearerAuth)
    authSetupMFA(
      @path @doc("User ID") userid: uuid
    ): ApiResponse<{setup_data: MFASetupResponse, next_step: "verify_mfa_setup", message: string}> | BadRequest | Conflict | InternalServerError;

    @doc("Verify MFA setup")
    @post
    @route("/mfa/verify/{userid}")
    @useAuth(BearerAuth)
    authVerifyMFASetup(
      @path @doc("User ID") userid: uuid,
      @body @doc("MFA verification details") request: MFAVerificationRequest
    ): ApiResponse<never> | BadRequest | NotFound | InternalServerError;

    @doc("Login with MFA")
    @post
    @route("/login/mfa/{userid}")
    authLoginWithMFA(
      @path @doc("User ID") userid: uuid,
      @body @doc("MFA login details") request: MFALoginRequest
    ): ApiResponse<AuthResponse> | BadRequest | Unauthorized | InternalServerError;
  }
}