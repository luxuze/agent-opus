import * as jspb from 'google-protobuf'

import * as common_pb from './common_pb'; // proto import: "common.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb'; // proto import: "google/protobuf/empty.proto"
import * as google_api_annotations_pb from './google/api/annotations_pb'; // proto import: "google/api/annotations.proto"


export class Message extends jspb.Message {
  getId(): string;
  setId(value: string): Message;

  getRole(): string;
  setRole(value: string): Message;

  getContent(): string;
  setContent(value: string): Message;

  getMetadata(): google_protobuf_struct_pb.Struct | undefined;
  setMetadata(value?: google_protobuf_struct_pb.Struct): Message;
  hasMetadata(): boolean;
  clearMetadata(): Message;

  getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): Message;
  hasTimestamp(): boolean;
  clearTimestamp(): Message;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Message.AsObject;
  static toObject(includeInstance: boolean, msg: Message): Message.AsObject;
  static serializeBinaryToWriter(message: Message, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Message;
  static deserializeBinaryFromReader(message: Message, reader: jspb.BinaryReader): Message;
}

export namespace Message {
  export type AsObject = {
    id: string;
    role: string;
    content: string;
    metadata?: google_protobuf_struct_pb.Struct.AsObject;
    timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject;
  };
}

export class Conversation extends jspb.Message {
  getId(): string;
  setId(value: string): Conversation;

  getAgentId(): string;
  setAgentId(value: string): Conversation;

  getUserId(): string;
  setUserId(value: string): Conversation;

  getTitle(): string;
  setTitle(value: string): Conversation;

  getMessagesList(): Array<Message>;
  setMessagesList(value: Array<Message>): Conversation;
  clearMessagesList(): Conversation;
  addMessages(value?: Message, index?: number): Message;

  getContext(): google_protobuf_struct_pb.Struct | undefined;
  setContext(value?: google_protobuf_struct_pb.Struct): Conversation;
  hasContext(): boolean;
  clearContext(): Conversation;

  getStatus(): string;
  setStatus(value: string): Conversation;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Conversation;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Conversation;

  getUpdatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Conversation;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Conversation;

  getLastMessageAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastMessageAt(value?: google_protobuf_timestamp_pb.Timestamp): Conversation;
  hasLastMessageAt(): boolean;
  clearLastMessageAt(): Conversation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Conversation.AsObject;
  static toObject(includeInstance: boolean, msg: Conversation): Conversation.AsObject;
  static serializeBinaryToWriter(message: Conversation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Conversation;
  static deserializeBinaryFromReader(message: Conversation, reader: jspb.BinaryReader): Conversation;
}

export namespace Conversation {
  export type AsObject = {
    id: string;
    agentId: string;
    userId: string;
    title: string;
    messagesList: Array<Message.AsObject>;
    context?: google_protobuf_struct_pb.Struct.AsObject;
    status: string;
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
    updatedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
    lastMessageAt?: google_protobuf_timestamp_pb.Timestamp.AsObject;
  };
}

export class CreateConversationRequest extends jspb.Message {
  getAgentId(): string;
  setAgentId(value: string): CreateConversationRequest;

  getTitle(): string;
  setTitle(value: string): CreateConversationRequest;

  getContext(): google_protobuf_struct_pb.Struct | undefined;
  setContext(value?: google_protobuf_struct_pb.Struct): CreateConversationRequest;
  hasContext(): boolean;
  clearContext(): CreateConversationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateConversationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateConversationRequest): CreateConversationRequest.AsObject;
  static serializeBinaryToWriter(message: CreateConversationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateConversationRequest;
  static deserializeBinaryFromReader(message: CreateConversationRequest, reader: jspb.BinaryReader): CreateConversationRequest;
}

export namespace CreateConversationRequest {
  export type AsObject = {
    agentId: string;
    title: string;
    context?: google_protobuf_struct_pb.Struct.AsObject;
  };
}

export class SendMessageRequest extends jspb.Message {
  getConversationId(): string;
  setConversationId(value: string): SendMessageRequest;

