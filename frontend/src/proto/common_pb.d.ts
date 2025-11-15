import * as jspb from 'google-protobuf'

import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"


export class Response extends jspb.Message {
  getCode(): number;
  setCode(value: number): Response;

  getMessage(): string;
  setMessage(value: string): Response;

  getData(): google_protobuf_struct_pb.Struct | undefined;
  setData(value?: google_protobuf_struct_pb.Struct): Response;
  hasData(): boolean;
  clearData(): Response;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Response.AsObject;
  static toObject(includeInstance: boolean, msg: Response): Response.AsObject;
  static serializeBinaryToWriter(message: Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Response;
  static deserializeBinaryFromReader(message: Response, reader: jspb.BinaryReader): Response;
}

export namespace Response {
  export type AsObject = {
    code: number;
    message: string;
    data?: google_protobuf_struct_pb.Struct.AsObject;
  };
}

export class PaginationRequest extends jspb.Message {
  getPage(): number;
  setPage(value: number): PaginationRequest;

  getPageSize(): number;
  setPageSize(value: number): PaginationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PaginationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PaginationRequest): PaginationRequest.AsObject;
  static serializeBinaryToWriter(message: PaginationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PaginationRequest;
  static deserializeBinaryFromReader(message: PaginationRequest, reader: jspb.BinaryReader): PaginationRequest;
}

export namespace PaginationRequest {
  export type AsObject = {
    page: number;
    pageSize: number;
  };
}

export class PaginationResponse extends jspb.Message {
  getItemsList(): Array<google_protobuf_struct_pb.Struct>;
  setItemsList(value: Array<google_protobuf_struct_pb.Struct>): PaginationResponse;
  clearItemsList(): PaginationResponse;
  addItems(value?: google_protobuf_struct_pb.Struct, index?: number): google_protobuf_struct_pb.Struct;

  getPage(): number;
  setPage(value: number): PaginationResponse;

  getPageSize(): number;
  setPageSize(value: number): PaginationResponse;

  getTotal(): number;
  setTotal(value: number): PaginationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PaginationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PaginationResponse): PaginationResponse.AsObject;
  static serializeBinaryToWriter(message: PaginationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PaginationResponse;
  static deserializeBinaryFromReader(message: PaginationResponse, reader: jspb.BinaryReader): PaginationResponse;
}

export namespace PaginationResponse {
  export type AsObject = {
    itemsList: Array<google_protobuf_struct_pb.Struct.AsObject>;
    page: number;
    pageSize: number;
    total: number;
  };
}

export class IDRequest extends jspb.Message {
  getId(): string;
  setId(value: string): IDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: IDRequest): IDRequest.AsObject;
  static serializeBinaryToWriter(message: IDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDRequest;
  static deserializeBinaryFromReader(message: IDRequest, reader: jspb.BinaryReader): IDRequest;
}

export namespace IDRequest {
  export type AsObject = {
    id: string;
  };
}

export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  };
}

