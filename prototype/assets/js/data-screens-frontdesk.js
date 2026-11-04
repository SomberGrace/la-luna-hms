/* Screens: front desk workstation */
(function () {
  'use strict';
  const A = window.A, P = window.PROTO;

  Object.assign(P.screens, {
    '/frontdesk': {
      dept: 'frontdesk', chrome: 'brand', title: 'Front Desk Hub',
      blocks: [
        {
          k: 'hero', kicker: 'Morning shift · 06:00 – 14:00', title: 'Front Desk Hub',
          sub: 'Reception A · Grace M. · 12 rooms in house', icon: 'ph-device-mobile', tone: 'amber',
          chips: ['Online', 'Queue clear', 'Handover 14:00']
        },
        {
          k: 'stats', items: [
            { label: 'Occupancy', value: '78%', note: '9 of 12 rooms', tone: 'amber' },
            { label: 'Arrivals', value: '4', note: '2 expected now', tone: 'emerald' },
            { label: 'Departures', value: '3', note: '1 balance due', tone: 'orange' },
            { label: 'In house', value: '12', note: '4 VIP guests', tone: 'indigo' }
          ]
        },
        { k: 'section', title: 'Expected arrivals', action: { text: '4 today', tone: 'muted' } },
        {
          k: 'list', items: [
            { icon: 'ph-user', title: 'Alice Banda', sub: 'Room 104 · Standard Double · 3 nights', tone: 'amber', right: 'Check in', action: A.nav('/frontdesk/checkin', 'Check in Alice Banda') },
            { icon: 'ph-user', title: 'Thandiwe Phiri', sub: 'Room 201 · Executive · 2 nights', tone: 'emerald', right: 'Check in', action: A.nav('/frontdesk/checkin', 'Check in Thandiwe Phiri') },
            { icon: 'ph-user', title: 'James Mwale', sub: 'Room 301 · Apartment · 5 nights', tone: 'indigo', right: 'Check in', action: A.nav('/frontdesk/checkin', 'Check in James Mwale') }
          ]
        },
        { k: 'section', title: 'Active guest stays', action: { text: 'View rooms', tone: 'link', action: A.nav('/frontdesk/rooms', 'View all rooms') } },
        {
          k: 'list', items: [
            { icon: 'ph-bed', title: 'Room 102 · Chisomo Nyirenda', sub: 'In house · Departure 02 Nov · Balance 18,000', tone: 'orange', right: 'Folio', action: A.nav('/frontdesk/folio', 'Open folio for room 102') },
            { icon: 'ph-bed', title: 'Room 205 · Peter Kamanga', sub: 'In house · Departure 04 Nov · Balance 0', tone: 'emerald', right: 'Folio', action: A.nav('/frontdesk/folio', 'Open folio for room 205') }
          ]
        },
        { k: 'section', title: 'Quick actions' },
        {
          k: 'quick', items: [
            { icon: 'ph-door-open', label: 'Rooms', tone: 'amber', action: A.nav('/frontdesk/rooms', 'Rooms') },
            { icon: 'ph-user-plus', label: 'Walk-in', tone: 'emerald', action: A.nav('/frontdesk/walkin', 'Walk-in') },
            { icon: 'ph-receipt', label: 'Folio', tone: 'indigo', action: A.nav('/frontdesk/folio', 'Folio') },
            { icon: 'ph-credit-card', label: 'Checkout', tone: 'orange', action: A.nav('/frontdesk/checkout', 'Checkout') },
            { icon: 'ph-hand-coins', label: 'Payments', tone: 'emerald', action: A.nav('/frontdesk/payments', 'Payments') },
            { icon: 'ph-broom', label: 'Housekeeping', tone: 'indigo', action: A.nav('/frontdesk/housekeeping', 'Housekeeping') },
            { icon: 'ph-clock', label: 'Extension', tone: 'amber', action: A.nav('/frontdesk/extension', 'Stay extension') },
            { icon: 'ph-cloud-slash', label: 'Offline', tone: 'rose', action: A.nav('/frontdesk/offline', 'Offline sync') },
            { icon: 'ph-chart-line-up', label: 'Reports', tone: 'slate', action: A.nav('/frontdesk/reports', 'Reception reports') }
          ]
        },
        { k: 'note', icon: 'ph-warning-circle', tone: 'orange', title: '1 folio needs attention', text: 'Room 102 has an outstanding balance of MWK 18,000 and departs tomorrow at 10:00.', action: A.nav('/frontdesk/folio', 'Review room 102 folio') }
      ]
    },

    '/frontdesk/rooms': {
      dept: 'frontdesk', chrome: 'sub', title: 'Room Inventory', back: '/frontdesk',
      blocks: [
        { k: 'hero', kicker: 'Live housekeeping status', title: 'Room Inventory', sub: 'Tap a room to change its status. Availability syncs to the guest booking page.', icon: 'ph-door-open', tone: 'amber', chips: ['12 rooms', '9 occupied', '1 out of service'] },
        {
          k: 'chips', items: [
            { label: 'All 12', active: true, action: A.toast('Showing all rooms') },
            { label: 'Available 2', action: A.toast('Filtered: available') },
            { label: 'Occupied 9', action: A.toast('Filtered: occupied') },
            { label: 'Out of service 1', action: A.toast('Filtered: out of service') }
          ]
        },
        {
          k: 'rooms', items: [
            { no: '101', type: 'Standard', bed: 'Single', rate: 45000, status: 'Occupied', tone: 'orange', action: A.modal('room-status', 'Change status for room 101') },
            { no: '102', type: 'Standard', bed: 'Double', rate: 45000, status: 'Occupied', tone: 'orange', action: A.modal('room-status', 'Change status for room 102') },
            { no: '103', type: 'Standard', bed: 'Single', rate: 45000, status: 'Available', tone: 'emerald', action: A.modal('room-status', 'Change status for room 103') },
            { no: '104', type: 'Standard', bed: 'Double', rate: 45000, status: 'Occupied', tone: 'orange', action: A.modal('room-status', 'Change status for room 104') },
            { no: '201', type: 'Executive', bed: 'Queen', rate: 85000, status: 'Occupied', tone: 'indigo', action: A.modal('room-status', 'Change status for room 201') },
            { no: '205', type: 'Executive', bed: 'Queen', rate: 85000, status: 'Occupied', tone: 'indigo', action: A.modal('room-status', 'Change status for room 205') },
            { no: '301', type: 'Apartment', bed: 'King', rate: 120000, status: 'Available', tone: 'emerald', action: A.modal('room-status', 'Change status for room 301') },
            { no: '304', type: 'Apartment', bed: 'King', rate: 120000, status: 'Occupied', tone: 'orange', action: A.modal('room-status', 'Change status for room 304') },
            { no: '401', type: 'Sector Suite', bed: 'King', rate: 180000, status: 'Out of service', tone: 'rose', action: A.modal('room-status', 'Change status for room 401') }
          ]
        },
        { k: 'note', icon: 'ph-info', tone: 'slate', title: 'Status rules', text: 'Out of service rooms are hidden from guest booking. Dirty rooms can be sold but not checked in until housekeeping marks them clean.' }
      ]
    },

    '/frontdesk/walkin': {
      dept: 'frontdesk', chrome: 'sub', title: 'Rapid Walk-In', back: '/frontdesk',
      blocks: [
        { k: 'hero', kicker: 'No reservation required', title: 'Rapid Walk-In', sub: 'Register a guest, assign a room and open a folio in one pass.', icon: 'ph-user-plus', tone: 'emerald', chips: ['Deposit required', 'ID capture'] },
        {
          k: 'form', title: 'Guest registration',
          fields: [
            { label: 'Guest full name', type: 'text', value: 'Memory Tembo' },
            { label: 'Phone number', type: 'tel', value: '+265 991 234 567' },
            { label: 'ID / passport', type: 'text', value: 'MW 448392' },
            { label: 'Room', type: 'select', value: '103 · Standard Single', options: ['103 · Standard Single · MWK 45,000', '301 · Apartment King · MWK 120,000'] },
            { label: 'Nights', type: 'number', value: '2' },
            { label: 'Payment method', type: 'select', value: 'Cash', options: ['Cash', 'Visa', 'Cheque', 'Credit account'] }
          ],
          submit: { label: 'Check in guest', icon: 'ph-check-circle', action: A.modal('walkin-success', 'Complete walk-in check-in') },
          note: 'A folio is opened automatically with the room charge for the stay.'
        },
        { k: 'note', icon: 'ph-shield-check', tone: 'emerald', title: 'Deposit policy', text: 'Collect one night deposit for walk-ins without a corporate account. The deposit is applied to the folio immediately.' }
      ]
    },

    '/frontdesk/checkin': {
      dept: 'frontdesk', chrome: 'sub', title: 'Reservation Check-In', back: '/frontdesk',
      blocks: [
        { k: 'hero', kicker: 'Arrival · 01 Nov 2026', title: 'Reservation Check-In', sub: 'Verify the reservation, capture payment and issue the room card.', icon: 'ph-key', tone: 'amber', chips: ['4 arrivals', '2 VIP'] },
        {
          k: 'chips', items: [
            { label: 'Alice Banda', active: true, action: A.toast('Reservation for Alice Banda selected') },
            { label: 'Thandiwe Phiri', action: A.toast('Reservation for Thandiwe Phiri selected') },
            { label: 'James Mwale', action: A.toast('Reservation for James Mwale selected') }
          ]
        },
        {
          k: 'form', title: 'Reservation details',
          fields: [
            { label: 'Reservation', type: 'select', value: 'RES-2026-0184 · Alice Banda', options: ['RES-2026-0184 · Alice Banda', 'RES-2026-0185 · Thandiwe Phiri', 'RES-2026-0186 · James Mwale'] },
            { label: 'Room', type: 'select', value: '104 · Standard Double', options: ['104 · Standard Double · MWK 45,000', '201 · Executive Queen · MWK 85,000', '301 · Apartment King · MWK 120,000'] },
            { label: 'Nights', type: 'number', value: '3' },
            { label: 'Deposit (MWK)', type: 'number', value: '45000' },
            { label: 'Payment method', type: 'select', value: 'Visa', options: ['Visa', 'Cash', 'Cheque', 'Credit account'] },
            { label: 'Card / reference', type: 'text', value: 'VISA ···· 4242' }
          ],
          submit: { label: 'Confirm check-in', icon: 'ph-check-circle', action: A.modal('checkin-success', 'Confirm check-in') },
          note: 'Check-in locks the room, creates the guest record and opens the folio.'
        }
      ]
    },

    '/frontdesk/folio': {
      dept: 'frontdesk', chrome: 'sub', title: 'Guest Billing Folio', back: '/frontdesk',
      blocks: [
        { k: 'hero', kicker: 'Room 104 · Alice Banda', title: 'MWK 62,500', sub: 'Outstanding balance · 2 nights in house', icon: 'ph-receipt', tone: 'emerald', chips: ['Folio F-1042', 'Departure 03 Nov'] },
        {
          k: 'statement',
          rows: [
            { label: 'Accommodation — 2 nights', sub: '30 Oct – 01 Nov · Standard Double', amount: 90000 },
            { label: 'In-room dining', sub: 'Tilapia, juice, mandasi', amount: 6000 },
            { label: 'Bar — poolside', sub: '2 × Carlsberg, 1 × gin & tonic', amount: 6500 },
            { label: 'Payments received', sub: 'Visa · 30 Oct', amount: -40000, tone: 'emerald' }
          ],
          total: { label: 'Balance due', amount: 62500 }
        },
        { k: 'section', title: 'Folio actions' },
        {
          k: 'quick', items: [
            { icon: 'ph-hand-coins', label: 'Post payment', tone: 'emerald', action: A.nav('/frontdesk/payments', 'Post payment') },
            { icon: 'ph-clock', label: 'Extend stay', tone: 'amber', action: A.nav('/frontdesk/extension', 'Extend stay') },
            { icon: 'ph-credit-card', label: 'Check out', tone: 'orange', action: A.nav('/frontdesk/checkout', 'Check out room 104') }
          ]
        },
        { k: 'section', title: 'Charge history' },
        {
          k: 'timeline', items: [
            { time: 'Today · 19:02', title: 'In-room dining posted', sub: 'Tilapia & juice · MWK 6,000', tone: 'orange' },
            { time: 'Today · 14:20', title: 'Bar charge posted', sub: 'Poolside · MWK 6,500', tone: 'purple' },
            { time: '30 Oct · 15:10', title: 'Deposit received', sub: 'Visa · MWK 40,000', tone: 'emerald' }
          ]
        }
      ]
    },

    '/frontdesk/payments': {
      dept: 'frontdesk', chrome: 'sub', title: 'Payment Terminal', back: '/frontdesk',
      blocks: [
        { k: 'hero', kicker: 'Folio F-1042 · Room 104', title: 'MWK 62,500', sub: 'Post a payment against the guest folio. Receipts are generated instantly.', icon: 'ph-hand-coins', tone: 'emerald', chips: ['Card on file', 'No refunds after checkout'] },
        {
          k: 'form', title: 'New payment',
          fields: [
            { label: 'Amount (MWK)', type: 'number', value: '40000' },
            { label: 'Method', type: 'select', value: 'Visa', options: ['Visa', 'Cash', 'Cheque', 'Credit account'] },
            { label: 'Reference', type: 'text', value: 'VISA ···· 4242' },
            { label: 'Note', type: 'text', value: 'Part settlement of room balance' }
          ],
          submit: { label: 'Post payment', icon: 'ph-check-circle', action: A.modal('payment-receipt', 'Post payment') }
        },
        { k: 'section', title: 'Recent payments' },
        {
          k: 'list', items: [
            { icon: 'ph-credit-card', title: 'MWK 40,000 · Visa', sub: '30 Oct · Deposit · Grace M.', tone: 'emerald', right: 'Receipt', action: A.modal('payment-receipt', 'View receipt') },
            { icon: 'ph-money', title: 'MWK 15,000 · Cash', sub: '30 Oct · Bar settlement · Peter K.', tone: 'amber', right: 'Receipt', action: A.modal('payment-receipt', 'View receipt') }
          ]
        }
      ]
    },

    '/frontdesk/checkout': {
      dept: 'frontdesk', chrome: 'sub', title: 'Check-Out Terminal', back: '/frontdesk',
      blocks: [
        { k: 'hero', kicker: 'Room 104 · Alice Banda', title: 'Check-Out Terminal', sub: 'Settle the balance, release the room and notify housekeeping.', icon: 'ph-credit-card', tone: 'orange', chips: ['Balance 62,500', 'Departs today'] },
        {
          k: 'statement',
          rows: [
            { label: 'Total charges', sub: 'Accommodation, dining and bar', amount: 102500 },
            { label: 'Payments received', sub: 'Visa deposit', amount: -40000, tone: 'emerald' }
          ],
          total: { label: 'Balance to settle', amount: 62500 }
        },
        { k: 'section', title: 'Pre-departure checks' },
        {
          k: 'list', items: [
            { icon: 'ph-check-circle', title: 'Minibar checked', sub: 'No missing items', tone: 'emerald', right: 'Clear' },
            { icon: 'ph-key', title: 'Room card returned', sub: '2 cards issued', tone: 'emerald', right: 'Clear' },
            { icon: 'ph-broom', title: 'Housekeeping notified', sub: 'Room marked dirty on checkout', tone: 'indigo', right: 'Auto' }
          ]
        },
        { k: 'form', compact: true, fields: [{ label: 'Settlement method', type: 'select', value: 'Visa', options: ['Visa', 'Cash', 'Credit account'] }],
          submit: { label: 'Complete check-out', icon: 'ph-check-circle', action: A.modal('checkout-success', 'Complete check-out') },
          note: 'The folio is closed and the room becomes available after housekeeping.' }
      ]
    },

    '/frontdesk/extension': {
      dept: 'frontdesk', chrome: 'sub', title: 'Stay Extension', back: '/frontdesk',
      blocks: [
        { k: 'hero', kicker: 'Room 104 · Alice Banda', title: 'Stay Extension', sub: 'Current departure 03 Nov 2026 · Standard Double · MWK 45,000 / night', icon: 'ph-clock', tone: 'amber', chips: ['Room available', 'Rate held'] },
        {
          k: 'form', title: 'Extend the stay',
          fields: [
            { label: 'Additional nights', type: 'number', value: '1' },
            { label: 'New departure', type: 'date', value: '2026-11-04', readonly: true },
            { label: 'Rate per night (MWK)', type: 'number', value: '45000' },
            { label: 'Additional charge (MWK)', type: 'number', value: '45000', readonly: true }
          ],
          submit: { label: 'Confirm extension', icon: 'ph-check-circle', action: A.modal('extension-success', 'Confirm extension') },
          note: 'The extra night is posted to the folio and the room is held until the new departure date.'
        },
        { k: 'note', icon: 'ph-calendar-check', tone: 'slate', title: 'Availability note', text: 'Room 104 is free until 05 Nov. A second extension can be requested on the folio screen.' }
      ]
    },

    '/frontdesk/housekeeping': {
      dept: 'frontdesk', chrome: 'sub', title: 'Housekeeping', back: '/frontdesk',
      blocks: [
        { k: 'hero', kicker: 'Live board · 08:15', title: 'Housekeeping', sub: 'Track cleaning, inspection and out-of-service rooms across the property.', icon: 'ph-broom', tone: 'indigo', chips: ['5 clean', '3 dirty', '1 inspect'] },
        {
          k: 'stats', items: [
            { label: 'Clean', value: '5', note: 'Ready to sell', tone: 'emerald' },
            { label: 'Dirty', value: '3', note: 'Awaiting service', tone: 'orange' },
            { label: 'Inspect', value: '1', note: 'Supervisor check', tone: 'indigo' }
          ]
        },
        { k: 'section', title: 'Room board' },
        {
          k: 'list', items: [
            { icon: 'ph-broom', title: 'Room 102', sub: 'Dirty · Guest departed 10:00', tone: 'orange', right: 'Mark clean', action: A.toast('Room 102 marked clean') },
            { icon: 'ph-broom', title: 'Room 205', sub: 'Dirty · Stayover service', tone: 'orange', right: 'Mark clean', action: A.toast('Room 205 marked clean') },
            { icon: 'ph-sparkle', title: 'Room 201', sub: 'Inspect · Supervisor assigned', tone: 'indigo', right: 'Pass', action: A.toast('Room 201 passed inspection') },
            { icon: 'ph-check-circle', title: 'Room 301', sub: 'Clean · Ready for walk-in', tone: 'emerald', right: 'Ready' },
            { icon: 'ph-warning', title: 'Room 401', sub: 'Out of service · Plumbing repair', tone: 'rose', right: 'Blocked' }
          ]
        }
      ]
    },

    '/frontdesk/offline': {
      dept: 'frontdesk', chrome: 'sub', title: 'Offline Sync & Queue', back: '/frontdesk',
      blocks: [
        { k: 'hero', kicker: 'Connection restored', title: 'Offline Sync & Queue', sub: 'Actions taken during the outage are queued locally and replay in order.', icon: 'ph-cloud-check', tone: 'emerald', chips: ['3 queued', 'Last sync 08:02'] },
        {
          k: 'stats', items: [
            { label: 'Queued', value: '3', note: 'Waiting to sync', tone: 'orange' },
            { label: 'Synced', value: '18', note: 'Since 06:00', tone: 'emerald' },
            { label: 'Conflicts', value: '0', note: 'All clear', tone: 'slate' }
          ]
        },
        { k: 'section', title: 'Pending sync queue' },
        {
          k: 'list', items: [
            { icon: 'ph-user-plus', title: 'Walk-in check-in · Memory Tembo', sub: 'Room 103 · 08:01 · queued', tone: 'amber', right: 'Retry', action: A.toast('Retrying walk-in check-in') },
            { icon: 'ph-receipt', title: 'Folio charge · Room 102', sub: 'Minibar · MWK 3,500 · 07:58', tone: 'indigo', right: 'Retry', action: A.toast('Retrying folio charge') },
            { icon: 'ph-hand-coins', title: 'Cash payment · Room 205', sub: 'MWK 18,000 · 07:54', tone: 'emerald', right: 'Retry', action: A.toast('Retrying cash payment') }
          ]
        },
        { k: 'section', title: 'Outage standard operating procedure' },
        {
          k: 'timeline', items: [
            { time: '1 · Detect', title: 'Switch to paper log', sub: 'Record guest name, room and action on the outage sheet', tone: 'amber' },
            { time: '2 · Capture', title: 'Photograph ID and receipt', sub: 'Attach to the queued action when the connection returns', tone: 'indigo' },
            { time: '3 · Replay', title: 'Sync in order', sub: 'Never re-enter an action that is already queued', tone: 'emerald' }
          ]
        },
        { k: 'note', icon: 'ph-arrows-clockwise', tone: 'orange', title: 'Queue replays automatically', text: 'Keep this screen open while the connection is unstable. Actions replay oldest first.', action: A.toast('Sync started — 3 actions queued') }
      ]
    },

    '/frontdesk/reports': {
      dept: 'frontdesk', chrome: 'sub', title: 'Reception Reports', back: '/frontdesk',
      blocks: [
        { k: 'hero', kicker: 'Shift 06:00 – 14:00', title: 'Reception Reports', sub: 'Morning shift performance, settlement mix and handover log.', icon: 'ph-chart-line-up', tone: 'amber', chips: ['Grace M.', '01 Nov 2026'] },
        {
          k: 'stats', items: [
            { label: 'Room revenue', value: '315,000', note: 'MWK this shift', tone: 'amber' },
            { label: 'Collected', value: '164,000', note: 'MWK settled', tone: 'emerald' },
            { label: 'ADR', value: '45,000', note: 'Average daily rate', tone: 'indigo' },
            { label: 'RevPAR', value: '35,100', note: '78% occupancy', tone: 'orange' }
          ]
        },
        { k: 'section', title: 'Settlement method' },
        {
          k: 'bars', items: [
            { label: 'Visa', value: 'MWK 92,000', pct: 56, tone: 'indigo' },
            { label: 'Cash', value: 'MWK 48,000', pct: 29, tone: 'emerald' },
            { label: 'Cheque', value: 'MWK 16,000', pct: 10, tone: 'amber' },
            { label: 'Credit account', value: 'MWK 8,000', pct: 5, tone: 'orange' }
          ]
        },
        { k: 'section', title: 'Shift handover log' },
        {
          k: 'table', head: ['Time', 'Event', 'By'],
          rows: [
            ['06:05', 'Shift opened · float MWK 20,000', 'Grace M.'],
            ['07:58', 'Outage logged · 3 actions queued', 'Grace M.'],
            ['08:02', 'Connection restored · queue replayed', 'System'],
            ['09:40', 'Room 401 blocked · plumbing', 'Maintenance']
          ]
        },
        { k: 'note', icon: 'ph-printer', tone: 'slate', title: 'Shift close', text: 'Closing the shift locks the cash float and sends the report to the manager consolidation queue.', action: A.modal('shift-close', 'Close shift') }
      ]
    }
  });
})();