  getContent(): string;
  setContent(value: string): SendMessageRequest;

  getRole(): string;
  setRole(value: string): SendMessageRequest;

  getMetadata(): google_protobuf_struct_pb.Struct | undefined;
  setMetadata(value?: google_protobuf_struct_pb.Struct): SendMessageRequest;
  hasMetadata(): boolean;
  clearMetadata(): SendMessageRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendMessageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SendMessageRequest): SendMessageRequest.AsObject;
  static serializeBinaryToWriter(message: SendMessageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendMessageRequest;
  static deserializeBinaryFromReader(message: SendMessageRequest, reader: jspb.BinaryReader): SendMessageRequest;
}

export namespace SendMessageRequest {
  export type AsObject = {
    conversationId: string;
    content: string;
    role: string;
    metadata?: google_protobuf_struct_pb.Struct.AsObject;
  };
}

export class SendMessageResponse extends jspb.Message {
  getConversationId(): string;
  setConversationId(value: string): SendMessageResponse;

  getMessagesList(): Array<Message>;
  setMessagesList(value: Array<Message>): SendMessageResponse;
  clearMessagesList(): SendMessageResponse;
  addMessages(value?: Message, index?: number): Message;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendMessageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SendMessageResponse): SendMessageResponse.AsObject;
  static serializeBinaryToWriter(message: SendMessageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendMessageResponse;
  static deserializeBinaryFromReader(message: SendMessageResponse, reader: jspb.BinaryReader): SendMessageResponse;
}

export namespace SendMessageResponse {
  export type AsObject = {
    conversationId: string;
    messagesList: Array<Message.AsObject>;
  };
}

export class GetConversationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetConversationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConversationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetConversationRequest): GetConversationRequest.AsObject;
  static serializeBinaryToWriter(message: GetConversationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConversationRequest;
  static deserializeBinaryFromReader(message: GetConversationRequest, reader: jspb.BinaryReader): GetConversationRequest;
}

export namespace GetConversationRequest {
  export type AsObject = {
    id: string;
  };
}

export class ListConversationsRequest extends jspb.Message {
  getAgentId(): string;
  setAgentId(value: string): ListConversationsRequest;

  getPage(): number;
  setPage(value: number): ListConversationsRequest;

  getPageSize(): number;
  setPageSize(value: number): ListConversationsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListConversationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListConversationsRequest): ListConversationsRequest.AsObject;
  static serializeBinaryToWriter(message: ListConversationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListConversationsRequest;
  static deserializeBinaryFromReader(message: ListConversationsRequest, reader: jspb.BinaryReader): ListConversationsRequest;
}

export namespace ListConversationsRequest {
  export type AsObject = {
    agentId: string;
    page: number;
    pageSize: number;
  };
}

export class ListConversationsResponse extends jspb.Message {
  getItemsList(): Array<Conversation>;
  setItemsList(value: Array<Conversation>): ListConversationsResponse;
  clearItemsList(): ListConversationsResponse;
  addItems(value?: Conversation, index?: number): Conversation;

  getPage(): number;
  setPage(value: number): ListConversationsResponse;

  getPageSize(): number;
  setPageSize(value: number): ListConversationsResponse;

  getTotal(): number;
  setTotal(value: number): ListConversationsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListConversationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListConversationsResponse): ListConversationsResponse.AsObject;
  static serializeBinaryToWriter(message: ListConversationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListConversationsResponse;
  static deserializeBinaryFromReader(message: ListConversationsResponse, reader: jspb.BinaryReader): ListConversationsResponse;
}

export namespace ListConversationsResponse {
  export type AsObject = {
    itemsList: Array<Conversation.AsObject>;
    page: number;
    pageSize: number;
    total: number;
  };
}

