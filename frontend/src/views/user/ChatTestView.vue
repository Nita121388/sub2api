<template>
  <AppLayout>
    <div class="space-y-6">
      <div class="card p-6">
        <div class="mb-5 flex items-center justify-between gap-3">
          <div>
            <h1 class="text-xl font-semibold text-gray-900 dark:text-white">对话</h1>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              Responses first, Chat Completions compatible. Built for fast local verification.
            </p>
          </div>
          <button class="btn btn-secondary" @click="showSettings = !showSettings">
            {{ showSettings ? '收起设置' : '展开设置' }}
          </button>
        </div>

        <div v-if="!showSettings" class="rounded-xl border border-dashed border-gray-200 bg-gray-50/70 px-4 py-3 text-sm text-gray-600 dark:border-dark-700 dark:bg-dark-900/30 dark:text-gray-300">
          {{ keyMode === 'manual' ? '手动 Key' : '已选 Key' }} · {{ protocol }} · {{ model || DEFAULT_MODEL }} ·
          {{ stream ? '流式输出' : '非流式输出' }}
        </div>

        <div v-else class="grid grid-cols-1 gap-4 lg:grid-cols-2">
          <div>
            <label class="input-label">Key Source</label>
            <div class="mt-1 flex gap-3">
              <label class="inline-flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                <input v-model="keyMode" type="radio" value="select" />
                Select Existing Key
              </label>
              <label class="inline-flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                <input v-model="keyMode" type="radio" value="manual" />
                Paste Key
              </label>
            </div>
          </div>

          <div>
            <label class="input-label">Protocol</label>
            <select v-model="protocol" class="input mt-1">
              <option value="responses">Responses (default)</option>
              <option value="chat-completions">Chat Completions</option>
            </select>
          </div>

          <div v-if="keyMode === 'select'">
            <label class="input-label">API Key</label>
            <select v-model.number="selectedKeyId" class="input mt-1">
              <option :value="0" disabled>Select one key...</option>
              <option v-for="key in apiKeys" :key="key.id" :value="key.id">
                {{ key.name }} · {{ key.group?.platform || 'default' }} · {{ maskKey(key.key) }}
              </option>
            </select>
            <p v-if="keyLoading" class="mt-1 text-xs text-gray-500 dark:text-gray-400">Loading keys...</p>
          </div>

          <div v-else>
            <label class="input-label">Manual API Key</label>
            <input v-model.trim="manualKey" type="password" class="input mt-1" placeholder="sk-..." />
          </div>

          <div>
            <label class="input-label">Model</label>
            <input v-model.trim="model" type="text" class="input mt-1" placeholder="gpt-4.1-mini" />
          </div>

          <div class="flex items-end">
            <label class="inline-flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
              <input v-model="stream" type="checkbox" />
              Stream output
            </label>
          </div>

          <div>
            <label class="input-label">Base URL</label>
            <input v-model.trim="baseUrl" type="text" class="input mt-1" placeholder="http://127.0.0.1:15173" />
          </div>

          <div>
            <label class="input-label">Endpoint Path</label>
            <div class="mt-1 flex gap-2">
              <input
                :value="endpointPath"
                type="text"
                class="input"
                placeholder="/responses"
                @input="onEndpointInput"
              />
              <button class="btn btn-secondary shrink-0" @click="resetEndpointPath">Reset</button>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-6 xl:grid-cols-3">
        <div class="card p-6 xl:col-span-2">
          <div class="mb-3 flex items-center justify-between">
            <h2 class="text-base font-semibold text-gray-900 dark:text-white">Conversation</h2>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ requestUrl }}</p>
          </div>

          <div class="mb-4 h-[440px] overflow-y-auto rounded-xl border border-gray-200 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-900/40">
            <div v-if="messages.length === 0" class="text-sm text-gray-500 dark:text-gray-400">
              Send a message to start testing.
            </div>
            <div v-for="msg in messages" :key="msg.id" class="mb-3 last:mb-0">
              <div class="mb-1 text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">
                {{ msg.role }}
              </div>
              <div
                class="whitespace-pre-wrap rounded-xl px-3 py-2 text-sm"
                :class="msg.role === 'user'
                  ? 'bg-blue-100 text-blue-900 dark:bg-blue-900/40 dark:text-blue-100'
                  : msg.error
                    ? 'bg-rose-100 text-rose-900 dark:bg-rose-900/40 dark:text-rose-100'
                    : 'bg-white text-gray-900 dark:bg-dark-800 dark:text-gray-100'"
              >
                {{ msg.content || (msg.pending ? 'Streaming...' : '-') }}
              </div>
            </div>
          </div>

          <div class="space-y-3">
            <textarea
              v-model="userInput"
              rows="4"
              class="input min-h-[110px]"
              placeholder="Type your message..."
            />
            <div class="flex flex-wrap items-center gap-2">
              <button class="btn btn-primary" :disabled="sending" @click="sendMessage">Send</button>
              <button class="btn btn-secondary" :disabled="!sending" @click="stopRequest">Stop</button>
              <button class="btn btn-secondary" :disabled="sending" @click="clearConversation">Clear</button>
            </div>
          </div>
        </div>

        <div class="space-y-6">
          <div class="card p-6">
            <div class="mb-3 flex items-center justify-between gap-3">
              <div>
                <h2 class="text-base font-semibold text-gray-900 dark:text-white">Session History</h2>
                <p class="text-xs text-gray-500 dark:text-gray-400">
                  Local only. Reopen or delete previous Chat Test sessions.
                </p>
              </div>
              <button class="btn btn-secondary px-3 py-1 text-xs" :disabled="sending" @click="startNewSession">
                New Session
              </button>
            </div>

            <div v-if="sessionHistory.length === 0" class="rounded-lg bg-gray-50 px-3 py-4 text-sm text-gray-500 dark:bg-dark-900/40 dark:text-gray-400">
              No saved sessions yet.
            </div>

            <div v-else class="space-y-2">
              <div
                v-for="session in sessionHistory"
                :key="session.id"
                class="rounded-xl border p-3 transition"
                :class="session.id === activeSessionId
                  ? 'border-blue-300 bg-blue-50/80 dark:border-blue-700 dark:bg-blue-950/20'
                  : 'border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900/30'"
              >
                <div class="flex items-start justify-between gap-3">
                  <button class="min-w-0 flex-1 text-left" @click="loadSession(session.id)">
                    <div class="truncate text-sm font-semibold text-gray-900 dark:text-white">
                      {{ session.title }}
                    </div>
                    <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                      {{ formatSessionTime(session.updatedAt) }}
                    </div>
                    <div class="mt-2 text-xs text-gray-500 dark:text-gray-400">
                      {{ session.protocol }} · {{ session.model }} · {{ session.messages.length }} messages
                    </div>
                  </button>
                  <button
                    class="rounded-lg px-2 py-1 text-xs font-medium text-rose-600 transition hover:bg-rose-50 dark:text-rose-300 dark:hover:bg-rose-950/30"
                    :disabled="sending"
                    @click.stop="deleteSession(session.id)"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div class="card p-6">
            <h2 class="mb-3 text-base font-semibold text-gray-900 dark:text-white">Debug</h2>
            <div class="space-y-3 text-xs">
              <div>
                <p class="mb-1 font-semibold text-gray-700 dark:text-gray-300">Request URL</p>
                <pre class="rounded bg-gray-100 p-2 text-gray-800 dark:bg-dark-800 dark:text-gray-200">{{ debug.requestUrl || '-' }}</pre>
              </div>
              <div>
                <p class="mb-1 font-semibold text-gray-700 dark:text-gray-300">Headers</p>
                <pre class="rounded bg-gray-100 p-2 text-gray-800 dark:bg-dark-800 dark:text-gray-200">{{ debug.headersSummary || '-' }}</pre>
              </div>
              <div>
                <p class="mb-1 font-semibold text-gray-700 dark:text-gray-300">Request Body</p>
                <pre class="max-h-[180px] overflow-auto rounded bg-gray-100 p-2 text-gray-800 dark:bg-dark-800 dark:text-gray-200">{{ debug.requestBody || '-' }}</pre>
              </div>
              <div>
                <p class="mb-1 font-semibold text-gray-700 dark:text-gray-300">Status / Error</p>
                <pre class="rounded bg-gray-100 p-2 text-gray-800 dark:bg-dark-800 dark:text-gray-200">{{ debug.statusText || '-' }}</pre>
              </div>
              <div>
                <p class="mb-1 font-semibold text-gray-700 dark:text-gray-300">Raw Response Excerpt</p>
                <pre class="max-h-[220px] overflow-auto rounded bg-gray-100 p-2 text-gray-800 dark:bg-dark-800 dark:text-gray-200">{{ debug.rawResponse || '-' }}</pre>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { keysAPI } from '@/api'
