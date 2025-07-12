import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import App from './App.vue'
import { CalendarDate } from '@internationalized/date'

const mockLogs = {
  lines: [
    '01/01/2024 12:00:00 Log entry one',
    '02/01/2024 12:01:00 Log entry two',
    '03/01/2024 12:02:00 Log entry three',
  ],
}

vi.stubGlobal(
  'fetch',
  vi.fn(() => Promise.resolve({ json: () => Promise.resolve(mockLogs) })),
)

describe('App.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('filters logs by search', async () => {
    const wrapper = mount(App)
    await flushPromises()

    const input = wrapper.find('input[type="text"]')
    await input.setValue('entry two')
    await flushPromises()

    expect(wrapper.text()).toContain('Log entry two')
    expect(wrapper.text()).not.toContain('Log entry one')
    expect(wrapper.text()).not.toContain('Log entry three')
  })

  it('filters logs by date range', async () => {
    const wrapper = mount(App)
    await flushPromises()

    // Initially all logs should be visible
    expect(wrapper.text()).toContain('Log entry one')
    expect(wrapper.text()).toContain('Log entry two')
    expect(wrapper.text()).toContain('Log entry three')

    // Set date range to filter logs from 01/01/2024 to 02/01/2024
    const vm = wrapper.vm as any
    vm.dateRange = {
      start: new CalendarDate(2024, 1, 1),
      end: new CalendarDate(2024, 1, 2),
    }
    await flushPromises()

    // Should show logs from 01/01/2024 and 02/01/2024
    expect(wrapper.text()).toContain('Log entry one')
    expect(wrapper.text()).toContain('Log entry two')
    // Should not show logs from 03/01/2024
    expect(wrapper.text()).not.toContain('Log entry three')
  })

  it('combines search and date range filters', async () => {
    const wrapper = mount(App)
    await flushPromises()

    // Set search filter
    const input = wrapper.find('input[type="text"]')
    await input.setValue('entry')
    await flushPromises()

    // Set date range to filter logs from 01/01/2024 to 01/01/2024
    const vm = wrapper.vm as any
    vm.dateRange = {
      start: new CalendarDate(2024, 1, 1),
      end: new CalendarDate(2024, 1, 1),
    }
    await flushPromises()

    // Should only show logs from 01/01/2024 that contain "entry"
    expect(wrapper.text()).toContain('Log entry one')
    // Should not show logs from other dates
    expect(wrapper.text()).not.toContain('Log entry two')
    expect(wrapper.text()).not.toContain('Log entry three')
  })

  it('shows "No results" when filter matches nothing', async () => {
    const wrapper = mount(App)
    await flushPromises()

    const input = wrapper.find('input[type="text"]')
    await input.setValue('no-such-log')
    await flushPromises()

    expect(wrapper.text()).toContain('No results')
  })
})
