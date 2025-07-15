import "@typespec/http";
import "@typespec/openapi3";
import "../lib/common.tsp";



using Common;

namespace Models.MFA {
  /** MFA setup response */
  model MFASetupResponse {
    @doc("Secret key for authenticator app")
    secret: string;
    
    @doc("QR code URL for easy setup")
    qr_code_url: string;
    
    @doc("Base64 encoded QR code image")
    qr_code_image: string;
    
    @doc("Token ID for verification")
    token_id: int32;
    
    @doc("Manual entry key for authenticator apps")
    manual_entry: string;
  }

  /** MFA verification request */
  model MFAVerificationRequest {
    @doc("Token ID from MFA setup")
    token_id: int32;
    
    @doc("6-digit verification code from authenticator app")
    @minLength(6)
    @maxLength(6)
    @pattern("^[0-9]{6}$")
    code: string;
  }

  /** MFA login request */
  model MFALoginRequest {
    @doc("6-digit MFA verification code")
    @minLength(6)
    @maxLength(6)
    @pattern("^[0-9]{6}$")
    mfa_code: string;
  }
}