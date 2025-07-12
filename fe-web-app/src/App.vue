<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { RangeCalendar } from '@/components/ui/range-calendar'
import { getLocalTimeZone } from '@internationalized/date'
import { useQuery } from '@tanstack/vue-query'
import { endOfDay, isWithinInterval, parse, startOfDay, type Interval } from 'date-fns'
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
const popoverOpen = ref(false)
const dateRange = ref<DateRange>({ start: undefined, end: undefined })

const filteredLines = computed(() => {
  const trimmedSearch = search.value.trim()
  const { start, end } = dateRange.value

  if (!trimmedSearch && !start && !end) return mergedLines.value

  let dateInterval: Interval | undefined
  if (start || end) {
    dateInterval = {
      start: startOfDay((start || end)?.toDate(getLocalTimeZone()) || new Date()),
      end: endOfDay((end || start)?.toDate(getLocalTimeZone()) || new Date()),
    }
  }

  return mergedLines.value.filter(
    (l) =>
      l.text.includes(trimmedSearch) && (!dateInterval || isWithinInterval(l.date, dateInterval)),
  )
})

function onDateChange() {
  popoverOpen.value = false
}

function clearFilters() {
  search.value = ''
  dateRange.value = { start: undefined, end: undefined }
}
</script>

<template>
  <main class="px-4 py-2 flex flex-col h-screen items-start gap-2">
    <h1 class="text-2xl">Sidero Logs</h1>

    <div class="flex gap-2">
      <Popover v-model:open="popoverOpen">
        <PopoverTrigger as-child>
          <Button variant="outline">
            <template v-if="!dateRange.start && !dateRange.end">Date range</template>
            <template v-else-if="!dateRange.start!.compare(dateRange.end!)">{{
              dateRange.start
            }}</template>
            <template v-else>{{ dateRange.start }} - {{ dateRange.end }}</template>
          </Button>
        </PopoverTrigger>

        <PopoverContent class="w-80">
          <RangeCalendar
            v-model="dateRange"
            @update:model-value="onDateChange"
            class="rounded-md border"
          />
        </PopoverContent>
      </Popover>

      <Input type="text" placeholder="Search" v-model="search" />

      <Button
        v-if="search || dateRange.start || dateRange.end"
        variant="ghost"
        @click="clearFilters"
        >Clear filters</Button
      >
    </div>

    <pre class="border-gray-400 border-1 rounded-md px-4 py-2 w-full flex-1 overflow-y-auto">{{
      filteredLines.map((l) => l.text).join('\n')
    }}</pre>
  </main>
</template>
