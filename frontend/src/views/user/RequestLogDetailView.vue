<template>
  <AppLayout>
    <div class="space-y-6">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.detailTitle') }}</p>
          <h1 class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">
            {{ requestLog?.request_id || `#${requestLog?.id || route.params.id}` }}
          </h1>
        </div>
        <div class="flex items-center gap-3">
          <button class="btn btn-secondary" @click="router.back()">
            {{ t('common.back') }}
          </button>
          <button class="btn btn-secondary" :disabled="loading" @click="loadRequestLog">
            {{ t('common.refresh') }}
          </button>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-4">
        <div class="card p-4">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
            {{ t('requestLogs.statusCode') }}
          </p>
          <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
            {{ requestLog?.status_code ?? '-' }}
          </p>
        </div>
        <div class="card p-4">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
            {{ t('requestLogs.tokens') }}
          </p>
          <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
            {{ totalTokens.toLocaleString() }}
          </p>
        </div>
        <div class="card p-4">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
            {{ t('requestLogs.cost') }}
          </p>
          <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
            ${{ formatCost(requestLog?.total_cost) }}
          </p>
        </div>
        <div class="card p-4">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
            {{ t('requestLogs.duration') }}
          </p>
          <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
            {{ formatDuration(requestLog?.duration_ms) }}
          </p>
        </div>
      </div>

      <div v-if="loading" class="card p-10 text-center text-sm text-gray-500 dark:text-gray-400">
        {{ t('common.loading') }}
      </div>

      <template v-else-if="requestLog">
        <section class="card p-6">
          <div class="flex flex-wrap items-start justify-between gap-3">
            <div>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('requestLogs.visualizedContent') }}
              </h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {{ t('requestLogs.visualizedDescription') }}
              </p>
            </div>
            <div class="flex flex-wrap items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
              <span class="rounded-full bg-gray-100 px-3 py-1 dark:bg-dark-700">
                {{ t('requestLogs.requestConversation') }}: {{ requestConversation.length }}
              </span>
              <span class="rounded-full bg-gray-100 px-3 py-1 dark:bg-dark-700">
                {{ t('requestLogs.responseConversation') }}: {{ responseConversation.length }}
              </span>
            </div>
          </div>

          <div
            v-if="conversationEntries.length === 0"
            class="mt-4 rounded-2xl border border-dashed border-gray-200 bg-gray-50 px-4 py-6 text-sm text-gray-500 dark:border-dark-700 dark:bg-dark-900/40 dark:text-gray-400"
          >
            {{ t('requestLogs.noConversationBlocks') }}
          </div>

          <div v-else class="mt-6">
            <div class="rounded-3xl border border-gray-200/80 bg-gradient-to-br from-white via-slate-50 to-blue-50/60 p-4 dark:border-dark-700 dark:from-dark-900 dark:via-dark-900 dark:to-slate-950">
              <div class="mb-4 flex items-center justify-between gap-3">
                <div>
                  <p class="text-xs font-semibold uppercase tracking-[0.18em] text-gray-500 dark:text-gray-400">
                    {{ t('requestLogs.sessionTimeline') }}
                  </p>
                  <p class="mt-1 text-sm text-gray-600 dark:text-gray-300">
                    {{ t('requestLogs.sessionTimelineDescription') }}
                  </p>
                </div>
                <div class="flex flex-wrap items-center gap-2">
                  <span class="rounded-full bg-white/80 px-3 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-gray-300">
                    {{ t('requestLogs.turnCount', { count: turnGroups.length }) }}
                  </span>
                  <span class="rounded-full bg-white/80 px-3 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-gray-300">
                    {{ t('requestLogs.stepCount', { count: sessionTimeline.length }) }}
                  </span>
                </div>
              </div>

              <div class="space-y-4">
                <div
                  v-for="turn in turnGroups"
                  :key="turn.id"
                  class="grid grid-cols-[44px_minmax(0,1fr)] gap-3"
                >
                  <div class="relative flex justify-center">
                    <div class="flex h-11 w-11 items-center justify-center rounded-2xl border text-sm font-semibold shadow-sm" :class="turnStepClass(turn)">
                      {{ turn.step }}
                    </div>
                    <div
                      v-if="turn.step < turnGroups.length"
                      class="absolute top-12 h-[calc(100%+0.75rem)] w-px bg-gradient-to-b from-gray-300 to-transparent dark:from-dark-600"
                    ></div>
                  </div>

                  <article class="rounded-2xl border p-4 shadow-sm" :class="turnCardClass(turn)">
                    <div class="flex flex-wrap items-start justify-between gap-3">
                      <div>
                        <div class="flex flex-wrap items-center gap-2">
                          <span class="rounded-full px-2.5 py-1 text-xs font-semibold" :class="roleBadgeClass(turn.primaryRole)">
                            {{ roleLabel(turn.primaryRole) }}
                          </span>
                          <span class="rounded-full bg-gray-100 px-2.5 py-1 text-xs font-medium text-gray-700 dark:bg-dark-700 dark:text-gray-300">
                            {{ turnSourceLabel(turn) }}
                          </span>
                          <span class="rounded-full bg-white/80 px-2.5 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-gray-300">
                            {{ turn.protocolSummary }}
                          </span>
                          <span
                            v-if="turn.entries.length > 1"
                            class="rounded-full bg-white/80 px-2.5 py-1 text-xs font-medium text-gray-600 dark:bg-dark-800 dark:text-gray-300"
                          >
                            {{ t('requestLogs.turnStepSummary', { count: turn.entries.length }) }}
                          </span>
                        </div>
                        <h3 class="mt-3 text-base font-semibold text-gray-900 dark:text-white">
                          {{ turn.title }}
                        </h3>
                      </div>
                    </div>

                    <div class="mt-4 space-y-4">
                      <section
                        v-for="entry in turn.entries"
                        :key="entry.id"
                        class="rounded-2xl border border-black/5 bg-white/70 p-3 dark:border-white/5 dark:bg-black/10"
                      >
                        <div class="mb-3 flex flex-wrap items-center gap-2">
                          <span class="rounded-full px-2.5 py-1 text-[11px] font-semibold" :class="roleBadgeClass(entry.role)">
                            {{ roleLabel(entry.role) }}
                          </span>
                          <span class="rounded-full bg-gray-100 px-2.5 py-1 text-[11px] font-medium text-gray-700 dark:bg-dark-700 dark:text-gray-300">
                            {{ entry.source === 'request' ? t('requestLogs.sourceRequest') : t('requestLogs.sourceResponse') }}
                          </span>
                          <span class="text-xs font-medium text-gray-500 dark:text-gray-400">
                            {{ timelineEntryTitle(entry) }}
                          </span>
                        </div>

                        <div class="space-y-3">
                          <section
                            v-for="block in entry.blocks"
                            :key="block.id"
                            class="rounded-xl border border-black/5 bg-white/70 p-3 dark:border-white/5 dark:bg-black/10"
                          >
                            <div
                              v-if="block.label"
                              class="mb-2 text-[11px] font-semibold uppercase tracking-[0.14em] text-gray-500 dark:text-gray-400"
                            >
                              {{ block.label }}
                            </div>

                            <div
                              v-if="block.kind === 'markdown'"
                              class="markdown-body prose prose-sm max-w-none break-words dark:prose-invert"
                              v-html="renderMarkdown(block.text)"
                            ></div>
                            <pre
                              v-else-if="block.kind === 'json'"
                              class="overflow-x-auto whitespace-pre-wrap break-words rounded-lg bg-gray-50 p-3 text-xs leading-6 text-gray-800 dark:bg-dark-900 dark:text-gray-100"
                            >{{ block.text }}</pre>
                            <div
                              v-else
                              class="whitespace-pre-wrap break-words text-sm leading-6 text-gray-800 dark:text-gray-100"
                            >
                              {{ block.text }}
                            </div>
                          </section>
                        </div>
                      </section>
                    </div>
                  </article>
                </div>
              </div>
            </div>
          </div>
        </section>

        <div class="grid grid-cols-1 gap-6 xl:grid-cols-2">
          <section class="card p-6">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('requestLogs.metadata') }}
            </h2>
            <dl class="mt-4 space-y-4">
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.requestId') }}</dt>
                <dd class="text-sm break-all font-medium text-gray-900 dark:text-white">{{ requestLog.request_id || missingValue }}</dd>
              </div>
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.modelFilter') }}</dt>
                <dd class="text-sm font-medium text-gray-900 dark:text-white">{{ requestLog.model || missingValue }}</dd>
              </div>
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.apiKeyFilter') }}</dt>
                <dd class="text-sm font-medium text-gray-900 dark:text-white">#{{ requestLog.api_key_id }}</dd>
              </div>
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.method') }}</dt>
                <dd class="text-sm font-medium text-gray-900 dark:text-white">{{ requestLog.method || missingValue }}</dd>
              </div>
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.stream') }}</dt>
                <dd class="text-sm font-medium text-gray-900 dark:text-white">
                  {{ requestLog.stream ? t('requestLogs.streamLabel') : t('requestLogs.syncLabel') }}
                </dd>
              </div>
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.createdAt') }}</dt>
                <dd class="text-sm font-medium text-gray-900 dark:text-white">{{ formatDateTime(requestLog.created_at) }}</dd>
              </div>
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.userAgent') }}</dt>
                <dd class="text-sm break-all text-gray-900 dark:text-white">{{ requestLog.user_agent || missingValue }}</dd>
              </div>
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.ipAddress') }}</dt>
                <dd class="text-sm font-medium text-gray-900 dark:text-white">{{ requestLog.ip_address || missingValue }}</dd>
              </div>
            </dl>
          </section>

          <section class="card p-6">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('requestLogs.routing') }}
            </h2>
            <dl class="mt-4 space-y-4">
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.inboundEndpoint') }}</dt>
                <dd class="text-sm break-all text-gray-900 dark:text-white">{{ requestLog.inbound_endpoint || missingValue }}</dd>
              </div>
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.upstreamEndpoint') }}</dt>
                <dd class="text-sm break-all text-gray-900 dark:text-white">{{ requestLog.upstream_endpoint || missingValue }}</dd>
              </div>
            </dl>

            <h2 class="mt-8 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('requestLogs.performance') }}
            </h2>
            <dl class="mt-4 space-y-4">
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.tokens') }}</dt>
                <dd class="text-sm font-medium text-gray-900 dark:text-white">
                  {{ requestLog.input_tokens.toLocaleString() }} / {{ requestLog.output_tokens.toLocaleString() }}
                </dd>
              </div>
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.firstToken') }}</dt>
                <dd class="text-sm font-medium text-gray-900 dark:text-white">{{ formatDuration(requestLog.first_token_ms) }}</dd>
              </div>
              <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
                <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.duration') }}</dt>
                <dd class="text-sm font-medium text-gray-900 dark:text-white">{{ formatDuration(requestLog.duration_ms) }}</dd>
              </div>
            </dl>
          </section>
        </div>

        <section class="card p-6">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('requestLogs.errorInfo') }}
          </h2>
          <dl class="mt-4 space-y-4">
            <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
              <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.errorCode') }}</dt>
              <dd class="text-sm font-medium text-gray-900 dark:text-white">{{ requestLog.error_code || missingValue }}</dd>
            </div>
            <div class="grid grid-cols-1 gap-1 sm:grid-cols-[160px_1fr]">
              <dt class="text-sm text-gray-500 dark:text-gray-400">{{ t('requestLogs.errorMessage') }}</dt>
              <dd class="text-sm whitespace-pre-wrap break-words text-gray-900 dark:text-white">{{ requestLog.error_message || missingValue }}</dd>
            </div>
          </dl>
        </section>

        <div class="grid grid-cols-1 gap-6 xl:grid-cols-2">
          <section class="card p-6">
            <div class="flex items-center justify-between gap-3">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('requestLogs.requestBody') }}
              </h2>
              <div class="text-xs text-gray-500 dark:text-gray-400">
                <span>{{ t('requestLogs.bodyBytes') }}: {{ requestLog.payload?.request_body_bytes ?? missingValue }}</span>
                <span class="ml-3">{{ t('requestLogs.truncated') }}: {{ requestLog.payload?.request_body_truncated ? t('requestLogs.yes') : t('requestLogs.no') }}</span>
              </div>
            </div>
            <pre class="mt-4 overflow-x-auto rounded-lg bg-gray-50 p-4 text-xs leading-6 text-gray-800 dark:bg-gray-900 dark:text-gray-100">{{ requestLog.payload?.request_body || t('requestLogs.payloadUnavailable') }}</pre>
          </section>

          <section class="card p-6">
            <div class="flex items-center justify-between gap-3">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('requestLogs.responseBody') }}
              </h2>
              <div class="text-xs text-gray-500 dark:text-gray-400">
                <span>{{ t('requestLogs.bodyBytes') }}: {{ requestLog.payload?.response_body_bytes ?? missingValue }}</span>
                <span class="ml-3">{{ t('requestLogs.truncated') }}: {{ requestLog.payload?.response_body_truncated ? t('requestLogs.yes') : t('requestLogs.no') }}</span>
              </div>
            </div>
            <pre class="mt-4 overflow-x-auto rounded-lg bg-gray-50 p-4 text-xs leading-6 text-gray-800 dark:bg-gray-900 dark:text-gray-100">{{ requestLog.payload?.response_body || t('requestLogs.payloadUnavailable') }}</pre>
          </section>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { requestLogsAPI } from '@/api'
