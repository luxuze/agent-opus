import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Button, Spin, message } from 'antd'
import { agentService } from '@/services/agent'
import type { Agent } from '@/types'

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
        </Descriptions>
      </Card>
    </div>
  )
}

export default AgentDetail
