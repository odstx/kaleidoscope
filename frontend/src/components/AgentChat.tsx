import { useState, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { Bot, Maximize2, PanelRight, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/contexts/AuthContext'
import { apiCall } from '@/utils/api'

type Layout = 'floating' | 'fullscreen'

interface Message {
  role: 'user' | 'assistant'
  content: string
}

interface AgentChatProps {
  onLayoutChange?: (layout: Layout, isOpen: boolean) => void
  embedded?: boolean
}

const defaultWelcome = 'Hi! I am your AI assistant. How can I help you today?'

export function AgentChat({ onLayoutChange, embedded = false }: AgentChatProps) {
  const { isAuthenticated } = useAuth()
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)
  const [layout, setLayout] = useState<Layout>('floating')
  const [position, setPosition] = useState({ right: 16, bottom: 64 })
  const [messages, setMessages] = useState<Message[]>([
    { role: 'assistant', content: defaultWelcome },
  ])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const dragData = useRef({ isDragging: false, startX: 0, startY: 0, startRight: 16, startBottom: 64 })
  const messagesEndRef = useRef<HTMLDivElement>(null)

  if (!isAuthenticated) return null

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  const handleOpenChange = (newOpen: boolean) => {
    setOpen(newOpen)
    onLayoutChange?.(layout, newOpen)
  }

  const handleClose = () => {
    if (layout === 'fullscreen') {
      setLayout('floating')
      setOpen(false)
      onLayoutChange?.('floating', false)
    } else if (!embedded) {
      setOpen(false)
      onLayoutChange?.(layout, false)
    } else {
      onLayoutChange?.('floating', false)
    }
  }

  const handleMouseDown = (e: React.MouseEvent) => {
    dragData.current = {
      isDragging: true,
      startX: e.clientX,
      startY: e.clientY,
      startRight: position.right,
      startBottom: position.bottom,
    }
    e.preventDefault()
  }

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!dragData.current.isDragging) return
    const dx = dragData.current.startX - e.clientX
    const dy = dragData.current.startY - e.clientY
    setPosition({
      right: dragData.current.startRight + dx,
      bottom: dragData.current.startBottom + dy,
    })
  }

  const handleMouseUp = () => {
    dragData.current.isDragging = false
  }

  const toggleLayout = () => {
    const layouts: Layout[] = ['floating', 'fullscreen']
    const currentIndex = layouts.indexOf(layout)
    const nextIndex = (currentIndex + 1) % layouts.length
    const newLayout = layouts[nextIndex]
    setLayout(newLayout)
    onLayoutChange?.(newLayout, open)
    if (newLayout === 'fullscreen') {
      setOpen(true)
    }
  }

  const getNextLayout = (): Layout => {
    const layouts: Layout[] = ['floating', 'fullscreen']
    const currentIndex = layouts.indexOf(layout)
    const nextIndex = (currentIndex + 1) % layouts.length
    return layouts[nextIndex]
  }

  const getLayoutIcon = () => {
    const nextLayout = getNextLayout()
    switch (nextLayout) {
      case 'floating':
        return <PanelRight className="h-4 w-4" />
      case 'fullscreen':
        return <Maximize2 className="h-4 w-4" />
    }
  }

  const getContainerStyle = (): React.CSSProperties => {
    if (layout === 'floating') {
      return {
        right: `${position.right}px`,
        bottom: `${position.bottom}px`,
        width: '400px',
        height: '500px',
      }
    }
    if (layout === 'fullscreen') {
      return {}
    }
    return {}
  }

  const handleSend = async () => {
    if (!input.trim() || loading) return
    const userMessage = input.trim()
    setInput('')
    setMessages((prev) => [...prev, { role: 'user', content: userMessage }])
    setLoading(true)

    try {
      const resp = await apiCall<{ message: string }>('/api/v1/agent/chat', {
        method: 'POST',
        body: JSON.stringify({ message: userMessage }),
      })
      setMessages((prev) => [...prev, { role: 'assistant', content: resp.message }])
    } catch (err) {
      setMessages((prev) => [
        ...prev,
        { role: 'assistant', content: t('agent.error') },
      ])
    } finally {
      setLoading(false)
      setTimeout(scrollToBottom, 100)
    }
  }

  const renderToggleButton = () => {
    if (open) return null
    return (
      <div
        className="fixed z-50"
        style={{ right: `${position.right}px`, bottom: `${position.bottom}px` }}
      >
        <div
          className="cursor-move"
          onMouseDown={handleMouseDown}
          onMouseMove={handleMouseMove}
          onMouseUp={handleMouseUp}
          onMouseLeave={handleMouseUp}
        >
          <Button
            onClick={() => handleOpenChange(!open)}
            className="h-12 w-12 rounded-full shadow-lg"
            size="icon"
          >
            <Bot className="h-6 w-6" />
          </Button>
        </div>
      </div>
    )
  }

  const renderChat = () => {
    if (!open && !embedded) return null
    return (
      <div
        className={`flex flex-col ${embedded || layout === 'fullscreen' ? 'h-full' : ''} ${layout === 'floating' && !embedded ? 'fixed z-50 bg-background border border-input rounded-lg shadow-xl' : ''} ${layout === 'fullscreen' ? 'fixed inset-0 z-50 bg-background' : ''}`}
        style={!embedded ? getContainerStyle() : {}}
      >
        <div
          className={`h-14 px-4 border-b border-input flex items-center justify-between select-none ${layout === 'floating' ? 'cursor-move' : ''}`}
          onMouseDown={layout === 'floating' ? handleMouseDown : undefined}
          onMouseMove={layout === 'floating' ? handleMouseMove : undefined}
          onMouseUp={layout === 'floating' ? handleMouseUp : undefined}
          onMouseLeave={layout === 'floating' ? handleMouseUp : undefined}
        >
          <div className="flex items-center gap-2">
            <Bot className="h-5 w-5" />
            <span className="font-semibold">{t('agent.title')}</span>
          </div>
          <div className="flex items-center gap-1">
            {!embedded && (
              <Button variant="ghost" size="sm" onClick={toggleLayout} title={getNextLayout()}>
                {getLayoutIcon()}
              </Button>
            )}
            <Button variant="ghost" size="sm" onClick={handleClose}>
              <X className="h-4 w-4" />
            </Button>
          </div>
        </div>
        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {messages.map((msg, i) => (
            <div
              key={i}
              className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-[80%] rounded-lg p-3 ${
                  msg.role === 'user'
                    ? 'bg-primary text-primary-foreground'
                    : 'bg-muted'
                }`}
              >
                {msg.content}
              </div>
            </div>
          ))}
          {loading && (
            <div className="flex justify-start">
              <div className="max-w-[80%] rounded-lg p-3 bg-muted">
                {t('agent.thinking')}
              </div>
            </div>
          )}
          <div ref={messagesEndRef} />
        </div>
        <div className="flex gap-2 p-4 border-t border-input">
          <input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSend()}
            placeholder={t('agent.placeholder')}
            className="flex-1 px-3 py-2 rounded-md border bg-background"
            disabled={loading}
          />
          <Button onClick={handleSend} disabled={loading || !input.trim()}>
            {t('agent.send')}
          </Button>
        </div>
      </div>
    )
  }

  return (
    <>
      {renderToggleButton()}
      {renderChat()}
    </>
  )
}

export function useAgentChat() {
  const [agentChatOpen, setAgentChatOpen] = useState(false)
  const [agentChatLayout, setAgentChatLayout] = useState<Layout>('floating')

  const handleLayoutChange = (layout: Layout, isOpen: boolean) => {
    setAgentChatLayout(layout)
    setAgentChatOpen(isOpen)
  }

  return {
    agentChatOpen,
    agentChatLayout,
    handleLayoutChange,
  }
}