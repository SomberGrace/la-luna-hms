// frontend/src/lib/types/room.ts

export type RoomStatus = 'Available' | 'Occupied' | 'Cleaning' | 'OutOfService' | 'Reserved';

export type RoomType = 'Standard' | 'Executive' | 'Apartment' | 'Sector Suite';

export interface Room {
  id: string;
  room_number: string;
  type: RoomType;
  price_per_night: number;
  status: RoomStatus;
  floor?: number;
  housekeeper?: string;
  priority?: 'Normal' | 'High' | 'Urgent';
}

export type FolioCategory = 'Accommodation' | 'Restaurant' | 'Bar' | 'Incidental' | 'Payment';

export interface FolioItem {
  id: string;
  description: string;
  amount: number;
  category: FolioCategory;
  timestamp: string;
}

export interface Folio {
  id: string;
  card_no?: string;
  guest_id: string;
  guest_name: string;
  guest_phone?: string;
  guest_id_doc?: string;
  room_number: string;
  check_in_date: string;
  check_out_date: string;
  total_charges: number;
  total_paid: number;
  balance: number;
  items: FolioItem[];
}

export interface OfflineAction {
  id: string;
  action_type: 'checkin' | 'payment' | 'checkout' | 'extension' | 'room_status';
  payload: any;
  timestamp: string;
  status: 'pending' | 'syncing' | 'failed' | 'synced';
  error_message?: string;
}

export interface ReceptionReport {
  total_rooms: number;
  occupied_rooms: number;
  available_rooms: number;
  cleaning_rooms: number;
  out_of_service: number;
  occupancy_rate: number;
  total_revenue_mwk: number;
  total_balance_mwk: number;
  active_stays_count: number;
  timestamp: string;
}
