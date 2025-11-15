import json from '@rollup/plugin-json';
import { nodeResolve } from '@rollup/plugin-node-resolve';
import replace from '@rollup/plugin-replace';
import terser from '@rollup/plugin-terser';
import typescript from '@rollup/plugin-typescript';
import copy from 'rollup-plugin-copy';
import del from 'rollup-plugin-delete';
import css from 'rollup-plugin-import-css';

const isDev = process.env.BUILD === 'dev';
const isProduction = process.env.BUILD === 'production';

export default [
	{
		input: './src/client/ts/application.ts',
		output: {
			file: './build/client/js/application.js',
			format: 'esm'
		},
		plugins: [
			replace({
				preventAssignment: true,
				TESTING: isDev,
			}),
			css(),
			isProduction ? del({ targets: 'build/*' }) : null,
			json({
				compact: true,
			}),
			typescript(),
			nodeResolve({
				dedupe: ['harmony-ui', 'harmony-browser-utils'],
			}),
			isProduction ? terser() : null,
			copy({
				copyOnce: true,
				targets: [
					{ src: 'src/client/index.html', dest: 'build/client/' },
					{ src: 'src/client/img/', dest: 'build/client/' },
				]
			}),
		],
	},
];
