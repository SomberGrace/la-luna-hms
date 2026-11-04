/* Modal / lightbox definitions */
(function () {
  'use strict';
  const A = window.A, P = window.PROTO;

  Object.assign(P.modals, {
    'book-room': {
      title: 'Book Room 104', subtitle: 'Standard Double · MWK 45,000 / night', icon: 'ph-bed', tone: 'amber', size: 'sheet',
      blocks: [
        { k: 'form', fields: [
          { label: 'Guest name', type: 'text', value: 'Alice Banda' },
          { label: 'Email', type: 'email', value: 'alice.banda@example.mw' },
          { label: 'Phone', type: 'tel', value: '+265 991 234 567' },
          { label: 'Nights', type: 'number', value: '3' }
        ], submit: { label: 'Confirm booking', icon: 'ph-check-circle', action: A.modal('booking-confirmed', 'Confirm booking') },
        note: '10 – 13 Nov · 2 adults · No deposit required today.' }
      ],
      actions: [{ label: 'Cancel', tone: 'ghost', action: A.close('Cancel booking') }]
    },
    'booking-confirmed': {
      title: 'Booking Confirmed', subtitle: 'Reference BK-2026-0412', icon: 'ph-seal-check', tone: 'emerald', size: 'center',
      blocks: [
        { k: 'success', title: 'You are booked, Alice', text: 'A confirmation has been sent to alice.banda@example.mw. The front desk will issue your room card on arrival.', rows: [
          ['Room', '104 · Standard Double'],
          ['Dates', '10 – 13 Nov 2026'],
          ['Total', 'MWK 135,000'],
          ['Reference', 'BK-2026-0412']
        ] }
      ],
      actions: [
        { label: 'Open guest portal', tone: 'primary', action: A.navc('/book/login', 'Open guest portal') },
        { label: 'Book another room', tone: 'ghost', action: A.close('Book another room') }
      ]
    },
    'room-status': {
      title: 'Room 104', subtitle: 'Standard Double · Floor 1 · Guest: Alice Banda', icon: 'ph-door-open', tone: 'amber', size: 'sheet',
      blocks: [
        { k: 'chips', items: [
          { label: 'Available', action: A.toast('Status set to available') },
          { label: 'Occupied', active: true, action: A.toast('Status set to occupied') },
          { label: 'Dirty', action: A.toast('Status set to dirty') },
          { label: 'Out of service', action: A.toast('Status set to out of service') }
        ] },
        { k: 'form', fields: [
          { label: 'Housekeeping note', type: 'text', value: 'Stayover service completed at 09:15' }
        ], submit: { label: 'Save status', icon: 'ph-check-circle', action: A.close('Save room status') } },
        { k: 'note', icon: 'ph-info', tone: 'slate', title: 'Guest impact', text: 'Marking a room out of service removes it from the public booking page immediately.' }
      ],
      actions: [{ label: 'Cancel', tone: 'ghost', action: A.close('Cancel room status') }]
    },
    'walkin-success': {
      title: 'Registration Complete', subtitle: 'Folio F-1043 opened', icon: 'ph-check-circle', tone: 'emerald', size: 'center',
      blocks: [
        { k: 'success', title: 'Memory Tembo is checked in', text: 'Room 103 is now occupied and the deposit has been applied to the folio.', rows: [
          ['Room', '103 · Standard Single'],
          ['Nights', '2'],
          ['Deposit', 'MWK 45,000'],
          ['Folio', 'F-1043']
        ] }
      ],
      actions: [
        { label: 'Open folio', tone: 'primary', action: A.navc('/frontdesk/folio', 'Open new folio') },
        { label: 'New walk-in', tone: 'ghost', action: A.close('Start another walk-in') }
      ]
    },
    'checkin-success': {
      title: 'Check-In Complete', subtitle: 'Room 104 · Alice Banda', icon: 'ph-check-circle', tone: 'emerald', size: 'center',
      blocks: [
        { k: 'success', title: 'Room 104 is now occupied', text: 'The guest record, booking and folio were created together. Two room cards were issued.', rows: [
          ['Reservation', 'RES-2026-0184'],
          ['Nights', '3'],
          ['Deposit', 'MWK 45,000 · Visa'],
          ['Folio', 'F-1042']
        ] }
      ],
      actions: [
        { label: 'Open folio', tone: 'primary', action: A.navc('/frontdesk/folio', 'Open checked-in folio') },
        { label: 'Back to hub', tone: 'ghost', action: A.navc('/frontdesk', 'Back to front desk hub') }
      ]
    },
    'payment-receipt': {
      title: 'Payment Receipt', subtitle: 'Folio F-1042 · Room 104', icon: 'ph-receipt', tone: 'emerald', size: 'sheet',
      blocks: [
        { k: 'success', compact: true, title: 'MWK 40,000 received', text: 'Visa ···· 4242 · approved 19:12', rows: [
          ['Guest', 'Alice Banda'],
          ['Method', 'Visa'],
          ['Reference', 'PAY-2026-0912'],
          ['Balance after', 'MWK 22,500']
        ] },
        { k: 'note', icon: 'ph-printer', tone: 'slate', title: 'Receipt', text: 'A copy is emailed to the guest and attached to the folio.', action: A.toast('Receipt emailed to the guest') }
      ],
      actions: [
        { label: 'Proceed to checkout', tone: 'primary', action: A.navc('/frontdesk/checkout', 'Proceed to checkout') },
        { label: 'Done', tone: 'ghost', action: A.close('Close receipt') }
      ]
    },
    'checkout-success': {
      title: 'Check-Out Successful', subtitle: 'Room 104 released', icon: 'ph-check-circle', tone: 'emerald', size: 'center',
      blocks: [
        { k: 'success', title: 'Alice Banda has checked out', text: 'The folio is closed and housekeeping has been notified. The room will be sellable once cleaned.', rows: [
          ['Folio', 'F-1042 · closed'],
          ['Settled', 'MWK 62,500'],
          ['Room status', 'Dirty'],
          ['Next arrival', '04 Nov · 14:00']
        ] }
      ],
      actions: [
        { label: 'Open housekeeping', tone: 'primary', action: A.navc('/frontdesk/housekeeping', 'Open housekeeping') },
        { label: 'Back to hub', tone: 'ghost', action: A.navc('/frontdesk', 'Back to front desk hub') }
      ]
    },
    'extension-success': {
      title: 'Stay Extended', subtitle: 'Room 104 · new departure 04 Nov', icon: 'ph-calendar-check', tone: 'emerald', size: 'center',
      blocks: [
        { k: 'success', title: 'One extra night added', text: 'MWK 45,000 has been posted to the folio and the room is held for the new departure date.', rows: [
          ['Previous departure', '03 Nov 2026'],
          ['New departure', '04 Nov 2026'],
          ['Additional charge', 'MWK 45,000'],
          ['New balance', 'MWK 107,500']
        ] }
      ],
      actions: [
        { label: 'Back to folio', tone: 'primary', action: A.navc('/frontdesk/folio', 'Back to folio') },
        { label: 'Done', tone: 'ghost', action: A.close('Close extension') }
      ]
    },
    'add-menu-item': {
      title: 'Add Menu Item', subtitle: 'New dish · linked to kitchen BOM', icon: 'ph-book-open', tone: 'orange', size: 'sheet',
      blocks: [
        { k: 'form', fields: [
          { label: 'Item name', type: 'text', value: 'Lake Malawi Fish Curry' },
          { label: 'Category', type: 'select', value: 'Mains', options: ['Mains', 'Breakfast', 'Snacks', 'Drinks'] },
          { label: 'Selling price (MWK)', type: 'number', value: '6500' },
          { label: 'Food cost (MWK)', type: 'number', value: '2600' }
        ], submit: { label: 'Add to menu', icon: 'ph-plus', action: A.close('Add menu item') },
        note: 'The item stays unavailable until ingredients are linked to a recipe.' }
      ],
      actions: [{ label: 'Cancel', tone: 'ghost', action: A.close('Cancel menu item') }]
    },
    'new-docket': {
      title: 'New Docket', subtitle: 'Table or room service', icon: 'ph-clipboard-text', tone: 'orange', size: 'sheet',
      blocks: [
        { k: 'form', fields: [
          { label: 'Table / room', type: 'select', value: 'Table 7', options: ['Table 2', 'Table 5', 'Table 7', 'Room 104', 'Room 205'] },
          { label: 'Guests', type: 'number', value: '2' },
          { label: 'Items', type: 'textarea', value: '2 × Grilled Tilapia, 1 × Nsima & Chicken, 2 × Gin & Tonic' }
        ], submit: { label: 'Fire to kitchen', icon: 'ph-fire', action: A.close('Fire docket') },
        note: 'Room-service dockets post to the guest folio; table dockets settle at the till.' }
      ],
      actions: [{ label: 'Cancel', tone: 'ghost', action: A.close('Cancel docket') }]
    },
    'add-bar-item': {
      title: 'Update Bar Stock', subtitle: 'Carlsberg Green Lager 500ml', icon: 'ph-package', tone: 'purple', size: 'sheet',
      blocks: [
        { k: 'form', fields: [
          { label: 'Item name', type: 'text', value: 'Carlsberg Green Lager 500ml' },
          { label: 'Category', type: 'select', value: 'Beer', options: ['Beer', 'Spirits', 'Wine', 'Mixer'] },
          { label: 'Quantity on hand', type: 'number', value: '24' },
          { label: 'Reorder at', type: 'number', value: '12' }
        ], submit: { label: 'Save stock', icon: 'ph-check-circle', action: A.close('Save bar stock') },
        note: 'Stock movements are written to the bar ledger for the nightly variance check.' }
      ],
      actions: [{ label: 'Cancel', tone: 'ghost', action: A.close('Cancel bar stock') }]
    },
    'add-ingredient': {
      title: 'Add Raw Ingredient', subtitle: 'Kitchen stock · cold room or dry store', icon: 'ph-carrot', tone: 'emerald', size: 'sheet',
      blocks: [
        { k: 'form', fields: [
          { label: 'Ingredient', type: 'text', value: 'Fresh Tilapia (Lake Malawi)' },
          { label: 'Storage', type: 'select', value: 'Cold room', options: ['Cold room', 'Freezer', 'Dry store'] },
          { label: 'Quantity', type: 'number', value: '10' },
          { label: 'Unit', type: 'select', value: 'pieces', options: ['kg', 'litres', 'pieces', 'birds'] },
          { label: 'Reorder at', type: 'number', value: '6' }
        ], submit: { label: 'Save ingredient', icon: 'ph-check-circle', action: A.close('Save ingredient') },
        note: 'Ingredient counts drive menu availability through the recipe bill of materials.' }
      ],
      actions: [{ label: 'Cancel', tone: 'ghost', action: A.close('Cancel ingredient') }]
    },
    'guest-tray': {
      title: 'In-Room Dining Tray', subtitle: 'Charged to Room 104', icon: 'ph-basket', tone: 'orange', size: 'sheet',
      blocks: [
        { k: 'statement', rows: [
          { label: '2 × Grilled Lake Tilapia', sub: 'MWK 4,500 each', amount: 9000 },
          { label: '1 × Fresh Mango Juice', sub: 'MWK 1,500', amount: 1500 }
        ], total: { label: 'Tray total', amount: 10500 } },
        { k: 'form', compact: true, fields: [
          { label: 'Delivery note', type: 'text', value: 'Please knock quietly — baby sleeping.' }
        ], submit: { label: 'Place order', icon: 'ph-paper-plane-tilt', action: A.close('Place dining order') },
        note: 'Average delivery 20–35 minutes. You can track the order from the guest portal.' }
      ],
      actions: [{ label: 'Add more', tone: 'ghost', action: A.close('Add more dishes') }]
    },
    'guest-request': {
      title: 'Concierge Request', subtitle: 'Average response 8 minutes', icon: 'ph-bell-ringing', tone: 'indigo', size: 'sheet',
      blocks: [
        { k: 'chips', items: [
          { label: 'Bath towels', active: true, action: A.toast('Request type: bath towels') },
          { label: 'Room clean', action: A.toast('Request type: room clean') },
          { label: 'Drinking water', action: A.toast('Request type: drinking water') },
          { label: 'Late check-out', action: A.toast('Request type: late check-out') }
        ] },
        { k: 'form', fields: [
          { label: 'When', type: 'select', value: 'As soon as possible', options: ['As soon as possible', 'Within the hour', 'This evening', 'Tomorrow morning'] },
          { label: 'Notes for the team', type: 'textarea', value: 'Two extra towels please, and one hand towel.' }
        ], submit: { label: 'Send request', icon: 'ph-paper-plane-tilt', action: A.close('Send concierge request') },
        note: 'You will get a notification when the team accepts the request.' }
      ],
      actions: [{ label: 'Cancel', tone: 'ghost', action: A.close('Cancel request') }]
    },
    'wifi': {
      title: 'Hotel WiFi', subtitle: 'La Luna guest network', icon: 'ph-wifi-high', tone: 'amber', size: 'sheet',
      blocks: [
        { k: 'wifi', network: 'LaLuna-Guest-WiFi', password: 'laluna2026', note: 'Covers all rooms, the restaurant, the bar and the pool deck. No password required on the beach.' },
        { k: 'note', icon: 'ph-info', tone: 'slate', title: 'Streaming', text: 'Smart TVs support casting from the guest portal. Ask reception for a HDMI cable if needed.' }
      ],
      actions: [
        { label: 'Copy password', tone: 'primary', action: A.toast('WiFi password copied') },
        { label: 'Done', tone: 'ghost', action: A.close('Close WiFi') }
      ]
    },
    'guest-request-success': {
      title: 'Request Sent', subtitle: 'Housekeeping notified', icon: 'ph-check-circle', tone: 'emerald', size: 'center',
      blocks: [
        { k: 'success', title: 'We are on it', text: 'Your request is with the housekeeping team. You will get a notification when it is on the way.', rows: [
          ['Request', 'Extra bath towels'],
          ['Room', '104'],
          ['ETA', '8 – 12 minutes']
        ] }
      ],
      actions: [{ label: 'Done', tone: 'primary', action: A.close('Close request') }]
    },
    'audit-review': {
      title: 'Review Bar Shift Close', subtitle: 'Submitted by Peter K. · MWK 86,500', icon: 'ph-shield-check', tone: 'indigo', size: 'sheet',
      blocks: [
        { k: 'statement', rows: [
          { label: 'Cash on hand', sub: 'Counted at close', amount: 42000 },
          { label: 'Room charges', sub: '3 rooms', amount: 31500 },
          { label: 'Card takings', sub: 'Batch 8802', amount: 13000 },
          { label: 'Expected variance', sub: '0.4% · within 1% tolerance', amount: 0, tone: 'emerald' }
        ], total: { label: 'Total banked', amount: 86500 } },
        { k: 'form', compact: true, fields: [{ label: 'Approval note', type: 'text', value: 'Variance accepted. Ledger finalised.' }],
          submit: { label: 'Approve & sign', icon: 'ph-check-circle', action: A.close('Approve audit') } },
        { k: 'note', icon: 'ph-warning', tone: 'orange', title: 'Audit trail', text: 'Your name and the timestamp are written to the approval log. This cannot be edited later.' }
      ],
      actions: [
        { label: 'Reject', tone: 'danger', action: A.close('Reject audit') },
        { label: 'Cancel', tone: 'ghost', action: A.close('Cancel audit review') }
      ]
    },
    'shift-close': {
      title: 'Close Shift', subtitle: 'Morning shift · Grace M.', icon: 'ph-lock', tone: 'amber', size: 'sheet',
      blocks: [
        { k: 'form', fields: [
          { label: 'Cash float checked (MWK)', type: 'number', value: '20000' },
          { label: 'Handover note', type: 'textarea', value: 'Queue clear, one folio balance outstanding. Room 401 blocked for plumbing.' }
        ], submit: { label: 'Close and hand over', icon: 'ph-check-circle', action: A.close('Close shift') },
        note: 'Closing locks the float and sends the shift report to the manager consolidation queue.' }
      ],
      actions: [{ label: 'Cancel', tone: 'ghost', action: A.close('Cancel shift close') }]
    },
    'sale-recorded': {
      title: 'Sale Recorded', subtitle: 'Charged to Room 205', icon: 'ph-check-circle', tone: 'emerald', size: 'center',
      blocks: [
        { k: 'success', compact: true, title: 'MWK 5,000 posted', text: '2 × Guinness Stout · stock updated · folio charge posted.', rows: [
          ['Payment', 'Room 205'],
          ['Reference', 'BAR-2026-0731'],
          ['Stock left', '24 bottles']
        ] }
      ],
      actions: [{ label: 'Done', tone: 'primary', action: A.close('Close sale') }]
    }
  });
})();
