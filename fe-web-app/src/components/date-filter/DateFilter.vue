<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { RangeCalendar } from '@/components/ui/range-calendar'
import type { DateRange } from 'reka-ui'
import { ref } from 'vue'

const dateRange = defineModel<DateRange>()
const popoverOpen = ref(false)

function onDateChange() {
  popoverOpen.value = false
}
</script>

<template>
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

    <PopoverContent class="p-0 w-auto">
      <RangeCalendar v-model="dateRange" @update:model-value="onDateChange" />
    </PopoverContent>
  </Popover>
</template>
