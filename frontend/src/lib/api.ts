import { Capacitor } from '@capacitor/core';
import { get } from 'svelte/store';
import { API_BASE_URL, isNgrokHost, resolveApiRequestUrl } from '#lib/mobile';
import { authStore } from '#lib/stores/authStore';

export async function apiFetch(input: RequestInfo | URL, init: RequestInit = {}) {
	const headers = new Headers(init.headers);
	const token = get(authStore).token;
	if (token) {
		headers.set('Authorization', `Bearer ${token}`);
	}

	let request: RequestInfo | URL = input;
	if (typeof input === 'string') {
		const apiBaseUrl = API_BASE_URL.trim().replace(/\/+$/, '');
		request = resolveApiRequestUrl(input, apiBaseUrl, Capacitor.isNativePlatform());
		if (input.startsWith('/api/') && apiBaseUrl) {
			if (isNgrokHost(apiBaseUrl)) {
				headers.set('ngrok-skip-browser-warning', 'true');
			}
		}
	}

	return fetch(request, { ...init, headers });
}
