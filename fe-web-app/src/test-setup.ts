import { VueQueryPlugin } from '@tanstack/vue-query'
import { config } from '@vue/test-utils'

config.global.plugins = [VueQueryPlugin]
