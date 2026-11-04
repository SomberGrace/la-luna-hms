import { apiFetch } from '#lib/api';

type GuestBooking = {
	checked_in: boolean;
	folio_status: string;
	room_status: string;
};

function isGuestBooking(value: unknown): value is GuestBooking {
	if (typeof value !== 'object' || value === null) return false;
	const booking = value as Record<string, unknown>;
	return (
		typeof booking.checked_in === 'boolean' &&
		typeof booking.folio_status === 'string' &&
		typeof booking.room_status === 'string'
	);
}

export function hasActiveCurrentGuestBooking(bookings: readonly GuestBooking[]): boolean {
	return bookings.some(
		(booking) =>
			booking.checked_in &&
			booking.room_status.trim().toUpperCase() === 'OCCUPIED' &&
			booking.folio_status.trim().toUpperCase() !== 'CLOSED'
	);
}

export async function getGuestLandingPath(): Promise<'/guest' | '/book'> {
	const response = await apiFetch('/api/v1/guest/bookings');
	if (!response.ok) {
		throw new Error('Unable to verify your current reservation. Please try signing in again.');
	}

	const bookings: unknown = await response.json();
	if (!Array.isArray(bookings) || !bookings.every(isGuestBooking)) {
		throw new Error('The server returned an invalid reservation list.');
	}

	return hasActiveCurrentGuestBooking(bookings) ? '/guest' : '/book';
}
