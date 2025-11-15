import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import type { Agent } from '@/types'

interface AgentState {
  agents: Agent[]
  currentAgent: Agent | null
  loading: boolean
}

const initialState: AgentState = {
  agents: [],
  currentAgent: null,
  loading: false,
}

const agentSlice = createSlice({
  name: 'agent',
  initialState,
  reducers: {
    setAgents: (state, action: PayloadAction<Agent[]>) => {
      state.agents = action.payload
    },
    setCurrentAgent: (state, action: PayloadAction<Agent | null>) => {
      state.currentAgent = action.payload
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload
    },
  },
})

export const { setAgents, setCurrentAgent, setLoading } = agentSlice.actions
export default agentSlice.reducer
