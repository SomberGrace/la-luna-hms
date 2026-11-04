import { describe, expect, it } from 'vitest';
import { isNgrokHost, normalizeUserRole, resolveApiRequestUrl } from './mobile';

describe('normalizeUserRole', () => {
	it.each(['CHEF', 'KITCHEN', 'chef', 'kitchen'])('normalizes %s to kitchen', (role) => {
		expect(normalizeUserRole(role)).toBe('kitchen');
	});
});

describe('isNgrokHost', () => {
	it('matches ngrok free domains', () => {
		expect(isNgrokHost('https://rosy-lubricate-ranged.ngrok-free.dev')).toBe(true);
	});

	it('does not match other hosts', () => {
		expect(isNgrokHost('https://api.laluna.mw')).toBe(false);
	});

	it('returns false for malformed URLs', () => {
		expect(isNgrokHost('not a url')).toBe(false);
	});
});

describe('resolveApiRequestUrl', () => {
	it('keeps relative requests for web when no API origin is configured', () => {
		expect(resolveApiRequestUrl('/api/v1/health', '', false)).toBe('/api/v1/health');
	});

	it('prefixes API requests with the configured origin', () => {
		expect(resolveApiRequestUrl('/api/v1/health', 'https://api.example.com/', false)).toBe(
			'https://api.example.com/api/v1/health'
		);
	});

	it('requires an API origin for native app requests', () => {
		expect(() => resolveApiRequestUrl('/api/v1/health', '', true)).toThrow(
			'VITE_API_BASE_URL must be set'
		);
	});

	it('leaves non-API paths unchanged', () => {
		expect(resolveApiRequestUrl('/login', '', true)).toBe('/login');
	});
});
