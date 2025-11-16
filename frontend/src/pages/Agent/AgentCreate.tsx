import { useState } from 'react'
import { Form, Input, Select, Button, Card, message, Space, Tag } from 'antd'
import { useNavigate } from 'react-router-dom'
import { agentService } from '@/services/agent'

const { TextArea } = Input

// 可用的 AI 模型列表
const AI_MODELS = [
  { label: 'DeepSeek V3 (推荐)', value: 'deepseek-ai/DeepSeek-V3', provider: 'SiliconFlow', icon: '🚀' },
  { label: 'DeepSeek V3.1', value: 'deepseek-ai/DeepSeek-V3.1-Terminus', provider: 'SiliconFlow', icon: '🚀' },
  { label: 'DeepSeek R1', value: 'deepseek-ai/DeepSeek-R1', provider: 'SiliconFlow', icon: '🧠' },
  { label: 'GPT-4', value: 'gpt-4', provider: 'OpenAI', icon: '🤖' },
  { label: 'GPT-4 Turbo', value: 'gpt-4-turbo-preview', provider: 'OpenAI', icon: '⚡' },
  { label: 'GPT-3.5 Turbo', value: 'gpt-3.5-turbo', provider: 'OpenAI', icon: '💬' },
]

const AgentCreate = () => {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  const handleSubmit = async (values: any) => {
    setLoading(true)
    try {
      // 将 model 字段转换为 model_config
      const payload = {
        ...values,
        model_config: {
          model: values.model || 'deepseek-ai/DeepSeek-V3',
          temperature: 0.7,
        },
      }
      delete payload.model

      await agentService.createAgent(payload)
      message.success('创建成功')
      navigate('/agents')
    } catch (error) {
      message.error('创建失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <h1 style={{ marginBottom: 24 }}>创建 Agent</h1>
      <Card>
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{
            type: 'single',
            status: 'draft',
            model: 'deepseek-ai/DeepSeek-V3',
          }}
        >
          <Form.Item
            label="Agent 名称"
            name="name"
            rules={[{ required: true, message: '请输入 Agent 名称' }]}
          >
            <Input placeholder="请输入 Agent 名称" />
          </Form.Item>

          <Form.Item
            label="描述"
            name="description"
          >
            <TextArea rows={4} placeholder="请输入描述" />
          </Form.Item>

          <Form.Item
            label="AI 模型"
            name="model"
            rules={[{ required: true, message: '请选择 AI 模型' }]}
          >
            <Select placeholder="选择 AI 模型">
              {AI_MODELS.map((model) => (
                <Select.Option key={model.value} value={model.value}>
                  <Space>
                    <span>{model.icon}</span>
                    <span>{model.label}</span>
                    <Tag color={model.provider === 'SiliconFlow' ? 'green' : 'blue'} style={{ fontSize: 10 }}>
                      {model.provider}
                    </Tag>
                  </Space>
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            label="类型"
            name="type"
          >
            <Select>
              <Select.Option value="single">单 Agent</Select.Option>
              <Select.Option value="multi">多 Agent</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="提示词模板"
            name="prompt_template"
          >
            <TextArea rows={6} placeholder="请输入提示词模板" />
          </Form.Item>

          <Form.Item
            label="状态"
            name="status"
          >
            <Select>
              <Select.Option value="draft">草稿</Select.Option>
              <Select.Option value="published">已发布</Select.Option>
              <Select.Option value="archived">已归档</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              创建
            </Button>
            <Button style={{ marginLeft: 8 }} onClick={() => navigate('/agents')}>
              取消
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default AgentCreate
