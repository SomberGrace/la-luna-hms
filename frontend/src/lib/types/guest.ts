// frontend/src/lib/types/guest.ts

export interface Guest {
  id: string;
  fullName: string;
  phone: string;
  nationalIdOrPassport: string;
  email?: string;
  nationality?: string;
  emergencyContact?: string;
}

export interface ReservationItem {
  id: string;
  res_code: string;
  guest_name: string;
  guest_phone: string;
  guest_id_doc: string;
  room_type: 'Standard' | 'Executive' | 'Apartment';
  room_number: string;
  check_in_date: string;
  check_out_date: string;
  nights: number;
  deposit_paid: number;
  status: 'Confirmed' | 'CheckedIn' | 'Cancelled';
}