import type { RequestLog } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useAppStore } from '@/stores/app'
import { formatDateTime } from '@/utils/format'

interface ConversationBlock {
  id: string
  kind: 'markdown' | 'json' | 'plain'
  label?: string
  text: string
}

interface ConversationEntry {
  id: string
  role: string
  source: 'request' | 'response'
  protocol: string
  blocks: ConversationBlock[]
}

interface TimelineEntry extends ConversationEntry {
  step: number
}

interface TurnGroup {
  id: string
  step: number
  title: string
  primaryRole: string
  entries: TimelineEntry[]
  protocolSummary: string
  containsRequest: boolean
  containsResponse: boolean
}

interface ParsedSSEEvent {
  event: string
  payload: Record<string, unknown> | null
  rawData: string
}

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()

marked.setOptions({
  breaks: true,
  gfm: true
})

let abortController: AbortController | null = null

const loading = ref(false)
const requestLog = ref<RequestLog | null>(null)

const missingValue = computed(() => t('requestLogs.missingValue'))
const totalTokens = computed(() => {
  const log = requestLog.value
  return (log?.input_tokens || 0) + (log?.output_tokens || 0)
})

const parsedRequestPayload = computed(() => parseJSONPayload(requestLog.value?.payload?.request_body))
const parsedResponsePayload = computed(() => parseJSONPayload(requestLog.value?.payload?.response_body))
const parsedResponseEvents = computed(() => parseSSEEvents(requestLog.value?.payload?.response_body))

