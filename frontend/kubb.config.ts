import { defineConfig } from '@kubb/core'
import { pluginSwagger } from '@kubb/swagger'
import { pluginSwaggerClient } from '@kubb/swagger-client'
import { pluginTanstackQuery } from '@kubb/swagger-tanstack-query'

export default defineConfig({
  input: {
    path: '../backend/docs/swagger.json',
  },
  output: {
    path: './src/generated',
  },
  plugins: [
    pluginSwagger(),
    pluginSwaggerClient(),
    pluginTanstackQuery({
      output: {
        path: './src/generated/hooks',
      },
    }),
  ],
})
