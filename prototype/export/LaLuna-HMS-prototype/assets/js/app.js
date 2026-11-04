/* La Luna HMS — prototype runtime */
(function () {
  'use strict';
  const P = window.PROTO;
  const app = document.getElementById('app');
  const viewport = document.getElementById('viewport');
  const appbar = document.getElementById('appbar');
  const tabbar = document.getElementById('tabbar');
  const layer = document.getElementById('modal-layer');
  const drawer = document.getElementById('drawer');
  const toastWrap = document.getElementById('toasts');
  const params = new URLSearchParams(location.search);
  const capture = params.has('capture');
  const state = { route: '/', modal: null, hotspots: false, tray: 2 };

  const esc = (v) => String(v == null ? '' : v)
    .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');

  const money = (n) => {
    const v = Math.abs(Number(n) || 0).toLocaleString('en-US');
    return (Number(n) < 0 ? '− ' : '') + 'MWK ' + v;
  };

  function actAttrs(action, extra) {
    if (!action) return '';
    extra = extra || {};
    let a = ' data-do="' + esc(action.do) + '"';
    if (action.to) a += ' data-to="' + esc(action.to) + '"';
    if (action.text) a += ' data-text="' + esc(action.text) + '"';
    if (action.close) a += ' data-close="1"';
    if (action.label) a += ' data-label="' + esc(action.label) + '"';
    if (extra.cls) a += ' class="' + esc(extra.cls) + '"';
    return a;
  }

  function tag(action, inner, cls, extra) {
    if (action) {
      const el = (extra && extra.as === 'a') ? 'a' : 'button';
      return '<' + el + ' type="button" class="' + cls + '"' + actAttrs(action) + '>' + inner + '</' + el + '>';
    }
    return '<div class="' + cls + '">' + inner + '</div>';
  }

  function icon(name, cls) {
    return '<i class="ph ' + esc(name) + (cls ? ' ' + cls : '') + '" aria-hidden="true"></i>';
  }

  /* ---------------- blocks ---------------- */

  const blocks = {
    login() {
      const roles = [
        ['frontdesk', 'Front Desk', '/frontdesk', '11 views', 'ph-device-mobile'],
        ['restaurant', 'Restaurant', '/restaurant', 'Dockets & menu', 'ph-fork-knife'],
        ['bar', 'Bar', '/bar', 'Stock & sales', 'ph-wine'],
        ['kitchen', 'Kitchen', '/kitchen', 'Orders & BOM', 'ph-cooking-pot'],
        ['manager', 'Manager', '/manager', 'Executive desk', 'ph-briefcase'],
        ['guest', 'Guest Booking', '/book', 'Booking & portal', 'ph-user']
      ];
      return `
        <div class="login-brand">
          <div class="moon-mark">${icon('ph-moon')}</div>
          <h1>LA LUNA <span>HMS</span></h1>
          <p>${esc(P.tagline)}</p>
        </div>
        <div class="login-card">
          <div class="login-head">
            <h2>Station sign in</h2>
            <p>Select a workstation to explore the prototype</p>
          </div>
          <p class="login-label">Direct station access</p>
          <div class="role-grid">
            ${roles.map(([id, label, to, sub, ic]) => `
              <button type="button" class="role-card tone-${id}" ${actAttrs(A.nav(to, 'Open ' + label))}>
                <span class="role-icon">${icon(ic)}</span>
                <span class="role-text"><strong>${esc(label)}</strong><small>${esc(sub)}</small></span>
                ${icon('ph-caret-right', 'role-arrow')}
              </button>`).join('')}
          </div>
          <div class="or"><span>or sign in</span></div>
          <label class="field"><span>Username</span><input type="text" value="grace.m" autocomplete="username"></label>
          <label class="field"><span>Password</span><input type="password" value="laluna2026" autocomplete="current-password"></label>
          <button type="button" class="btn primary wide" ${actAttrs(A.nav('/frontdesk', 'Sign in to front desk'))}>Sign in to Front Desk ${icon('ph-arrow-right')}</button>
          <p class="login-foot">Demo prototype · credentials are not checked</p>
        </div>
        <p class="login-copy">${esc(P.version)} · Built for the La Luna Hotel & Resort degree project</p>`;
    },

    hero(b) {
      return `
        <section class="hero tone-${esc(b.tone || 'amber')}">
          <div class="hero-watermark">${icon(b.icon || 'ph-moon')}</div>
          <div class="hero-kicker">${esc(b.kicker || '')}</div>
          <h2>${esc(b.title || '')}</h2>
          <p>${esc(b.sub || '')}</p>
          ${b.chips && b.chips.length ? `<div class="hero-chips">${b.chips.map(c => `<span>${esc(c)}</span>`).join('')}</div>` : ''}
        </section>`;
    },

    stats(b) {
      return `<div class="stats">${(b.items || []).map(s => `
        <div class="stat tone-${esc(s.tone || 'amber')}">
          <span class="stat-label">${esc(s.label)}</span>
          <span class="stat-value">${esc(s.value)}</span>
          <span class="stat-note">${esc(s.note || '')}</span>
        </div>`).join('')}</div>`;
    },

    section(b) {
      const a = b.action;
      return `<div class="section-head"><h3>${esc(b.title)}</h3>${a ? (a.action
        ? `<button type="button" class="section-link" ${actAttrs(a.action)}>${esc(a.text)} ${icon('ph-caret-right')}</button>`
        : `<span class="section-note">${esc(a.text)}</span>`) : ''}</div>`;
    },

    list(b) {
      return `<div class="list">${(b.items || []).map(item => {
        const inner = `
          <span class="row-icon tone-${esc(item.tone || 'slate')}">${icon(item.icon || 'ph-dot')}</span>
          <span class="row-main">
            <span class="row-title">${esc(item.title)}</span>
            <span class="row-sub">${esc(item.sub || '')}</span>
          </span>
          <span class="row-right ${item.rightTone ? 'tone-text-' + esc(item.rightTone) : ''}">${item.right ? `<span>${esc(item.right)}</span>` : ''}${item.action ? icon('ph-caret-right', 'row-chev') : ''}</span>`;
        return item.action ? `<button type="button" class="row" ${actAttrs(item.action)}>${inner}</button>` : `<div class="row">${inner}</div>`;
      }).join('')}</div>`;
    },

    grid(b) {
      return `<div class="grid2">${(b.items || []).map(item => {
        const inner = `<span class="tile-icon tone-${esc(item.tone || 'amber')}">${icon(item.icon || 'ph-squares-four')}</span>
          <span class="tile-title">${esc(item.title)}</span><span class="tile-sub">${esc(item.sub || '')}</span>`;
        return `<button type="button" class="tile" ${actAttrs(item.action)}>${inner}</button>`;
      }).join('')}</div>`;
    },

    quick(b) {
      return `<div class="quick">${(b.items || []).map(item => `
        <button type="button" class="quick-item" ${actAttrs(item.action)}>
          <span class="quick-icon tone-${esc(item.tone || 'amber')}">${icon(item.icon || 'ph-plus')}</span>
          <span>${esc(item.label)}</span>
        </button>`).join('')}</div>`;
    },

    chips(b) {
      return `<div class="chips-row">${(b.items || []).map(c => `
        <button type="button" class="chip${c.active ? ' is-active' : ''}" ${actAttrs(c.action)}>${esc(c.label)}</button>`).join('')}</div>`;
    },

    rooms(b) {
      return `<div class="room-grid">${(b.items || []).map(r => `
        <button type="button" class="room tone-${esc(r.tone || 'amber')}" ${actAttrs(r.action)}>
          <span class="room-top"><span class="room-no">${esc(r.no)}</span><span class="badge tone-${esc(r.tone || 'amber')}">${esc(r.status)}</span></span>
          <span class="room-type">${esc(r.type)} · ${esc(r.bed)}</span>
          <span class="room-rate">${money(r.rate)}<small>/ night</small></span>
        </button>`).join('')}</div>`;
    },

    menu(b) {
      return `<div class="menu">${(b.items || []).map(d => `
        <article class="dish tone-${esc(d.tone || 'orange')}">
          <div class="dish-main">
            <div class="dish-top"><h4>${esc(d.name)}</h4><span class="dish-price">${money(d.price)}</span></div>
            <p>${esc(d.desc)}</p>
            <span class="tag tone-${esc(d.tone || 'orange')}">${esc(d.tag || '')}</span>
          </div>
          <button type="button" class="dish-add" ${actAttrs(d.action)} aria-label="${esc(d.action && d.action.label || 'Add ' + d.name)}">${icon('ph-plus')}</button>
        </article>`).join('')}</div>`;
    },

    form(b) {
      const allWide = (b.fields || []).length <= 2;
      const fields = (b.fields || []).map(f => {
        const wide = allWide || f.type === 'textarea' || f.type === 'select' || f.type === 'email' || f.type === 'tel' || f.wide;
        let input;
        if (f.type === 'select') {
          input = `<select>${(f.options || [f.value]).map(o => `<option${o === f.value ? ' selected' : ''}>${esc(o)}</option>`).join('')}</select>`;
        } else if (f.type === 'textarea') {
          input = `<textarea rows="2">${esc(f.value || '')}</textarea>`;
        } else {
          input = `<input type="${esc(f.type || 'text')}" value="${esc(f.value || '')}"${f.readonly ? ' readonly' : ''}>`;
        }
        return `<label class="field${wide ? ' wide' : ''}"><span>${esc(f.label)}</span>${input}</label>`;
      }).join('');
      const wide = (b.fields || []).filter(f => f.type === 'textarea').length > 0 || (b.fields || []).length > 4;
      return `<div class="card form-card${b.compact ? ' compact' : ''}${wide ? ' form-wide' : ''}">
        ${b.title ? `<div class="form-title">${esc(b.title)}</div>` : ''}
        <div class="fields">${fields}</div>
        ${b.submit ? `<button type="button" class="btn primary wide" ${actAttrs(b.submit.action)}>${b.submit.icon ? icon(b.submit.icon) : ''}${esc(b.submit.label)}</button>` : ''}
        ${b.note ? `<p class="form-note">${esc(b.note)}</p>` : ''}
      </div>`;
    },

    table(b) {
      return `<div class="card table-card"><table>
        <thead><tr>${(b.head || []).map(h => `<th>${esc(h)}</th>`).join('')}</tr></thead>
        <tbody>${(b.rows || []).map(r => `<tr>${r.map((c, i) => `<td${i === 0 ? ' class="strong"' : ''}>${esc(c)}</td>`).join('')}</tr>`).join('')}</tbody>
      </table></div>`;
    },

    bars(b) {
      return `<div class="card bars">${(b.items || []).map(x => `
        <div class="bar-row tone-${esc(x.tone || 'amber')}">
          <div class="bar-top"><span>${esc(x.label)}</span><span class="bar-value">${esc(x.value)}</span></div>
          <div class="bar-track"><i style="width:${Math.max(0, Math.min(100, Number(x.pct) || 0))}%"></i></div>
        </div>`).join('')}</div>`;
    },

    statement(b) {
      const row = (r) => `<div class="st-row">
        <div class="st-main"><span class="st-label">${esc(r.label)}</span>${r.sub ? `<span class="st-sub">${esc(r.sub)}</span>` : ''}</div>
        <span class="st-amount${r.tone === 'emerald' ? ' positive' : ''}">${money(r.amount)}</span>
      </div>`;
      return `<div class="card statement">
        ${(b.rows || []).map(row).join('')}
        ${b.total ? `<div class="st-total"><span>${esc(b.total.label)}</span><span>${money(b.total.amount)}</span></div>` : ''}
        ${b.note ? `<p class="st-note">${esc(b.note)}</p>` : ''}
      </div>`;
    },

    tickets(b) {
      return `<div class="tickets">${(b.items || []).map(t => `
        <article class="ticket tone-${esc(t.tone || 'orange')}">
          <header><span class="ticket-id">${esc(t.id)}</span><span class="badge tone-${esc(t.tone || 'orange')}">${esc(t.status)}</span></header>
          <p class="ticket-meta">${esc(t.meta)}</p>
          <ul class="ticket-lines">${(t.lines || []).map(l => `
            <li><span class="ticket-qty">${esc(l.qty)}×</span><span class="ticket-line"><strong>${esc(l.name)}</strong>${l.note ? `<small>${esc(l.note)}</small>` : ''}</span></li>`).join('')}</ul>
          <footer><span class="ticket-total">${money(t.total)}</span>
            <button type="button" class="btn small" ${actAttrs(t.action)}>${esc(t.actionLabel || 'Bump')} ${icon('ph-check')}</button>
          </footer>
        </article>`).join('')}</div>`;
    },

    timeline(b) {
      return `<div class="card timeline">${(b.items || []).map(t => `
        <div class="tl-item tone-${esc(t.tone || 'slate')}">
          <span class="tl-dot"></span>
          <div class="tl-main"><span class="tl-time">${esc(t.time)}</span><span class="tl-title">${esc(t.title)}</span><span class="tl-sub">${esc(t.sub)}</span></div>
        </div>`).join('')}</div>`;
    },

    note(b) {
      const inner = `<span class="note-icon tone-${esc(b.tone || 'slate')}">${icon(b.icon || 'ph-info')}</span>
        <span class="note-main"><span class="note-title">${esc(b.title)}</span><span class="note-text">${esc(b.text)}</span></span>
        ${b.action ? icon('ph-caret-right', 'note-arrow') : ''}`;
      return b.action
        ? `<button type="button" class="note tone-${esc(b.tone || 'slate')} is-action" ${actAttrs(b.action)}>${inner}</button>`
        : `<div class="note tone-${esc(b.tone || 'slate')}">${inner}</div>`;
    },

    tray(b) {
      return `<button type="button" class="tray" ${actAttrs(b.action)}>
        <span class="tray-count">${esc(b.count)}</span>
        <span class="tray-label">Review dining tray<small>${esc(b.items || 'Charged to room 104')}</small></span>
        <span class="tray-total">${money(b.total)}</span>
      </button>`;
    },

    success(b) {
      return `<div class="success tone-${esc(b.tone || 'emerald')}${b.compact ? ' compact' : ''}">
        <div class="success-icon">${icon('ph-check-circle')}</div>
        <h3>${esc(b.title)}</h3>
        <p>${esc(b.text)}</p>
        <div class="success-rows">${(b.rows || []).map(r => `<div><span>${esc(r[0])}</span><strong>${esc(r[1])}</strong></div>`).join('')}</div>
      </div>`;
    },

    wifi(b) {
      return `<div class="wifi-card">
        <div class="wifi-head">${icon('ph-wifi-high')}<span>Connected automatically</span></div>
        <div class="wifi-row"><span>Network</span><strong>${esc(b.network)}</strong></div>
        <div class="wifi-row"><span>Password</span><strong class="mono">${esc(b.password)}</strong></div>
        <p>${esc(b.note)}</p>
      </div>`;
    }
  };

  function renderBlocks(list) {
    return (list || []).map(b => {
      const fn = blocks[b.k];
      return fn ? fn(b) : '';
    }).join('');
  }

  /* ---------------- chrome ---------------- */

  function renderAppbar(screen) {
    const dept = P.depts[screen.dept] || P.depts.portal;
    if (screen.chrome === 'none') { appbar.innerHTML = ''; appbar.hidden = true; return; }
    appbar.hidden = false;
    const back = screen.back
      ? `<button type="button" class="icon-btn" ${actAttrs(A.nav(screen.back, 'Back to ' + screen.back))} aria-label="Back">${icon('ph-arrow-left')}</button>`
      : `<span class="brand-mark">${icon('ph-moon')}</span>`;
    const title = screen.chrome === 'brand'
      ? `<span class="appbar-brand">LA LUNA</span><span class="appbar-tag tone-${screen.dept}">${esc(dept.tag)}</span>`
      : `<span class="appbar-title">${esc(screen.title)}</span>`;
    appbar.innerHTML = `${back}<div class="appbar-mid">${title}</div>
      <button type="button" class="icon-btn" ${actAttrs(A.toast('No new alerts · you are all caught up', 'Notifications'))} aria-label="Notifications">${icon('ph-bell')}<span class="dot"></span></button>`;
  }

  function renderTabbar(screen) {
    const dept = P.depts[screen.dept] || P.depts.portal;
    const nav = dept.nav || [];
    if (!nav.length) { tabbar.innerHTML = ''; tabbar.hidden = true; return; }
    tabbar.hidden = false;
    const direct = nav.find(n => n.route === state.route);
    const activeRoute = direct ? direct.route : nav[0].route;
    tabbar.innerHTML = nav.map(n => `
      <button type="button" class="tab${n.route === activeRoute ? ' is-active' : ''}" ${actAttrs(A.nav(n.route, n.label))}>
        <span class="tab-icon">${icon(n.icon)}</span><span class="tab-label">${esc(n.label)}</span>
      </button>`).join('');
  }

  function renderDrawer(active) {
    const groups = Object.entries(P.depts).filter(([id]) => id !== 'portal');
    drawer.innerHTML = `
      <div class="drawer-head">
        <div><span class="drawer-kicker">La Luna HMS</span><h2>Screen map</h2></div>
        <button type="button" class="icon-btn" data-drawer-close aria-label="Close map">${icon('ph-x')}</button>
      </div>
      <p class="drawer-sub">34 screens · 7 workstations · 19 popups</p>
      <div class="drawer-scroll">
        ${groups.map(([id, d]) => {
          const screens = Object.entries(P.screens).filter(([, s]) => s.dept === id);
          return `<div class="drawer-group">
            <div class="drawer-group-head"><span class="drawer-dot tone-${id}"></span>${esc(d.label)}<small>${screens.length}</small></div>
            ${screens.map(([route, s]) => `<button type="button" class="drawer-link${route === active ? ' is-active' : ''}" ${actAttrs(A.nav(route, s.title))}>
              <span>${esc(s.title)}</span><small>${esc(route)}</small></button>`).join('')}
          </div>`;
        }).join('')}
        <div class="drawer-group">
          <div class="drawer-group-head"><span class="drawer-dot tone-slate"></span>Popups<small>${Object.keys(P.modals).length}</small></div>
          ${Object.entries(P.modals).map(([id, m]) => `<button type="button" class="drawer-link" data-do="modal" data-to="${esc(id)}" data-label="Open ${esc(m.title)}">
            <span>${esc(m.title)}</span><small>${esc(id)}</small></button>`).join('')}
        </div>
      </div>`;
  }

  /* ---------------- routing ---------------- */

  function currentHashRoute() {
    const h = decodeURIComponent(location.hash.replace(/^#/, ''));
    return h || '/';
  }

  function navigate(to, opts) {
    opts = opts || {};
    closeModal(true);
    if (to === currentHashRoute() && !opts.force) {
      render();
      return;
    }
    location.hash = to;
  }

  function render() {
    const route = currentHashRoute();
    const screen = P.screens[route] || P.screens['/'];
    state.route = P.screens[route] ? route : '/';
    const dept = P.depts[screen.dept] || P.depts.portal;
    app.dataset.dept = screen.dept;
    app.style.setProperty('--accent', dept.accent);
    app.style.setProperty('--accent-soft', dept.accent + '22');
    renderAppbar(screen);
    renderTabbar(screen);
    viewport.innerHTML = `<div class="screen tone-${screen.dept}" data-screen="${esc(state.route)}">${renderBlocks(screen.blocks)}</div>`;
    viewport.scrollTop = 0;
    viewport.classList.remove('anim');
    void viewport.offsetWidth;
    viewport.classList.add('anim');
    document.title = screen.title + ' · La Luna HMS Prototype';
    document.querySelectorAll('.drawer-link').forEach(a => a.classList.remove('is-active'));
    renderDrawer(state.route);
    document.querySelector('#route-label').textContent = state.route;
    history.replaceState(history.state || { route: state.route }, '', location.href);
  }

  /* ---------------- modals ---------------- */

  function openModal(id) {
    const m = P.modals[id];
    if (!m) { toast('Popup not found: ' + id); return; }
    state.modal = id;
    app.classList.add('modal-open');
    layer.hidden = false;
    layer.innerHTML = `
      <div class="scrim" data-modal-close></div>
      <section class="modal size-${esc(m.size || 'sheet')} tone-${esc(m.tone || 'amber')}" role="dialog" aria-modal="true" aria-label="${esc(m.title)}">
        <div class="modal-grab"></div>
        <header class="modal-head">
          <span class="modal-icon tone-${esc(m.tone || 'amber')}">${icon(m.icon || 'ph-info')}</span>
          <div class="modal-title"><h2>${esc(m.title)}</h2><p>${esc(m.subtitle || '')}</p></div>
          <button type="button" class="icon-btn" data-modal-close aria-label="Close popup">${icon('ph-x')}</button>
        </header>
        <div class="modal-body">${renderBlocks(m.blocks)}</div>
        ${m.actions && m.actions.length ? `<footer class="modal-actions">${m.actions.map(a => `<button type="button" class="btn ${esc(a.tone || 'ghost')}" ${actAttrs(a.action)}>${esc(a.label)}</button>`).join('')}</footer>` : ''}
      </section>`;
    requestAnimationFrame(() => layer.classList.add('is-open'));
    const first = layer.querySelector('input, textarea, select, button');
    if (first && !capture) setTimeout(() => first.focus({ preventScroll: true }), 120);
    if (capture) layer.classList.add('is-open');
  }

  function closeModal(instant) {
    if (!state.modal) return;
    state.modal = null;
    app.classList.remove('modal-open');
    if (instant) {
      layer.classList.remove('is-open');
      layer.hidden = true;
      layer.innerHTML = '';
      return;
    }
    layer.classList.remove('is-open');
    setTimeout(() => { layer.hidden = true; layer.innerHTML = ''; }, 220);
  }

  /* ---------------- toast ---------------- */

  let toastTimer;
  function toast(text) {
    if (capture) return;
    toastWrap.innerHTML = `<div class="toast">${icon('ph-check-circle')}<span>${esc(text)}</span></div>`;
    toastWrap.classList.add('is-open');
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => toastWrap.classList.remove('is-open'), 2600);
  }

  /* ---------------- interactions ---------------- */

  function runAction(a) {
    if (!a) return;
    if (a.do === 'nav') { closeDrawer(); navigate(a.to); return; }
    if (a.do === 'back') {
      closeDrawer();
      const screen = P.screens[state.route];
      navigate((screen && screen.back) || '/');
      return;
    }
    if (a.do === 'modal') {
      closeDrawer();
      if (a.close) closeModal(true);
      openModal(a.to);
      return;
    }
    if (a.do === 'close') { closeModal(); return; }
    if (a.do === 'toast') { closeModal(); toast(a.text || 'Done'); return; }
    if (a.do === 'tray') { closeDrawer(); openModal('guest-tray'); return; }
  }

  function closeDrawer() {
    drawer.classList.remove('is-open');
    document.body.classList.remove('drawer-open');
  }

  document.addEventListener('click', (e) => {
    const closer = e.target.closest('[data-modal-close]');
    if (closer) { closeModal(); return; }
    const el = e.target.closest('[data-do]');
    if (el) {
      e.preventDefault();
      const action = {
        do: el.dataset.do, to: el.dataset.to, text: el.dataset.text,
        close: el.dataset.close === '1', label: el.dataset.label
      };
      runAction(action);
      return;
    }
    if (e.target.closest('[data-drawer-close]')) { closeDrawer(); return; }
    if (e.target.closest('#drawer-scrim')) { closeDrawer(); return; }
    if (e.target.closest('#btn-screens')) { toggleDrawer(); return; }
    if (e.target.closest('#btn-hotspots')) { toggleHotspots(); return; }
    if (e.target.closest('#btn-reset')) { resetDemo(); return; }
    if (e.target.closest('#btn-focus')) { document.body.classList.toggle('focus-mode'); return; }
  });

  function toggleDrawer() {
    const open = !drawer.classList.contains('is-open');
    drawer.classList.toggle('is-open', open);
    document.body.classList.toggle('drawer-open', open);
  }

  function toggleHotspots() {
    state.hotspots = !state.hotspots;
    document.body.classList.toggle('hotspots', state.hotspots);
    const b = document.getElementById('btn-hotspots');
    if (b) b.classList.toggle('is-on', state.hotspots);
  }

  function resetDemo() {
    state.tray = 2;
    closeModal(true);
    location.hash = '/';
    if (currentHashRoute() === '/') render();
    toast('Prototype reset');
  }

  window.addEventListener('hashchange', render);
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
      if (state.modal) { closeModal(); return; }
      if (drawer.classList.contains('is-open')) { toggleDrawer(); return; }
      if (document.body.classList.contains('focus-mode')) document.body.classList.remove('focus-mode');
    }
    if (e.target.matches('input, textarea, select')) return;
    if (e.key === 'm' || e.key === 'M') toggleDrawer();
    if (e.key === 'h' || e.key === 'H') toggleHotspots();
  });

  /* exposed for the Axure kit generator */
  window.__proto = {
    hotspots() {
      const base = app.getBoundingClientRect();
      return Array.from(app.querySelectorAll('[data-do]')).map(el => {
        const r = el.getBoundingClientRect();
        const label = el.dataset.label || (el.innerText || el.getAttribute('aria-label') || '').replace(/\s+/g, ' ').trim().slice(0, 70);
        return {
          label, action: el.dataset.do, target: el.dataset.to || '',
          x: Math.round(r.left - base.left), y: Math.round(r.top - base.top),
          w: Math.round(r.width), h: Math.round(r.height)
        };
      }).filter(h => h.w > 0 && h.h > 0);
    },
    ready: true
  };

  /* ---------------- boot ---------------- */

  if (capture) {
    document.body.classList.add('capture');
    const screen = params.get('screen') || '/';
    if (location.hash !== '#' + screen) location.hash = screen;
  } else {
    const screen = params.get('screen');
    if (screen && P.screens[screen] && location.hash !== '#' + screen) location.hash = screen;
  }
  render();
  layer.hidden = true;
  const modalParam = params.get('modal');
  if (modalParam && P.modals[modalParam]) {
    if (capture) document.body.classList.add('has-modal');
    openModal(modalParam);
    layer.classList.add('is-open');
  }
  if (params.has('map') && !capture) {
    drawer.classList.add('is-open');
    document.body.classList.add('drawer-open');
  }
  window.__proto.ready = true;
})();
