// Enhanced mock API for the DnD Magical Items shop with authentication
const express = require('express');
const path = require('path');
const fs = require('fs');
const cors = require('cors');

const PORT = process.env.PORT || 8001;
const ITEMS_PATH = path.join(__dirname, 'data', 'items.json');
const USERS_PATH = path.join(__dirname, 'data', 'users.json');
const ORDERS_PATH = path.join(__dirname, 'data', 'orders.json');
const REAGENT_SITES_PATH = path.join(__dirname, 'data', 'reagent-sites.json');
const EXPEDITIONS_PATH = path.join(__dirname, 'data', 'expeditions.json');

// Aurelis is a perfect sphere; distances on it use its own radius, never Earth's.
const AURELIS_RADIUS_KM = 6000;
const MIN_EXPEDITION_DISTANCE_KM = 1;

const app = express();
app.use(cors());
app.use(express.json());
app.use('/images', express.static(path.join(__dirname, 'public', 'images')));

// In-memory session storage (hybrid approach - easy to upgrade to JWT later)
const sessions = new Map();

// Helper functions for data persistence
function readJSON(filePath) {
  try {
    const raw = fs.readFileSync(filePath, 'utf8');
    return JSON.parse(raw);
  } catch (err) {
    console.error(`Failed to read ${filePath}`, err);
    return [];
  }
}

function writeJSON(filePath, data) {
  try {
    fs.writeFileSync(filePath, JSON.stringify(data, null, 2), 'utf8');
    return true;
  } catch (err) {
    console.error(`Failed to write ${filePath}`, err);
    return false;
  }
}

function readItems() { return readJSON(ITEMS_PATH); }
function readUsers() { return readJSON(USERS_PATH); }
function readOrders() { return readJSON(ORDERS_PATH); }
function readReagentSites() { return readJSON(REAGENT_SITES_PATH); }
function readExpeditions() { return readJSON(EXPEDITIONS_PATH); }
function writeUsers(users) { return writeJSON(USERS_PATH, users); }
function writeOrders(orders) { return writeJSON(ORDERS_PATH, orders); }
function writeExpeditions(expeditions) { return writeJSON(EXPEDITIONS_PATH, expeditions); }

// Auth middleware
function requireAuth(req, res, next) {
  const token = req.headers.authorization?.replace('Bearer ', '');
  if (!token || !sessions.has(token)) {
    return res.status(401).json({ error: 'Unauthorized' });
  }
  req.userId = sessions.get(token);
  next();
}

// Looks up the session for a request without rejecting it. Unlike requireAuth, a missing or
// invalid token is not an error here — it just means no personalization. Never call res/next.
function resolveSessionUserId(req) {
  const token = req.headers.authorization?.replace('Bearer ', '');
  if (!token || !sessions.has(token)) {
    return null;
  }
  return sessions.get(token);
}

