import * as jspb from 'google-protobuf'

import * as common_pb from './common_pb'; // proto import: "common.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb'; // proto import: "google/protobuf/empty.proto"
import * as google_api_annotations_pb from './google/api/annotations_pb'; // proto import: "google/api/annotations.proto"


export class Agent extends jspb.Message {
  getId(): string;
  setId(value: string): Agent;

  getName(): string;
  setName(value: string): Agent;

  getDescription(): string;
  setDescription(value: string): Agent;

  getType(): string;
  setType(value: string): Agent;

  getModelConfig(): google_protobuf_struct_pb.Struct | undefined;
  setModelConfig(value?: google_protobuf_struct_pb.Struct): Agent;
  hasModelConfig(): boolean;
  clearModelConfig(): Agent;

  getToolsList(): Array<string>;
  setToolsList(value: Array<string>): Agent;
  clearToolsList(): Agent;
  addTools(value: string, index?: number): Agent;

  getKnowledgeBasesList(): Array<string>;
  setKnowledgeBasesList(value: Array<string>): Agent;
  clearKnowledgeBasesList(): Agent;
  addKnowledgeBases(value: string, index?: number): Agent;

  getPromptTemplate(): string;
  setPromptTemplate(value: string): Agent;

  getParameters(): google_protobuf_struct_pb.Struct | undefined;
  setParameters(value?: google_protobuf_struct_pb.Struct): Agent;
  hasParameters(): boolean;
  clearParameters(): Agent;

  getStatus(): string;
  setStatus(value: string): Agent;

  getVersion(): string;
  setVersion(value: string): Agent;

  getCreatedBy(): string;
  setCreatedBy(value: string): Agent;

  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): Agent;
  clearTagsList(): Agent;
  addTags(value: string, index?: number): Agent;

  getFolder(): string;
  setFolder(value: string): Agent;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Agent;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Agent;

  getUpdatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Agent;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Agent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Agent.AsObject;
  static toObject(includeInstance: boolean, msg: Agent): Agent.AsObject;
  static serializeBinaryToWriter(message: Agent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Agent;
  static deserializeBinaryFromReader(message: Agent, reader: jspb.BinaryReader): Agent;
}

export namespace Agent {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
    type: string;
    modelConfig?: google_protobuf_struct_pb.Struct.AsObject;
    toolsList: Array<string>;
    knowledgeBasesList: Array<string>;
    promptTemplate: string;
    parameters?: google_protobuf_struct_pb.Struct.AsObject;
    status: string;
    version: string;
    createdBy: string;
    tagsList: Array<string>;
    folder: string;
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
    updatedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
  };
}

export class CreateAgentRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateAgentRequest;

  getDescription(): string;
  setDescription(value: string): CreateAgentRequest;

  getType(): string;
  setType(value: string): CreateAgentRequest;

  getModelConfig(): google_protobuf_struct_pb.Struct | undefined;
  setModelConfig(value?: google_protobuf_struct_pb.Struct): CreateAgentRequest;
  hasModelConfig(): boolean;
  clearModelConfig(): CreateAgentRequest;

  getToolsList(): Array<string>;
  setToolsList(value: Array<string>): CreateAgentRequest;
  clearToolsList(): CreateAgentRequest;
  addTools(value: string, index?: number): CreateAgentRequest;

  getKnowledgeBasesList(): Array<string>;
  setKnowledgeBasesList(value: Array<string>): CreateAgentRequest;
  clearKnowledgeBasesList(): CreateAgentRequest;
  addKnowledgeBases(value: string, index?: number): CreateAgentRequest;

  getPromptTemplate(): string;
  setPromptTemplate(value: string): CreateAgentRequest;

  getParameters(): google_protobuf_struct_pb.Struct | undefined;
  setParameters(value?: google_protobuf_struct_pb.Struct): CreateAgentRequest;
  hasParameters(): boolean;
  clearParameters(): CreateAgentRequest;

  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): CreateAgentRequest;
  clearTagsList(): CreateAgentRequest;
  addTags(value: string, index?: number): CreateAgentRequest;

  getFolder(): string;
  setFolder(value: string): CreateAgentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAgentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAgentRequest): CreateAgentRequest.AsObject;
  static serializeBinaryToWriter(message: CreateAgentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAgentRequest;
  static deserializeBinaryFromReader(message: CreateAgentRequest, reader: jspb.BinaryReader): CreateAgentRequest;
}

export namespace CreateAgentRequest {
  export type AsObject = {
    name: string;
    description: string;
    type: string;
    modelConfig?: google_protobuf_struct_pb.Struct.AsObject;
    toolsList: Array<string>;
    knowledgeBasesList: Array<string>;
    promptTemplate: string;
    parameters?: google_protobuf_struct_pb.Struct.AsObject;
    tagsList: Array<string>;
    folder: string;
  };
}

