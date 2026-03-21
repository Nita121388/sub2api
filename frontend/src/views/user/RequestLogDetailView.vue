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
          <button @click="router.back()" class="btn btn-secondary">
            {{ t('common.back') }}
          </button>
          <button @click="loadRequestLog" :disabled="loading" class="btn btn-secondary">
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
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { requestLogsAPI } from '@/api'
import type { RequestLog } from '@/types'
import { useAppStore } from '@/stores/app'
import AppLayout from '@/components/layout/AppLayout.vue'
import { formatDateTime } from '@/utils/format'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()

let abortController: AbortController | null = null

const requestLog = ref<RequestLog | null>(null)
const loading = ref(false)
const missingValue = computed(() => t('requestLogs.missingValue'))
const totalTokens = computed(() => (requestLog.value?.input_tokens || 0) + (requestLog.value?.output_tokens || 0))

const formatCost = (value: number | null | undefined): string => {
  return Number(value || 0).toFixed(6)
}

const formatDuration = (value: number | null | undefined): string => {
  if (value == null) return missingValue.value
  if (value < 1000) return `${value}ms`
  return `${(value / 1000).toFixed(2)}s`
}

const loadRequestLog = async () => {
  const id = Number(route.params.id)
  if (!Number.isFinite(id) || id <= 0) {
    appStore.showError(t('requestLogs.failedToLoadDetail'))
    return
  }

  if (abortController) {
    abortController.abort()
  }

  const currentAbortController = new AbortController()
  abortController = currentAbortController
  loading.value = true

  try {
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
