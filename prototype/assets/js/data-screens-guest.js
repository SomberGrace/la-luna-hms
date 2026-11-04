/* Screens: portal, public booking, guest portal */
(function () {
  'use strict';
  const A = window.A, P = window.PROTO;

  Object.assign(P.screens, {
    '/': {
      dept: 'portal', chrome: 'none', title: 'Station Sign In',
      blocks: [{ k: 'login' }]
    },

    '/book': {
      dept: 'book', chrome: 'brand', title: 'Book Your Stay',
      blocks: [
        {
          k: 'hero', kicker: 'Lake Malawi · Malawi', title: 'Book Your Stay',
          sub: 'Twelve garden and lake-facing rooms. Instant confirmation, no deposit required.',
          icon: 'ph-moon-stars', tone: 'amber',
          chips: ['Free cancellation', 'Best rate', 'Airport transfer']
        },
        {
          k: 'form', compact: true, title: 'Find a room',
          fields: [
            { label: 'Check in', type: 'date', value: '2026-11-10' },
            { label: 'Check out', type: 'date', value: '2026-11-13' },
            { label: 'Guests', type: 'select', value: '2 adults', options: ['1 adult', '2 adults', '2 adults · 1 child', '3 adults'] },
            { label: 'Room type', type: 'select', value: 'Any room', options: ['Any room', 'Standard', 'Executive', 'Apartment', 'Sector Suite'] }
          ],
          submit: { label: 'Check availability', icon: 'ph-magnifying-glass', action: A.toast('Availability refreshed for 10–13 Nov', 'Check availability') }
        },
        { k: 'section', title: 'Available rooms', action: { text: '3 nights', tone: 'muted' } },
        {
          k: 'rooms', items: [
            { no: '101', type: 'Standard', bed: 'Single', rate: 45000, status: 'Available', tone: 'amber', action: A.modal('book-room', 'Book room 101') },
            { no: '104', type: 'Standard', bed: 'Double', rate: 45000, status: 'Available', tone: 'amber', action: A.modal('book-room', 'Book room 104') },
            { no: '201', type: 'Executive', bed: 'Queen', rate: 85000, status: 'Available', tone: 'orange', action: A.modal('book-room', 'Book room 201') },
            { no: '301', type: 'Apartment', bed: 'King', rate: 120000, status: 'Available', tone: 'emerald', action: A.modal('book-room', 'Book room 301') },
            { no: '401', type: 'Sector Suite', bed: 'King', rate: 180000, status: 'Available', tone: 'indigo', action: A.modal('book-room', 'Book room 401') }
          ]
        },
        { k: 'note', icon: 'ph-info', tone: 'amber', title: 'Direct booking promise', text: 'Book direct for the lowest rate. Reservations are held for 24 hours and confirmed by the front desk on arrival.' },
        { k: 'section', title: 'Already staying with us?' },
        { k: 'list', items: [
          { icon: 'ph-key', title: 'Guest portal login', sub: 'View your bill, order dining, request services', tone: 'amber', right: 'Sign in', action: A.nav('/book/login', 'Guest portal login') }
        ] }
      ]
    },

    '/book/login': {
      dept: 'book', chrome: 'brand', title: 'Guest Account', back: '/book',
      blocks: [
        { k: 'hero', kicker: 'Resident access', title: 'Guest Account', sub: 'Sign in to view your statement and reach every hotel service from your phone.', icon: 'ph-identification-card', tone: 'amber', chips: ['Room 104', 'Departure 03 Nov'] },
        {
          k: 'chips', items: [
            { label: 'Sign in', active: true, action: A.toast('Sign in selected') },
            { label: 'Create account', action: A.toast('Registration form selected') }
          ]
        },
        {
          k: 'form', title: 'Sign in',
          fields: [
            { label: 'Email or room card', type: 'text', value: 'alice.banda@example.mw' },
            { label: 'Password', type: 'password', value: 'laluna2026' }
          ],
          submit: { label: 'Enter guest portal', icon: 'ph-sign-in', action: A.nav('/guest', 'Enter guest portal') },
          note: 'Demo: any details will open the guest portal.'
        },
        { k: 'note', icon: 'ph-lifebuoy', tone: 'slate', title: 'Need help?', text: 'Call reception on extension 0, or ask at the front desk. The guest portal is available on the hotel WiFi.' }
      ]
    },

    '/guest': {
      dept: 'guest', chrome: 'brand', title: 'Your Stay',
      blocks: [
        {
          k: 'hero', kicker: 'Room 104 · Garden Wing', title: 'Good evening, Alice',
          sub: 'Checked in 30 Oct 2026 · Departure 03 Nov 2026 · 2 guests',
          icon: 'ph-bed', tone: 'amber',
          chips: ['3 nights', 'Standard Double', 'Card on file']
        },
        {
          k: 'stats', items: [
            { label: 'Nights left', value: '2', note: 'Departure 03 Nov', tone: 'amber' },
            { label: 'Folio balance', value: '62,500', note: 'MWK outstanding', tone: 'emerald' },
            { label: 'Requests', value: '1', note: 'In progress', tone: 'indigo' },
            { label: 'WiFi', value: 'Live', note: 'LaLuna-Guest', tone: 'slate' }
          ]
        },
        { k: 'section', title: 'Resident services' },
        {
          k: 'grid', items: [
            { icon: 'ph-bowl-food', title: 'In-room dining', sub: 'Order to your door', tone: 'orange', action: A.nav('/guest/dining', 'In-room dining') },
            { icon: 'ph-receipt', title: 'My bill', sub: 'Live statement', tone: 'emerald', action: A.nav('/guest/folio', 'My bill') },
            { icon: 'ph-bell-ringing', title: 'Concierge', sub: 'Requests & help', tone: 'indigo', action: A.nav('/guest/services', 'Concierge') },
            { icon: 'ph-map-pin', title: 'Resort guide', sub: 'Hours & excursions', tone: 'amber', action: A.nav('/guest/guide', 'Resort guide') }
          ]
        },
        { k: 'section', title: 'Active requests', action: { text: '1 open', tone: 'muted' } },
        {
          k: 'list', items: [
            { icon: 'ph-towel', title: 'Extra bath towels', sub: 'Requested 18:42 · Housekeeping', tone: 'indigo', right: 'On the way', rightTone: 'indigo', action: A.toast('Housekeeping is 10 minutes away') }
          ]
        },
        { k: 'note', icon: 'ph-wifi-high', tone: 'amber', title: 'Hotel WiFi', text: 'Network LaLuna-Guest-WiFi · tap to view the password', action: A.modal('wifi', 'View WiFi details') }
      ]
    },

    '/guest/dining': {
      dept: 'guest', chrome: 'sub', title: 'In-Room Dining', back: '/guest',
      blocks: [
        { k: 'hero', kicker: 'Kitchen open 07:00–22:00', title: 'In-Room Dining', sub: 'Fresh from the La Luna kitchen. Average delivery 20–35 minutes.', icon: 'ph-bowl-food', tone: 'orange', chips: ['Charge to room 104', 'No service fee'] },
        {
          k: 'chips', items: [
            { label: 'All', active: true, action: A.toast('Showing all dishes') },
            { label: 'Mains', action: A.toast('Mains selected') },
            { label: 'Breakfast', action: A.toast('Breakfast selected') },
            { label: 'Snacks', action: A.toast('Snacks selected') },
            { label: 'Drinks', action: A.toast('Drinks selected') }
          ]
        },
        {
          k: 'menu', items: [
            { name: 'Nsima with Stewed Chicken', desc: 'Traditional Malawian maize meal with rich tomato chicken stew and seasonal greens.', price: 3500, tag: 'Mains', tone: 'orange', action: A.tray('Add Nsima') },
            { name: 'Grilled Lake Tilapia (Chambo)', desc: 'Whole Lake Malawi chambo grilled with lemon butter and hand-cut chips.', price: 4500, tag: 'Mains', tone: 'emerald', action: A.tray('Add Tilapia') },
            { name: 'Prime Beef Tenderloin', desc: 'Pan-seared tenderloin, garlic mash, roasted vegetables and pepper sauce.', price: 8500, tag: 'Mains', tone: 'indigo', action: A.tray('Add Tenderloin') },
            { name: 'Warm Mandasi & Local Honey', desc: 'Fluffy African donuts with butter, Zomba honey and berry jam.', price: 800, tag: 'Snacks', tone: 'amber', action: A.tray('Add Mandasi') },
            { name: 'Fresh Mango & Passion Juice', desc: 'Freshly pressed seasonal tropical fruit, chilled over ice.', price: 1500, tag: 'Drinks', tone: 'amber', action: A.tray('Add Juice') },
            { name: 'Malawi Gin & Tonic', desc: 'Signature MGT with Malawi Gin, crisp tonic and fresh lime.', price: 2500, tag: 'Drinks', tone: 'purple', action: A.tray('Add Gin & Tonic') }
          ]
        },
        { k: 'tray', count: 2, total: 8000, action: A.modal('guest-tray', 'Review dining tray') },
        { k: 'section', title: 'Your last order' },
        { k: 'list', items: [
          { icon: 'ph-check-circle', title: 'Tilapia & Mango Juice', sub: 'Yesterday 19:20 · Delivered', tone: 'emerald', right: 'MWK 6,000', action: A.toast('Reorder added to your tray') }
        ] }
      ]
    },

    '/guest/folio': {
      dept: 'guest', chrome: 'sub', title: 'My Room Statement', back: '/guest',
      blocks: [
        { k: 'hero', kicker: 'Room 104 · Alice Banda', title: 'MWK 62,500', sub: 'Outstanding balance · statement updated at 19:05', icon: 'ph-receipt', tone: 'emerald', chips: ['3 nights', '2 guests'] },
        {
          k: 'statement',
          rows: [
            { label: 'Accommodation — 2 nights', sub: '30 Oct – 01 Nov · Standard Double', amount: 90000 },
            { label: 'In-room dining', sub: 'Tilapia, juice, mandasi', amount: 6000 },
            { label: 'Bar — poolside', sub: '2 × Carlsberg, 1 × gin & tonic', amount: 6500 },
            { label: 'Payments received', sub: 'Visa · 30 Oct', amount: -40000, tone: 'emerald' }
          ],
          total: { label: 'Balance due', amount: 62500 },
          note: 'Charges are live. Settle at reception or from the guest portal before check-out.'
        },
        { k: 'section', title: 'Need something on your bill?' },
        { k: 'list', items: [
          { icon: 'ph-chat-circle-text', title: 'Query a charge', sub: 'Message reception about this statement', tone: 'indigo', right: 'Message', action: A.modal('guest-request', 'Query a charge') },
          { icon: 'ph-printer', title: 'Email a copy', sub: 'Sent to alice.banda@example.mw', tone: 'slate', right: 'Send', action: A.toast('Statement emailed to you') }
        ] }
      ]
    },

    '/guest/services': {
      dept: 'guest', chrome: 'sub', title: 'Concierge & Requests', back: '/guest',
      blocks: [
        { k: 'hero', kicker: 'Average response 8 minutes', title: 'Concierge & Requests', sub: 'Housekeeping, amenities, transport and anything else. We are a tap away.', icon: 'ph-bell-ringing', tone: 'indigo', chips: ['24-hour desk', 'Ext. 0'] },
        { k: 'section', title: 'Fast requests' },
        {
          k: 'quick', items: [
            { icon: 'ph-towel', label: 'Bath towels', tone: 'indigo', action: A.modal('guest-request', 'Request bath towels') },
            { icon: 'ph-broom', label: 'Room clean', tone: 'emerald', action: A.modal('guest-request', 'Request room clean') },
            { icon: 'ph-drop', label: 'Drinking water', tone: 'amber', action: A.modal('guest-request', 'Request drinking water') },
            { icon: 'ph-clock', label: 'Late check-out', tone: 'orange', action: A.modal('guest-request', 'Request late check-out') },
            { icon: 'ph-car', label: 'Airport transfer', tone: 'slate', action: A.modal('guest-request', 'Request airport transfer') },
            { icon: 'ph-first-aid', label: 'First aid', tone: 'rose', action: A.modal('guest-request', 'Request first aid') }
          ]
        },
        { k: 'section', title: 'Something else?' },
        {
          k: 'form', compact: true,
          fields: [{ label: 'Tell the concierge', type: 'textarea', value: '' }],
          submit: { label: 'Send request', icon: 'ph-paper-plane-tilt', action: A.modal('guest-request', 'Send custom request') }
        },
        { k: 'section', title: 'My requests', action: { text: 'This stay', tone: 'muted' } },
        {
          k: 'list', items: [
            { icon: 'ph-towel', title: 'Extra bath towels', sub: '18:42 · Housekeeping', tone: 'indigo', right: 'On the way', rightTone: 'indigo', action: A.toast('Housekeeping is 10 minutes away') },
            { icon: 'ph-drop', title: 'Two bottles of water', sub: 'Yesterday 16:10 · Delivered', tone: 'emerald', right: 'Done', rightTone: 'emerald', action: A.toast('Delivered yesterday at 16:10') }
          ]
        }
      ]
    },

    '/guest/guide': {
      dept: 'guest', chrome: 'sub', title: 'Resort Guide & Directory', back: '/guest',
      blocks: [
        { k: 'hero', kicker: 'La Luna Hotel & Resort', title: 'Resort Guide', sub: 'Everything within reach — hours, phone numbers and the best of Lake Malawi.', icon: 'ph-map-pin', tone: 'amber', chips: ['Lake Malawi', 'Mangochi'] },
        { k: 'section', title: 'Dining & facility hours' },
        {
          k: 'list', items: [
            { icon: 'ph-bowl-food', title: 'In-room dining', sub: 'Daily · 07:00 – 22:00', tone: 'orange', right: 'Ext. 2', action: A.toast('Dial extension 2 for in-room dining') },
            { icon: 'ph-fork-knife', title: 'Restaurant', sub: 'Breakfast 07:00 · Dinner 18:30', tone: 'amber', right: 'Ext. 2', action: A.toast('Dial extension 2 for the restaurant') },
            { icon: 'ph-wine', title: 'Poolside bar', sub: 'Daily · 11:00 – 23:00', tone: 'purple', right: 'Ext. 1', action: A.toast('Dial extension 1 for the bar') },
            { icon: 'ph-waves', title: 'Pool & beach', sub: 'Open 06:00 – 19:00', tone: 'emerald', right: 'Ext. 0', action: A.toast('Dial extension 0 for pool service') }
          ]
        },
        { k: 'section', title: 'Room phone speed dial' },
        {
          k: 'list', items: [
            { icon: 'ph-phone', title: 'Front Desk Reception', sub: 'Housekeeping, billing and general help', tone: 'slate', right: 'Ext. 0', action: A.toast('Dial 0 on the room phone') },
            { icon: 'ph-cooking-pot', title: 'In-Room Dining Kitchen', sub: 'Orders and dietary questions', tone: 'orange', right: 'Ext. 2', action: A.toast('Dial 2 on the room phone') },
            { icon: 'ph-sparkle', title: 'Housekeeping & Linen', sub: 'Towels, laundry and turndown', tone: 'indigo', right: 'Ext. 3', action: A.toast('Dial 3 on the room phone') },
            { icon: 'ph-briefcase', title: 'General Manager On Duty', sub: 'Feedback and escalations', tone: 'amber', right: 'Ext. 4', action: A.toast('Dial 4 on the room phone') },
            { icon: 'ph-first-aid', title: 'Emergency & First Aid', sub: 'Urgent medical assistance', tone: 'rose', right: 'Ext. 9', action: A.toast('Dial 9 for emergency assistance') }
          ]
        },
        { k: 'section', title: 'Lake Malawi excursions' },
        {
          k: 'grid', items: [
            { icon: 'ph-boat', title: 'Sunset cruise', sub: '17:00 · 90 min · MWK 45,000', tone: 'amber', action: A.toast('Sunset cruise added to your requests') },
            { icon: 'ph-fish', title: 'Snorkel safari', sub: '09:00 · 2 hrs · MWK 38,000', tone: 'emerald', action: A.toast('Snorkel safari added to your requests') },
            { icon: 'ph-tree-palm', title: 'Village walk', sub: '14:00 · 2 hrs · MWK 20,000', tone: 'orange', action: A.toast('Village walk added to your requests') },
            { icon: 'ph-bird', title: 'Birdwatching', sub: '06:30 · 2 hrs · MWK 25,000', tone: 'indigo', action: A.toast('Birdwatching added to your requests') }
          ]
        },
        { k: 'note', icon: 'ph-info', tone: 'slate', title: 'Key hotel policies', text: 'Check-out is 10:00. Noise curfew 22:00. Smoking is permitted only in designated outdoor areas. Outside visitors must sign in at reception.' }
      ]
    }
  });
})();

