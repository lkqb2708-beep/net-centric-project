
/* DEBUG START */
console.log('[MangaHub] Script initializing...');
/* DEBUG END */

const API = {
  base: '',
  token: localStorage.getItem('mh_token') || '',

  headers() {
    const h = { 'Content-Type': 'application/json' };
    if (this.token) h['Authorization'] = `Bearer ${this.token}`;
    return h;
  },

  async get(path) {
    const r = await fetch(this.base + path, { headers: this.headers() });
    return r.json();
  },

  async post(path, body) {
    const r = await fetch(this.base + path, {
      method: 'POST',
      headers: this.headers(),
      body: JSON.stringify(body),
    });
    return r.json();
  },

  async put(path, body) {
    const r = await fetch(this.base + path, {
      method: 'PUT',
      headers: this.headers(),
      body: JSON.stringify(body),
    });
    return r.json();
  },

  async delete(path) {
    const r = await fetch(this.base + path, {
      method: 'DELETE',
      headers: this.headers(),
    });
    return r.json();
  },

  setToken(t) {
    this.token = t;
    if (t) localStorage.setItem('mh_token', t);
    else localStorage.removeItem('mh_token');
  },

  clearToken() {
    this.token = '';
    localStorage.removeItem('mh_token');
    localStorage.removeItem('mh_user');
  },

  setUser(u) { if (u) localStorage.setItem('mh_user', JSON.stringify(u)); },
  getUser()  { try { return JSON.parse(localStorage.getItem('mh_user')); } catch (e) { return null; } },
};
window.API = API; // Export early

/* ─── Toast notifications ──────────────────────────────────────────────── */
const Toast = {
  container: null,
  init() {
    this.container = document.getElementById('toast-container');
    if (!this.container) {
      this.container = document.createElement('div');
      this.container.className = 'toast-container';
      this.container.id = 'toast-container';
      document.body.appendChild(this.container);
    }
  },

  show(msg, type = 'info', duration = 3500) {
    const icons = { success: '✅', error: '❌', info: 'ℹ️', warning: '⚠️' };
    const t = document.createElement('div');
    t.className = `toast toast-${type}`;
    t.innerHTML = `<span class="toast-icon">${icons[type]||'📢'}</span>
                   <span class="toast-msg">${msg}</span>`;
    this.container.appendChild(t);
    setTimeout(() => {
      t.classList.add('removing');
      t.addEventListener('animationend', () => t.remove());
    }, duration);
  },

  success(m) { this.show(m, 'success'); },
  error(m)   { this.show(m, 'error',   4000); },
  info(m)    { this.show(m, 'info'); },
};

/* ─── Auth helpers ─────────────────────────────────────────────────────── */
const Auth = {
  isLoggedIn() { return !!API.token; },
  currentUser() { return API.getUser(); },

  async login(email, password) {
    const res = await API.post('/api/auth/login', { email, password });
    if (res.success) {
      API.setToken(res.data.token);
      API.setUser(res.data.user);
      return { ok: true, user: res.data.user };
    }
    return { ok: false, error: res.error };
  },

  async register(username, email, password) {
    const res = await API.post('/api/auth/register', { username, email, password });
    if (res.success) {
      API.setToken(res.data.token);
      API.setUser(res.data.user);
      return { ok: true, user: res.data.user };
    }
    return { ok: false, error: res.error };
  },

  async logout() {
    await API.post('/api/auth/logout', {});
    API.clearToken();
    window.location.href = '/login';
  },

  guard() {
    if (!this.isLoggedIn()) {
      window.location.href = '/login?next=' + encodeURIComponent(location.pathname);
      return false;
    }
    return true;
  },
};
window.Auth = Auth; // Export early
console.log('[MangaHub] Auth initialized, logged in:', Auth.isLoggedIn());

