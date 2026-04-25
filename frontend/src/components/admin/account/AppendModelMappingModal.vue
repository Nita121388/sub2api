<template>
  <BaseDialog
    :show="show"
    :title="t('admin.accounts.appendModelMapping.title')"
    width="wide"
    close-on-click-outside
    @close="handleClose"
  >
    <form id="append-model-mapping-form" class="space-y-5" @submit.prevent="handlePreview">
      <div class="grid gap-3 md:grid-cols-[1fr_auto]">
        <div>
          <label class="input-label">{{ t('admin.accounts.appendModelMapping.modelFamily') }}</label>
          <input
            v-model="modelFamily"
            type="text"
            class="input font-mono"
            :placeholder="t('admin.accounts.appendModelMapping.modelFamilyPlaceholder')"
          />
        </div>
        <div class="flex items-end">
          <button type="button" class="btn btn-secondary w-full md:w-auto" @click="generatePresetRows">
            {{ t('admin.accounts.appendModelMapping.generatePreset') }}
          </button>
        </div>
      </div>

      <div class="rounded-lg border border-gray-200 dark:border-dark-700">
        <div class="grid grid-cols-[1fr_1fr_40px] gap-2 border-b border-gray-100 px-3 py-2 text-xs font-medium text-gray-500 dark:border-dark-700 dark:text-dark-400">
          <span>{{ t('admin.accounts.appendModelMapping.requestModel') }}</span>
          <span>{{ t('admin.accounts.appendModelMapping.targetModel') }}</span>
          <span></span>
        </div>
        <div class="space-y-2 p-3">
          <div
            v-for="(row, index) in mappingRows"
            :key="row.id"
            class="grid grid-cols-[1fr_1fr_40px] gap-2"
          >
            <input
              v-model="row.from"
              type="text"
              class="input font-mono text-sm"
              :placeholder="t('admin.accounts.appendModelMapping.requestModelPlaceholder')"
            />
            <input
              v-model="row.to"
              type="text"
              class="input font-mono text-sm"
              :placeholder="row.from || t('admin.accounts.appendModelMapping.targetModelPlaceholder')"
            />
            <button
              type="button"
              class="btn btn-secondary px-2"
              :disabled="mappingRows.length === 1"
              :title="t('common.delete')"
              @click="removeRow(index)"
            >
              <Icon name="trash" size="sm" />
            </button>
          </div>
          <button type="button" class="btn btn-secondary btn-sm" @click="addRow">
            {{ t('admin.accounts.appendModelMapping.addRow') }}
          </button>
        </div>
      </div>

      <div class="grid gap-3 md:grid-cols-2">
        <label class="flex items-start gap-3 rounded-lg border border-gray-200 p-3 dark:border-dark-700">
          <input v-model="scope" type="radio" value="selected" class="mt-1" :disabled="selectedIds.length === 0" />
          <span>
            <span class="block text-sm font-medium text-gray-900 dark:text-white">
              {{ t('admin.accounts.appendModelMapping.scopeSelected', { count: selectedIds.length }) }}
            </span>
            <span class="block text-xs text-gray-500 dark:text-dark-400">
              {{ t('admin.accounts.appendModelMapping.scopeSelectedHint') }}
            </span>
          </span>
        </label>
        <label class="flex items-start gap-3 rounded-lg border border-gray-200 p-3 dark:border-dark-700">
          <input v-model="scope" type="radio" value="all" class="mt-1" />
          <span>
            <span class="block text-sm font-medium text-gray-900 dark:text-white">
              {{ t('admin.accounts.appendModelMapping.scopeAll') }}
            </span>
            <span class="block text-xs text-gray-500 dark:text-dark-400">
              {{ t('admin.accounts.appendModelMapping.scopeAllHint') }}
            </span>
          </span>
        </label>
      </div>

      <div class="space-y-2">
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-dark-200">
          <input v-model="onlyWithExistingMapping" type="checkbox" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
          <span>{{ t('admin.accounts.appendModelMapping.onlyExisting') }}</span>
        </label>
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-dark-200">
          <input v-model="overwriteExisting" type="checkbox" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
          <span>{{ t('admin.accounts.appendModelMapping.overwriteExisting') }}</span>
        </label>
      </div>

      <div
        v-if="result"
        class="space-y-3 rounded-lg border border-gray-200 p-4 dark:border-dark-700"
      >
        <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
          <div class="rounded-md bg-gray-50 p-3 dark:bg-dark-800">
            <div class="text-xs text-gray-500 dark:text-dark-400">{{ t('common.total') }}</div>
            <div class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">{{ result.total }}</div>
          </div>
          <div class="rounded-md bg-emerald-50 p-3 dark:bg-emerald-900/20">
            <div class="text-xs text-emerald-700 dark:text-emerald-300">{{ t('admin.accounts.appendModelMapping.changed') }}</div>
            <div class="mt-1 text-lg font-semibold text-emerald-700 dark:text-emerald-300">{{ result.changed }}</div>
          </div>
          <div class="rounded-md bg-amber-50 p-3 dark:bg-amber-900/20">
            <div class="text-xs text-amber-700 dark:text-amber-300">{{ t('admin.accounts.appendModelMapping.skipped') }}</div>
            <div class="mt-1 text-lg font-semibold text-amber-700 dark:text-amber-300">{{ result.skipped }}</div>
          </div>
          <div class="rounded-md bg-red-50 p-3 dark:bg-red-900/20">
            <div class="text-xs text-red-700 dark:text-red-300">{{ t('admin.accounts.appendModelMapping.conflicted') }}</div>
            <div class="mt-1 text-lg font-semibold text-red-700 dark:text-red-300">{{ result.conflicted }}</div>
          </div>
        </div>

        <div class="max-h-64 overflow-auto rounded-md bg-gray-50 p-3 text-xs dark:bg-dark-800">
          <div v-for="item in visibleResults" :key="item.account_id" class="grid gap-1 border-b border-gray-200 py-2 last:border-b-0 dark:border-dark-700 md:grid-cols-[90px_1fr_120px]">
            <span class="font-mono text-gray-500">#{{ item.account_id }}</span>
            <span class="truncate text-gray-700 dark:text-dark-200">{{ item.account_name || '-' }}</span>
            <span :class="resultClass(item.action)">{{ resultLabel(item) }}</span>
          </div>
          <div v-if="result.results.length > visibleResults.length" class="pt-2 text-gray-500">
            {{ t('admin.accounts.appendModelMapping.moreResults', { count: result.results.length - visibleResults.length }) }}
          </div>
        </div>
      </div>
    </form>

    <template #footer>
      <div class="flex flex-wrap justify-end gap-3">
        <button class="btn btn-secondary" type="button" :disabled="working" @click="handleClose">
          {{ t('common.cancel') }}
        </button>
        <button class="btn btn-secondary" type="submit" form="append-model-mapping-form" :disabled="working">
          {{ working ? t('common.processing') : t('admin.accounts.appendModelMapping.preview') }}
        </button>
        <button
          class="btn btn-primary"
          type="button"
          :disabled="working || !result || result.changed === 0"
          @click="handleApply"
        >
          {{ working ? t('common.processing') : t('admin.accounts.appendModelMapping.apply') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import { adminAPI } from '@/api/admin'
import { useAppStore } from '@/stores/app'
import type { AppendModelMappingAccountResult, AppendModelMappingsResult } from '@/api/admin/accounts'

const props = defineProps<{
  show: boolean
  selectedIds: number[]
}>()

const emit = defineEmits<{
  close: []
  updated: []
}>()

const { t } = useI18n()
const appStore = useAppStore()

type Scope = 'selected' | 'all'
type MappingRow = { id: number; from: string; to: string }

const working = ref(false)
const scope = ref<Scope>('all')
const modelFamily = ref('gpt-5.5')
const mappingRows = ref<MappingRow[]>([])
const overwriteExisting = ref(false)
const onlyWithExistingMapping = ref(true)
const result = ref<AppendModelMappingsResult | null>(null)
let nextRowID = 1

const visibleResults = computed(() => result.value?.results.slice(0, 30) ?? [])

const makeRow = (from = '', to = ''): MappingRow => ({ id: nextRowID++, from, to })

const generateRowsForModelFamily = (family: string): MappingRow[] => {
  const base = family.trim()
  if (!base) return [makeRow()]
  const pro = `${base}-pro`
  return [
    makeRow(base, base),
    makeRow(`${base}-none`, base),
    makeRow(`${base}-low`, base),
    makeRow(`${base}-medium`, base),
    makeRow(`${base}-high`, base),
    makeRow(`${base}-xhigh`, base),
    makeRow(`${base}-chat-latest`, base),
    makeRow(pro, pro),
    makeRow(`${pro}-none`, pro),
    makeRow(`${pro}-low`, pro),
    makeRow(`${pro}-medium`, pro),
    makeRow(`${pro}-high`, pro),
    makeRow(`${pro}-xhigh`, pro)
  ]
}

const resetState = () => {
  scope.value = props.selectedIds.length > 0 ? 'selected' : 'all'
  modelFamily.value = 'gpt-5.5'
  mappingRows.value = generateRowsForModelFamily(modelFamily.value)
  overwriteExisting.value = false
  onlyWithExistingMapping.value = true
  result.value = null
}

watch(
  () => props.show,
  (open) => {
    if (open) resetState()
  },
  { immediate: true }
)

watch([scope, overwriteExisting, onlyWithExistingMapping], () => {
  result.value = null
})

watch(
  mappingRows,
  () => {
    result.value = null
  },
  { deep: true }
)

const generatePresetRows = () => {
  mappingRows.value = generateRowsForModelFamily(modelFamily.value)
  result.value = null
}

const addRow = () => {
  mappingRows.value.push(makeRow())
  result.value = null
}

const removeRow = (index: number) => {
  if (mappingRows.value.length <= 1) return
  mappingRows.value.splice(index, 1)
  result.value = null
}

const buildMappings = (): Record<string, string> => {
  const out: Record<string, string> = {}
  for (const row of mappingRows.value) {
    const from = row.from.trim()
    const to = row.to.trim() || from
    if (!from) continue
    out[from] = to
  }
  return out
}

const buildPayload = (dryRun: boolean) => ({
  platform: 'openai',
  account_ids: scope.value === 'selected' ? props.selectedIds : undefined,
  mappings: buildMappings(),
  dry_run: dryRun,
  overwrite_existing: overwriteExisting.value,
  only_with_existing_mapping: onlyWithExistingMapping.value
})

const handlePreview = async () => {
  const mappings = buildMappings()
  if (Object.keys(mappings).length === 0) {
    appStore.showError(t('admin.accounts.appendModelMapping.emptyMappings'))
    return
  }
  if (scope.value === 'selected' && props.selectedIds.length === 0) {
    appStore.showError(t('admin.accounts.appendModelMapping.noSelectedAccounts'))
    return
  }

  working.value = true
  try {
    result.value = await adminAPI.accounts.appendModelMappings(buildPayload(true))
  } catch (error: any) {
    appStore.showError(error?.message || t('admin.accounts.appendModelMapping.previewFailed'))
  } finally {
    working.value = false
  }
}

const handleApply = async () => {
  const mappings = buildMappings()
  if (Object.keys(mappings).length === 0) return

  working.value = true
  try {
    result.value = await adminAPI.accounts.appendModelMappings(buildPayload(false))
    appStore.showSuccess(t('admin.accounts.appendModelMapping.applySuccess', { count: result.value.changed }))
    emit('updated')
  } catch (error: any) {
    appStore.showError(error?.message || t('admin.accounts.appendModelMapping.applyFailed'))
  } finally {
    working.value = false
  }
}

const handleClose = () => {
  if (working.value) return
  emit('close')
}

const resultClass = (action: string) => {
  if (action === 'changed') return 'text-emerald-700 dark:text-emerald-300'
  if (action === 'conflict') return 'text-red-700 dark:text-red-300'
  return 'text-gray-500 dark:text-dark-400'
}

const resultLabel = (item: AppendModelMappingAccountResult) => {
  if (item.action === 'changed') {
    return t('admin.accounts.appendModelMapping.resultChanged', {
      count: Object.keys(item.added || {}).length + Object.keys(item.overwritten || {}).length
    })
  }
  if (item.action === 'conflict') {
    return t('admin.accounts.appendModelMapping.resultConflict', { count: item.conflicts?.length || 0 })
  }
  if (item.action === 'unchanged') {
    return t('admin.accounts.appendModelMapping.resultUnchanged')
  }
  return t(`admin.accounts.appendModelMapping.reason.${item.reason || 'skipped'}`)
}
</script>
