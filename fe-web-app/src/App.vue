<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { ref, computed } from 'vue'

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
  return (data.value?.lines ?? []).reduce<string[]>((prev, curr) => {
    if (dateRegex.test(curr)) return [...prev, curr]

    const merged = prev.at(-1)?.concat(`\n${curr}`) ?? curr

    return [...prev.slice(0, prev.length - 1), merged]
  }, [])
})

const search = ref('')

const filteredLines = computed(() => {
  const trimmedSearch = search.value.trim()

  if (!trimmedSearch) return mergedLines.value

  return mergedLines.value.filter((line) => line.includes(trimmedSearch))
})
</script>

<template>
  <main class="px-4 py-2 flex flex-col h-screen items-start gap-2">
    <h1 class="text-2xl">Sidero Logs</h1>

    <input type="text" class="border-red-500 border-2 rounded-md" v-model="search" />
    <p>{{ search }}</p>

    <pre class="border-gray-400 border-1 rounded-md px-4 py-2 w-full flex-1 overflow-y-auto">{{
      filteredLines.join('\n')
    }}</pre>
  </main>
</template>
