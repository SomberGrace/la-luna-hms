import type { ReservationItem } from '#lib/types/guest';
import type { RoomType } from '#lib/types/room';

type ReservationResponse = {
	id: string;
	res_code: string;
	card_no: string;
	guest_name: string;
	guest_phone: string;
	guest_id_doc?: string;
	room_type: ReservationItem['room_type'];
	room_number: string;
	check_in: string;
	check_out: string;
	nights?: number;
	deposit_paid?: number;
	folio_status: string;
	checked_in: boolean;
};

function isRoomType(value: unknown): value is RoomType {
	return (
		value === 'Standard' ||
		value === 'Executive' ||
		value === 'Apartment' ||
		value === 'Sector Suite'
	);
}

function isReservationResponse(value: unknown): value is ReservationResponse {
	if (typeof value !== 'object' || value === null) return false;
	const row = value as Record<string, unknown>;

	return (
		typeof row.id === 'string' &&
		typeof row.res_code === 'string' &&
		typeof row.card_no === 'string' &&
		typeof row.guest_name === 'string' &&
		typeof row.guest_phone === 'string' &&
		isRoomType(row.room_type) &&
		typeof row.room_number === 'string' &&
		typeof row.check_in === 'string' &&
		typeof row.check_out === 'string' &&
		typeof row.folio_status === 'string' &&
		typeof row.checked_in === 'boolean' &&
		(row.guest_id_doc === undefined || typeof row.guest_id_doc === 'string') &&
		(row.nights === undefined || typeof row.nights === 'number') &&
		(row.deposit_paid === undefined || typeof row.deposit_paid === 'number')
	);
}

export function getPendingArrivals(value: unknown): ReservationItem[] | null {
	if (!Array.isArray(value)) return null;

	const rows = value.filter(isReservationResponse);
	if (rows.length !== value.length) return null;

	return rows
		.filter((row) => !row.checked_in && row.folio_status !== 'Closed')
		.map((row) => ({
			id: row.id || row.card_no,
			res_code: row.res_code || row.card_no,
			guest_name: row.guest_name,
			guest_phone: row.guest_phone,
			guest_id_doc: row.guest_id_doc || '',
			room_type: row.room_type,
			room_number: row.room_number,
			check_in_date: row.check_in,
			check_out_date: row.check_out,
			nights: row.nights || 1,
			deposit_paid: row.deposit_paid || 0,
			status: 'Pending Arrival'
		}));
}
