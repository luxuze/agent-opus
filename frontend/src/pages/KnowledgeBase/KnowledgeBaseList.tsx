import { Card, List, Statistic, Row, Col } from 'antd'
import { FileTextOutlined } from '@ant-design/icons'

const mockKnowledgeBases = [
  {
    id: '1',
    name: '产品文档',
    description: '产品相关文档知识库',
    document_count: 150,
    vector_count: 5000,
    type: 'document',
  },
  {
    id: '2',
    name: '客户FAQ',
    description: '常见问题知识库',
    document_count: 80,
    vector_count: 2500,
    type: 'document',
  },
]

const KnowledgeBaseList = () => {
  return (
    <div>
      <h1 style={{ marginBottom: 16 }}>知识库</h1>
      <List
        grid={{ gutter: 16, column: 2 }}
        dataSource={mockKnowledgeBases}
        renderItem={(kb) => (
          <List.Item>
            <Card title={kb.name}>
              <p>{kb.description}</p>
              <Row gutter={16} style={{ marginTop: 16 }}>
                <Col span={12}>
                  <Statistic
                    title="文档数"
                    value={kb.document_count}
                    prefix={<FileTextOutlined />}
                  />
                </Col>
                <Col span={12}>
                  <Statistic
                    title="向量数"
                    value={kb.vector_count}
                  />
                </Col>
              </Row>
            </Card>
          </List.Item>
        )}
      />
    </div>
  )
}

export default KnowledgeBaseList
