import commonjs from '@rollup/plugin-commonjs'
import json from '@rollup/plugin-json'
import resolve from '@rollup/plugin-node-resolve'
import typescript from '@rollup/plugin-typescript'

export default {
  input: 'server/index.ts',
  output: {
    file: 'distServer/index.js',
    inlineDynamicImports: true
  },
  plugins: [
    typescript({ tsconfig: './tsconfig.server.json' }),
    resolve({ preferBuiltins: true }),
    commonjs(),
    json()
  ],
  external: ['path', 'fs']
}
