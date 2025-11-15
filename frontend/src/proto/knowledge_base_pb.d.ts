import * as jspb from 'google-protobuf'

import * as common_pb from './common_pb'; // proto import: "common.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb'; // proto import: "google/protobuf/empty.proto"
import * as google_api_annotations_pb from './google/api/annotations_pb'; // proto import: "google/api/annotations.proto"


export class Document extends jspb.Message {
  getId(): string;
  setId(value: string): Document;

  getKnowledgeBaseId(): string;
  setKnowledgeBaseId(value: string): Document;

  getTitle(): string;
  setTitle(value: string): Document;

  getContent(): string;
  setContent(value: string): Document;

  getMetadata(): google_protobuf_struct_pb.Struct | undefined;
  setMetadata(value?: google_protobuf_struct_pb.Struct): Document;
  hasMetadata(): boolean;
  clearMetadata(): Document;

  getStatus(): string;
  setStatus(value: string): Document;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Document;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Document;

  getUpdatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Document;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Document;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Document.AsObject;
  static toObject(includeInstance: boolean, msg: Document): Document.AsObject;
  static serializeBinaryToWriter(message: Document, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Document;
  static deserializeBinaryFromReader(message: Document, reader: jspb.BinaryReader): Document;
}

export namespace Document {
  export type AsObject = {
    id: string;
    knowledgeBaseId: string;
    title: string;
    content: string;
    metadata?: google_protobuf_struct_pb.Struct.AsObject;
    status: string;
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
    updatedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
  };
}

export class KnowledgeBase extends jspb.Message {
  getId(): string;
  setId(value: string): KnowledgeBase;

  getName(): string;
  setName(value: string): KnowledgeBase;

  getDescription(): string;
  setDescription(value: string): KnowledgeBase;

  getType(): string;
  setType(value: string): KnowledgeBase;

  getEmbeddingModel(): string;
  setEmbeddingModel(value: string): KnowledgeBase;

  getChunkConfig(): google_protobuf_struct_pb.Struct | undefined;
  setChunkConfig(value?: google_protobuf_struct_pb.Struct): KnowledgeBase;
  hasChunkConfig(): boolean;
  clearChunkConfig(): KnowledgeBase;

  getCreatedBy(): string;
  setCreatedBy(value: string): KnowledgeBase;

  getDocumentCount(): number;
  setDocumentCount(value: number): KnowledgeBase;

  getVectorCount(): number;
  setVectorCount(value: number): KnowledgeBase;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): KnowledgeBase;
  hasCreatedAt(): boolean;
  clearCreatedAt(): KnowledgeBase;

  getUpdatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: google_protobuf_timestamp_pb.Timestamp): KnowledgeBase;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): KnowledgeBase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KnowledgeBase.AsObject;
  static toObject(includeInstance: boolean, msg: KnowledgeBase): KnowledgeBase.AsObject;
  static serializeBinaryToWriter(message: KnowledgeBase, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KnowledgeBase;
  static deserializeBinaryFromReader(message: KnowledgeBase, reader: jspb.BinaryReader): KnowledgeBase;
}

export namespace KnowledgeBase {
  export type AsObject = {
    id: string;
    name: string;
    description: string;
    type: string;
    embeddingModel: string;
    chunkConfig?: google_protobuf_struct_pb.Struct.AsObject;
    createdBy: string;
    documentCount: number;
    vectorCount: number;
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
    updatedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
  };
}

export class CreateKnowledgeBaseRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateKnowledgeBaseRequest;

  getDescription(): string;
  setDescription(value: string): CreateKnowledgeBaseRequest;

  getType(): string;
  setType(value: string): CreateKnowledgeBaseRequest;

  getEmbeddingModel(): string;
  setEmbeddingModel(value: string): CreateKnowledgeBaseRequest;

  getChunkConfig(): google_protobuf_struct_pb.Struct | undefined;
  setChunkConfig(value?: google_protobuf_struct_pb.Struct): CreateKnowledgeBaseRequest;
  hasChunkConfig(): boolean;
  clearChunkConfig(): CreateKnowledgeBaseRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateKnowledgeBaseRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateKnowledgeBaseRequest): CreateKnowledgeBaseRequest.AsObject;
  static serializeBinaryToWriter(message: CreateKnowledgeBaseRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateKnowledgeBaseRequest;
  static deserializeBinaryFromReader(message: CreateKnowledgeBaseRequest, reader: jspb.BinaryReader): CreateKnowledgeBaseRequest;
}

export namespace CreateKnowledgeBaseRequest {
  export type AsObject = {
    name: string;
    description: string;
    type: string;
    embeddingModel: string;
    chunkConfig?: google_protobuf_struct_pb.Struct.AsObject;
  };
}

