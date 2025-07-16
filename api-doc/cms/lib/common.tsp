import "@typespec/http";
import "@typespec/openapi3";

using Http;
namespace Common {

  model ApiResponse<T> {
    @doc("Indicates if the request was successful")
    success: boolean;

    @doc("Human-readable message describing the result")
    message: string;

    @doc("Response data (only present on success)")
    data?: T;

    @doc("Error details (only present on failure)")
    error?: unknown;
  }

  /** Paginated API response structure */
  model PaginatedResponse<T> {
    @doc("Indicates if the request was successful")
    success: boolean;

    @doc("Human-readable message describing the result")
    message: string;

    @doc("Response data array")
    data?: T[];

    @doc("Pagination information")
    pagination: Pagination;
  }

  /** Pagination metadata */
  model Pagination {
    @doc("Current page number")
    page: int32;

    @doc("Number of items per page")
    limit: int32;

    @doc("Total number of items")
    total: int64;

    @doc("Total number of pages")
    total_pages: int32;
  }

  /** Pagination query parameters */
  model PaginationParams {
    @doc("Page number (default: 1)")
    @query
    page?: int32 = 1;

    @doc("Items per page (default: 10)")
    @query
    limit?: int32 = 10;
  }

  /** UUID format */
  @format("uuid")
  scalar uuid extends string;

  /** Email format */
  @format("email")
  scalar email extends string;

  /** Base timestamps */
  model Timestamps {
    @doc("Creation timestamp")
    created_at: utcDateTime;

    @doc("Last update timestamp")
    updated_at: utcDateTime;
  }
}