import type { ApiKey } from '@/types'
import { useAppStore } from '@/stores/app'
import AppLayout from '@/components/layout/AppLayout.vue'

type ProtocolType = 'responses' | 'chat-completions'
type KeyMode = 'select' | 'manual'
type MessageRole = 'user' | 'assistant'

interface ChatMessage {
  id: number
  role: MessageRole
  content: string
  pending?: boolean
  error?: boolean
}

interface PersistedChatTestSettings {
  keyMode: KeyMode
  selectedKeyId: number
  manualKey: string
  protocol: ProtocolType
  model: string
  stream: boolean
  baseUrl: string
  endpointPath: string
}

interface StreamChunkExtraction {
  delta?: string
  finalText?: string
}

interface ChatDebugSnapshot {
  requestUrl: string
  headersSummary: string
  requestBody: string
  statusText: string
  rawResponse: string
}

interface ChatSessionSnapshot {
  id: string
  title: string
  createdAt: string
  updatedAt: string
  protocol: ProtocolType
  model: string
  stream: boolean
  baseUrl: string
  endpointPath: string
  messages: ChatMessage[]
  debug: ChatDebugSnapshot
}

const STORAGE_KEY = 'sub2api.chat-test.settings.v1'
const HISTORY_STORAGE_KEY = 'sub2api.chat-test.sessions.v1'
const DEFAULT_MODEL = 'gpt-4.1-mini'

