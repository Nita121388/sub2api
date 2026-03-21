<template>
  <AppLayout>
    <TablePageLayout>
      <template #actions>
        <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
          <div class="card p-4">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
              {{ t('common.total') }}
            </p>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
              {{ pagination.total.toLocaleString() }}
            </p>
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {{ t('requestLogs.pageSummary') }}
            </p>
          </div>
          <div class="card p-4">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
              {{ t('requestLogs.currentPageCount') }}
            </p>
            <p class="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
              {{ requestLogs.length.toLocaleString() }}
            </p>
            <p class="mt-1 text-xs text-emerald-600 dark:text-emerald-400">
              {{ t('requestLogs.successCount') }}: {{ successCount.toLocaleString() }}
            </p>
          </div>
          <div class="card p-4">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
              {{ t('requestLogs.errorCount') }}
            </p>
            <p class="mt-2 text-2xl font-bold text-rose-600 dark:text-rose-400">
              {{ errorCount.toLocaleString() }}
            </p>
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {{ t('requestLogs.timeRange') }}
            </p>
          </div>
        </div>
      </template>

      <template #filters>
        <div class="card">
          <div class="px-6 py-4">
            <div class="flex flex-wrap items-end gap-4">
              <div class="min-w-[180px]">
                <label class="input-label">{{ t('requestLogs.apiKeyFilter') }}</label>
                <Select
                  v-model="filters.api_key_id"
                  :options="apiKeyOptions"
                  :placeholder="t('requestLogs.allApiKeys')"
                />
              </div>

              <div class="min-w-[180px]">
                <label class="input-label">{{ t('requestLogs.modelFilter') }}</label>
                <input
                  v-model.trim="filters.model"
                  type="text"
                  class="input"
                  :placeholder="t('requestLogs.modelPlaceholder')"
                  @keyup.enter="applyFilters"
                />
              </div>

              <div class="w-[140px]">
                <label class="input-label">{{ t('requestLogs.statusCode') }}</label>
                <input
                  v-model.trim="statusCodeInput"
                  type="number"
                  min="0"
                  class="input"
                  :placeholder="t('requestLogs.statusCodePlaceholder')"
                  @keyup.enter="applyFilters"
                />
              </div>

              <div>
                <label class="input-label">{{ t('requestLogs.timeRange') }}</label>
                <DateRangePicker
                  v-model:start-date="startDate"
                  v-model:end-date="endDate"
                  @change="onDateRangeChange"
                />
              </div>

              <div class="ml-auto flex items-center gap-3">
                <button @click="applyFilters" :disabled="loading" class="btn btn-secondary">
                  {{ t('common.refresh') }}
                </button>
                <button @click="resetFilters" class="btn btn-secondary">
                  {{ t('common.reset') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </template>

      <template #table>
        <DataTable
          :columns="columns"
          :data="requestLogs"
          :loading="loading"
          row-key="id"
          default-sort-key="created_at"
          default-sort-order="desc"
        >
          <template #cell-api_key_id="{ value }">
            <span class="text-sm text-gray-900 dark:text-white">
              {{ resolveApiKeyLabel(Number(value)) }}
            </span>
          </template>

          <template #cell-model="{ row }">
            <div class="min-w-[180px]">
              <p class="font-medium text-gray-900 dark:text-white">{{ row.model || '-' }}</p>
              <p v-if="row.request_id" class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                {{ row.request_id }}
              </p>
            </div>
          </template>

          <template #cell-endpoint="{ row }">
            <div class="max-w-[280px] space-y-1 whitespace-normal break-all text-xs text-gray-600 dark:text-gray-300">
              <p>{{ row.inbound_endpoint || '-' }}</p>
              <p class="text-gray-400 dark:text-gray-500">{{ row.upstream_endpoint || '-' }}</p>
            </div>
          </template>

          <template #cell-status_code="{ value }">
            <span
              class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium"
              :class="getStatusCodeBadgeClass(value)"
            >
              {{ value ?? '-' }}
            </span>
          </template>

          <template #cell-stream="{ row }">
            <span
              class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium"
              :class="row.stream ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300' : 'bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-gray-300'"
            >
              {{ row.stream ? t('requestLogs.streamLabel') : t('requestLogs.syncLabel') }}
            </span>
          </template>

          <template #cell-tokens="{ row }">
            <div class="space-y-1 text-sm">
              <div class="flex items-center gap-2">
                <span class="text-emerald-600 dark:text-emerald-400">{{ row.input_tokens.toLocaleString() }}</span>
                <span class="text-gray-400">/</span>
                <span class="text-violet-600 dark:text-violet-400">{{ row.output_tokens.toLocaleString() }}</span>
              </div>
            </div>
          </template>

          <template #cell-total_cost="{ value }">
            <span class="font-medium text-gray-900 dark:text-white">${{ formatCost(value) }}</span>
          </template>

          <template #cell-duration_ms="{ row }">
            <div class="space-y-1 text-xs text-gray-600 dark:text-gray-300">
              <p>{{ formatDuration(row.duration_ms) }}</p>
              <p>{{ t('requestLogs.firstToken') }}: {{ formatDuration(row.first_token_ms) }}</p>
            </div>
          </template>

          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-900 dark:text-white">
              {{ formatDateTime(value) }}
            </span>
          </template>

          <template #cell-actions="{ row }">
            <router-link
              :to="{ name: 'RequestLogDetail', params: { id: row.id } }"
              class="btn btn-secondary px-3 py-1 text-xs"
            >
              {{ t('requestLogs.viewDetail') }}
            </router-link>
          </template>

          <template #empty>
            <EmptyState :message="t('requestLogs.noRecords')" />
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          :total="pagination.total"
          :page="pagination.page"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { keysAPI, requestLogsAPI } from '@/api'
import type { ApiKey, RequestLog, RequestLogQueryParams } from '@/types'
import { useAppStore } from '@/stores/app'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Select from '@/components/common/Select.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import type { Column } from '@/components/common/types'
import { formatDateTime } from '@/utils/format'

const { t } = useI18n()
const appStore = useAppStore()

let abortController: AbortController | null = null

const requestLogs = ref<RequestLog[]>([])
const apiKeys = ref<ApiKey[]>([])
const loading = ref(false)
const statusCodeInput = ref('')

const columns = computed<Column[]>(() => [
  { key: 'api_key_id', label: t('requestLogs.apiKeyFilter') },
  { key: 'model', label: t('usage.model'), sortable: true },
  { key: 'endpoint', label: t('usage.endpoint') },
  { key: 'status_code', label: t('requestLogs.statusCode'), sortable: true },
  { key: 'stream', label: t('requestLogs.stream') },
  { key: 'tokens', label: t('requestLogs.tokens') },
  { key: 'total_cost', label: t('requestLogs.cost'), sortable: true },
  { key: 'duration_ms', label: t('requestLogs.duration') },
  { key: 'created_at', label: t('requestLogs.createdAt'), sortable: true },
  { key: 'actions', label: t('requestLogs.actions') }
])

const apiKeyOptions = computed(() => [
  { value: null, label: t('requestLogs.allApiKeys') },
  ...apiKeys.value.map((key) => ({
    value: key.id,
    label: key.name
  }))
])

const formatLocalDate = (date: Date): string => {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

const now = new Date()
const weekAgo = new Date(now)
weekAgo.setDate(weekAgo.getDate() - 6)

const startDate = ref(formatLocalDate(weekAgo))
const endDate = ref(formatLocalDate(now))

const filters = ref<RequestLogQueryParams>({
  start_date: startDate.value,
  end_date: endDate.value
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  pages: 0
})

const successCount = computed(() => requestLogs.value.filter((log) => (log.status_code ?? 0) < 400).length)
const errorCount = computed(() => requestLogs.value.filter((log) => (log.status_code ?? 0) >= 400).length)

const formatCost = (value: unknown): string => {
  const amount = Number(value || 0)
  return amount.toFixed(6)
}

const formatDuration = (value: number | null | undefined): string => {
  if (value == null) return t('requestLogs.missingValue')
  if (value < 1000) return `${value}ms`
  return `${(value / 1000).toFixed(2)}s`
}

const getStatusCodeBadgeClass = (value: unknown): string => {
  const statusCode = Number(value)
  if (!statusCode) return 'bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-gray-300'
  if (statusCode < 400) return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300'
  if (statusCode < 500) return 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300'
  return 'bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-300'
}

const resolveApiKeyLabel = (apiKeyID: number): string => {
  const apiKey = apiKeys.value.find((item) => item.id === apiKeyID)
  return apiKey?.name || `#${apiKeyID}`
}

const loadApiKeys = async () => {
  try {
    const response = await keysAPI.list(1, 100)
    apiKeys.value = response.items
  } catch (error) {
    console.error('Failed to load API keys:', error)
  }
}

const loadRequestLogs = async () => {
  if (abortController) {
    abortController.abort()
  }

  const currentAbortController = new AbortController()
  abortController = currentAbortController
  loading.value = true

  try {
    const response = await requestLogsAPI.query(
      {
        page: pagination.page,
        page_size: pagination.page_size,
        ...filters.value
      },
      { signal: currentAbortController.signal }
    )
    requestLogs.value = response.items
    pagination.total = response.total
    pagination.pages = response.pages
  } catch (error) {
    const abortError = error as { name?: string; code?: string }
    if (abortError?.name === 'AbortError' || abortError?.code === 'ERR_CANCELED') {
      return
    }
    appStore.showError(t('requestLogs.failedToLoad'))
  } finally {
    if (abortController === currentAbortController) {
      loading.value = false
    }
  }
}

const onDateRangeChange = (range: {
  startDate: string
  endDate: string
  preset: string | null
}) => {
  filters.value.start_date = range.startDate
  filters.value.end_date = range.endDate
  applyFilters()
}

const applyFilters = () => {
  if (statusCodeInput.value) {
    const parsed = Number(statusCodeInput.value)
    if (Number.isNaN(parsed)) {
      appStore.showError(t('requestLogs.invalidStatusCode'))
      return
    }
    filters.value.status_code = parsed
  } else {
    delete filters.value.status_code
  }

  if (!filters.value.model) {
    delete filters.value.model
  }

  pagination.page = 1
  loadRequestLogs()
}

const resetFilters = () => {
  const resetNow = new Date()
  const resetWeekAgo = new Date(resetNow)
  resetWeekAgo.setDate(resetWeekAgo.getDate() - 6)

  startDate.value = formatLocalDate(resetWeekAgo)
  endDate.value = formatLocalDate(resetNow)
  statusCodeInput.value = ''
  filters.value = {
    start_date: startDate.value,
    end_date: endDate.value
  }
  pagination.page = 1
  loadRequestLogs()
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadRequestLogs()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.page_size = pageSize
  pagination.page = 1
  loadRequestLogs()
}

onMounted(() => {
  loadApiKeys()
  loadRequestLogs()
})
</script>