const requestConversation = computed(() => normalizeRequestConversation(parsedRequestPayload.value))
const responseConversation = computed(() => {
  if (parsedResponsePayload.value) {
    return normalizeResponseConversation(parsedResponsePayload.value)
  }
  if (parsedResponseEvents.value.length > 0) {
    return normalizeResponseStreamConversation(parsedResponseEvents.value)
  }
  return []
})
const conversationEntries = computed(() => [...requestConversation.value, ...responseConversation.value])
const sessionTimeline = computed<TimelineEntry[]>(() =>
  conversationEntries.value.map((entry, index) => ({
    ...entry,
    step: index + 1
  }))
)
const turnGroups = computed<TurnGroup[]>(() => groupTimelineEntries(sessionTimeline.value))

function parseJSONPayload(value: string | null | undefined): Record<string, unknown> | null {
  if (!value) {
    return null
  }
  try {
    const parsed = JSON.parse(value) as unknown
    return parsed && typeof parsed === 'object' && !Array.isArray(parsed)
      ? (parsed as Record<string, unknown>)
      : null
  } catch {
    return null
  }
}

function parseSSEEvents(value: string | null | undefined): ParsedSSEEvent[] {
  if (!value) {
    return []
  }

  const trimmed = value.trim()
  if (!trimmed.startsWith('data:') && !trimmed.startsWith('event:')) {
    return []
  }

  const events: ParsedSSEEvent[] = []
  const blocks = trimmed.split(/\n\s*\n/)
  for (const block of blocks) {
    const lines = block.split('\n')
    let eventName = ''
    const dataLines: string[] = []

    for (const rawLine of lines) {
      const line = rawLine.trim()
      if (!line) {
        continue
      }
      if (line.startsWith('event:')) {
        eventName = line.slice(6).trim()
        continue
      }
      if (line.startsWith('data:')) {
        dataLines.push(line.slice(5).trim())
      }
    }

    if (dataLines.length === 0) {
      continue
    }

    const rawData = dataLines.join('\n')
    if (!rawData || rawData === '[DONE]') {
      continue
    }

    let payload: Record<string, unknown> | null = null
    try {
      const parsed = JSON.parse(rawData) as unknown
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
        payload = parsed as Record<string, unknown>
      }
    } catch {
      payload = null
    }

    events.push({
      event: eventName,
      payload,
      rawData
    })
  }

  return events
}