const appStore = useAppStore()

const apiKeys = ref<ApiKey[]>([])
const keyLoading = ref(false)
const keyMode = ref<KeyMode>('select')
const selectedKeyId = ref<number>(0)
const manualKey = ref('')

const protocol = ref<ProtocolType>('responses')
const model = ref(DEFAULT_MODEL)
const stream = ref(true)
const baseUrl = ref('')
const endpointPath = ref('/responses')
const endpointEdited = ref(false)
const showSettings = ref(false)

const sending = ref(false)
const userInput = ref('')
const messages = ref<ChatMessage[]>([])
const sessionHistory = ref<ChatSessionSnapshot[]>([])
const activeSessionId = ref<string | null>(null)
const requestAbortController = ref<AbortController | null>(null)
let messageId = 0

const debug = reactive<ChatDebugSnapshot>({
  requestUrl: '',
  headersSummary: '',
  requestBody: '',
  statusText: '',
  rawResponse: ''
})

const settingsHydrated = ref(false)

const requestUrl = computed(() => {
  const base = normalizeBaseUrl(baseUrl.value)
  const path = normalizePath(endpointPath.value)
  return `${base}${path}`
})

const resolvedApiKey = computed(() => {
  if (keyMode.value === 'manual') {
    return manualKey.value.trim()
  }
  return apiKeys.value.find((k) => k.id === selectedKeyId.value)?.key || ''
})

