import { useState } from 'react'
import { Form, Input, Select, Button, Card, message } from 'antd'
import { useNavigate } from 'react-router-dom'
import { agentService } from '@/services/agent'

const { TextArea } = Input

const AgentCreate = () => {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  const handleSubmit = async (values: any) => {
    setLoading(true)
    try {
      await agentService.createAgent(values)
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