function detectProtocol(payload: Record<string, unknown> | null): string {
  if (!payload) {
    return 'unknown'
  }
  if (payload.input !== undefined || typeof payload.output_text === 'string' || Array.isArray(payload.output)) {
    return 'responses'
  }
  if (Array.isArray(payload.choices)) {
    return 'chat-completions'
  }
  if (Array.isArray(payload.messages) || Array.isArray(payload.content)) {
    return 'anthropic'
  }
  return 'generic'
}

function normalizeRequestConversation(payload: Record<string, unknown> | null): ConversationEntry[] {
  if (!payload) {
    return []
  }

  const protocol = detectProtocol(payload)
  const entries: ConversationEntry[] = []

  pushSystemEntries(entries, payload.system, protocol)

  normalizeRequestInput(payload.input, protocol).forEach((entry) => {
    entries.push(entry)
  })

  const messages = payload.messages
  if (Array.isArray(messages)) {
    messages.forEach((item, index) => {
      const entry = normalizeInputItem(item, protocol, 'request', index)
      if (entry) {
        entries.push(entry)
      }
    })
  }

  return entries
}

function normalizeRequestInput(value: unknown, protocol: string): ConversationEntry[] {
  if (typeof value === 'string') {
    const entry = normalizeInputItem(value, protocol, 'request', 0)
    return entry ? [entry] : []
  }

  if (Array.isArray(value)) {
    return value
      .map((item, index) => normalizeInputItem(item, protocol, 'request', index))
      .filter((entry): entry is ConversationEntry => entry !== null)
  }

  if (value && typeof value === 'object') {
    const entry = normalizeInputItem(value, protocol, 'request', 0)
    return entry ? [entry] : []
  }

  return []
}

function normalizeResponseConversation(payload: Record<string, unknown> | null): ConversationEntry[] {
  if (!payload) {
    return []
  }

  const protocol = detectProtocol(payload)
  const entries: ConversationEntry[] = []

  const output = payload.output
  if (Array.isArray(output)) {
    output.forEach((item, index) => {
      const entry = normalizeOutputItem(item, protocol, index)
      if (entry) {
        entries.push(entry)
      }
    })
  }

  const choices = payload.choices
  if (Array.isArray(choices)) {
    choices.forEach((choice, index) => {
      if (!choice || typeof choice !== 'object') {
        return
      }
      const choiceObject = choice as Record<string, unknown>
      const message = choiceObject.message
      if (message && typeof message === 'object') {
        const entry = normalizeInputItem(message, protocol, 'response', index)
        if (entry) {
          entries.push(entry)
        }
      }
      const toolCalls = choiceObject.tool_calls
      if (Array.isArray(toolCalls) && toolCalls.length > 0) {
        const blocks = toolCalls
          .map((toolCall, blockIndex) => normalizeJSONBlock(toolCall, blockIndex, t('requestLogs.toolCall')))
          .filter((block): block is ConversationBlock => block !== null)
        if (blocks.length > 0) {
          entries.push({
            id: `response-tool-${index}`,
            role: 'assistant',
            source: 'response',
            protocol,
            blocks
          })
        }
      }
    })
  }

  const content = payload.content
  if (Array.isArray(content) && entries.length === 0) {
    const blocks = normalizeContentBlocks(content)
    if (blocks.length > 0) {
      entries.push({
        id: 'response-anthropic-content',
        role: String(payload.role || 'assistant'),
        source: 'response',
        protocol,
        blocks
      })
    }
  }

  const outputText = payload.output_text
  if (typeof outputText === 'string' && outputText.trim() && entries.length === 0) {
    entries.push({
      id: 'response-output-text',
      role: 'assistant',
      source: 'response',
      protocol,
      blocks: [createMarkdownBlock('response-output-text', outputText)]
    })
  }

  return entries
}