watch(protocol, () => {
  if (!endpointEdited.value) {
    endpointPath.value = getDefaultPath(protocol.value)
  }
})

watch(
  [keyMode, selectedKeyId, manualKey, protocol, model, stream, baseUrl, endpointPath],
  () => {
    if (!settingsHydrated.value) {
      return
    }
    persistSettings()
    persistActiveSession()
  }
)

onMounted(async () => {
  await appStore.fetchPublicSettings()
  baseUrl.value = normalizeBaseUrl(appStore.apiBaseUrl || window.location.origin)
  hydrateSettings()
  await loadApiKeys()
  hydrateSessionHistory()
  settingsHydrated.value = true
})

async function loadApiKeys(): Promise<void> {
  keyLoading.value = true
  try {
    const result = await keysAPI.list(1, 100)
    apiKeys.value = result.items || []
    if (!selectedKeyId.value && apiKeys.value.length > 0) {
      selectedKeyId.value = apiKeys.value[0].id
    }
  } catch (error) {
    appStore.showError((error as { message?: string }).message || 'Failed to load API keys')
  } finally {
    keyLoading.value = false
  }
}

function hydrateSettings(): void {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) {
      endpointPath.value = getDefaultPath(protocol.value)
      return
    }

    const saved = JSON.parse(raw) as Partial<PersistedChatTestSettings>
    if (saved.keyMode === 'manual' || saved.keyMode === 'select') {
      keyMode.value = saved.keyMode
    }
    if (typeof saved.selectedKeyId === 'number' && saved.selectedKeyId >= 0) {
      selectedKeyId.value = saved.selectedKeyId
    }
    if (typeof saved.manualKey === 'string') {
      manualKey.value = saved.manualKey
    }
    if (saved.protocol === 'responses' || saved.protocol === 'chat-completions') {
      protocol.value = saved.protocol
    }
    if (typeof saved.model === 'string' && saved.model.trim()) {
      model.value = saved.model.trim()
    }
    if (typeof saved.stream === 'boolean') {
      stream.value = saved.stream
    }
    if (typeof saved.baseUrl === 'string' && saved.baseUrl.trim()) {
      baseUrl.value = normalizeBaseUrl(saved.baseUrl)
    }
    if (typeof saved.endpointPath === 'string' && saved.endpointPath.trim()) {
      endpointPath.value = normalizePath(saved.endpointPath)
      endpointEdited.value = saved.endpointPath !== getDefaultPath(protocol.value)
    } else {
      endpointPath.value = getDefaultPath(protocol.value)
    }
  } catch {
    endpointPath.value = getDefaultPath(protocol.value)
  }
}

function persistSettings(): void {
  const payload: PersistedChatTestSettings = {
    keyMode: keyMode.value,
    selectedKeyId: selectedKeyId.value,
    manualKey: manualKey.value,
    protocol: protocol.value,
    model: model.value.trim() || DEFAULT_MODEL,
    stream: stream.value,
    baseUrl: normalizeBaseUrl(baseUrl.value),
    endpointPath: normalizePath(endpointPath.value)
  }
  localStorage.setItem(STORAGE_KEY, JSON.stringify(payload))
}

function getDefaultPath(value: ProtocolType): string {
  return value === 'responses' ? '/responses' : '/chat/completions'
}

function normalizeBaseUrl(value: string): string {
  const trimmed = value.trim() || window.location.origin
  return trimmed.replace(/\/+$/, '')
}

function normalizePath(value: string): string {
  const trimmed = value.trim()
  if (!trimmed) {
    return getDefaultPath(protocol.value)
  }
  return trimmed.startsWith('/') ? trimmed : `/${trimmed}`
}

function maskKey(key: string): string {
  if (!key) return ''
  if (key.length <= 8) return key
  return `${key.slice(0, 4)}...${key.slice(-4)}`
}

function onEndpointInput(event: Event): void {
  const target = event.target as HTMLInputElement
  endpointPath.value = target.value
  endpointEdited.value = true
}

