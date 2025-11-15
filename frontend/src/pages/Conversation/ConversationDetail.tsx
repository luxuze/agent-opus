import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { Card, List, Input, Button, message, Avatar } from 'antd'
import { UserOutlined, RobotOutlined, SendOutlined } from '@ant-design/icons'
import { conversationService } from '@/services/conversation'
import type { Conversation, Message } from '@/types'

const ConversationDetail = () => {
  const { id } = useParams<{ id: string }>()
  const [conversation, setConversation] = useState<Conversation | null>(null)
  const [inputMessage, setInputMessage] = useState('')
  const [sending, setSending] = useState(false)

  useEffect(() => {
    if (id) {
      fetchConversation(id)
    }
  }, [id])

  const fetchConversation = async (conversationId: string) => {
    try {
      const response = await conversationService.getConversation(conversationId)
      setConversation(response.data)
    } catch (error) {
      message.error('获取对话详情失败')
    }
  }

  const handleSend = async () => {
    if (!inputMessage.trim() || !id) return

    setSending(true)
    try {
      const response = await conversationService.sendMessage(id, inputMessage)
      if (conversation) {
        setConversation({
          ...conversation,
          messages: [...conversation.messages, ...response.data.messages],
        })
      }
      setInputMessage('')
    } catch (error) {
      message.error('发送消息失败')
    } finally {
      setSending(false)
    }
  }

  if (!conversation) {
    return <div>加载中...</div>
  }

  return (
    <div>
      <h1 style={{ marginBottom: 16 }}>{conversation.title || '对话详情'}</h1>
      <Card style={{ height: 'calc(100vh - 250px)', display: 'flex', flexDirection: 'column' }}>
        <div style={{ flex: 1, overflow: 'auto', marginBottom: 16 }}>
          <List
            dataSource={conversation.messages}
            renderItem={(msg: Message) => (
              <List.Item style={{ border: 'none' }}>
                <List.Item.Meta
                  avatar={
                    <Avatar icon={msg.role === 'user' ? <UserOutlined /> : <RobotOutlined />} />
                  }
                  title={msg.role === 'user' ? '用户' : 'Agent'}
                  description={msg.content}
                />
              </List.Item>
            )}
          />
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          <Input
            placeholder="输入消息..."
            value={inputMessage}
            onChange={(e) => setInputMessage(e.target.value)}
            onPressEnter={handleSend}
          />
          <Button
            type="primary"
            icon={<SendOutlined />}
            loading={sending}
            onClick={handleSend}
          >
            发送
          </Button>
        </div>
      </Card>
    </div>
  )
}

export default ConversationDetail
