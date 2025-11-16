import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Button, Spin, message, Tag, Space } from 'antd'
import { agentService } from '@/services/agent'
import type { Agent } from '@/types'

// 模型信息映射
const MODEL_INFO: Record<string, { label: string; provider: string; icon: string }> = {
  'deepseek-ai/DeepSeek-V3': { label: 'DeepSeek V3', provider: 'SiliconFlow', icon: '🚀' },
  'deepseek-ai/DeepSeek-V3.1-Terminus': { label: 'DeepSeek V3.1', provider: 'SiliconFlow', icon: '🚀' },
  'deepseek-ai/DeepSeek-R1': { label: 'DeepSeek R1', provider: 'SiliconFlow', icon: '🧠' },
  'gpt-4': { label: 'GPT-4', provider: 'OpenAI', icon: '🤖' },
  'gpt-4-turbo-preview': { label: 'GPT-4 Turbo', provider: 'OpenAI', icon: '⚡' },
  'gpt-3.5-turbo': { label: 'GPT-3.5 Turbo', provider: 'OpenAI', icon: '💬' },
}

const AgentDetail = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [agent, setAgent] = useState<Agent | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (id) {
      fetchAgent(id)
    }
  }, [id])

  const fetchAgent = async (agentId: string) => {
    setLoading(true)
    try {
      const response = await agentService.getAgent(agentId)
      setAgent(response.data)
    } catch (error) {
      message.error('获取 Agent 详情失败')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <Spin size="large" />
  }

  if (!agent) {
    return <div>Agent 不存在</div>
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h1>Agent 详情</h1>
        <Button onClick={() => navigate('/agents')}>返回列表</Button>
      </div>
      <Card>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="ID">{agent.id}</Descriptions.Item>
          <Descriptions.Item label="名称">{agent.name}</Descriptions.Item>
          <Descriptions.Item label="类型">{agent.type}</Descriptions.Item>
          <Descriptions.Item label="状态">{agent.status}</Descriptions.Item>
          <Descriptions.Item label="AI 模型">
            {agent.model_config?.model ? (
              <Space>
                <span>{MODEL_INFO[agent.model_config.model]?.icon || '🤖'}</span>
                <span>{MODEL_INFO[agent.model_config.model]?.label || agent.model_config.model}</span>
                <Tag color={MODEL_INFO[agent.model_config.model]?.provider === 'SiliconFlow' ? 'green' : 'blue'}>
                  {MODEL_INFO[agent.model_config.model]?.provider || 'Unknown'}
                </Tag>
              </Space>
            ) : (
              <Tag color="orange">未配置</Tag>
            )}
          </Descriptions.Item>
          <Descriptions.Item label="Temperature">
            {agent.model_config?.temperature || 0.7}
          </Descriptions.Item>
          <Descriptions.Item label="版本">{agent.version}</Descriptions.Item>
          <Descriptions.Item label="创建者">{agent.created_by}</Descriptions.Item>
          <Descriptions.Item label="描述" span={2}>
            {agent.description || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="提示词模板" span={2}>
            <pre style={{ whiteSpace: 'pre-wrap' }}>
              {agent.prompt_template || '-'}
            </pre>
          </Descriptions.Item>
          <Descriptions.Item label="工具" span={2}>
            {agent.tools && agent.tools.length > 0 ? (
              agent.tools.map((tool) => (
                <Tag key={tool} color="green">
                  {tool}
                </Tag>
              ))
            ) : (
              '-'
            )}
          </Descriptions.Item>
          <Descriptions.Item label="知识库" span={2}>
            {agent.knowledge_bases && agent.knowledge_bases.length > 0 ? (
              agent.knowledge_bases.map((kb) => (
                <Tag key={kb} color="blue">
                  {kb}
                </Tag>
              ))
            ) : (
              '-'
            )}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}

export default AgentDetail