function resetEndpointPath(): void {
  endpointPath.value = getDefaultPath(protocol.value)
  endpointEdited.value = false
}

function clearConversation(): void {
  messages.value = []
  resetDebug()
  persistActiveSession()
}

function stopRequest(): void {
  requestAbortController.value?.abort()
}

async function sendMessage(): Promise<void> {
  if (sending.value) {
    return
  }

  const content = userInput.value.trim()
  if (!content) {
    appStore.showWarning('Please enter a message')
    return
  }

  const apiKey = resolvedApiKey.value
  if (!apiKey) {
    appStore.showWarning('Please select or input an API key')
    return
  }

  const normalizedModel = model.value.trim() || DEFAULT_MODEL
  model.value = normalizedModel
  baseUrl.value = normalizeBaseUrl(baseUrl.value)
  endpointPath.value = normalizePath(endpointPath.value)

  const userMessage: ChatMessage = {
    id: ++messageId,
    role: 'user',
    content
  }
  const assistantMessage: ChatMessage = {
    id: ++messageId,
    role: 'assistant',
    content: '',
    pending: true
  }

  messages.value.push(userMessage, assistantMessage)
  userInput.value = ''
  sending.value = true
  persistSettings()
  ensureActiveSession()
  persistActiveSession()

  const controller = new AbortController()
  requestAbortController.value = controller

  const headers: Record<string, string> = {
    Authorization: `Bearer ${apiKey}`,
    'Content-Type': 'application/json'
  }

  const payload =
    protocol.value === 'responses'
      ? {
          model: normalizedModel,
          input: messages.value
            .filter((msg) => msg.role !== 'assistant' || !msg.pending)
            .map((msg) => ({
              role: msg.role,
              content: msg.content
            })),
          stream: stream.value
        }
      : {
          model: normalizedModel,
          messages: messages.value
            .filter((msg) => msg.role !== 'assistant' || !msg.pending)
            .map((msg) => ({
              role: msg.role,
              content: msg.content
            })),
          stream: stream.value
        }

  debug.requestUrl = requestUrl.value
  debug.headersSummary = `Authorization: Bearer ${maskKey(apiKey)}\nContent-Type: application/json`
  debug.requestBody = JSON.stringify(payload, null, 2)
  debug.statusText = ''
  debug.rawResponse = ''

  try {
    const response = await fetch(requestUrl.value, {
      method: 'POST',
      headers,
      body: JSON.stringify(payload),
      signal: controller.signal
    })

    if (stream.value && response.body) {
      await consumeStreamResponse(response, assistantMessage.id)
    } else {
      await consumeJsonResponse(response, assistantMessage.id)
    }
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Request failed'
    markAssistantError(assistantMessage.id, message)
    debug.statusText = message
  } finally {
    sending.value = false
    requestAbortController.value = null
    const target = messages.value.find((msg) => msg.id === assistantMessage.id)
    if (target) {
      target.pending = false
    }
    persistActiveSession()
  }
}

async function consumeJsonResponse(response: Response, assistantMessageId: number): Promise<void> {
  const text = await response.text()
  debug.statusText = `${response.status} ${response.statusText}`
  debug.rawResponse = clip(text)

  let body: Record<string, unknown> = {}
  try {
    body = text ? (JSON.parse(text) as Record<string, unknown>) : {}
  } catch {
    body = {}
  }

  if (!response.ok) {
    throw new Error(extractErrorMessage(body, text) || `${response.status} ${response.statusText}`)
  }

  const content = extractNonStreamContent(body, protocol.value) || '(empty response)'
  updateAssistantContent(assistantMessageId, content)
}