// Generate simple session token
function generateToken() {
  return `session-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

// Great-circle (haversine) distance in km between two {lat, lng} points on a sphere of the
// given radius. Pure and side-effect-free so it can be unit-tested independently of the HTTP layer.
function geodesicDistanceKm(from, to, radiusKm) {
  const toRadians = (degrees) => (degrees * Math.PI) / 180;
  const dLat = toRadians(to.lat - from.lat);
  const dLng = toRadians(to.lng - from.lng);
  const lat1 = toRadians(from.lat);
  const lat2 = toRadians(to.lat);

  const a = Math.sin(dLat / 2) ** 2 + Math.cos(lat1) * Math.cos(lat2) * Math.sin(dLng / 2) ** 2;
  const centralAngle = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

  return radiusKm * centralAngle;
}

// Number.isFinite (not typeof) excludes NaN/Infinity — both pass `typeof x === 'number'` but
// make distanceKm NaN, and NaN < MIN_EXPEDITION_DISTANCE_KM is false, so the zero-length guard
// would silently accept and persist them (JSON.stringify then writes them to disk as null).
function isCoordinateInRange(value, limit) {
  return Number.isFinite(value) && Math.abs(value) <= limit;
}

function isValidPoint(point) {
  return !!point && isCoordinateInRange(point.lat, 90) && isCoordinateInRange(point.lng, 180);
}

// ============ HEALTH CHECK ============
app.get('/api/health', (req, res) => {
  res.json({ status: 'ok' });
});

// ============ ITEMS ENDPOINTS ============
app.get('/api/items', (req, res) => {
  const items = readItems();
  res.json(items);
});

app.get('/api/items/:id', (req, res) => {
  const id = req.params.id;
  const items = readItems();
  const item = items.find(i => i.id === id);
  if (!item) return res.status(404).json({ error: 'Not found' });
  res.json(item);
});

// ============ REAGENT SITES ENDPOINTS ============
app.get('/api/reagent-sites', (req, res) => {
  res.json(readReagentSites());
});

// ============ EXPEDITIONS ENDPOINTS ============
app.get('/api/expeditions', (req, res) => {
  const expeditions = readExpeditions();
  const userId = resolveSessionUserId(req);
  // userId guard: without it, a no-session request (userId === null) would match any
  // stored record whose owner is also null/undefined, as if it were the caller's own.
  const visible = expeditions.filter((e) => e.owner === 'shop' || (userId !== null && e.owner === userId));
  res.json(visible);
});

app.post('/api/expeditions', requireAuth, (req, res) => {
  const { from, to } = req.body;

  if (!isValidPoint(from) || !isValidPoint(to)) {
    return res.status(400).json({ error: 'from and to are required as { lat, lng }' });
  }

  const distanceKm = geodesicDistanceKm(from, to, AURELIS_RADIUS_KM);
  if (distanceKm < MIN_EXPEDITION_DISTANCE_KM) {
    return res.status(422).json({ error: 'Zero-length line: from and to are the same place' });
  }

  const expeditions = readExpeditions();
  const newExpedition = {
    id: `expedition-${Date.now()}`,
    owner: req.userId,
    from,
    to,
    savedAt: new Date().toISOString()
  };

  expeditions.push(newExpedition);
  writeExpeditions(expeditions);

  console.log({ expeditionId: newExpedition.id, owner: req.userId }, 'Expedition saved');
  res.json(newExpedition);
});

// ============ AUTH ENDPOINTS ============
app.post('/api/auth/register', (req, res) => {
  const { email, password, name } = req.body;
  
  if (!email || !password || !name) {
    return res.status(400).json({ error: 'Email, password, and name are required' });
  }

  const users = readUsers();
  
  // Check if user already exists
  if (users.find(u => u.email === email)) {
    return res.status(409).json({ error: 'User already exists' });
  }

  // Create new user
  const newUser = {
    id: `user-${Date.now()}`,
    email,
    password, // In production, this would be hashed
    name,
    createdAt: new Date().toISOString(),
    address: null
  };

  users.push(newUser);
  writeUsers(users);

  // Create session
  const token = generateToken();
  sessions.set(token, newUser.id);

  // Return user without password
  const { password: _, ...userWithoutPassword } = newUser;
  res.json({ user: userWithoutPassword, token });
});

app.post('/api/auth/login', (req, res) => {
  const { email, password } = req.body;

  if (!email || !password) {
    return res.status(400).json({ error: 'Email and password are required' });
  }

  const users = readUsers();
  const user = users.find(u => u.email === email && u.password === password);

  if (!user) {
    return res.status(401).json({ error: 'Invalid credentials' });
  }

  // Create session
  const token = generateToken();
  sessions.set(token, user.id);

  // Return user without password
  const { password: _, ...userWithoutPassword } = user;
  res.json({ user: userWithoutPassword, token });
});

app.post('/api/auth/logout', requireAuth, (req, res) => {
  const token = req.headers.authorization?.replace('Bearer ', '');
  if (token) {
    sessions.delete(token);
  }
  res.json({ success: true });
});

app.get('/api/auth/me', requireAuth, (req, res) => {
  const users = readUsers();
  const user = users.find(u => u.id === req.userId);
  
  if (!user) {
    return res.status(404).json({ error: 'User not found' });
  }

  const { password: _, ...userWithoutPassword } = user;
  res.json(userWithoutPassword);
});

// ============ USER PROFILE ENDPOINTS ============
app.get('/api/users/profile', requireAuth, (req, res) => {
  const users = readUsers();
  const user = users.find(u => u.id === req.userId);
  
  if (!user) {
    return res.status(404).json({ error: 'User not found' });
  }

  const { password: _, ...userWithoutPassword } = user;
  res.json(userWithoutPassword);
});

app.put('/api/users/profile', requireAuth, (req, res) => {
  const { name, address } = req.body;
  const users = readUsers();
  const userIndex = users.findIndex(u => u.id === req.userId);
  
  if (userIndex === -1) {
    return res.status(404).json({ error: 'User not found' });
  }

  // Update user
  if (name) users[userIndex].name = name;
  if (address) users[userIndex].address = address;

  writeUsers(users);

  const { password: _, ...userWithoutPassword } = users[userIndex];
  res.json(userWithoutPassword);
});

// ============ ORDER ENDPOINTS ============
app.get('/api/orders', requireAuth, (req, res) => {
  const orders = readOrders();
  const userOrders = orders.filter(o => o.userId === req.userId);
  res.json(userOrders);
});

app.get('/api/orders/:id', requireAuth, (req, res) => {
  const orders = readOrders();
  const order = orders.find(o => o.id === req.params.id && o.userId === req.userId);
  
  if (!order) {
    return res.status(404).json({ error: 'Order not found' });
  }

  res.json(order);
});

app.post('/api/orders', requireAuth, (req, res) => {
  const { items, total, shippingAddress } = req.body;

  if (!items || !Array.isArray(items) || items.length === 0) {
    return res.status(400).json({ error: 'Items are required' });
  }

  const orders = readOrders();
  const newOrder = {
    id: `order-${Date.now()}`,
    userId: req.userId,
    items,
    total: total || items.reduce((sum, item) => sum + (item.priceSnapshot * item.quantity), 0),
    status: 'pending',
    createdAt: new Date().toISOString(),
    shippingAddress: shippingAddress || null
  };

  orders.push(newOrder);
  writeOrders(orders);

  console.log('Created order:', newOrder.id, 'for user:', req.userId);
  res.json(newOrder);
});

// Note: mock-api intentionally does NOT serve Cucumber feature/index information.
// The real test-runner-service (Go) serves Cucumber index data on port 9001.

if (require.main === module) {
  app.listen(PORT, () => {
    console.log(`Mock API listening on http://localhost:${PORT}`);
    console.log('Available endpoints:');
    console.log('  - POST /api/auth/register');
    console.log('  - POST /api/auth/login');
    console.log('  - POST /api/auth/logout');
    console.log('  - GET  /api/auth/me');
    console.log('  - GET  /api/users/profile');
    console.log('  - PUT  /api/users/profile');
    console.log('  - GET  /api/orders');
    console.log('  - GET  /api/orders/:id');
    console.log('  - POST /api/orders');
    console.log('  - GET  /api/items');
    console.log('  - GET  /api/items/:id');
    console.log('  - GET  /api/reagent-sites');
    console.log('  - GET  /api/expeditions');
    console.log('  - POST /api/expeditions');
  });
}

// Exposed on the app object (not a separate export) so tests can exercise this pure function
// directly without changing the module's existing `module.exports = app` contract.
app.geodesicDistanceKm = geodesicDistanceKm;

module.exports = app;
