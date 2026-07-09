import { existsSync } from 'node:fs';
import { spawnSync } from 'node:child_process';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const scriptDir = dirname(fileURLToPath(import.meta.url));
const webRoot = resolve(scriptDir, '..');
const svelteKitBin = resolve(webRoot, 'node_modules', '.bin', 'svelte-kit');
const svelteKitCommand = process.platform === 'win32' ? 'npx.cmd' : 'npx';

if (existsSync(svelteKitBin) || existsSync(`${svelteKitBin}.cmd`)) {
	const sync = spawnSync(svelteKitCommand, ['--no-install', 'svelte-kit', 'sync'], {
		cwd: webRoot,
		stdio: 'inherit',
		shell: process.platform === 'win32'
	});

	if (sync.error) {
		console.warn(sync.error.message);
	}
}

const hooks = spawnSync(process.execPath, [resolve(scriptDir, 'install-root-hooks.mjs')], {
	cwd: webRoot,
	stdio: 'inherit',
	shell: false
});

if (hooks.error) {
	console.error(hooks.error.message);
	process.exit(1);
}

process.exit(hooks.status ?? 1);