async function consumeStreamResponse(response: Response, assistantMessageId: number): Promise<void> {
  debug.statusText = `${response.status} ${response.statusText}`
  if (!response.ok || !response.body) {
    const text = await response.text()
    debug.rawResponse = clip(text)
    throw new Error(`${response.status} ${response.statusText}`)
  }

  const reader = response.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let sawDelta = false
  let terminalText = ''

  while (true) {
    const { done, value } = await reader.read()
    if (done) {
      break
    }
    buffer += decoder.decode(value, { stream: true })
    const chunks = buffer.split('\n')
    buffer = chunks.pop() ?? ''

    for (const line of chunks) {
      const trimmed = line.trim()
      if (!trimmed.startsWith('data:')) {
        continue
      }
      const rawData = trimmed.slice(5).trim()
      if (!rawData || rawData === '[DONE]') {
        continue
      }

      debug.rawResponse = clip(`${debug.rawResponse}\n${rawData}`)

      try {
        const parsed = JSON.parse(rawData) as Record<string, unknown>
        const chunk = extractStreamChunk(parsed, protocol.value)
        if (chunk.delta) {
          appendAssistantContent(assistantMessageId, chunk.delta)
          sawDelta = true
        }
        if (chunk.finalText) {
          terminalText = chunk.finalText
        }
      } catch {
        // Ignore invalid stream chunks and keep reading.
      }
    }
  }

  reconcileAssistantContent(assistantMessageId, terminalText, sawDelta)
}

function extractStreamChunk(event: Record<string, unknown>, value: ProtocolType): StreamChunkExtraction {
  if (value === 'chat-completions') {
    const choices = event.choices
    if (!Array.isArray(choices)) {
      return {}
    }
    const first = choices[0]
    if (!first || typeof first !== 'object') {
      return {}
    }
    const deltaObj = (first as { delta?: unknown }).delta
    if (deltaObj && typeof deltaObj === 'object') {
      const delta = (deltaObj as { content?: unknown }).content
      if (typeof delta === 'string' && delta.length > 0) {
        return { delta }
      }
    }
    const messageObj = (first as { message?: unknown }).message
    if (messageObj && typeof messageObj === 'object') {
      const content = (messageObj as { content?: unknown }).content
      if (typeof content === 'string' && content.length > 0) {
        return { finalText: content }
      }
    }
    return {}
  }

  const eventType = typeof event.type === 'string' ? event.type : ''
  if (eventType === 'response.output_text.delta') {
    const deltaField = event.delta
    return typeof deltaField === 'string' && deltaField.length > 0 ? { delta: deltaField } : {}
  }
  if (eventType === 'response.output_text.done') {
    const textField = event.text
    return typeof textField === 'string' && textField.length > 0 ? { finalText: textField } : {}
  }
  if (isResponsesTerminalEvent(eventType)) {
    return { finalText: extractResponsesTerminalText(event) }
  }

  return {}
}

function isResponsesTerminalEvent(eventType: string): boolean {
  return (
    eventType === 'response.completed' ||
    eventType === 'response.done' ||
    eventType === 'response.failed' ||
    eventType === 'response.incomplete'
  )
}

function extractResponsesTerminalText(event: Record<string, unknown>): string {
  const responseField = event.response
  if (responseField && typeof responseField === 'object') {
    return extractNonStreamContent(responseField as Record<string, unknown>, 'responses')
  }

  const deltaField = event.delta
  if (typeof deltaField === 'string') {
    return deltaField
  }
  const outputTextField = event.output_text
  if (typeof outputTextField === 'string') {
    return outputTextField
  }
  const textField = event.text
  if (typeof textField === 'string') {
    return textField
  }
  return ''
}

