const API_BASE = window.location.hostname === 'localhost' && window.location.port === '8081'
  ? 'http://localhost:8080'
  : '';

const tokenKey = 'chengetai_token';

function parseToken(token) {
  try {
    const [, payload] = token.split('.');
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
  } catch {
    return null;
  }
}

function getToken() {
  return localStorage.getItem(tokenKey);
}

function saveToken(token) {
  localStorage.setItem(tokenKey, token);
}

function clearToken() {
  localStorage.removeItem(tokenKey);
}

function isTokenValid(token) {
  const payload = parseToken(token);
  if (!payload || !payload.exp) return false;
  return payload.exp * 1000 > Date.now();
}

async function apiRequest(path, options = {}) {
  const token = getToken();
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) };
  if (token) {
    headers.Authorization = 'Bearer ' + token;
  }
  const response = await fetch(`${API_BASE}${path}`, { ...options, headers });
  const body = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(body?.error?.message || 'Request failed');
  }
  return body;
}

function protectPage() {
  const token = getToken();
  if (!token || !isTokenValid(token)) {
    clearToken();
    window.location.href = 'login.html';
    return null;
  }
  return parseToken(token);
}

function setMessage(message, isError = false) {
  const node = document.getElementById('form-message');
  if (!node) return;
  node.textContent = message;
  node.style.color = isError ? '#c62828' : '#5b6c87';
}

function bindLogin() {
  const form = document.getElementById('login-form');
  if (!form) return;
  form.addEventListener('submit', async (event) => {
    event.preventDefault();
    const formData = new FormData(form);
    try {
      const payload = Object.fromEntries(formData.entries());
      const result = await apiRequest('/api/v1/auth/login', {
        method: 'POST',
        body: JSON.stringify(payload)
      });
      saveToken(result.token);
      window.location.href = 'dashboard.html';
    } catch (error) {
      setMessage(error.message, true);
    }
  });
}

function toggleLearnerFields() {
  const role = document.querySelector('input[name="role"]:checked')?.value;
  const learnerFields = document.getElementById('learner-fields');
  if (learnerFields) {
    learnerFields.style.display = role === 'learner' ? 'grid' : 'none';
  }
}

function bindRegister() {
  const form = document.getElementById('register-form');
  if (!form) return;
  form.addEventListener('change', toggleLearnerFields);
  toggleLearnerFields();
  form.addEventListener('submit', async (event) => {
    event.preventDefault();
    const formData = new FormData(form);
    try {
      const payload = Object.fromEntries(formData.entries());
      const result = await apiRequest('/api/v1/auth/register', {
        method: 'POST',
        body: JSON.stringify(payload)
      });
      saveToken(result.token);
      window.location.href = 'dashboard.html';
    } catch (error) {
      setMessage(error.message, true);
    }
  });
}

async function bindDashboard() {
  if (!protectPage()) return;

  const title = document.getElementById('dashboard-title');
  const copy = document.getElementById('dashboard-copy');
  const panels = document.getElementById('dashboard-panels');
  const logoutButton = document.getElementById('logout-button');

  const roleContent = {
    learner: {
      title: 'Learner dashboard',
      copy: 'Review lessons, practise activities, and progress snapshots.',
      cards: ['Daily lesson', 'Practice queue', 'Progress streak']
    },
    teacher: {
      title: 'Teacher dashboard',
      copy: 'Track learners, review classroom activity, and plan support.',
      cards: ['Class activity', 'Assignments', 'Support insights']
    },
    parent: {
      title: 'Parent dashboard',
      copy: 'See recent learner activity and encourage study routines.',
      cards: ['Recent progress', 'Attendance prompts', 'Support tips']
    }
  };

  try {
    const me = await apiRequest('/api/v1/auth/me');
    const content = roleContent[me.role] || roleContent.parent;
    title.textContent = content.title;
    copy.textContent = `${content.copy} Signed in as ${me.first_name} ${me.last_name}.`;
    panels.innerHTML = content.cards.map((card) => `<section class="card"><h2>${card}</h2></section>`).join('');
  } catch {
    clearToken();
    window.location.href = 'login.html';
    return;
  }

  logoutButton?.addEventListener('click', async () => {
    try {
      await apiRequest('/api/v1/auth/logout', { method: 'POST' });
    } finally {
      clearToken();
      window.location.href = 'login.html';
    }
  });
}

document.addEventListener('DOMContentLoaded', () => {
  const page = document.body.dataset.page;
  if (page === 'login') bindLogin();
  if (page === 'register') bindRegister();
  if (page === 'dashboard') bindDashboard();
});
