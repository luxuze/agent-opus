import api from './api'
import type { Conversation, Message, APIResponse } from '@/types'

export const conversationService = {
  // Get conversation list
  getConversations: (params?: any) => {
    return api.get<any, APIResponse<{ items: Conversation[]; total: number }>>('/conversations', { params })
  },

  // Get conversation detail
  getConversation: (id: string) => {
    return api.get<any, APIResponse<Conversation>>(`/conversations/${id}`)
  },

  // Create conversation
  createConversation: (data: { agent_id: string; title?: string }) => {
    return api.post<any, APIResponse<Conversation>>('/conversations', data)
  },

  // Send message
  sendMessage: (conversationId: string, content: string) => {
    return api.post<any, APIResponse<{ messages: Message[] }>>(`/conversations/${conversationId}/messages`, {
      content,
      role: 'user',
    })
  },
}
