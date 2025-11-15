import { Row, Col, Card, Statistic } from 'antd'
import {
  RobotOutlined,
  MessageOutlined,
  ToolOutlined,
  DatabaseOutlined,
} from '@ant-design/icons'

const Dashboard = () => {
  return (
    <div>
      <h1 style={{ marginBottom: 24 }}>Dashboard</h1>
      <Row gutter={16}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总 Agent 数"
              value={12}
              prefix={<RobotOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="对话数"
              value={156}
              prefix={<MessageOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="工具数"
              value={28}
              prefix={<ToolOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="知识库"
              value={5}
              prefix={<DatabaseOutlined />}
              valueStyle={{ color: '#eb2f96' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginTop: 24 }}>
        <Col span={12}>
          <Card title="最近活跃 Agent" bordered={false}>
            <p>Customer Service Agent - 活跃</p>
            <p>Data Analysis Agent - 活跃</p>
            <p>Content Generator - 空闲</p>
          </Card>
        </Col>
        <Col span={12}>
          <Card title="系统状态" bordered={false}>
            <p>API 响应时间: 45ms</p>
            <p>Token 使用: 15,234 / 100,000</p>
            <p>系统可用性: 99.9%</p>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
