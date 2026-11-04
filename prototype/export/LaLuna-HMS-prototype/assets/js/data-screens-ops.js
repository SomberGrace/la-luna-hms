/* Screens: restaurant, bar, kitchen, manager */
(function () {
  'use strict';
  const A = window.A, P = window.PROTO;

  Object.assign(P.screens, {
    /* ---------------- Restaurant ---------------- */
    '/restaurant': {
      dept: 'restaurant', chrome: 'brand', title: 'Restaurant Hub',
      blocks: [
        { k: 'hero', kicker: 'Lunch service · 42 covers', title: 'Dining & Dockets', sub: 'Restaurant floor, kitchen tickets and menu costing in one station.', icon: 'ph-fork-knife', tone: 'orange', chips: ['Service open', '5 open dockets'] },
        { k: 'stats', items: [
          { label: 'Covers', value: '42', note: 'Lunch service', tone: 'orange' },
          { label: 'Open dockets', value: '5', note: '2 in kitchen', tone: 'amber' },
          { label: 'Revenue', value: '168,000', note: 'MWK today', tone: 'emerald' },
          { label: 'Avg ticket', value: '4,000', note: 'MWK per cover', tone: 'indigo' }
        ] },
        { k: 'section', title: 'Recent dockets', action: { text: 'View all', tone: 'link', action: A.nav('/restaurant/dockets', 'View all dockets') } },
        { k: 'list', items: [
          { icon: 'ph-clipboard-text', title: 'Docket #0184 · Table 7', sub: '2 × Tilapia, 1 × Nsima · MWK 12,500', tone: 'orange', right: 'Preparing', rightTone: 'amber', action: A.nav('/restaurant/dockets', 'Open docket 0184') },
          { icon: 'ph-clipboard-text', title: 'Docket #0183 · Room 104', sub: 'Tilapia, juice, mandasi · MWK 8,000', tone: 'indigo', right: 'Served', rightTone: 'emerald', action: A.nav('/restaurant/dockets', 'Open docket 0183') },
          { icon: 'ph-clipboard-text', title: 'Docket #0182 · Table 2', sub: '2 × Steak, 2 × Gin & tonic · MWK 22,000', tone: 'orange', right: 'Paid', rightTone: 'slate', action: A.nav('/restaurant/dockets', 'Open docket 0182') }
        ] },
        { k: 'section', title: 'Quick actions' },
        { k: 'quick', items: [
          { icon: 'ph-plus', label: 'New docket', tone: 'orange', action: A.modal('new-docket', 'New docket') },
          { icon: 'ph-book-open', label: 'Menu', tone: 'amber', action: A.nav('/restaurant/menu', 'Menu') },
          { icon: 'ph-clipboard-text', label: 'Dockets', tone: 'indigo', action: A.nav('/restaurant/dockets', 'Dockets') },
          { icon: 'ph-chart-line-up', label: 'Reports', tone: 'emerald', action: A.nav('/restaurant/reports', 'Restaurant reports') }
        ] },
        { k: 'note', icon: 'ph-warning-circle', tone: 'orange', title: 'Chambo stock is low', text: 'Only 4 portions of Lake Malawi chambo remain. The kitchen has been notified.', action: A.nav('/kitchen/stock', 'Open kitchen stock') }
      ]
    },

    '/restaurant/menu': {
      dept: 'restaurant', chrome: 'sub', title: 'Menu Items', back: '/restaurant',
      blocks: [
        { k: 'hero', kicker: 'Live menu · kitchen linked', title: 'Menu Items', sub: 'Availability follows the ingredient stock count from the kitchen.', icon: 'ph-book-open', tone: 'amber', chips: ['9 dishes', '2 low stock'] },
        { k: 'chips', items: [
          { label: 'All 9', active: true, action: A.toast('Showing all menu items') },
          { label: 'Mains', action: A.toast('Mains selected') },
          { label: 'Breakfast', action: A.toast('Breakfast selected') },
          { label: 'Snacks', action: A.toast('Snacks selected') },
          { label: 'Drinks', action: A.toast('Drinks selected') }
        ] },
        { k: 'menu', items: [
          { name: 'Nsima with Stewed Chicken', desc: 'Maize meal, tomato chicken stew, seasonal greens.', price: 3500, tag: 'Mains', tone: 'orange', action: A.toast('Nsima available · 24 portions') },
          { name: 'Grilled Lake Tilapia (Chambo)', desc: 'Whole chambo, lemon butter, hand-cut chips.', price: 4500, tag: 'Low stock', tone: 'rose', action: A.toast('Only 4 chambo portions left') },
          { name: 'Prime Beef Tenderloin', desc: 'Pan-seared tenderloin, garlic mash, pepper sauce.', price: 8500, tag: 'Mains', tone: 'indigo', action: A.toast('Steak available · 12 portions') },
          { name: 'Full English Resort Breakfast', desc: 'Eggs, bacon, sausage, beans, grilled tomato, sourdough.', price: 5500, tag: 'Breakfast', tone: 'amber', action: A.toast('Breakfast available until 10:30') },
          { name: 'Warm Mandasi & Local Honey', desc: 'Fluffy donuts, Zomba honey, berry jam.', price: 800, tag: 'Snacks', tone: 'amber', action: A.toast('Mandasi available · 30 portions') },
          { name: 'Malawi Gin & Tonic', desc: 'Malawi Gin, tonic, fresh lime.', price: 2500, tag: 'Drinks', tone: 'purple', action: A.toast('Bar stock linked · 18 serves') }
        ] },
        { k: 'form', compact: true, fields: [{ label: 'Add a menu item', type: 'text', value: '' }],
          submit: { label: 'Add menu item', icon: 'ph-plus', action: A.modal('add-menu-item', 'Add menu item') } }
      ]
    },

    '/restaurant/dockets': {
      dept: 'restaurant', chrome: 'sub', title: 'Dockets', back: '/restaurant',
      blocks: [
        { k: 'hero', kicker: 'Live service', title: 'Restaurant Operations', sub: 'Every table and room-service ticket, from order to settlement.', icon: 'ph-clipboard-text', tone: 'orange', chips: ['5 open', '2 preparing'] },
        { k: 'stats', items: [
          { label: 'Open', value: '5', note: 'Not settled', tone: 'amber' },
          { label: 'Preparing', value: '2', note: 'In the kitchen', tone: 'orange' },
          { label: 'Served', value: '14', note: 'Today', tone: 'emerald' }
        ] },
        { k: 'section', title: 'Open dockets', action: { text: 'New docket', tone: 'link', action: A.modal('new-docket', 'New docket') } },
        { k: 'tickets', items: [
          { id: '#0184', meta: 'Table 7 · 12:42', status: 'Preparing', tone: 'amber', total: 12500,
            lines: [{ qty: 2, name: 'Grilled Tilapia', note: 'No chilli' }, { qty: 1, name: 'Nsima & Chicken', note: '' }],
            action: A.toast('Docket #0184 sent to kitchen') },
          { id: '#0185', meta: 'Room 205 · 12:48', status: 'New', tone: 'orange', total: 8000,
            lines: [{ qty: 1, name: 'Beef Tenderloin', note: 'Medium rare' }, { qty: 1, name: 'Mango Juice', note: '' }],
            action: A.toast('Docket #0185 accepted') },
          { id: '#0186', meta: 'Table 2 · 12:55', status: 'New', tone: 'orange', total: 22000,
            lines: [{ qty: 2, name: 'Prime Steak', note: 'One well done' }, { qty: 2, name: 'Gin & Tonic', note: '' }],
            action: A.toast('Docket #0186 accepted') }
        ] },
        { k: 'note', icon: 'ph-info', tone: 'slate', title: 'Room service', text: 'Dockets raised against a room post directly to the guest folio. Table dockets settle at the restaurant till.' }
      ]
    },

    '/restaurant/reports': {
      dept: 'restaurant', chrome: 'sub', title: 'Restaurant Report', back: '/restaurant',
      blocks: [
        { k: 'hero', kicker: 'Lunch service · 01 Nov', title: 'Restaurant Report', sub: 'Covers, category mix and the dishes carrying the service.', icon: 'ph-chart-line-up', tone: 'orange', chips: ['42 covers', 'MWK 168,000'] },
        { k: 'stats', items: [
          { label: 'Covers', value: '42', note: 'Lunch', tone: 'orange' },
          { label: 'Revenue', value: '168,000', note: 'MWK', tone: 'emerald' },
          { label: 'Avg ticket', value: '4,000', note: 'MWK', tone: 'amber' },
          { label: 'Room charges', value: '31,500', note: 'Charged to folios', tone: 'indigo' }
        ] },
        { k: 'section', title: 'Category breakdown' },
        { k: 'bars', items: [
          { label: 'Mains', value: 'MWK 98,000', pct: 58, tone: 'orange' },
          { label: 'Drinks', value: 'MWK 38,000', pct: 23, tone: 'purple' },
          { label: 'Breakfast', value: 'MWK 20,000', pct: 12, tone: 'amber' },
          { label: 'Snacks', value: 'MWK 12,000', pct: 7, tone: 'emerald' }
        ] },
        { k: 'section', title: 'Top dishes' },
        { k: 'table', head: ['Dish', 'Sold', 'Revenue'], rows: [
          ['Grilled Tilapia', '14', '63,000'],
          ['Nsima & Chicken', '11', '38,500'],
          ['Beef Tenderloin', '6', '51,000'],
          ['Gin & Tonic', '9', '22,500']
        ] },
        { k: 'note', icon: 'ph-printer', tone: 'slate', title: 'Send to manager', text: 'Restaurant totals feed the manager consolidation report at shift close.', action: A.toast('Restaurant report sent to manager queue') }
      ]
    },

    /* ---------------- Bar ---------------- */
    '/bar': {
      dept: 'bar', chrome: 'brand', title: 'Bar Hub',
      blocks: [
        { k: 'hero', kicker: 'Poolside bar · Open', title: 'Bar Inventory & Audit', sub: 'Sales, stock movement and the nightly variance check.', icon: 'ph-wine', tone: 'purple', chips: ['Open 11:00 – 23:00', '1 low stock'] },
        { k: 'stats', items: [
          { label: 'Sales', value: '86,500', note: 'MWK today', tone: 'purple' },
          { label: 'Stock value', value: '412,000', note: 'MWK on hand', tone: 'emerald' },
          { label: 'Low stock', value: '1', note: 'Carlsberg 500ml', tone: 'rose' },
          { label: 'Variance', value: '0.4%', note: 'Within tolerance', tone: 'amber' }
        ] },
        { k: 'section', title: 'Low stock alerts', action: { text: 'Manage stock', tone: 'link', action: A.nav('/bar/stock', 'Manage bar stock') } },
        { k: 'list', items: [
          { icon: 'ph-beer-stein', title: 'Carlsberg Green Lager 500ml', sub: '8 bottles left · reorder at 12', tone: 'rose', right: 'Restock', action: A.toast('Restock request sent to supplier') },
          { icon: 'ph-wine', title: 'Malawi Gin 750ml', sub: '3 bottles left · reorder at 4', tone: 'orange', right: 'Restock', action: A.toast('Restock request sent to supplier') }
        ] },
        { k: 'section', title: 'Quick actions' },
        { k: 'quick', items: [
          { icon: 'ph-cash-register', label: 'Record sale', tone: 'purple', action: A.nav('/bar/sales', 'Record bar sale') },
          { icon: 'ph-package', label: 'Stock audit', tone: 'emerald', action: A.nav('/bar/stock', 'Stock audit') },
          { icon: 'ph-chart-bar', label: 'Reports', tone: 'amber', action: A.nav('/bar/reports', 'Bar reports') }
        ] },
        { k: 'note', icon: 'ph-info', tone: 'slate', title: 'Tonight', text: 'Live band at 20:00. Expect a busy late service — pre-chill two crates of lager.' }
      ]
    },

    '/bar/stock': {
      dept: 'bar', chrome: 'sub', title: 'Stock Audit', back: '/bar',
      blocks: [
        { k: 'hero', kicker: 'Counted 07:45', title: 'Stock Audit', sub: 'Opening count, live levels and reorder points for every bar item.', icon: 'ph-package', tone: 'emerald', chips: ['18 items', '1 low', '0 variance'] },
        { k: 'chips', items: [
          { label: 'All items', active: true, action: A.toast('Showing all bar items') },
          { label: 'Low stock', action: A.toast('Filtered: low stock') },
          { label: 'Spirits', action: A.toast('Filtered: spirits') },
          { label: 'Beer', action: A.toast('Filtered: beer') }
        ] },
        { k: 'list', items: [
          { icon: 'ph-beer-stein', title: 'Carlsberg Green Lager 500ml', sub: 'Beer · reorder at 12', tone: 'rose', right: '8 left', rightTone: 'rose', action: A.modal('add-bar-item', 'Update lager stock') },
          { icon: 'ph-wine', title: 'Malawi Gin 750ml', sub: 'Spirits · reorder at 4', tone: 'orange', right: '3 left', rightTone: 'orange', action: A.modal('add-bar-item', 'Update gin stock') },
          { icon: 'ph-wine', title: 'KWV Sauvignon Blanc', sub: 'Wine · reorder at 6', tone: 'emerald', right: '14 left', rightTone: 'emerald', action: A.modal('add-bar-item', 'Update wine stock') },
          { icon: 'ph-beer-stein', title: 'Guinness Stout 440ml', sub: 'Beer · reorder at 12', tone: 'emerald', right: '24 left', rightTone: 'emerald', action: A.modal('add-bar-item', 'Update stout stock') },
          { icon: 'ph-flask', title: 'Coca-Cola 300ml', sub: 'Mixer · reorder at 24', tone: 'amber', right: '36 left', rightTone: 'emerald', action: A.modal('add-bar-item', 'Update mixer stock') },
          { icon: 'ph-flask', title: 'Tonic Water 300ml', sub: 'Mixer · reorder at 24', tone: 'amber', right: '9 left', rightTone: 'orange', action: A.modal('add-bar-item', 'Update tonic stock') }
        ] },
        { k: 'form', compact: true, fields: [{ label: 'Add a bar item', type: 'text', value: '' }],
          submit: { label: 'Add bar item', icon: 'ph-plus', action: A.modal('add-bar-item', 'Add bar item') } },
        { k: 'note', icon: 'ph-warning', tone: 'orange', title: 'Variance check', text: 'A variance above 1% requires a manager sign-off before the shift can close.' }
      ]
    },

    '/bar/sales': {
      dept: 'bar', chrome: 'sub', title: 'Record Sales', back: '/bar',
      blocks: [
        { k: 'hero', kicker: 'Session total · MWK 86,500', title: 'Record Sales', sub: 'Post a drink sale, charge it to a room or take payment at the bar.', icon: 'ph-cash-register', tone: 'purple', chips: ['Open tab', 'Stock updates live'] },
        { k: 'form', title: 'New sale entry', fields: [
          { label: 'Item', type: 'select', value: 'Guinness Stout 440ml · MWK 2,500', options: ['Guinness Stout 440ml · MWK 2,500', 'Carlsberg Green Lager 500ml · MWK 2,000', 'Malawi Gin & Tonic · MWK 3,500', 'KWV Sauvignon Blanc glass · MWK 4,000'] },
          { label: 'Quantity', type: 'number', value: '2' },
          { label: 'Payment', type: 'select', value: 'Room 104', options: ['Cash', 'Room 104', 'Room 205', 'Card'] },
          { label: 'Note', type: 'text', value: 'Poolside · table 3' }
        ], submit: { label: 'Record sale · MWK 5,000', icon: 'ph-check-circle', action: A.modal('sale-recorded', 'Record sale') } },
        { k: 'section', title: 'Sales recorded this session' },
        { k: 'list', items: [
          { icon: 'ph-beer-stein', title: '2 × Guinness Stout', sub: '12:51 · Room 205', tone: 'purple', right: 'MWK 5,000', action: A.toast('Sale posted to room 205 folio') },
          { icon: 'ph-wine', title: '1 × Gin & Tonic', sub: '12:38 · Cash', tone: 'emerald', right: 'MWK 3,500', action: A.toast('Cash sale recorded') },
          { icon: 'ph-wine', title: '2 × Sauvignon Blanc', sub: '12:20 · Card', tone: 'indigo', right: 'MWK 8,000', action: A.toast('Card sale recorded') }
        ] }
      ]
    },

    '/bar/reports': {
      dept: 'bar', chrome: 'sub', title: 'Bar Daily Report', back: '/bar',
      blocks: [
        { k: 'hero', kicker: '01 Nov 2026 · Poolside', title: 'Bar Daily Report', sub: 'Sales, stock movement and the nightly ledger in one view.', icon: 'ph-chart-bar', tone: 'purple', chips: ['MWK 86,500', '0.4% variance'] },
        { k: 'stats', items: [
          { label: 'Sales', value: '86,500', note: 'MWK today', tone: 'purple' },
          { label: 'Cost of sales', value: '31,200', note: '36% pour cost', tone: 'amber' },
          { label: 'Profit', value: '55,300', note: 'MWK gross', tone: 'emerald' }
        ] },
        { k: 'section', title: 'Item ledger' },
        { k: 'table', head: ['Item', 'Open', 'Sold', 'Close'], rows: [
          ['Carlsberg 500ml', '24', '16', '8'],
          ['Guinness 440ml', '36', '12', '24'],
          ['Malawi Gin', '5', '2', '3'],
          ['KWV Sauvignon', '16', '2', '14']
        ] },
        { k: 'section', title: 'Settlement mix' },
        { k: 'bars', items: [
          { label: 'Cash', value: 'MWK 42,000', pct: 49, tone: 'emerald' },
          { label: 'Room charge', value: 'MWK 31,500', pct: 36, tone: 'indigo' },
          { label: 'Card', value: 'MWK 13,000', pct: 15, tone: 'purple' }
        ] },
        { k: 'note', icon: 'ph-shield-check', tone: 'slate', title: 'Manager sign-off', text: 'Variance is within tolerance. Send the ledger to the manager audit queue at close.', action: A.toast('Bar ledger sent to manager queue') }
      ]
    },

    /* ---------------- Kitchen ---------------- */
    '/kitchen': {
      dept: 'kitchen', chrome: 'brand', title: 'Kitchen Orders',
      blocks: [
        { k: 'hero', kicker: 'Live KDS · 3 tickets', title: 'Kitchen Orders', sub: 'Tickets land here the moment a docket is fired. Bump when the plate leaves.', icon: 'ph-kanban', tone: 'emerald', chips: ['3 new', '2 preparing', '4 ready'] },
        { k: 'stats', items: [
          { label: 'New', value: '3', note: 'Awaiting pickup', tone: 'orange' },
          { label: 'Preparing', value: '2', note: 'On the pass', tone: 'amber' },
          { label: 'Ready', value: '4', note: 'Call service', tone: 'emerald' }
        ] },
        { k: 'section', title: 'Current tickets', action: { text: 'Auto-refresh 10s', tone: 'muted' } },
        { k: 'tickets', items: [
          { id: '#0184', meta: 'Table 7 · fired 12:42', status: 'Preparing', tone: 'amber', total: 12500,
            lines: [{ qty: 2, name: 'Grilled Tilapia', note: 'No chilli · extra lemon' }, { qty: 1, name: 'Nsima & Chicken', note: 'Extra gravy' }],
            action: A.toast('Ticket #0184 bumped to ready') },
          { id: '#0185', meta: 'Room 205 · fired 12:48', status: 'New', tone: 'orange', total: 8000,
            lines: [{ qty: 1, name: 'Beef Tenderloin', note: 'Medium rare' }, { qty: 1, name: 'Mango Juice', note: 'No ice' }],
            action: A.toast('Ticket #0185 started') },
          { id: '#0186', meta: 'Table 2 · fired 12:55', status: 'New', tone: 'orange', total: 22000,
            lines: [{ qty: 2, name: 'Prime Steak', note: 'One well done' }, { qty: 2, name: 'Gin & Tonic', note: '' }],
            action: A.toast('Ticket #0186 started') }
        ] },
        { k: 'note', icon: 'ph-fire', tone: 'orange', title: 'Low stock: chambo', text: 'Four portions remain. Tell the restaurant before the next lake order is taken.', action: A.nav('/kitchen/stock', 'Open kitchen stock') }
      ]
    },

    '/kitchen/stock': {
      dept: 'kitchen', chrome: 'sub', title: 'Kitchen Stock', back: '/kitchen',
      blocks: [
        { k: 'hero', kicker: 'Counted 07:30', title: 'Kitchen Stock', sub: 'Raw ingredients, reorder points and the cost behind every plate.', icon: 'ph-carrot', tone: 'emerald', chips: ['6 ingredients', '2 low'] },
        { k: 'list', items: [
          { icon: 'ph-package', title: 'Maize Flour (Gramil)', sub: 'Dry store · reorder at 10 kg', tone: 'emerald', right: '42 kg', rightTone: 'emerald', action: A.modal('add-ingredient', 'Update maize flour') },
          { icon: 'ph-fish', title: 'Fresh Tilapia (Lake Malawi)', sub: 'Cold room · reorder at 6 pieces', tone: 'rose', right: '4 pcs', rightTone: 'rose', action: A.modal('add-ingredient', 'Update tilapia') },
          { icon: 'ph-package', title: 'Whole Chickens', sub: 'Freezer · reorder at 8 birds', tone: 'orange', right: '7 birds', rightTone: 'orange', action: A.modal('add-ingredient', 'Update chicken') },
          { icon: 'ph-flask', title: 'Cooking Oil (Sunseed)', sub: 'Dry store · reorder at 10 L', tone: 'emerald', right: '18 L', rightTone: 'emerald', action: A.modal('add-ingredient', 'Update cooking oil') },
          { icon: 'ph-carrot', title: 'Potatoes (Irish)', sub: 'Cold room · reorder at 15 kg', tone: 'emerald', right: '32 kg', rightTone: 'emerald', action: A.modal('add-ingredient', 'Update potatoes') },
          { icon: 'ph-carrot', title: 'Tomatoes & Onions', sub: 'Cold room · reorder at 10 kg', tone: 'amber', right: '9 kg', rightTone: 'orange', action: A.modal('add-ingredient', 'Update tomatoes') }
        ] },
        { k: 'form', compact: true, fields: [{ label: 'Add a raw ingredient', type: 'text', value: '' }],
          submit: { label: 'Add raw ingredient', icon: 'ph-plus', action: A.modal('add-ingredient', 'Add raw ingredient') } },
        { k: 'note', icon: 'ph-info', tone: 'slate', title: 'Stock feeds the menu', text: 'Menu availability is derived from these counts through the recipe bill of materials.' }
      ]
    },

    '/kitchen/recipes': {
      dept: 'kitchen', chrome: 'sub', title: 'Recipe BOM Costing', back: '/kitchen',
      blocks: [
        { k: 'hero', kicker: 'Bill of materials', title: 'Recipe BOM Costing', sub: 'Every plate costed from live ingredient prices, with the margin left on the pass.', icon: 'ph-flask', tone: 'indigo', chips: ['4 recipes', 'Avg margin 62%'] },
        { k: 'list', items: [
          { icon: 'ph-bowl-food', title: 'Nsima & Stewed Chicken', sub: 'Flour 0.35 kg · chicken 0.25 bird · tomato 0.15 kg · oil 0.05 L', tone: 'orange', right: 'Cost 1,850', action: A.toast('Nsima recipe opened') },
          { icon: 'ph-fish', title: 'Grilled Chambo', sub: 'Tilapia 1 pc · oil & spices 0.08 L · salad 0.1 kg', tone: 'emerald', right: 'Cost 2,200', action: A.toast('Chambo recipe opened') },
          { icon: 'ph-fire', title: 'Crispy Chips', sub: 'Potatoes 0.45 kg · oil 0.1 L · sauce 0.05 L', tone: 'amber', right: 'Cost 780', action: A.toast('Chips recipe opened') },
          { icon: 'ph-knife', title: 'Beef Tenderloin', sub: 'Beef 0.35 kg · potatoes 0.3 kg · butter & herbs 0.05 kg', tone: 'indigo', right: 'Cost 3,500', action: A.toast('Tenderloin recipe opened') }
        ] },
        { k: 'section', title: 'Margin check' },
        { k: 'bars', items: [
          { label: 'Nsima & Chicken', value: '47% food cost', pct: 47, tone: 'orange' },
          { label: 'Grilled Chambo', value: '49% food cost', pct: 49, tone: 'emerald' },
          { label: 'Crispy Chips', value: '65% food cost', pct: 65, tone: 'rose' },
          { label: 'Beef Tenderloin', value: '41% food cost', pct: 41, tone: 'indigo' }
        ] },
        { k: 'note', icon: 'ph-warning', tone: 'orange', title: 'Chips margin is thin', text: 'Potato cost rose 12% this week. Review the portion or the menu price before the weekend.' }
      ]
    },

    '/kitchen/reports': {
      dept: 'kitchen', chrome: 'sub', title: 'Kitchen Shift Report', back: '/kitchen',
      blocks: [
        { k: 'hero', kicker: 'Shift 06:00 – 14:00', title: 'Kitchen Shift Report', sub: 'Covers served, stock consumed and the closing count for the next shift.', icon: 'ph-chart-line-up', tone: 'emerald', chips: ['42 covers', '18 tickets'] },
        { k: 'stats', items: [
          { label: 'Covers', value: '42', note: 'Lunch service', tone: 'emerald' },
          { label: 'Tickets', value: '18', note: 'Avg 11 min', tone: 'indigo' },
          { label: 'Waste', value: '0.8%', note: 'Within target', tone: 'amber' }
        ] },
        { k: 'section', title: 'Stock closing summary' },
        { k: 'table', head: ['Ingredient', 'Open', 'Used', 'Close'], rows: [
          ['Maize flour', '46 kg', '4 kg', '42 kg'],
          ['Tilapia', '9 pcs', '5 pcs', '4 pcs'],
          ['Chicken', '10 birds', '3 birds', '7 birds'],
          ['Cooking oil', '19 L', '1 L', '18 L']
        ] },
        { k: 'note', icon: 'ph-printer', tone: 'slate', title: 'Handover', text: 'Send the closing count to the evening chef and the manager consolidation queue.', action: A.toast('Kitchen report sent to manager queue') }
      ]
    },

    /* ---------------- Manager ---------------- */
    '/manager': {
      dept: 'manager', chrome: 'brand', title: 'Executive Desk',
      blocks: [
        { k: 'hero', kicker: '01 Nov 2026 · Executive desk', title: 'Good morning, Mr. Banda', sub: 'The property at a glance. All departments reporting, one approval waiting.', icon: 'ph-buildings', tone: 'indigo', chips: ['All systems live', '1 approval'] },
        { k: 'stats', items: [
          { label: 'Gross revenue', value: '569,500', note: 'MWK today', tone: 'indigo' },
          { label: 'Occupancy', value: '78%', note: '9 of 12 rooms', tone: 'amber' },
          { label: 'ADR', value: '45,000', note: 'MWK average', tone: 'emerald' },
          { label: 'RevPAR', value: '35,100', note: 'MWK', tone: 'orange' }
        ] },
        { k: 'section', title: 'Revenue by settlement' },
        { k: 'bars', items: [
          { label: 'Room sales', value: 'MWK 315,000', pct: 55, tone: 'amber' },
          { label: 'Restaurant', value: 'MWK 168,000', pct: 30, tone: 'orange' },
          { label: 'Bar', value: 'MWK 86,500', pct: 15, tone: 'purple' }
        ] },
        { k: 'section', title: 'Department performance', action: { text: 'Consolidated', tone: 'link', action: A.nav('/manager/reports', 'Open consolidated report') } },
        { k: 'list', items: [
          { icon: 'ph-device-mobile', title: 'Front desk', sub: 'Collected MWK 164,000 · 3 departures', tone: 'amber', right: 'On target', rightTone: 'emerald', action: A.nav('/manager/reports', 'Front desk report') },
          { icon: 'ph-fork-knife', title: 'Restaurant', sub: '42 covers · MWK 168,000', tone: 'orange', right: '+8%', rightTone: 'emerald', action: A.nav('/manager/reports', 'Restaurant report') },
          { icon: 'ph-wine', title: 'Bar', sub: 'MWK 86,500 · 0.4% variance', tone: 'purple', right: 'On target', rightTone: 'emerald', action: A.nav('/manager/reports', 'Bar report') },
          { icon: 'ph-cooking-pot', title: 'Kitchen', sub: '18 tickets · 0.8% waste', tone: 'emerald', right: 'Efficient', rightTone: 'emerald', action: A.nav('/manager/reports', 'Kitchen report') }
        ] },
        { k: 'section', title: 'Management' },
        { k: 'quick', items: [
          { icon: 'ph-shield-check', label: 'Audits', tone: 'indigo', action: A.nav('/manager/audits', 'Shift audits') },
          { icon: 'ph-chart-pie-slice', label: 'Consolidated', tone: 'amber', action: A.nav('/manager/reports', 'Consolidated report') },
          { icon: 'ph-note-pencil', label: 'Shift note', tone: 'emerald', action: A.modal('shift-close', 'Add shift note') }
        ] },
        { k: 'note', icon: 'ph-warning-circle', tone: 'indigo', title: '1 approval waiting', text: 'The bar shift close needs a manager signature before the ledger can be finalised.', action: A.nav('/manager/audits', 'Review pending approvals') }
      ]
    },

    '/manager/audits': {
      dept: 'manager', chrome: 'sub', title: 'Shift Audits & Approvals', back: '/manager',
      blocks: [
        { k: 'hero', kicker: '3 pending · 2 approved', title: 'Shift Audits & Approvals', sub: 'Sign off cash-ups, stock variances and folio adjustments before they post.', icon: 'ph-shield-check', tone: 'indigo', chips: ['1 urgent', 'Cut-off 18:00'] },
        { k: 'stats', items: [
          { label: 'Pending', value: '3', note: 'Awaiting sign-off', tone: 'orange' },
          { label: 'Approved', value: '2', note: 'Today', tone: 'emerald' },
          { label: 'Rejected', value: '0', note: 'No exceptions', tone: 'slate' }
        ] },
        { k: 'section', title: 'Approval queue' },
        { k: 'list', items: [
          { icon: 'ph-cash-register', title: 'Bar shift close · MWK 86,500', sub: 'Variance 0.4% · submitted by Peter K.', tone: 'purple', right: 'Review', action: A.modal('audit-review', 'Review bar shift close') },
          { icon: 'ph-receipt', title: 'Folio adjustment · Room 102', sub: 'Waive MWK 3,500 minibar · by Grace M.', tone: 'amber', right: 'Review', action: A.modal('audit-review', 'Review folio adjustment') },
          { icon: 'ph-package', title: 'Stock variance · Kitchen', sub: 'Tilapia −1 pc · submitted by Chef A.', tone: 'emerald', right: 'Review', action: A.modal('audit-review', 'Review stock variance') }
        ] },
        { k: 'section', title: 'Recent decisions' },
        { k: 'timeline', items: [
          { time: 'Today · 07:05', title: 'Front desk cash-up approved', sub: 'Float MWK 20,000 · Grace M.', tone: 'emerald' },
          { time: 'Yesterday · 18:10', title: 'Restaurant report approved', sub: 'MWK 168,000 · 42 covers', tone: 'indigo' }
        ] },
        { k: 'note', icon: 'ph-lock', tone: 'slate', title: 'Segregation of duties', text: 'The manager who approves a cash-up cannot also be the officer who submitted it.' }
      ]
    },

    '/manager/reports': {
      dept: 'manager', chrome: 'sub', title: 'Consolidated Report', back: '/manager',
      blocks: [
        { k: 'hero', kicker: 'La Luna Hotel & Resort', title: 'Consolidated Report', sub: 'Room, restaurant, bar and kitchen performance for 01 Nov 2026.', icon: 'ph-chart-pie-slice', tone: 'indigo', chips: ['Gross 569,500', 'Net 402,300'] },
        { k: 'stats', items: [
          { label: 'Gross', value: '569,500', note: 'MWK revenue', tone: 'indigo' },
          { label: 'Net', value: '402,300', note: 'After costs', tone: 'emerald' },
          { label: 'Occupancy', value: '78%', note: '9 of 12 rooms', tone: 'amber' },
          { label: 'Guest score', value: '4.7', note: 'of 5 today', tone: 'orange' }
        ] },
        { k: 'section', title: 'Department summary' },
        { k: 'table', head: ['Department', 'Revenue', 'Cost', 'Net'], rows: [
          ['Rooms', '315,000', '42,000', '273,000'],
          ['Restaurant', '168,000', '61,000', '107,000'],
          ['Bar', '86,500', '31,200', '55,300'],
          ['Kitchen', '—', '33,000', '−33,000']
        ] },
        { k: 'section', title: 'Revenue mix' },
        { k: 'bars', items: [
          { label: 'Rooms', value: 'MWK 315,000', pct: 55, tone: 'amber' },
          { label: 'Restaurant', value: 'MWK 168,000', pct: 30, tone: 'orange' },
          { label: 'Bar', value: 'MWK 86,500', pct: 15, tone: 'purple' }
        ] },
        { k: 'note', icon: 'ph-download-simple', tone: 'slate', title: 'Board pack', text: 'Export the consolidated report as a signed PDF for the weekly board pack.', action: A.toast('Consolidated PDF queued for export') }
      ]
    }
  });
})();
