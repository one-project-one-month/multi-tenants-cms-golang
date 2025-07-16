import "@typespec/http";
import "@typespec/openapi3";
import "../lib/common.tsp";

using Common;

namespace Models.User {
  /** User role enumeration */
  enum UserRole {
    @doc("Administrator role")
    admin: "admin",
    
    @doc("Regular user role")
    user: "user",
    
    @doc("Owner role")
    owner: "owner"
  }

  /** User profile information */
  model UserResponse {
    @doc("Unique user identifier")
    id: uuid;
    
    @doc("User's full name")
    @minLength(2)
    name: string;
    
    @doc("User's email address")
    email: email;
    
    @doc("User's role in the system")
    role: UserRole;
    
    @doc("Whether the user's email is verified")
    verified: boolean;
    
    ...Timestamps;
  }

  /** User update request */
  model UserUpdateRequest {
    @doc("User's full name")
    @minLength(2)
    name: string;
  }

  /** Owner-specific response */
  model OwnerResponse {
    @doc("Unique owner identifier")
    id: uuid;
    
    @doc("Owner's full name")
    name: string;
    
    @doc("Owner's email address")
    email: email;
    
    @doc("Owner's role")
    role: UserRole;
    
    @doc("Owner's namespace")
    name_space: string;
    
    @doc("Whether the owner's email is verified")
    verified: boolean;
  }
}