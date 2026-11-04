import { describe, expect, it } from 'vitest';
import { hasActiveCurrentGuestBooking } from './guestBookings';

describe('hasActiveCurrentGuestBooking', () => {
	it('requires the guest to be checked into an occupied, non-closed room', () => {
		expect(
			hasActiveCurrentGuestBooking([
				{ checked_in: true, room_status: 'Occupied', folio_status: 'Open' }
			])
		).toBe(true);
	});

	it.each([
		{ checked_in: false, room_status: 'Reserved', folio_status: 'Open' },
		{ checked_in: true, room_status: 'Available', folio_status: 'Open' },
		{ checked_in: true, room_status: 'Occupied', folio_status: 'Closed' }
	])('rejects a reservation that is not an active current stay: %o', (booking) => {
		expect(hasActiveCurrentGuestBooking([booking])).toBe(false);
	});
});
