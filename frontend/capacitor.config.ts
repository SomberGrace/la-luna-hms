import type { CapacitorConfig } from '@capacitor/cli';

const appNameMap = {
	shared: 'La Luna HMS',
	guest: 'La Luna Guest',
	frontdesk: 'La Luna Front Desk',
	kitchen: 'La Luna Kitchen',
	manager: 'La Luna Manager',
	restaurant: 'La Luna Restaurant',
	bar: 'La Luna Bar'
} as const;

const appRole = process.env.VITE_APP_ROLE ?? 'shared';
if (!(appRole in appNameMap)) {
	throw new Error(`Unknown VITE_APP_ROLE: ${appRole}`);
}

const config: CapacitorConfig = {
	appId: 'com.lalunahms',
	appName: appNameMap[appRole as keyof typeof appNameMap],
	webDir: 'build',
	server: {
		androidScheme: 'https'
	}
};

export default config;
