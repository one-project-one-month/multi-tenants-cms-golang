import "@typespec/http";
import "@typespec/openapi3";
import "../lib/common.tsp";
import "../lib/error.tsp";
import "../model/owner.tsp";

using TypeSpec.Http;
using Common;
using Errors;
using Models.Owner;

namespace Services.Owner {
  @route("/cms/owners")
  @tag("Owner Management")
  interface OwnerOperations {
    @doc("Create a new owner")
    @post
    @route("/")
    create(
      @body @doc("Owner creation details") request: OwnerCreateRequest
    ): ApiResponse<OwnerResponse> | BadRequest | InternalServerError;

    @doc("Update an existing owner")
    @put
    @route("/{id}")
    update(
      @path @doc("Owner ID") id: string,
      @body @doc("Owner update details") request: OwnerUpdateRequest
    ): ApiResponse<OwnerResponse> | BadRequest | InternalServerError;

    @doc("Get all owners")
    @get
    @route("/")
    getAll(
      @query @doc("Page number") page?: int32,
      @query @doc("Items per page") limit?: int32
    ): ApiResponse<OwnerResponse[]> | InternalServerError;

    @doc("Get owner by ID")
    @get
    @route("/{id}")
    getById(
      @path @doc("Owner ID") id: string
    ): ApiResponse<OwnerResponse> | InternalServerError;

    @doc("Delete owners")
    @delete
    @route("/")
    delete(
      @body @doc("Owner deletion details") request: OwnerDeleteRequest
    ): ApiResponse<never> | BadRequest | InternalServerError;
  }
}