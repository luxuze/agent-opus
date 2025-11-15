import { useState, useEffect } from 'react'
import { Table, Button, Space, Tag, message } from 'antd'
import { MessageOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { conversationService } from '@/services/conversation'
import type { Conversation } from '@/types'

const ConversationList = () => {
  const navigate = useNavigate()
  const [conversations, setConversations] = useState<Conversation[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    fetchConversations()
  }, [])

  const fetchConversations = async () => {
    setLoading(true)
    try {
      const response = await conversationService.getConversations()
      setConversations(response.data.items || [])
    } catch (error) {
      message.error('获取对话列表失败')
    } finally {
      setLoading(false)
    }
  }

  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
    },
    {
      title: 'Agent ID',
      dataIndex: 'agent_id',
      key: 'agent_id',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'success' : 'default'}>
          {status}
        </Tag>
      ),
    },
    {
      title: '最后消息时间',
      dataIndex: 'last_message_at',
      key: 'last_message_at',
      render: (time: string) => time ? new Date(time).toLocaleString() : '-',
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Conversation) => (
        <Button
          type="link"
          icon={<MessageOutlined />}
          onClick={() => navigate(`/conversations/${record.id}`)}
        >
          查看
        </Button>
      ),
    },
  ]

  return (
    <div>
      <h1 style={{ marginBottom: 16 }}>对话管理</h1>
      <Table
        columns={columns}
        dataSource={conversations}
        rowKey="id"
        loading={loading}
      />
    </div>
  )
}

export default ConversationList
