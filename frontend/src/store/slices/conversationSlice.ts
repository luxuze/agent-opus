import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import type { Conversation } from '@/types'

interface ConversationState {
  conversations: Conversation[]
  currentConversation: Conversation | null
  loading: boolean
}

const initialState: ConversationState = {
  conversations: [],
  currentConversation: null,
  loading: false,
}

const conversationSlice = createSlice({
  name: 'conversation',
  initialState,
  reducers: {
    setConversations: (state, action: PayloadAction<Conversation[]>) => {
      state.conversations = action.payload
    },
    setCurrentConversation: (state, action: PayloadAction<Conversation | null>) => {
      state.currentConversation = action.payload
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload
    },
  },
})

export const { setConversations, setCurrentConversation, setLoading } = conversationSlice.actions
export default conversationSlice.reducer
