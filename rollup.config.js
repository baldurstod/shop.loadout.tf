import json from '@rollup/plugin-json';
import {nodeResolve} from '@rollup/plugin-node-resolve';
import replace from '@rollup/plugin-replace';
import terser from '@rollup/plugin-terser';
import copy from 'rollup-plugin-copy';
import del from 'rollup-plugin-delete';
import css from 'rollup-plugin-import-css';

const isDev = process.env.BUILD === 'dev';
const isProduction = process.env.BUILD === 'production';

export default [
	{
		input: './src/client/js/application.js',
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
			isProduction ? del({targets: 'build/*'}) : null,
			json({
				compact: true,
			}),
			nodeResolve({
				dedupe: ['harmony-ui', 'harmony-browser-utils'],
			}),
			isProduction ? terser() : null,
			copy({
				copyOnce: true,
				targets: [
					{src: 'src/client/index.html', dest: 'build/client/'},
				]
			}),
		],
	},
];