function normalizeResponseStreamConversation(events: ParsedSSEEvent[]): ConversationEntry[] {
  const terminalResponse = findOpenAITerminalResponse(events)
  if (terminalResponse) {
    return normalizeResponseConversation(terminalResponse)
  }

  const entries: ConversationEntry[] = []
  const anthropicBlocks = new Map<number, {
    type: string
    text: string
    name?: string
    input?: unknown
    partialJSON: string
  }>()

  for (const event of events) {
    const payload = event.payload
    if (!payload) {
      continue
    }

    const eventType = typeof payload.type === 'string' ? payload.type : event.event
    if (eventType === 'content_block_start') {
      const index = Number(payload.index ?? anthropicBlocks.size)
      const contentBlock = payload.content_block
      if (!contentBlock || typeof contentBlock !== 'object' || Array.isArray(contentBlock)) {
        continue
      }
      const block = contentBlock as Record<string, unknown>
      const blockType = typeof block.type === 'string' ? block.type : 'unknown'
      anthropicBlocks.set(index, {
        type: blockType,
        text: typeof block.text === 'string' ? block.text : '',
        name: typeof block.name === 'string' ? block.name : undefined,
        input: block.input,
        partialJSON: ''
      })
      continue
    }

    if (eventType === 'content_block_delta') {
      const index = Number(payload.index ?? 0)
      const delta = payload.delta
      if (!delta || typeof delta !== 'object' || Array.isArray(delta)) {
        continue
      }
      const deltaObject = delta as Record<string, unknown>
      const current = anthropicBlocks.get(index) ?? {
        type: 'text',
        text: '',
        partialJSON: ''
      }
      const deltaType = typeof deltaObject.type === 'string' ? deltaObject.type : ''
      if (deltaType === 'text_delta') {
        current.text += typeof deltaObject.text === 'string' ? deltaObject.text : ''
      } else if (deltaType === 'thinking_delta') {
        current.type = 'thinking'
        current.text += typeof deltaObject.thinking === 'string' ? deltaObject.thinking : ''
      } else if (deltaType === 'input_json_delta') {
        current.type = 'tool_use'
        current.partialJSON += typeof deltaObject.partial_json === 'string' ? deltaObject.partial_json : ''
      }
      anthropicBlocks.set(index, current)
    }
  }

  if (anthropicBlocks.size > 0) {
    Array.from(anthropicBlocks.entries())
      .sort(([left], [right]) => left - right)
      .forEach(([index, block]) => {
        if (block.type === 'tool_use') {
          const toolEntry: ConversationEntry = {
            id: `response-stream-tool-${index}`,
            role: 'tool',
            source: 'response',
            protocol: 'anthropic-stream',
            blocks: []
          }
          if (block.name) {
            toolEntry.blocks.push({
              id: `response-stream-tool-name-${index}`,
              kind: 'plain',
              label: t('requestLogs.toolName'),
              text: block.name
            })
          }
          const toolInput = block.partialJSON ? parseMaybeJSON(block.partialJSON) ?? block.partialJSON : block.input
          const argsBlock = normalizeStructuredBlock(toolInput, `response-stream-tool-args-${index}`, t('requestLogs.toolArguments'))
          if (argsBlock) {
            toolEntry.blocks.push(argsBlock)
          }
          entries.push(toolEntry)
          return
        }

        if (block.text.trim()) {
          entries.push({
            id: `response-stream-assistant-${index}`,
            role: 'assistant',
            source: 'response',
            protocol: 'anthropic-stream',
            blocks: [
              {
                id: `response-stream-block-${index}`,
                kind: block.type === 'thinking' ? 'plain' : 'markdown',
                label: block.type === 'thinking' ? t('requestLogs.reasoning') : undefined,
                text: block.text
              }
            ]
          })
        }
      })
  }

  if (entries.length === 0) {
    const transcript = events.map((event) => event.rawData).join('\n')
    if (transcript) {
      entries.push({
        id: 'response-stream-fallback',
        role: 'assistant',
        source: 'response',
        protocol: 'sse',
        blocks: [
          {
            id: 'response-stream-fallback-block',
            kind: 'plain',
            label: t('requestLogs.responseBody'),
            text: transcript
          }
        ]
      })
    }
  }

  return entries
}

function findOpenAITerminalResponse(events: ParsedSSEEvent[]): Record<string, unknown> | null {
  for (let index = events.length - 1; index >= 0; index -= 1) {
    const payload = events[index].payload
    if (!payload) {
      continue
    }
    const eventType = typeof payload.type === 'string' ? payload.type : ''
    if (!isResponsesTerminalEvent(eventType)) {
      continue
    }
    const response = payload.response
    if (response && typeof response === 'object' && !Array.isArray(response)) {
      return response as Record<string, unknown>
    }
  }
  return null
}

function isResponsesTerminalEvent(eventType: string): boolean {
  return (
    eventType === 'response.completed' ||
    eventType === 'response.done' ||
    eventType === 'response.failed' ||
    eventType === 'response.incomplete'
  )
}

function pushSystemEntries(entries: ConversationEntry[], systemValue: unknown, protocol: string): void {
  const blocks = normalizeContentValue(systemValue)
  if (blocks.length === 0) {
    return
  }
  entries.push({
    id: 'request-system',
    role: 'system',
    source: 'request',
    protocol,
    blocks
  })
}

