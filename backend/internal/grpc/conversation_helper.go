package grpc

import (
	pb "agent-platform/gen/go"
	"agent-platform/internal/model/ent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// entConversationToProto converts ent.Conversation to pb.Conversation
func entConversationToProto(conv *ent.Conversation) *pb.Conversation {
	pbConv := &pb.Conversation{
		Id:            conv.ID,
		AgentId:       conv.AgentID,
		UserId:        conv.UserID,
		Title:         conv.Title,
		Status:        conv.Status,
		CreatedAt:     timestamppb.New(conv.CreatedAt),
		UpdatedAt:     timestamppb.New(conv.UpdatedAt),
	}

	// Convert messages
	if conv.Messages != nil {
		pbMessages := make([]*pb.Message, 0)
		for _, msg := range conv.Messages {
			if msgMap, ok := msg.(map[string]interface{}); ok {
				pbMsg := &pb.Message{}
				if id, ok := msgMap["id"].(string); ok {
					pbMsg.Id = id
				}
				if role, ok := msgMap["role"].(string); ok {
					pbMsg.Role = role
				}
				if content, ok := msgMap["content"].(string); ok {
					pbMsg.Content = content
				}
				// Handle timestamp - could be Unix timestamp (int64) or time.Time
				if _, ok := msgMap["timestamp"].(int64); ok {
					pbMsg.Timestamp = timestamppb.New(conv.CreatedAt.Add(0)) // Use created time as base
				}
				pbMessages = append(pbMessages, pbMsg)
			}
		}
		pbConv.Messages = pbMessages
	}

	// Set last message time if available
	if !conv.LastMessageAt.IsZero() {
		pbConv.LastMessageAt = timestamppb.New(conv.LastMessageAt)
	}

	return pbConv
}
