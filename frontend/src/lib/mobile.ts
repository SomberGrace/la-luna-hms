declare const __APP_ROLE__: string;
declare const __API_BASE_URL__: string;

export type UserRole =
	'frontdesk' | 'manager' | 'housekeeping' | 'cashier' | 'restaurant' | 'bar' | 'kitchen' | 'guest';

export const APP_ROLE = __APP_ROLE__ as
	'shared' | 'guest' | 'frontdesk' | 'kitchen' | 'manager' | 'restaurant' | 'bar';
export const API_BASE_URL = __API_BASE_URL__;

export function resolveApiRequestUrl(path: string, apiBaseUrl: string, native: boolean) {
	if (!path.startsWith('/api/')) {
		return path;
	}

	const baseUrl = apiBaseUrl.trim().replace(/\/+$/, '');
	if (native && !baseUrl) {
		throw new Error('VITE_API_BASE_URL must be set to the HTTPS backend URL for mobile builds.');
	}

	return baseUrl ? `${baseUrl}${path}` : path;
}

export function isNgrokHost(url: string) {
	try {
		const host = new URL(url).hostname;
		return ['.ngrok-free.dev', '.ngrok-free.app', '.ngrok.app', '.ngrok.dev', '.ngrok.io'].some(
			(suffix) => host.endsWith(suffix)
		);
	} catch {
		return false;
	}
}

export const ROLE_ROUTE_MAP: Record<typeof APP_ROLE, string> = {
	shared: '/',
	guest: '/book',
	frontdesk: '/frontdesk',
	kitchen: '/kitchen',
	manager: '/manager',
	restaurant: '/restaurant',
	bar: '/bar'
};

const APP_PATHS: Record<Exclude<typeof APP_ROLE, 'shared'>, string[]> = {
	guest: ['/guest', '/book'],
	frontdesk: ['/frontdesk'],
	kitchen: ['/kitchen'],
	manager: ['/manager'],
	restaurant: ['/restaurant'],
	bar: ['/bar']
};

export function normalizeUserRole(role: string): UserRole | null {
	switch (role.trim().toUpperCase()) {
		case 'FRONT_DESK':
		case 'FRONTDESK':
			return 'frontdesk';
		case 'CHEF':
		case 'KITCHEN':
			return 'kitchen';
		case 'BAR_ATTENDANT':
		case 'BAR':
			return 'bar';
		case 'WAITER':
		case 'RESTAURANT':
			return 'restaurant';
		case 'MANAGER':
		case 'ADMIN':
			return 'manager';
		case 'GUEST':
			return 'guest';
		case 'HOUSEKEEPING':
			return 'housekeeping';
		case 'CASHIER':
			return 'cashier';
		default:
			return null;
	}
}

export function isRoleAllowedForApp(role: string) {
	const normalizedRole = normalizeUserRole(role);
	return APP_ROLE === 'shared' || normalizedRole === APP_ROLE;
}

export function isPathAllowedForApp(pathname: string) {
	if (APP_ROLE === 'shared' || pathname === '/') {
		return true;
	}
	return APP_PATHS[APP_ROLE].some(
		(prefix) => pathname === prefix || pathname.startsWith(`${prefix}/`)
	);
}
