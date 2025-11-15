import { Card, List, Tag } from 'antd'

const mockTools = [
  {
    id: '1',
    name: 'Web Search',
    description: '搜索网络信息',
    type: 'api',
    category: 'search',
  },
  {
    id: '2',
    name: 'Send Email',
    description: '发送邮件',
    type: 'api',
    category: 'communication',
  },
  {
    id: '3',
    name: 'Database Query',
    description: '数据库查询',
    type: 'function',
    category: 'data',
  },
]

const ToolList = () => {
  return (
    <div>
      <h1 style={{ marginBottom: 16 }}>工具库</h1>
      <List
        grid={{ gutter: 16, column: 3 }}
        dataSource={mockTools}
        renderItem={(tool) => (
          <List.Item>
            <Card title={tool.name}>
              <p>{tool.description}</p>
              <div style={{ marginTop: 12 }}>
                <Tag color="blue">{tool.type}</Tag>
                <Tag>{tool.category}</Tag>
              </div>
            </Card>
          </List.Item>
        )}
      />
    </div>
  )
}

export default ToolList
