// Enhanced mock API for the DnD Magical Items shop with authentication
const express = require('express');
const path = require('path');
const fs = require('fs');
const cors = require('cors');

const PORT = process.env.PORT || 8001;
const ITEMS_PATH = path.join(__dirname, 'data', 'items.json');
const USERS_PATH = path.join(__dirname, 'data', 'users.json');
const ORDERS_PATH = path.join(__dirname, 'data', 'orders.json');

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
function writeUsers(users) { return writeJSON(USERS_PATH, users); }
function writeOrders(orders) { return writeJSON(ORDERS_PATH, orders); }

// Auth middleware
function requireAuth(req, res, next) {
  const token = req.headers.authorization?.replace('Bearer ', '');
  if (!token || !sessions.has(token)) {
    return res.status(401).json({ error: 'Unauthorized' });
  }
  req.userId = sessions.get(token);
  next();
}

// Generate simple session token
function generateToken() {
  return `session-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
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
});