/* ─── WebSocket client ─────────────────────────────────────────────────── */
const WS = {
  conn: null,
  handlers: {},
  reconnectDelay: 1000,
  maxDelay: 30000,

  connect() {
    if (!Auth.isLoggedIn()) return;
    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
    const url = `${proto}//${location.host}/ws?token=${API.token}`;
    this.conn = new WebSocket(url);

    this.conn.onopen = () => {
      console.log('[WS] connected');
      this.reconnectDelay = 1000;
      this.emit('connected', {});
    };

    this.conn.onmessage = (e) => {
      try {
        const evt = JSON.parse(e.data);
        this.emit(evt.type, evt.payload);
      } catch (e) {}
    };

    this.conn.onclose = () => {
      console.log('[WS] disconnected, reconnecting in', this.reconnectDelay, 'ms');
      setTimeout(() => this.connect(), this.reconnectDelay);
      this.reconnectDelay = Math.min(this.reconnectDelay * 2, this.maxDelay);
    };

    this.conn.onerror = () => this.conn.close();
  },

  on(type, fn) {
    if (!this.handlers[type]) this.handlers[type] = [];
    this.handlers[type].push(fn);
  },

  emit(type, payload) {
    (this.handlers[type] || []).forEach(fn => fn(payload));
  },

  send(type, payload) {
    if (this.conn && this.conn.readyState === WebSocket.OPEN) {
      this.conn.send(JSON.stringify({ type, payload }));
    }
  },

  joinRoom(roomID) {
    this.send('join_room', { room_id: roomID });
  },
};

/* ─── Navigation & Sidebar ─────────────────────────────────────────────── */
const Nav = {
  init() {
    const path = location.pathname;
    document.querySelectorAll('.nav-item[data-href]').forEach(el => {
      if (path === el.dataset.href || (path.startsWith(el.dataset.href) && el.dataset.href !== '/')) {
        el.classList.add('active');
      }
      el.addEventListener('click', () => { location.href = el.dataset.href; });
    });

    // Mobile sidebar toggle
    const toggleBtn = document.getElementById('sidebar-toggle');
    const sidebar   = document.getElementById('sidebar');
    if (toggleBtn && sidebar) {
      toggleBtn.addEventListener('click', () => sidebar.classList.toggle('open'));
      document.addEventListener('click', (e) => {
        if (!sidebar.contains(e.target) && !toggleBtn.contains(e.target)) {
          sidebar.classList.remove('open');
        }
      });
    }

    // User info
    const user = Auth.currentUser();
    if (user) {
      document.querySelectorAll('[data-user-name]').forEach(el => el.textContent = user.username);
      document.querySelectorAll('[data-user-role]').forEach(el => el.textContent = user.role);
      document.querySelectorAll('[data-user-avatar]').forEach(el => {
        if (user.avatar_url) {
          el.innerHTML = `<img src="${user.avatar_url}" alt="${user.username}">`;
        } else {
          el.textContent = user.username[0].toUpperCase();
        }
      });
    }

    // Logout button
    document.getElementById('logout-btn')?.addEventListener('click', () => Auth.logout());

    // Load unread notification count
    if (Auth.isLoggedIn()) {
      API.get('/api/notifications?unread=true').then(res => {
        if (res.success && res.data.unread_count > 0) {
          document.querySelectorAll('.notif-badge').forEach(el => {
            el.textContent = res.data.unread_count;
            el.classList.remove('hidden');
          });
        }
      }).catch(() => {});
    }
  },
};

/* ─── Modals ───────────────────────────────────────────────────────────── */
const Modal = {
  open(id) {
    const m = document.getElementById(id);
    if (m) { m.classList.add('open'); m.style.display = 'flex'; }
  },
  close(id) {
    const m = document.getElementById(id);
    if (m) {
      m.classList.remove('open');
      setTimeout(() => m.style.display = '', 200);
    }
  },
  closeAll() {
    document.querySelectorAll('.modal-backdrop.open').forEach(m => {
      m.classList.remove('open');
      setTimeout(() => m.style.display = '', 200);
    });
  },
};

/* ─── Skeleton Loaders ─────────────────────────────────────────────────── */
function skeletonMangaCards(n = 8) {
  return Array.from({ length: n }, () => `
    <div class="manga-card">
      <div class="skeleton" style="aspect-ratio:3/4;width:100%"></div>
      <div style="padding:10px">
        <div class="skeleton" style="height:14px;margin-bottom:6px;border-radius:4px"></div>
        <div class="skeleton" style="height:12px;width:60%;border-radius:4px"></div>
      </div>
    </div>
  `).join('');
}

/* ─── Rating Stars ─────────────────────────────────────────────────────── */
function renderStars(rating, max = 10) {
  const filled = Math.round((rating / max) * 5);
  return Array.from({ length: 5 }, (_, i) =>
    `<span class="star ${i < filled ? 'filled' : ''}">★</span>`
  ).join('');
}

