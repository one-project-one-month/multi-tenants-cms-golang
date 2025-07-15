import "@typespec/compiler";
import "@typespec/http";
import "@typespec/openapi3";
import "../lib/common.tsp";
import "../lib/error.tsp";

using Common;

namespace Models.Owner {
    model OwnerCreateRequest {
        @minLength(2)
        @doc("Owner's full name")
        name: string;
        
        @doc("Owner's email address")
        @pattern("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
        email: string;
        
        @doc("Owner's namespace for multi-tenancy")
        @minLength(3)
        namespaces: string;
        
        @doc("Owner's password")
        @minLength(6)
        password: string;
    }

    model OwnerUpdateRequest {
        @doc("Owner's full name")
        @minLength(2)
        name: string;
        
        @doc("Owner's namespace for multi-tenancy")
        @minLength(3)
        namespaces: string;
    
    }

    model OwnerDeleteRequest {
        @doc("Array of owner IDs to delete")
        ids: string[];
        
        @doc("Force delete flag")
        forceDelete?: boolean;
    }

    model OwnerResponse {
        @doc("Owner's unique identifier")
        id: uuid;
        
        @doc("Owner's full name")
        name: string;
        
        @doc("Owner's email address")
        email: string;
        
        @doc("Owner's role")
        role: string;
        
        @doc("Owner's namespace")
        name_space: string;
        
        @doc("Email verification status")
        verified: boolean;
    }
}