function normalizeInputItem(
  value: unknown,
  protocol: string,
  source: 'request' | 'response',
  index: number
): ConversationEntry | null {
  if (!value) {
    return null
  }
  if (typeof value === 'string') {
    return {
      id: `${source}-${protocol}-${index}`,
      role: source === 'request' ? 'user' : 'assistant',
      source,
      protocol,
      blocks: [createMarkdownBlock(`${source}-${protocol}-${index}-text`, value)]
    }
  }
  if (typeof value !== 'object' || Array.isArray(value)) {
    return null
  }

  const item = value as Record<string, unknown>
  const role = typeof item.role === 'string' && item.role.trim() ? item.role : source === 'request' ? 'user' : 'assistant'
  const contentBlocks = normalizeContentValue(item.content, role)

  if (contentBlocks.length === 0) {
    const fallback = normalizeToolBlock(item, 0) ?? normalizeJSONBlock(item, 0, t('requestLogs.structuredData'))
    if (!fallback) {
      return null
    }
    contentBlocks.push(fallback)
  }

  return {
    id: `${source}-${protocol}-${index}`,
    role,
    source,
    protocol,
    blocks: contentBlocks
  }
}

function normalizeOutputItem(value: unknown, protocol: string, index: number): ConversationEntry | null {
  if (!value || typeof value !== 'object' || Array.isArray(value)) {
    return null
  }
  const item = value as Record<string, unknown>
  const itemType = typeof item.type === 'string' ? item.type : ''

  if (itemType === 'message') {
    return {
      id: `response-${protocol}-${index}`,
      role: typeof item.role === 'string' && item.role.trim() ? item.role : 'assistant',
      source: 'response',
      protocol,
      blocks: normalizeContentValue(item.content, typeof item.role === 'string' ? item.role : 'assistant')
    }
  }

  const blocks = normalizeToolBlocks(item, 0)
  if (blocks.length === 0) {
    const label = outputItemLabel(itemType)
    const fallback = normalizeJSONBlock(item, 0, label)
    if (!fallback) {
      return null
    }
    blocks.push(fallback)
  }
  if (blocks.length === 0) {
    return null
  }
  return {
    id: `response-${protocol}-${index}`,
    role: itemType.includes('tool') ? 'tool' : 'assistant',
    source: 'response',
    protocol,
    blocks
  }
}

function normalizeContentValue(value: unknown, role = ''): ConversationBlock[] {
  if (typeof value === 'string') {
    return value.trim() ? [createMarkdownBlock(`text-${value.length}`, value)] : []
  }
  if (!Array.isArray(value)) {
    return []
  }
  return normalizeContentBlocks(value, role)
}

function normalizeContentBlocks(parts: unknown[], role = ''): ConversationBlock[] {
  const blocks: ConversationBlock[] = []

  parts.forEach((part, index) => {
    if (typeof part === 'string') {
      if (part.trim()) {
        blocks.push(createMarkdownBlock(`part-${index}`, part))
      }
      return
    }
    if (!part || typeof part !== 'object' || Array.isArray(part)) {
      return
    }

    const item = part as Record<string, unknown>
    const partType = typeof item.type === 'string' ? item.type : ''

    if (partType === 'text' || partType === 'input_text' || partType === 'output_text') {
      const text = typeof item.text === 'string' ? item.text : ''
      if (text.trim()) {
        blocks.push(createMarkdownBlock(`part-${index}`, text))
      }
      return
    }

    if (partType === 'reasoning' || partType === 'thinking') {
      const text = typeof item.text === 'string' ? item.text : safeJSONStringify(item)
      if (text) {
        blocks.push({
          id: `part-${index}`,
          kind: 'plain',
          label: t('requestLogs.reasoning'),
          text
        })
      }
      return
    }

    if (partType === 'input_image' || partType === 'image' || partType === 'image_url') {
      blocks.push({
        id: `part-${index}`,
        kind: 'plain',
        label: t('requestLogs.media'),
        text: safeJSONStringify(item)
      })
      return
    }

    if (partType.includes('tool') || partType.includes('function')) {
      const toolBlocks = normalizeToolBlocks(item, index)
      if (toolBlocks.length > 0) {
        blocks.push(...toolBlocks)
      }
      return
    }

    const text = typeof item.text === 'string' ? item.text : ''
    if (text.trim()) {
      blocks.push(createMarkdownBlock(`part-${index}`, text))
      return
    }

    const fallback = normalizeJSONBlock(item, index, t('requestLogs.structuredData'))
    if (fallback) {
      blocks.push(fallback)
    }
  })

  if (blocks.length === 0 && role === 'tool') {
    const text = parts
      .map((part) => {
        if (typeof part === 'string') {
          return part
        }
        if (part && typeof part === 'object' && !Array.isArray(part)) {
          return safeJSONStringify(part)
        }
        return ''
      })
      .filter(Boolean)
      .join('\n')
    if (text) {
      blocks.push({
        id: 'tool-fallback-output',
        kind: 'plain',
        label: t('requestLogs.toolOutput'),
        text
      })
    }
  }

  return blocks
}

function normalizeToolBlocks(value: Record<string, unknown>, index: number): ConversationBlock[] {
  const blocks: ConversationBlock[] = []
  const toolBlock = normalizeToolBlock(value, index)
  if (toolBlock) {
    blocks.push(toolBlock)
  }
  const argumentsBlock = normalizeToolArgumentsBlock(value, index)
  if (argumentsBlock) {
    blocks.push(argumentsBlock)
  }
  const outputBlock = normalizeToolOutputBlock(value, index)
  if (outputBlock) {
    blocks.push(outputBlock)
  }
  return blocks
}