function extractNonStreamContent(body: Record<string, unknown>, value: ProtocolType): string {
  if (value === 'chat-completions') {
    const choices = body.choices
    if (!Array.isArray(choices)) {
      return ''
    }
    const first = choices[0]
    if (!first || typeof first !== 'object') {
      return ''
    }
    const messageObj = (first as { message?: unknown }).message
    if (!messageObj || typeof messageObj !== 'object') {
      return ''
    }
    const message = (messageObj as { content?: unknown }).content
    return typeof message === 'string' ? message : ''
  }

  const outputText = body.output_text
  if (typeof outputText === 'string') {
    return outputText
  }
  const output = body.output
  if (Array.isArray(output)) {
    const texts: string[] = []
    for (const item of output) {
      if (!item || typeof item !== 'object') {
        continue
      }
      const content = (item as { content?: unknown }).content
      if (!Array.isArray(content)) {
        continue
      }
      for (const contentItem of content) {
        if (!contentItem || typeof contentItem !== 'object') {
          continue
        }
        const text = (contentItem as { text?: unknown }).text
        if (typeof text === 'string' && text.length > 0) {
          texts.push(text)
        }
      }
    }
    return texts.join('\n')
  }
  return ''
}

function extractErrorMessage(body: Record<string, unknown>, fallback: string): string {
  const message = body.message
  const fromTopLevel = typeof message === 'string' ? message : ''
  const errorField = body.error
  let fromNested = ''
  if (errorField && typeof errorField === 'object') {
    const nestedMessage = (errorField as { message?: unknown }).message
    fromNested = typeof nestedMessage === 'string' ? nestedMessage : ''
  }
  return fromTopLevel || fromNested || fallback
}

function updateAssistantContent(messageIdValue: number, content: string): void {
  const target = messages.value.find((msg) => msg.id === messageIdValue)
  if (!target) {
    return
  }
  target.content = content
  target.error = false
}

function appendAssistantContent(messageIdValue: number, delta: string): void {
  const target = messages.value.find((msg) => msg.id === messageIdValue)
  if (!target) {
    return
  }
  target.content += delta
}

function reconcileAssistantContent(messageIdValue: number, finalText: string, sawDelta: boolean): void {
  if (!finalText) {
    return
  }
  const target = messages.value.find((msg) => msg.id === messageIdValue)
  if (!target) {
    return
  }
  if (!sawDelta || !target.content) {
    target.content = finalText
    return
  }
  if (finalText.startsWith(target.content) && finalText !== target.content) {
    target.content = finalText
  }
}

function markAssistantError(messageIdValue: number, content: string): void {
  const target = messages.value.find((msg) => msg.id === messageIdValue)
  if (!target) {
    return
  }
  target.content = content
  target.error = true
}

function hydrateSessionHistory(): void {
  try {
    const raw = localStorage.getItem(HISTORY_STORAGE_KEY)
    if (!raw) {
      return
    }
    const parsed = JSON.parse(raw) as { activeSessionId?: string | null; sessions?: ChatSessionSnapshot[] }
    sessionHistory.value = Array.isArray(parsed.sessions) ? parsed.sessions : []
    activeSessionId.value = typeof parsed.activeSessionId === 'string' ? parsed.activeSessionId : null

    const activeSession =
      sessionHistory.value.find((session) => session.id === activeSessionId.value) ?? sessionHistory.value[0]
    if (activeSession) {
      applySessionSnapshot(activeSession)
    }
  } catch {
    sessionHistory.value = []
    activeSessionId.value = null
  }
}

function persistSessionHistory(): void {
  localStorage.setItem(
    HISTORY_STORAGE_KEY,
    JSON.stringify({
      activeSessionId: activeSessionId.value,
      sessions: sessionHistory.value
    })
  )
}

