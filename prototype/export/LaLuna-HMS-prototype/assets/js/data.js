/* La Luna HMS — interactive prototype data
   Shared by the browser prototype and the Axure kit generator. */
(function () {
  'use strict';

  const D = {
    portal: {
      label: 'Portal', tag: 'Access', icon: 'ph-moon', accent: '#f5b942',
      subtitle: 'Integrated Hotel Management System', nav: []
    },
    frontdesk: {
      label: 'Front Desk', tag: 'Mobile FD', icon: 'ph-device-mobile', accent: '#f5b942',
      subtitle: 'Front Desk Reception',
      nav: [
        { label: 'Hub', icon: 'ph-house', route: '/frontdesk' },
        { label: 'Rooms', icon: 'ph-door-open', route: '/frontdesk/rooms' },
        { label: 'Walk-In', icon: 'ph-user-plus', route: '/frontdesk/walkin' },
        { label: 'Folio', icon: 'ph-receipt', route: '/frontdesk/folio' },
        { label: 'Checkout', icon: 'ph-credit-card', route: '/frontdesk/checkout' }
      ]
    },
    restaurant: {
      label: 'Restaurant', tag: 'Restaurant', icon: 'ph-fork-knife', accent: '#fb923c',
      subtitle: 'Dining & Dockets',
      nav: [
        { label: 'Hub', icon: 'ph-house', route: '/restaurant' },
        { label: 'Menu', icon: 'ph-book-open', route: '/restaurant/menu' },
        { label: 'Dockets', icon: 'ph-clipboard-text', route: '/restaurant/dockets' },
        { label: 'Reports', icon: 'ph-chart-line-up', route: '/restaurant/reports' }
      ]
    },
    bar: {
      label: 'Bar', tag: 'Bar', icon: 'ph-wine', accent: '#a78bfa',
      subtitle: 'Bar Inventory & Audit',
      nav: [
        { label: 'Hub', icon: 'ph-house', route: '/bar' },
        { label: 'Stock', icon: 'ph-package', route: '/bar/stock' },
        { label: 'Sales', icon: 'ph-cash-register', route: '/bar/sales' },
        { label: 'Reports', icon: 'ph-chart-bar', route: '/bar/reports' }
      ]
    },
    kitchen: {
      label: 'Kitchen', tag: 'Kitchen KDS', icon: 'ph-cooking-pot', accent: '#34d399',
      subtitle: 'Head Chef & Stock Station',
      nav: [
        { label: 'Orders', icon: 'ph-kanban', route: '/kitchen' },
        { label: 'Stock', icon: 'ph-carrot', route: '/kitchen/stock' },
        { label: 'Recipes', icon: 'ph-flask', route: '/kitchen/recipes' },
        { label: 'Reports', icon: 'ph-chart-line-up', route: '/kitchen/reports' }
      ]
    },
    manager: {
      label: 'Manager', tag: 'Executive', icon: 'ph-briefcase', accent: '#818cf8',
      subtitle: 'Hotel Management & Audits',
      nav: [
        { label: 'Overview', icon: 'ph-buildings', route: '/manager' },
        { label: 'Audits', icon: 'ph-shield-check', route: '/manager/audits' },
        { label: 'Consolidated', icon: 'ph-chart-pie-slice', route: '/manager/reports' }
      ]
    },
    guest: {
      label: 'Guest', tag: 'Resident', icon: 'ph-user', accent: '#fbbf24',
      subtitle: 'Your stay at La Luna',
      nav: [
        { label: 'Stay', icon: 'ph-bed', route: '/guest' },
        { label: 'Dining', icon: 'ph-bowl-food', route: '/guest/dining' },
        { label: 'My Bill', icon: 'ph-receipt', route: '/guest/folio' },
        { label: 'Services', icon: 'ph-bell-ringing', route: '/guest/services' }
      ]
    },
    book: {
      label: 'Booking', tag: 'Guest Web', icon: 'ph-moon-stars', accent: '#f5b942',
      subtitle: 'Book Your Stay', nav: []
    }
  };

  const A = {
    nav: (to, label) => ({ do: 'nav', to, label }),
    modal: (to, label) => ({ do: 'modal', to, label }),
    close: (label) => ({ do: 'close', label }),
    back: (label) => ({ do: 'back', label }),
    toast: (text, label) => ({ do: 'toast', text, label }),
    tray: (label) => ({ do: 'tray', label }),
    navc: (to, label) => ({ do: 'nav', to, close: true, label })
  };

  window.PROTO = {
    name: 'La Luna HMS',
    tagline: 'Integrated Hotel Management System',
    version: 'Prototype 1.0',
    depts: D,
    screens: {},
    modals: {},
    shortcuts: [
      ['M', 'Screen map'],
      ['H', 'Reveal hotspots'],
      ['Esc', 'Close overlay']
    ]
  };
  window.A = A;
})();