export class ListAgentsRequest extends jspb.Message {
  getPage(): number;
  setPage(value: number): ListAgentsRequest;

  getPageSize(): number;
  setPageSize(value: number): ListAgentsRequest;

  getStatus(): string;
  setStatus(value: string): ListAgentsRequest;

  getType(): string;
  setType(value: string): ListAgentsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAgentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAgentsRequest): ListAgentsRequest.AsObject;
  static serializeBinaryToWriter(message: ListAgentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAgentsRequest;
  static deserializeBinaryFromReader(message: ListAgentsRequest, reader: jspb.BinaryReader): ListAgentsRequest;
}

export namespace ListAgentsRequest {
  export type AsObject = {
    page: number;
    pageSize: number;
    status: string;
    type: string;
  };
}

export class ListAgentsResponse extends jspb.Message {
  getItemsList(): Array<Agent>;
  setItemsList(value: Array<Agent>): ListAgentsResponse;
  clearItemsList(): ListAgentsResponse;
  addItems(value?: Agent, index?: number): Agent;

  getPage(): number;
  setPage(value: number): ListAgentsResponse;

  getPageSize(): number;
  setPageSize(value: number): ListAgentsResponse;

  getTotal(): number;
  setTotal(value: number): ListAgentsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAgentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAgentsResponse): ListAgentsResponse.AsObject;
  static serializeBinaryToWriter(message: ListAgentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAgentsResponse;
  static deserializeBinaryFromReader(message: ListAgentsResponse, reader: jspb.BinaryReader): ListAgentsResponse;
}

export namespace ListAgentsResponse {
  export type AsObject = {
    itemsList: Array<Agent.AsObject>;
    page: number;
    pageSize: number;
    total: number;
  };
}

export class GetAgentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetAgentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAgentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAgentRequest): GetAgentRequest.AsObject;
  static serializeBinaryToWriter(message: GetAgentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAgentRequest;
  static deserializeBinaryFromReader(message: GetAgentRequest, reader: jspb.BinaryReader): GetAgentRequest;
}

export namespace GetAgentRequest {
  export type AsObject = {
    id: string;
  };
}

export class UpdateAgentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateAgentRequest;

  getName(): string;
  setName(value: string): UpdateAgentRequest;

  getDescription(): string;
  setDescription(value: string): UpdateAgentRequest;

  getModelConfig(): google_protobuf_struct_pb.Struct | undefined;
  setModelConfig(value?: google_protobuf_struct_pb.Struct): UpdateAgentRequest;
  hasModelConfig(): boolean;
  clearModelConfig(): UpdateAgentRequest;

  getToolsList(): Array<string>;
  setToolsList(value: Array<string>): UpdateAgentRequest;
  clearToolsList(): UpdateAgentRequest;
  addTools(value: string, index?: number): UpdateAgentRequest;

  getKnowledgeBasesList(): Array<string>;
  setKnowledgeBasesList(value: Array<string>): UpdateAgentRequest;
  clearKnowledgeBasesList(): UpdateAgentRequest;
  addKnowledgeBases(value: string, index?: number): UpdateAgentRequest;

  getPromptTemplate(): string;
  setPromptTemplate(value: string): UpdateAgentRequest;

  getParameters(): google_protobuf_struct_pb.Struct | undefined;
  setParameters(value?: google_protobuf_struct_pb.Struct): UpdateAgentRequest;
  hasParameters(): boolean;
  clearParameters(): UpdateAgentRequest;

  getStatus(): string;
  setStatus(value: string): UpdateAgentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAgentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAgentRequest): UpdateAgentRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateAgentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAgentRequest;
  static deserializeBinaryFromReader(message: UpdateAgentRequest, reader: jspb.BinaryReader): UpdateAgentRequest;
}

export namespace UpdateAgentRequest {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
    modelConfig?: google_protobuf_struct_pb.Struct.AsObject;
    toolsList: Array<string>;
    knowledgeBasesList: Array<string>;
    promptTemplate: string;
    parameters?: google_protobuf_struct_pb.Struct.AsObject;
    status: string;
  };
}

export class DeleteAgentRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteAgentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAgentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAgentRequest): DeleteAgentRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteAgentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAgentRequest;
  static deserializeBinaryFromReader(message: DeleteAgentRequest, reader: jspb.BinaryReader): DeleteAgentRequest;
}

export namespace DeleteAgentRequest {
  export type AsObject = {
    id: string;
  };
}

export class DeleteAgentResponse extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteAgentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAgentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAgentResponse): DeleteAgentResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteAgentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAgentResponse;
  static deserializeBinaryFromReader(message: DeleteAgentResponse, reader: jspb.BinaryReader): DeleteAgentResponse;
}

export namespace DeleteAgentResponse {
  export type AsObject = {
    id: string;
  };
}