export class ListKnowledgeBasesRequest extends jspb.Message {
  getType(): string;
  setType(value: string): ListKnowledgeBasesRequest;

  getPage(): number;
  setPage(value: number): ListKnowledgeBasesRequest;

  getPageSize(): number;
  setPageSize(value: number): ListKnowledgeBasesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListKnowledgeBasesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListKnowledgeBasesRequest): ListKnowledgeBasesRequest.AsObject;
  static serializeBinaryToWriter(message: ListKnowledgeBasesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListKnowledgeBasesRequest;
  static deserializeBinaryFromReader(message: ListKnowledgeBasesRequest, reader: jspb.BinaryReader): ListKnowledgeBasesRequest;
}

export namespace ListKnowledgeBasesRequest {
  export type AsObject = {
    type: string;
    page: number;
    pageSize: number;
  };
}

export class ListKnowledgeBasesResponse extends jspb.Message {
  getItemsList(): Array<KnowledgeBase>;
  setItemsList(value: Array<KnowledgeBase>): ListKnowledgeBasesResponse;
  clearItemsList(): ListKnowledgeBasesResponse;
  addItems(value?: KnowledgeBase, index?: number): KnowledgeBase;

  getPage(): number;
  setPage(value: number): ListKnowledgeBasesResponse;

  getPageSize(): number;
  setPageSize(value: number): ListKnowledgeBasesResponse;

  getTotal(): number;
  setTotal(value: number): ListKnowledgeBasesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListKnowledgeBasesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListKnowledgeBasesResponse): ListKnowledgeBasesResponse.AsObject;
  static serializeBinaryToWriter(message: ListKnowledgeBasesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListKnowledgeBasesResponse;
  static deserializeBinaryFromReader(message: ListKnowledgeBasesResponse, reader: jspb.BinaryReader): ListKnowledgeBasesResponse;
}

export namespace ListKnowledgeBasesResponse {
  export type AsObject = {
    itemsList: Array<KnowledgeBase.AsObject>;
    page: number;
    pageSize: number;
    total: number;
  };
}

export class GetKnowledgeBaseRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetKnowledgeBaseRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetKnowledgeBaseRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetKnowledgeBaseRequest): GetKnowledgeBaseRequest.AsObject;
  static serializeBinaryToWriter(message: GetKnowledgeBaseRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetKnowledgeBaseRequest;
  static deserializeBinaryFromReader(message: GetKnowledgeBaseRequest, reader: jspb.BinaryReader): GetKnowledgeBaseRequest;
}

export namespace GetKnowledgeBaseRequest {
  export type AsObject = {
    id: string;
  };
}

export class UploadDocumentRequest extends jspb.Message {
  getKnowledgeBaseId(): string;
  setKnowledgeBaseId(value: string): UploadDocumentRequest;

  getTitle(): string;
  setTitle(value: string): UploadDocumentRequest;

  getContent(): string;
  setContent(value: string): UploadDocumentRequest;

  getMetadata(): google_protobuf_struct_pb.Struct | undefined;
  setMetadata(value?: google_protobuf_struct_pb.Struct): UploadDocumentRequest;
  hasMetadata(): boolean;
  clearMetadata(): UploadDocumentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadDocumentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UploadDocumentRequest): UploadDocumentRequest.AsObject;
  static serializeBinaryToWriter(message: UploadDocumentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadDocumentRequest;
  static deserializeBinaryFromReader(message: UploadDocumentRequest, reader: jspb.BinaryReader): UploadDocumentRequest;
}

export namespace UploadDocumentRequest {
  export type AsObject = {
    knowledgeBaseId: string;
    title: string;
    content: string;
    metadata?: google_protobuf_struct_pb.Struct.AsObject;
  };
}

export class DeleteKnowledgeBaseRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteKnowledgeBaseRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteKnowledgeBaseRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteKnowledgeBaseRequest): DeleteKnowledgeBaseRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteKnowledgeBaseRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteKnowledgeBaseRequest;
  static deserializeBinaryFromReader(message: DeleteKnowledgeBaseRequest, reader: jspb.BinaryReader): DeleteKnowledgeBaseRequest;
}

export namespace DeleteKnowledgeBaseRequest {
  export type AsObject = {
    id: string;
  };
}

export class DeleteKnowledgeBaseResponse extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteKnowledgeBaseResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteKnowledgeBaseResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteKnowledgeBaseResponse): DeleteKnowledgeBaseResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteKnowledgeBaseResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteKnowledgeBaseResponse;
  static deserializeBinaryFromReader(message: DeleteKnowledgeBaseResponse, reader: jspb.BinaryReader): DeleteKnowledgeBaseResponse;
}

export namespace DeleteKnowledgeBaseResponse {
  export type AsObject = {
    id: string;
  };
}

