import { Routes, Route, Navigate } from 'react-router-dom'
import { Layout } from 'antd'
import MainLayout from './components/Layout/MainLayout'
import Dashboard from './pages/Dashboard'
import AgentList from './pages/Agent/AgentList'
import AgentDetail from './pages/Agent/AgentDetail'
import AgentCreate from './pages/Agent/AgentCreate'
import ConversationList from './pages/Conversation/ConversationList'
import ConversationDetail from './pages/Conversation/ConversationDetail'
import ToolList from './pages/Tool/ToolList'
import KnowledgeBaseList from './pages/KnowledgeBase/KnowledgeBaseList'

function App() {
  return (
    <Routes>
      <Route path="/" element={<MainLayout />}>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="agents" element={<AgentList />} />
        <Route path="agents/create" element={<AgentCreate />} />
        <Route path="agents/:id" element={<AgentDetail />} />
        <Route path="conversations" element={<ConversationList />} />
        <Route path="conversations/:id" element={<ConversationDetail />} />
        <Route path="tools" element={<ToolList />} />
        <Route path="knowledge-bases" element={<KnowledgeBaseList />} />
      </Route>
    </Routes>
  )
}

export default App