function normalizeToolBlock(value: Record<string, unknown>, index: number): ConversationBlock | null {
  const functionObject =
    value.function && typeof value.function === 'object' && !Array.isArray(value.function)
      ? (value.function as Record<string, unknown>)
      : null
  const nameCandidates = [
    value.name,
    value.tool_name,
    value.call_name,
    functionObject?.name
  ]
  const name = nameCandidates.find((candidate) => typeof candidate === 'string' && candidate.trim()) as string | undefined
  if (!name) {
    return null
  }
  return {
    id: `tool-name-${index}`,
    kind: 'plain',
    label: t('requestLogs.toolName'),
    text: name
  }
}

function normalizeToolArgumentsBlock(value: Record<string, unknown>, index: number): ConversationBlock | null {
  const functionObject =
    value.function && typeof value.function === 'object' && !Array.isArray(value.function)
      ? (value.function as Record<string, unknown>)
      : null
  const candidates = [
    functionObject?.arguments,
    value.arguments,
    value.input,
    value.parameters
  ]
  const args = candidates.find((candidate) => candidate != null)
  return normalizeStructuredBlock(args, `tool-arguments-${index}`, t('requestLogs.toolArguments'))
}

function normalizeToolOutputBlock(value: Record<string, unknown>, index: number): ConversationBlock | null {
  const candidates = [
    value.output,
    value.result,
    value.response,
    value.content,
    value.output_text
  ]
  const output = candidates.find((candidate) => candidate != null)
  return normalizeStructuredBlock(output, `tool-output-${index}`, t('requestLogs.toolOutput'))
}

function normalizeStructuredBlock(value: unknown, id: string, label: string): ConversationBlock | null {
  if (value == null) {
    return null
  }
  if (typeof value === 'string') {
    const trimmed = value.trim()
    if (!trimmed) {
      return null
    }
    const parsed = parseMaybeJSON(trimmed)
    if (parsed !== null) {
      return {
        id,
        kind: 'json',
        label,
        text: safeJSONStringify(parsed)
      }
    }
    return {
      id,
      kind: 'plain',
      label,
      text: trimmed
    }
  }
  return {
    id,
    kind: 'json',
    label,
    text: safeJSONStringify(value)
  }
}

function parseMaybeJSON(value: string): unknown | null {
  if (!value) {
    return null
  }
  const first = value[0]
  if (first !== '{' && first !== '[') {
    return null
  }
  try {
    return JSON.parse(value) as unknown
  } catch {
    return null
  }
}

function normalizeJSONBlock(value: unknown, index: number, label: string): ConversationBlock | null {
  const text = safeJSONStringify(value)
  if (!text) {
    return null
  }
  return {
    id: `json-${index}`,
    kind: 'json',
    label,
    text
  }
}

function createMarkdownBlock(id: string, text: string): ConversationBlock {
  return {
    id,
    kind: 'markdown',
    text
  }
}

function safeJSONStringify(value: unknown): string {
  try {
    return JSON.stringify(value, null, 2)
  } catch {
    return ''
  }
}

function outputItemLabel(type: string): string {
  if (type.includes('tool') || type.includes('function')) {
    return t('requestLogs.toolCall')
  }
  if (type.includes('reasoning') || type.includes('thinking')) {
    return t('requestLogs.reasoning')
  }
  if (type.includes('image')) {
    return t('requestLogs.media')
  }
  return t('requestLogs.structuredData')
}

function renderMarkdown(content: string): string {
  const html = marked.parse(content) as string
  return DOMPurify.sanitize(html)
}

function groupTimelineEntries(entries: TimelineEntry[]): TurnGroup[] {
  const groups: TurnGroup[] = []
  let currentConversationGroup: TurnGroup | null = null

  entries.forEach((entry) => {
    if (isMetaEntry(entry)) {
      const standaloneGroup = createTurnGroup(groups.length + 1, [entry])
      groups.push(standaloneGroup)
      currentConversationGroup = null
      return
    }

    if (entry.source === 'request') {
      if (entry.role === 'user' || currentConversationGroup == null) {
        currentConversationGroup = createTurnGroup(groups.length + 1, [entry])
        groups.push(currentConversationGroup)
        return
      }

      currentConversationGroup.entries.push(entry)
      refreshTurnGroup(currentConversationGroup)
      return
    }

    if (currentConversationGroup == null) {
      currentConversationGroup = createTurnGroup(groups.length + 1, [entry])
      groups.push(currentConversationGroup)
      return
    }

    currentConversationGroup.entries.push(entry)
    refreshTurnGroup(currentConversationGroup)
  })

  return groups
}

function createTurnGroup(step: number, entries: TimelineEntry[]): TurnGroup {
  const group: TurnGroup = {
    id: `turn-${step}-${entries.map((entry) => entry.id).join('-')}`,
    step,
    title: '',
    primaryRole: entries[0]?.role || 'assistant',
    entries,
    protocolSummary: '',
    containsRequest: false,
    containsResponse: false
  }
  refreshTurnGroup(group)
  return group
}

