export interface Agent {
  id: string
  name: string
  description?: string
  type: 'single' | 'multi'
  model_config?: Record<string, any>
  tools?: string[]
  knowledge_bases?: string[]
  prompt_template?: string
  parameters?: Record<string, any>
  status: 'draft' | 'published' | 'archived'
  version: string
  created_by: string
  created_at?: string
  updated_at?: string
  tags?: string[]
  folder?: string
  is_public?: boolean
}

export interface Conversation {
  id: string
  agent_id: string
  user_id: string
  title?: string
  messages: Message[]
  context?: Record<string, any>
  metadata?: Record<string, any>
  status: 'active' | 'closed'
  created_at: string
  updated_at?: string
  last_message_at?: string
}

export interface Message {
  id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp: string
  metadata?: Record<string, any>
}

export interface Tool {
  id: string
  name: string
  description?: string
  type: 'function' | 'api' | 'plugin'
  schema?: Record<string, any>
  implementation?: string
  version: string
  is_public: boolean
  created_by: string
  category?: string
  tags?: string[]
}

export interface KnowledgeBase {
  id: string
  name: string
  description?: string
  type: 'document' | 'database' | 'api'
  embedding_model: string
  chunk_config?: Record<string, any>
  documents?: any[]
  metadata?: Record<string, any>
  created_by: string
  created_at: string
  document_count: number
  vector_count: number
}

export interface APIResponse<T = any> {
  code: number
  message: string
  data: T
  timestamp?: number
  request_id?: string
}