function createSessionId(): string {
  return `session-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

function createDebugSnapshot(): ChatDebugSnapshot {
  return {
    requestUrl: debug.requestUrl,
    headersSummary: debug.headersSummary,
    requestBody: debug.requestBody,
    statusText: debug.statusText,
    rawResponse: debug.rawResponse
  }
}

function cloneMessages(value: ChatMessage[]): ChatMessage[] {
  return value.map((message) => ({ ...message }))
}

function createSessionTitle(value: ChatMessage[]): string {
  const firstUserMessage = value.find((message) => message.role === 'user' && message.content.trim())
  if (!firstUserMessage) {
    return 'New Session'
  }
  return firstUserMessage.content.length > 48
    ? `${firstUserMessage.content.slice(0, 48)}...`
    : firstUserMessage.content
}

function buildSessionSnapshot(sessionId: string, createdAt?: string): ChatSessionSnapshot {
  const now = new Date().toISOString()
  return {
    id: sessionId,
    title: createSessionTitle(messages.value),
    createdAt: createdAt || now,
    updatedAt: now,
    protocol: protocol.value,
    model: model.value.trim() || DEFAULT_MODEL,
    stream: stream.value,
    baseUrl: normalizeBaseUrl(baseUrl.value),
    endpointPath: normalizePath(endpointPath.value),
    messages: cloneMessages(messages.value),
    debug: createDebugSnapshot()
  }
}

function persistActiveSession(): void {
  if (!activeSessionId.value) {
    return
  }

  const existing = sessionHistory.value.find((session) => session.id === activeSessionId.value)
  const snapshot = buildSessionSnapshot(activeSessionId.value, existing?.createdAt)
  const nextSessions = sessionHistory.value.filter((session) => session.id !== activeSessionId.value)
  sessionHistory.value = [snapshot, ...nextSessions].sort((left, right) => right.updatedAt.localeCompare(left.updatedAt))
  persistSessionHistory()
}

function ensureActiveSession(): string {
  if (activeSessionId.value) {
    return activeSessionId.value
  }
  const sessionId = createSessionId()
  activeSessionId.value = sessionId
  sessionHistory.value = [buildSessionSnapshot(sessionId), ...sessionHistory.value]
  persistSessionHistory()
  return sessionId
}

function resetDebug(): void {
  debug.requestUrl = ''
  debug.headersSummary = ''
  debug.requestBody = ''
  debug.statusText = ''
  debug.rawResponse = ''
}

function applySessionSnapshot(session: ChatSessionSnapshot): void {
  activeSessionId.value = session.id
  protocol.value = session.protocol
  model.value = session.model
  stream.value = session.stream
  baseUrl.value = normalizeBaseUrl(session.baseUrl)
  endpointPath.value = normalizePath(session.endpointPath)
  endpointEdited.value = session.endpointPath !== getDefaultPath(session.protocol)
  messages.value = cloneMessages(session.messages)
  debug.requestUrl = session.debug.requestUrl
  debug.headersSummary = session.debug.headersSummary
  debug.requestBody = session.debug.requestBody
  debug.statusText = session.debug.statusText
  debug.rawResponse = session.debug.rawResponse
  messageId = messages.value.reduce((maxId, message) => Math.max(maxId, message.id), 0)
}

function loadSession(sessionId: string): void {
  const session = sessionHistory.value.find((item) => item.id === sessionId)
  if (!session) {
    return
  }
  applySessionSnapshot(session)
  persistSessionHistory()
}

function startNewSession(): void {
  const sessionId = createSessionId()
  activeSessionId.value = sessionId
  messages.value = []
  userInput.value = ''
  resetDebug()
  messageId = 0
  sessionHistory.value = [buildSessionSnapshot(sessionId), ...sessionHistory.value]
  persistSessionHistory()
}

function deleteSession(sessionId: string): void {
  const nextSessions = sessionHistory.value.filter((session) => session.id !== sessionId)
  sessionHistory.value = nextSessions

  if (activeSessionId.value === sessionId) {
    activeSessionId.value = nextSessions[0]?.id || null
    if (nextSessions[0]) {
      applySessionSnapshot(nextSessions[0])
    } else {
      messages.value = []
      userInput.value = ''
      resetDebug()
      messageId = 0
    }
  }

  persistSessionHistory()
}

function formatSessionTime(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString()
}

function clip(value: string, maxLength = 16000): string {
  if (value.length <= maxLength) {
    return value
  }
  return value.slice(value.length - maxLength)
}
</script>
