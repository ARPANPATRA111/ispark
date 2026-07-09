import { existsSync } from 'node:fs';
import { spawnSync } from 'node:child_process';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const scriptDir = dirname(fileURLToPath(import.meta.url));
const repoRoot = resolve(scriptDir, '..', '..');
const huskyBin = resolve(repoRoot, 'web', 'node_modules', 'husky', 'bin.js');

if (!existsSync(resolve(repoRoot, '.git')) || !existsSync(resolve(repoRoot, '.husky'))) {
	process.exit(0);
}

if (!existsSync(huskyBin)) {
	console.log('Skipping root Husky hook install: web/node_modules/husky is not available yet.');
	process.exit(0);
}

const result = spawnSync(process.execPath, [huskyBin], {
	cwd: repoRoot,
	stdio: 'inherit'
});

if (result.error) {
	console.error(result.error.message);
	process.exit(1);
}

process.exit(result.status ?? 1);