function refreshTurnGroup(group: TurnGroup): void {
  group.primaryRole = group.entries[0]?.role || 'assistant'
  group.containsRequest = group.entries.some((entry) => entry.source === 'request')
  group.containsResponse = group.entries.some((entry) => entry.source === 'response')
  group.protocolSummary =
    Array.from(
      new Set(group.entries.map((entry) => entry.protocol).filter((protocol) => protocol && protocol !== 'unknown'))
    ).join(' / ') || t('requestLogs.protocolUnknown')

  const leadUser = group.entries.find((entry) => entry.role === 'user')
  if (leadUser) {
    group.title = t('requestLogs.turnWithUserMessage')
    return
  }

  const leadEntry = group.entries[0]
  group.title = leadEntry ? timelineEntryTitle(leadEntry) : t('requestLogs.assistantReplyTitle')
}

function isMetaEntry(entry: ConversationEntry): boolean {
  return entry.role === 'system' || entry.role === 'developer'
}

function timelineEntryTitle(entry: ConversationEntry): string {
  if (entry.role === 'tool') {
    return entry.source === 'request' ? t('requestLogs.toolResultTitle') : t('requestLogs.toolCallTitle')
  }
  if (entry.role === 'assistant') {
    return entry.source === 'request' ? t('requestLogs.assistantDraftTitle') : t('requestLogs.assistantReplyTitle')
  }
  if (entry.role === 'system') {
    return t('requestLogs.systemPromptTitle')
  }
  if (entry.role === 'developer') {
    return t('requestLogs.developerPromptTitle')
  }
  return t('requestLogs.userMessageTitle')
}

function turnSourceLabel(turn: TurnGroup): string {
  if (turn.containsRequest && turn.containsResponse) {
    return t('requestLogs.sourceMerged')
  }
  if (turn.containsRequest) {
    return t('requestLogs.sourceRequest')
  }
  return t('requestLogs.sourceResponse')
}

function roleLabel(role: string): string {
  switch (role) {
    case 'system':
      return t('requestLogs.roleSystem')
    case 'assistant':
      return t('requestLogs.roleAssistant')
    case 'tool':
      return t('requestLogs.roleTool')
    case 'developer':
      return t('requestLogs.roleDeveloper')
    default:
      return t('requestLogs.roleUser')
  }
}

function timelineStepClass(entry: ConversationEntry): string {
  if (entry.role === 'assistant') {
    return 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900/50 dark:bg-emerald-950/30 dark:text-emerald-300'
  }
  if (entry.role === 'tool') {
    return 'border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-900/50 dark:bg-amber-950/30 dark:text-amber-300'
  }
  if (entry.role === 'system' || entry.role === 'developer') {
    return 'border-slate-200 bg-slate-50 text-slate-700 dark:border-slate-800 dark:bg-slate-950/40 dark:text-slate-300'
  }
  return 'border-blue-200 bg-blue-50 text-blue-700 dark:border-blue-900/50 dark:bg-blue-950/30 dark:text-blue-300'
}

function turnStepClass(turn: TurnGroup): string {
  return timelineStepClass({
    id: turn.id,
    role: turn.primaryRole,
    source: 'request',
    protocol: '',
    blocks: []
  })
}

function roleBadgeClass(role: string): string {
  switch (role) {
    case 'system':
      return 'bg-slate-100 text-slate-700 dark:bg-slate-900/50 dark:text-slate-300'
    case 'assistant':
      return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300'
    case 'tool':
      return 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300'
    case 'developer':
      return 'bg-violet-100 text-violet-700 dark:bg-violet-900/40 dark:text-violet-300'
    default:
      return 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-200'
  }
}

function entryCardClass(entry: ConversationEntry): string {
  if (entry.role === 'assistant') {
    return 'border-emerald-200/70 bg-emerald-50/70 dark:border-emerald-900/40 dark:bg-emerald-950/20'
  }
  if (entry.role === 'tool') {
    return 'border-amber-200/70 bg-amber-50/70 dark:border-amber-900/40 dark:bg-amber-950/20'
  }
  if (entry.role === 'system' || entry.role === 'developer') {
    return 'border-slate-200/70 bg-slate-50/70 dark:border-slate-800 dark:bg-slate-950/20'
  }
  return 'border-blue-200/70 bg-blue-50/70 dark:border-blue-900/40 dark:bg-blue-950/20'
}

function turnCardClass(turn: TurnGroup): string {
  return entryCardClass({
    id: turn.id,
    role: turn.primaryRole,
    source: 'request',
    protocol: '',
    blocks: []
  })
}

const formatCost = (value: unknown): string => {
  const amount = Number(value || 0)
  return amount.toFixed(6)
}

const formatDuration = (value: number | null | undefined): string => {
  if (value == null) return t('requestLogs.missingValue')
  if (value < 1000) return `${value}ms`
  return `${(value / 1000).toFixed(2)}s`
}

async function loadRequestLog(): Promise<void> {
  if (abortController) {
    abortController.abort()
  }

  const currentAbortController = new AbortController()
  abortController = currentAbortController
  loading.value = true

  try {
    const id = Number(route.params.id)
    requestLog.value = await requestLogsAPI.getById(id, { signal: currentAbortController.signal })
  } catch (error) {
    const abortError = error as { name?: string; code?: string }
    if (abortError?.name === 'AbortError' || abortError?.code === 'ERR_CANCELED') {
      return
    }
    appStore.showError(t('requestLogs.failedToLoadDetail'))
  } finally {
    if (abortController === currentAbortController) {
      loading.value = false
    }
  }
}

watch(
  () => route.params.id,
  () => {
    loadRequestLog()
  }
)

onMounted(() => {
  loadRequestLog()
})
</script>
