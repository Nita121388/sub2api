import { describe, expect, it, vi, beforeEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { nextTick } from 'vue'

import RequestLogDetailView from '../RequestLogDetailView.vue'

const { getById, showError } = vi.hoisted(() => {
  const storage = new Map<string, string>()
  vi.stubGlobal('localStorage', {
    getItem: vi.fn((key: string) => storage.get(key) ?? null),
    setItem: vi.fn((key: string, value: string) => {
      storage.set(key, value)
    }),
    removeItem: vi.fn((key: string) => {
      storage.delete(key)
    }),
    clear: vi.fn(() => {
      storage.clear()
    }),
  })

  return {
    getById: vi.fn(),
    showError: vi.fn(),
  }
})

vi.mock('@/api', () => ({
  requestLogsAPI: {
    getById,
  },
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError,
  }),
}))

vi.mock('vue-router', () => ({
  useRoute: () => ({
    params: {
      id: '290',
    },
  }),
  useRouter: () => ({
    back: vi.fn(),
  }),
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key,
    }),
  }
})

describe('RequestLogDetailView', () => {
  beforeEach(() => {
    getById.mockReset()
    showError.mockReset()
  })

  it('visualizes responses request input when payload.input is a string', async () => {
    getById.mockResolvedValue({
      id: 290,
      request_id: 'req_290',
      model: 'gpt-5.4',
      input_tokens: 12,
      output_tokens: 4,
      total_cost: 0.01,
      stream: false,
      created_at: '2026-03-24T07:00:00Z',
      payload: {
        request_body: JSON.stringify({
          model: 'gpt-5.4',
          input: 'hello from api key',
        }),
        response_body: JSON.stringify({
          id: 'resp_290',
          output_text: 'ok',
        }),
      },
    })

    const wrapper = mount(RequestLogDetailView, {
      global: {
        stubs: {
          AppLayout: {
            template: '<div><slot /></div>',
          },
        },
      },
    })

    await flushPromises()
    await nextTick()

    const setupState = (wrapper.vm as any).$?.setupState
    expect(setupState.requestConversation).toHaveLength(1)
    expect(setupState.requestConversation[0].blocks[0].text).toBe('hello from api key')
    expect(showError).not.toHaveBeenCalled()
  })
})
