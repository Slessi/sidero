<script setup lang="ts">
import DateFilter from '@/components/date-filter/DateFilter.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { getLocalTimeZone } from '@internationalized/date'
import { useQuery } from '@tanstack/vue-query'
import { endOfDay, isWithinInterval, parse, startOfDay, type Interval } from 'date-fns'
import type { DateRange } from 'reka-ui'
import { computed, ref } from 'vue'

const { data, error, isLoading, refetch } = useQuery({
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

const logOutput = computed(() => filteredLines.value.map((l) => l.text).join('\n'))
const hasFilters = computed(() => search || dateRange.value.start || dateRange.value.end)

function clearFilters() {
  search.value = ''
  dateRange.value = { start: undefined, end: undefined }
}
</script>

<template>
  <main class="px-16 py-4 flex flex-col h-screen items-start gap-4">
    <h1 class="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4 mx-auto">
      Sidero Logs
    </h1>

    <div class="flex gap-2">
      <DateFilter v-model="dateRange" />
      <Input type="text" placeholder="Search" v-model="search" />
      <Button v-if="hasFilters" variant="ghost" @click="clearFilters">Clear filters</Button>
    </div>

    <div
      class="flex flex-col items-center gap-2 dark:bg-accent/30 border rounded-md px-4 py-2 w-full flex-1 bg-transparent overflow-auto scrollbar-thin scrollbar-thumb-accent-foreground scrollbar-track-[rgba(0,0,0,0)]"
    >
      <p :class="{ 'text-red-800': error }" v-if="!logOutput">
        <template v-if="isLoading">Loading ...</template>
        <template v-else-if="error">{{ error.message }}</template>
        <template v-else="hasFilters">No results</template>
      </p>

      <Button @click="refetch" variant="outline" v-if="error">Retry</Button>

      <pre class="text-xs w-full">{{ logOutput }}</pre>
    </div>
  </main>
</template>
