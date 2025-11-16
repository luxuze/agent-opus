import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { message, Space, Select, Tag } from 'antd'
import { DatabaseOutlined, ThunderboltOutlined } from '@ant-design/icons'
import { ProChat } from '@ant-design/pro-chat'
import { conversationService } from '@/services/conversation'
import { agentService } from '@/services/agent'
import type { Conversation, Agent } from '@/types'

// 可用的 AI 模型列表
const AI_MODELS = [
  { label: 'DeepSeek V3 (推荐)', value: 'deepseek-ai/DeepSeek-V3', provider: 'SiliconFlow', icon: '🚀' },
  { label: 'DeepSeek V3.1', value: 'deepseek-ai/DeepSeek-V3.1-Terminus', provider: 'SiliconFlow', icon: '🚀' },
  { label: 'DeepSeek R1', value: 'deepseek-ai/DeepSeek-R1', provider: 'SiliconFlow', icon: '🧠' },
  { label: 'GPT-4', value: 'gpt-4', provider: 'OpenAI', icon: '🤖' },
  { label: 'GPT-4 Turbo', value: 'gpt-4-turbo-preview', provider: 'OpenAI', icon: '⚡' },
  { label: 'GPT-3.5 Turbo', value: 'gpt-3.5-turbo', provider: 'OpenAI', icon: '💬' },
]

const ConversationDetail = () => {
  const { id } = useParams<{ id: string }>()
  const [conversation, setConversation] = useState<Conversation | null>(null)
  const [agent, setAgent] = useState<Agent | null>(null)
  const [selectedModel, setSelectedModel] = useState<string>('deepseek-ai/DeepSeek-V3')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (id) {
      fetchConversation(id)
    }
  }, [id])

  useEffect(() => {
    if (conversation?.agent_id) {
      fetchAgent(conversation.agent_id)
    }
  }, [conversation?.agent_id])

  // 从 Agent 配置中读取模型
  useEffect(() => {
    if (agent?.model_config?.model) {
      setSelectedModel(agent.model_config.model)
    }
  }, [agent])

  const fetchConversation = async (conversationId: string) => {
    try {
      const response = await conversationService.getConversation(conversationId)
      // 确保 messages 字段存在
      const conversationData = {
        ...response.data,
        messages: response.data.messages || [],
      }
      setConversation(conversationData)
    } catch (error) {
      message.error('获取对话详情失败')
    }
  }

  const fetchAgent = async (agentId: string) => {
    try {
      const response = await agentService.getAgent(agentId)
      setAgent(response.data)
    } catch (error) {
      console.error('获取 Agent 信息失败', error)
    }
  }

  const handleModelChange = async (model: string) => {
    setSelectedModel(model)

    // 更新 Agent 的 model_config
    if (agent) {
      try {
        await agentService.updateAgent(agent.id, {
          model_config: {
            ...agent.model_config,
            model,
          },
        })
        message.success('模型切换成功')
        // 更新本地 agent 状态
        setAgent({
          ...agent,
          model_config: {
            ...agent.model_config,
            model,
          },
        })
      } catch (error) {
        message.error('模型切换失败')
      }
    }
  }

  const handleSendMessage = async (content: string) => {
    if (!id) return

    setLoading(true)
    try {
      const response = await conversationService.sendMessage(id, content)
      if (conversation) {
        const newMessages = response.data.messages || []
        setConversation({
          ...conversation,
          messages: [...(conversation.messages || []), ...newMessages],
        })
      }
    } catch (error: any) {
      console.error('发送消息失败:', error)
      message.error(error.response?.data?.message || '发送消息失败')
      throw error
    } finally {
      setLoading(false)
    }
  }

  if (!conversation) {
    return <div style={{ padding: 24 }}>加载中...</div>
  }

  // 转换消息格式为 ProChat 需要的格式
  const chatMessages = (conversation.messages || []).map((msg) => ({
    id: msg.id,
    content: msg.content,
    role: msg.role as 'user' | 'assistant',
    createAt: new Date(msg.timestamp).getTime(),
    updateAt: new Date(msg.timestamp).getTime(),
  }))

  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      {/* 顶部工具栏 */}
      <div style={{
        padding: '16px 24px',
        borderBottom: '1px solid #f0f0f0',
        background: '#fff'
      }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
          <h2 style={{ margin: 0 }}>{conversation.title || '对话详情'}</h2>
          {agent && agent.knowledge_bases && agent.knowledge_bases.length > 0 && (
            <Space>
              <DatabaseOutlined style={{ color: '#1890ff' }} />
              <span style={{ fontSize: 14, color: '#666' }}>已启用知识库：</span>
              {agent.knowledge_bases.map((kbId) => (
                <Tag key={kbId} color="blue">
                  {kbId}
                </Tag>
              ))}
            </Space>
          )}
        </div>

        {/* 模型选择器 */}
        <Space>
          <ThunderboltOutlined style={{ color: '#52c41a' }} />
          <span style={{ fontSize: 14, color: '#666' }}>AI 模型：</span>
          <Select
            value={selectedModel}
            onChange={handleModelChange}
            style={{ width: 280 }}
            options={AI_MODELS.map((model) => ({
              label: (
                <Space>
                  <span>{model.icon}</span>
                  <span>{model.label}</span>
                  <Tag color={model.provider === 'SiliconFlow' ? 'green' : 'blue'} style={{ fontSize: 10 }}>
                    {model.provider}
                  </Tag>
                </Space>
              ),
              value: model.value,
            }))}
          />
        </Space>
      </div>

      {/* ProChat 对话区域 */}
      <div style={{ flex: 1, overflow: 'hidden' }}>
        <ProChat
          chats={chatMessages}
          onChatsChange={(chats) => {
            // ProChat 内部状态管理
            console.log('Chats changed:', chats)
          }}
          request={async (messages) => {
            const userMessage = messages[messages.length - 1]
            const content = typeof userMessage.content === 'string' ? userMessage.content : ''
            await handleSendMessage(content)

            // 返回 AI 响应（已经在 handleSendMessage 中处理）
            const lastMessage = conversation.messages?.[conversation.messages.length - 1]
            if (lastMessage && lastMessage.role === 'assistant') {
              return new Response(lastMessage.content)
            }
            return new Response('')
          }}
          loading={loading}
          locale="zh-CN"
          placeholder="输入消息..."
          style={{ height: '100%' }}
          assistantMeta={{
            avatar: '🤖',
            title: agent?.name || 'AI Agent',
            backgroundColor: '#f0f0f0',
          }}
          userMeta={{
            avatar: '👤',
            title: '用户',
          }}
          helloMessage={
            chatMessages.length === 0
              ? `你好！我是 ${agent?.name || 'AI Agent'}，有什么可以帮助你的吗？`
              : undefined
          }
        />
      </div>
    </div>
  )
}

export default ConversationDetail
