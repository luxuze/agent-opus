import { configureStore } from '@reduxjs/toolkit'
import agentReducer from './slices/agentSlice'
import conversationReducer from './slices/conversationSlice'

export const store = configureStore({
  reducer: {
    agent: agentReducer,
    conversation: conversationReducer,
  },
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
