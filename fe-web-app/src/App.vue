<script setup lang="ts">
import { Input } from '@/components/ui/input'
import { RangeCalendar } from '@/components/ui/range-calendar'
import { getLocalTimeZone } from '@internationalized/date'
import { useQuery } from '@tanstack/vue-query'
import { endOfDay, isWithinInterval, parse, startOfDay } from 'date-fns'
import type { DateRange } from 'reka-ui'
import { computed, ref } from 'vue'

const { data } = useQuery({
  queryKey: ['logs'],
  async queryFn() {
    const response = await fetch('//localhost:12000/v1/logs')

    const data = await response.json()
    return data as { lines: string[] }
  },
})

const dateRegex = /^\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}:\d{2}/

// This will merge lines of the same log entry into the same line
const mergedLines = computed(() => {
  return (data.value?.lines ?? []).reduce<{ date: Date; text: string }[]>((prev, curr) => {
    const dateMatch = curr.match(dateRegex)
    if (dateMatch) {
      return [
        ...prev,
        {
          date: parse(dateMatch[0], 'dd/MM/yyyy HH:mm:ss', new Date()),
          text: curr,
        },
      ]
    }

    const lastEntry = prev.at(-1)
    if (!lastEntry) return prev

    return [
      ...prev.slice(0, prev.length - 1),
      {
        ...lastEntry,
        text: `${lastEntry.text}\n${curr}`,
      },
    ]
  }, [])
})

const search = ref('')
const dateRange = ref<DateRange>({ start: undefined, end: undefined })

const filteredLines = computed(() => {
  const trimmedSearch = search.value.trim()
  const { start, end } = dateRange.value

  if (!trimmedSearch && !start && !end) return mergedLines.value

  const dateInterval = {
    start: startOfDay((start || end)?.toDate(getLocalTimeZone()) || new Date()),
    end: endOfDay((end || start)?.toDate(getLocalTimeZone()) || new Date()),
  }

  return mergedLines.value.filter(
    (l) => l.text.includes(trimmedSearch) && isWithinInterval(l.date, dateInterval),
  )
})
</script>

<template>
  <main class="px-4 py-2 flex flex-col h-screen items-start gap-2">
    <h1 class="text-2xl">Sidero Logs</h1>

    <RangeCalendar v-model="dateRange" class="rounded-md border" />
    <Input type="text" placeholder="Search" v-model="search" />

    <p>{{ dateRange.start?.toString() }}</p>
    <p>{{ dateRange.end?.toString() }}</p>

    <pre class="border-gray-400 border-1 rounded-md px-4 py-2 w-full flex-1 overflow-y-auto">{{
      filteredLines.map((l) => l.text).join('\n')
    }}</pre>
  </main>
</template>
