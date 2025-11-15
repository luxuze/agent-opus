import * as jspb from 'google-protobuf'

import * as common_pb from './common_pb'; // proto import: "common.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb'; // proto import: "google/protobuf/empty.proto"
import * as google_api_annotations_pb from './google/api/annotations_pb'; // proto import: "google/api/annotations.proto"


export class Tool extends jspb.Message {
  getId(): string;
  setId(value: string): Tool;

  getName(): string;
  setName(value: string): Tool;

  getDescription(): string;
  setDescription(value: string): Tool;

  getType(): string;
  setType(value: string): Tool;

  getSchema(): google_protobuf_struct_pb.Struct | undefined;
  setSchema(value?: google_protobuf_struct_pb.Struct): Tool;
  hasSchema(): boolean;
  clearSchema(): Tool;

  getImplementation(): string;
  setImplementation(value: string): Tool;

  getVersion(): string;
  setVersion(value: string): Tool;

  getIsPublic(): boolean;
  setIsPublic(value: boolean): Tool;

  getCreatedBy(): string;
  setCreatedBy(value: string): Tool;

  getCategory(): string;
  setCategory(value: string): Tool;

  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): Tool;
  clearTagsList(): Tool;
  addTags(value: string, index?: number): Tool;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Tool;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Tool;

  getUpdatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Tool;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Tool;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Tool.AsObject;
  static toObject(includeInstance: boolean, msg: Tool): Tool.AsObject;
  static serializeBinaryToWriter(message: Tool, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Tool;
  static deserializeBinaryFromReader(message: Tool, reader: jspb.BinaryReader): Tool;
}

export namespace Tool {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
    type: string;
    schema?: google_protobuf_struct_pb.Struct.AsObject;
    implementation: string;
    version: string;
    isPublic: boolean;
    createdBy: string;
    category: string;
    tagsList: Array<string>;
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
    updatedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
  };
}

export class CreateToolRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateToolRequest;

  getDescription(): string;
  setDescription(value: string): CreateToolRequest;

  getType(): string;
  setType(value: string): CreateToolRequest;

  getSchema(): google_protobuf_struct_pb.Struct | undefined;
  setSchema(value?: google_protobuf_struct_pb.Struct): CreateToolRequest;
  hasSchema(): boolean;
  clearSchema(): CreateToolRequest;

  getImplementation(): string;
  setImplementation(value: string): CreateToolRequest;

  getCategory(): string;
  setCategory(value: string): CreateToolRequest;

  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): CreateToolRequest;
  clearTagsList(): CreateToolRequest;
  addTags(value: string, index?: number): CreateToolRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateToolRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateToolRequest): CreateToolRequest.AsObject;
  static serializeBinaryToWriter(message: CreateToolRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateToolRequest;
  static deserializeBinaryFromReader(message: CreateToolRequest, reader: jspb.BinaryReader): CreateToolRequest;
}

export namespace CreateToolRequest {
  export type AsObject = {
    name: string;
    description: string;
    type: string;
    schema?: google_protobuf_struct_pb.Struct.AsObject;
    implementation: string;
    category: string;
    tagsList: Array<string>;
  };
}

export class ListToolsRequest extends jspb.Message {
  getType(): string;
  setType(value: string): ListToolsRequest;

  getCategory(): string;
  setCategory(value: string): ListToolsRequest;

  getPage(): number;
  setPage(value: number): ListToolsRequest;

  getPageSize(): number;
  setPageSize(value: number): ListToolsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListToolsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListToolsRequest): ListToolsRequest.AsObject;
  static serializeBinaryToWriter(message: ListToolsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListToolsRequest;
  static deserializeBinaryFromReader(message: ListToolsRequest, reader: jspb.BinaryReader): ListToolsRequest;
}

export namespace ListToolsRequest {
  export type AsObject = {
    type: string;
    category: string;
    page: number;
    pageSize: number;
  };
}

export class ListToolsResponse extends jspb.Message {
  getItemsList(): Array<Tool>;
  setItemsList(value: Array<Tool>): ListToolsResponse;
  clearItemsList(): ListToolsResponse;
  addItems(value?: Tool, index?: number): Tool;

  getPage(): number;
  setPage(value: number): ListToolsResponse;

  getPageSize(): number;
  setPageSize(value: number): ListToolsResponse;

  getTotal(): number;
  setTotal(value: number): ListToolsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListToolsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListToolsResponse): ListToolsResponse.AsObject;
  static serializeBinaryToWriter(message: ListToolsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListToolsResponse;
  static deserializeBinaryFromReader(message: ListToolsResponse, reader: jspb.BinaryReader): ListToolsResponse;
}

export namespace ListToolsResponse {
  export type AsObject = {
    itemsList: Array<Tool.AsObject>;
    page: number;
    pageSize: number;
    total: number;
  };
}

export class GetToolRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetToolRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetToolRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetToolRequest): GetToolRequest.AsObject;
  static serializeBinaryToWriter(message: GetToolRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetToolRequest;
  static deserializeBinaryFromReader(message: GetToolRequest, reader: jspb.BinaryReader): GetToolRequest;
}

export namespace GetToolRequest {
  export type AsObject = {
    id: string;
  };
}

export class DeleteToolRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteToolRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteToolRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteToolRequest): DeleteToolRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteToolRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteToolRequest;
  static deserializeBinaryFromReader(message: DeleteToolRequest, reader: jspb.BinaryReader): DeleteToolRequest;
}

export namespace DeleteToolRequest {
  export type AsObject = {
    id: string;
  };
}

export class DeleteToolResponse extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteToolResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteToolResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteToolResponse): DeleteToolResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteToolResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteToolResponse;
  static deserializeBinaryFromReader(message: DeleteToolResponse, reader: jspb.BinaryReader): DeleteToolResponse;
}

export namespace DeleteToolResponse {
  export type AsObject = {
    id: string;
  };
}