/* ─── Status badge ─────────────────────────────────────────────────────── */
function statusBadge(status) {
  const map = {
    reading:     ['badge-reading',   '📖 Reading'],
    completed:   ['badge-completed', '✅ Completed'],
    plan_to_read:['badge-plan',      '📌 Plan to Read'],
    on_hold:     ['badge-on_hold',   '⏸ On Hold'],
    dropped:     ['badge-dropped',   '❌ Dropped'],
    ongoing:     ['badge-ongoing',   '🔄 Ongoing'],
    hiatus:      ['badge-hiatus',    '⏸ Hiatus'],
    completed_m: ['badge-completed', '✅ Completed'],
    cancelled:   ['badge-dropped',   '🚫 Cancelled'],
  };
  const [cls, label] = map[status] || ['', status];
  return `<span class="badge ${cls}">${label}</span>`;
}

/* ─── Tabs ─────────────────────────────────────────────────────────────── */
function initTabs(containerId) {
  const container = document.getElementById(containerId);
  if (!container) return;
  const tabs   = container.querySelectorAll('.tab-btn');
  const panels = container.querySelectorAll('.tab-panel');
  tabs.forEach(tab => {
    tab.addEventListener('click', () => {
      tabs.forEach(t => t.classList.remove('active'));
      panels.forEach(p => p.classList.remove('active'));
      tab.classList.add('active');
      container.querySelector(`#${tab.dataset.tab}`)?.classList.add('active');
    });
  });
}

/* ─── Format helpers ───────────────────────────────────────────────────── */
function timeAgo(dateStr) {
  const diff = Date.now() - new Date(dateStr);
  const s = Math.floor(diff / 1000);
  if (s < 60) return `${s}s ago`;
  const m = Math.floor(s / 60);
  if (m < 60) return `${m}m ago`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h}h ago`;
  const d = Math.floor(h / 24);
  return `${d}d ago`;
}

function formatDate(dateStr) {
  return new Date(dateStr).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
}

/* ─── Global Manga Card Renderer ───────────────────────────────────────── */
function renderMangaCard(m) {
  const cover = m.cover_url || '/static/img/placeholder.svg';
  return `
  <div class="manga-card" onclick="location.href='/manga/${m.id}'" id="manga-${m.id}">
    <div style="position:relative;overflow:hidden;">
      <img class="manga-card-cover" src="${cover}" alt="${m.title}"
           onerror="this.onerror=null; this.src='/static/img/placeholder.svg';" loading="lazy">
      <span class="manga-card-badge">${m.status === 'completed' ? '✅' : '🔄'}</span>
      <div class="manga-card-overlay">
        <div class="manga-card-title">${m.title}</div>
        <div class="manga-card-meta">★ ${m.rating} · Ch.${m.chapter_count}</div>
      </div>
    </div>
  </div>`;
}

/* ─── Init on every page ───────────────────────────────────────────────── */
document.addEventListener('DOMContentLoaded', () => {
  Toast.init();
  Nav.init();

  // Page enter animation
  document.querySelector('.page-body')?.classList.add('page-enter-active');

  // Connect WebSocket if logged in
  if (Auth.isLoggedIn()) {
    WS.connect();
    // Handle real-time notifications
    WS.on('notification', (payload) => {
      Toast.info(`🔔 ${payload.title}`);
      const badge = document.querySelector('.notif-badge');
      if (badge) {
        const count = parseInt(badge.textContent || '0') + 1;
        badge.textContent = count;
        badge.classList.remove('hidden');
      }
    });
  }

  // Close modals on backdrop click
  document.querySelectorAll('.modal-backdrop').forEach(backdrop => {
    backdrop.addEventListener('click', (e) => {
      if (e.target === backdrop) Modal.closeAll();
    });
  });

  // Close modals with Escape
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') Modal.closeAll();
  });
});

/* Export globals EARLY to prevent reference errors */
window.API = API;
window.Auth = Auth;
window.Toast = Toast;
window.WS = WS;
window.Nav = Nav;
window.Modal = Modal;
window.renderMangaCard = renderMangaCard;
window.statusBadge = (s) => (s === 'ongoing' ? '🔄' : '✅');
window.renderStars = (r) => '★ ' + (r || 0);
window.timeAgo = timeAgo;
window.formatDate = formatDate;
window.skeletonMangaCards = skeletonMangaCards;
window.initTabs = initTabs;
