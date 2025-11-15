import api from './api'
import type { Agent, APIResponse } from '@/types'

export const agentService = {
  // Get agent list
  getAgents: (params?: any) => {
    return api.get<any, APIResponse<{ items: Agent[]; total: number }>>('/agents', { params })
  },

  // Get agent detail
  getAgent: (id: string) => {
    return api.get<any, APIResponse<Agent>>(`/agents/${id}`)
  },

  // Create agent
  createAgent: (data: Partial<Agent>) => {
    return api.post<any, APIResponse<Agent>>('/agents', data)
  },

  // Update agent
  updateAgent: (id: string, data: Partial<Agent>) => {
    return api.put<any, APIResponse<Agent>>(`/agents/${id}`, data)
  },

  // Delete agent
  deleteAgent: (id: string) => {
    return api.delete<any, APIResponse<any>>(`/agents/${id}`)
  },
}
