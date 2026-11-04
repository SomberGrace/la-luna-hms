import { spawnSync } from 'node:child_process';
import { cp, mkdir, readdir } from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const frontendDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const roles = ['guest', 'frontdesk', 'kitchen', 'manager', 'restaurant', 'bar'];
const args = new Set(process.argv.slice(2));
const selectedRole = process.argv.slice(2).find((arg) => !arg.startsWith('--'));
const syncOnly = args.has('--sync-only');
const release = args.has('--release');

if (release) {
	const requiredSigning = [
		'ANDROID_KEYSTORE_PATH',
		'ANDROID_KEYSTORE_PASSWORD',
		'ANDROID_KEY_ALIAS',
		'ANDROID_KEY_PASSWORD'
	];
	const missing = requiredSigning.filter((key) => !process.env[key]);
	if (missing.length) {
		throw new Error(`Release signing is not configured. Set: ${missing.join(', ')}`);
	}
}

function run(command, commandArgs, env, cwd = frontendDir) {
	const result = spawnSync(command, commandArgs, {
		cwd,
		env,
		stdio: 'inherit',
		shell: process.platform === 'win32'
	});

	if (result.error) {
		throw result.error;
	}
	if (result.status !== 0) {
		throw new Error(`${command} ${commandArgs.join(' ')} failed with exit code ${result.status}`);
	}
}

function validateBackendUrl(value) {
	if (!value) {
		throw new Error(
			'Set VITE_API_BASE_URL to the deployed HTTPS API origin before building mobile apps.'
		);
	}
	let url;
	try {
		url = new URL(value);
	} catch {
		throw new Error('VITE_API_BASE_URL must be a valid HTTPS URL.');
	}
	if (url.protocol !== 'https:' || url.username || url.password || url.search || url.hash) {
		throw new Error(
			'VITE_API_BASE_URL must be an HTTPS origin without credentials, query, or fragment.'
		);
	}
	return url.origin + url.pathname.replace(/\/+$/, '');
}

async function verifyUnifiedDatabaseApi(apiBaseUrl) {
	const healthUrl = new URL('api/v1/health', `${apiBaseUrl}/`);
	let response;
	try {
		response = await fetch(healthUrl, {
			signal: AbortSignal.timeout(90000),
			headers: {
				'User-Agent':
					'Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36 Chrome/125.0 Mobile Safari/537.36',
				'ngrok-skip-browser-warning': 'true',
				Accept: 'application/json'
			}
		});
	} catch {
		throw new Error(
			'The configured API is unreachable. Mobile builds require the deployed unified API.'
		);
	}
	const contentType = response.headers.get('content-type') ?? '';
	if (!contentType.toLowerCase().includes('application/json')) {
		throw new Error(
			`The configured API returned non-JSON content (${contentType || 'no content type'}). The tunnel may be serving ngrok's browser warning page.`
		);
	}
	if (!response.ok) {
		throw new Error('The configured API does not expose the unified HMS health endpoint.');
	}

	let health;
	try {
		health = await response.json();
	} catch {
		throw new Error('The configured API returned an invalid health response.');
	}
	if (health.status !== 'ok' || health.service !== 'La Luna HMS Unified API') {
		throw new Error(
			'VITE_API_BASE_URL must point to the unified HMS API, not the standalone front-desk service.'
		);
	}
	if (health.database !== 'Connected to Supabase PostgreSQL') {
		throw new Error('The unified HMS API is not connected to its PostgreSQL database.');
	}
}

function findArtifact(directory, extension) {
	return readdir(directory, { withFileTypes: true }).then(async (entries) => {
		for (const entry of entries) {
			const entryPath = path.join(directory, entry.name);
			if (entry.isDirectory()) {
				const found = await findArtifact(entryPath, extension).catch(() => null);
				if (found) return found;
			} else if (entry.name.endsWith(extension)) {
				return entryPath;
			}
		}
		return null;
	});
}

async function buildRole(role, apiBaseUrl) {
	console.log(`\nBuilding La Luna ${role} Android app`);
	const env = { ...process.env, VITE_APP_ROLE: role, VITE_API_BASE_URL: apiBaseUrl };
	run('npm', ['run', 'build'], env);
	run('npx', ['cap', 'sync', 'android'], env);

	if (syncOnly) return;

	const variant = release ? 'Release' : 'Debug';
	const flavor = role[0].toUpperCase() + role.slice(1);
	const task = `:app:assemble${flavor}${variant}`;

	const wrapper = process.platform === 'win32' ? 'gradlew.bat' : './gradlew';
	const androidDir = path.resolve(frontendDir, 'android');
	run(wrapper, [task, ...(release ? [`:app:bundle${flavor}Release`] : [])], env, androidDir);

	const outputDir = path.resolve(frontendDir, 'dist', 'mobile');
	await mkdir(outputDir, { recursive: true });
	const apk = await findArtifact(
		path.join(androidDir, 'app', 'build', 'outputs', 'apk', role, variant.toLowerCase()),
		'.apk'
	);
	if (!apk) {
		throw new Error(`Could not find the ${role} ${variant.toLowerCase()} APK output.`);
	}
	await cp(apk, path.join(outputDir, `la-luna-${role}-${variant.toLowerCase()}.apk`));

	if (release) {
		const aab = await findArtifact(
			path.join(androidDir, 'app', 'build', 'outputs', 'bundle', `${role}Release`),
			'.aab'
		);
		if (!aab) {
			throw new Error(`Could not find the ${role} release AAB output.`);
		}
		await cp(aab, path.join(outputDir, `la-luna-${role}-release.aab`));
	}
}

if (args.has('--help')) {
	console.log('Usage: npm run mobile:build -- <role|--all> [--release] [--sync-only]');
	process.exit(0);
}

const unknownArgs = [...args].filter(
	(arg) => !['--all', '--release', '--sync-only'].includes(arg) && arg !== selectedRole
);
if (unknownArgs.length || (!selectedRole && !args.has('--all'))) {
	throw new Error('Choose a role or --all. Available roles: ' + roles.join(', '));
}
if (selectedRole && !roles.includes(selectedRole)) {
	throw new Error(`Unknown role "${selectedRole}". Available roles: ${roles.join(', ')}`);
}
if (args.has('--all') && selectedRole) {
	throw new Error('Use either --all or one role, not both.');
}
if (release && syncOnly) {
	throw new Error('--release cannot be combined with --sync-only.');
}

const apiBaseUrl = validateBackendUrl(process.env.VITE_API_BASE_URL);
await verifyUnifiedDatabaseApi(apiBaseUrl);
const targetRoles = args.has('--all') ? roles : [selectedRole];
const outputDir = path.resolve(frontendDir, 'dist', 'mobile');

for (const role of targetRoles) {
	await buildRole(role, apiBaseUrl);
}

console.log(syncOnly ? 'Capacitor assets synced.' : `Android outputs are in ${outputDir}`);
