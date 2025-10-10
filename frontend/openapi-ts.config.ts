import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  client: {
    bundle: true,
  },
  input: '../schemas/openapi.yaml',
  output: './src/api',
  schemas: {
    export: true,
  },
  services: {
    export: true,
    name: '{{name}}Service',
  },
  types: {
    export: true,
  },
  // baseURLは実行時に設定するため、ここでは指定しない
});
