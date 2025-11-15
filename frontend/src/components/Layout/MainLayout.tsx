import { useState } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, theme } from 'antd'
import {
  DashboardOutlined,
  RobotOutlined,
  MessageOutlined,
  ToolOutlined,
  DatabaseOutlined,
} from '@ant-design/icons'

const { Header, Content, Sider } = Layout

const MainLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [collapsed, setCollapsed] = useState(false)
  const {
    token: { colorBgContainer },
  } = theme.useToken()

  const menuItems = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: 'Dashboard',
    },
    {
      key: '/agents',
      icon: <RobotOutlined />,
      label: 'Agent 管理',
    },
    {
      key: '/conversations',
      icon: <MessageOutlined />,
      label: '对话管理',
    },
    {
      key: '/tools',
      icon: <ToolOutlined />,
      label: '工具库',
    },
    {
      key: '/knowledge-bases',
      icon: <DatabaseOutlined />,
      label: '知识库',
    },
  ]

  const selectedKey = '/' + location.pathname.split('/')[1]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}>
        <div style={{
          height: 32,
          margin: 16,
          color: '#fff',
          fontSize: 18,
          fontWeight: 'bold',
          textAlign: 'center'
        }}>
          {collapsed ? 'AP' : 'Agent Platform'}
        </div>
        <Menu
          theme="dark"
          selectedKeys={[selectedKey]}
          mode="inline"
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ padding: 0, background: colorBgContainer }}>
          <div style={{
            paddingLeft: 24,
            fontSize: 18,
            fontWeight: 'bold'
          }}>
            Agent 平台管理系统
          </div>
        </Header>
        <Content style={{ margin: '16px' }}>
          <div style={{
            padding: 24,
            minHeight: 360,
            background: colorBgContainer,
            borderRadius: 8,
          }}>
            <Outlet />
          </div>
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
