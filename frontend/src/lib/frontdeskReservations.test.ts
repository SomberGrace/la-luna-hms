import { describe, expect, it } from 'vitest';
import { getPendingArrivals } from './frontdeskReservations';

describe('getPendingArrivals', () => {
	it('maps active reservations to the front desk list shape', () => {
		expect(
			getPendingArrivals([
				{
					id: 'card-1',
					res_code: 'card-1',
					card_no: 'card-1',
					guest_name: 'Amina Banda',
					guest_phone: '12345',
					guest_id_doc: 'P123',
					room_type: 'Standard',
					room_number: '101',
					check_in: '2026-11-04',
					check_out: '2026-11-06',
					nights: 2,
					deposit_paid: 100,
					folio_status: 'Open',
					checked_in: false
				}
			])
		).toEqual([
			{
				id: 'card-1',
				res_code: 'card-1',
				guest_name: 'Amina Banda',
				guest_phone: '12345',
				guest_id_doc: 'P123',
				room_type: 'Standard',
				room_number: '101',
				check_in_date: '2026-11-04',
				check_out_date: '2026-11-06',
				nights: 2,
				deposit_paid: 100,
				status: 'Pending Arrival'
			}
		]);
	});

	it('skips checked-in and closed reservations', () => {
		const active = {
			id: 'card-1',
			res_code: 'card-1',
			card_no: 'card-1',
			guest_name: 'Amina Banda',
			guest_phone: '12345',
			room_type: 'Standard',
			room_number: '101',
			check_in: '2026-11-04',
			check_out: '2026-11-06',
			folio_status: 'Open',
			checked_in: false
		};

		expect(
			getPendingArrivals([
				{ ...active, checked_in: true },
				{ ...active, folio_status: 'Closed' }
			])
		).toEqual([]);
	});

	it('returns null for malformed responses', () => {
		expect(getPendingArrivals([{ room_number: '101' }])).toBeNull();
	});

	it('uses the existing defaults for missing stay values', () => {
		const arrivals = getPendingArrivals([
			{
				id: 'card-1',
				res_code: 'card-1',
				card_no: 'card-1',
				guest_name: 'Amina Banda',
				guest_phone: '12345',
				room_type: 'Standard',
				room_number: '101',
				check_in: '2026-11-04',
				check_out: '2026-11-06',
				folio_status: 'Open',
				checked_in: false
			}
		]);

		expect(arrivals).toEqual([
			expect.objectContaining({
				nights: 1,
				deposit_paid: 0,
				guest_id_doc: ''
			})
		]);
	});

	it('keeps reservations for suite rooms', () => {
		expect(
			getPendingArrivals([
				{
					id: 'card-1',
					res_code: 'card-1',
					card_no: 'card-1',
					guest_name: 'Amina Banda',
					guest_phone: '12345',
					room_type: 'Sector Suite',
					room_number: '401',
					check_in: '2026-11-04',
					check_out: '2026-11-06',
					folio_status: 'Open',
					checked_in: false
				}
			])
		).toEqual([
			expect.objectContaining({
				room_type: 'Sector Suite'
			})
		]);
	});

	it('returns null for a non-array response', () => {
		expect(getPendingArrivals(null)).toBeNull();
	});
});
