const test = require('node:test');
const assert = require('node:assert/strict');
const fs = require('fs');
const path = require('path');
const app = require('./server.js');

const SEED_PATH = path.join(
  __dirname,
  '..',
  '..',
  'docs',
  'specs',
  'reagent-map-spec',
  'acceptance-criteria.json'
);
const USERS_PATH = path.join(__dirname, 'data', 'users.json');
const EXPEDITIONS_PATH = path.join(__dirname, 'data', 'expeditions.json');

function startServer() {
  return new Promise((resolve) => {
    const server = app.listen(0, () => resolve(server));
  });
}

function stopServer(server) {
  return new Promise((resolve) => server.close(resolve));
}

function get(baseUrl, routePath, headers) {
  return fetch(`${baseUrl}${routePath}`, { headers }).then(async (res) => ({
    status: res.status,
    body: await res.json()
  }));
}

function post(baseUrl, routePath, body, headers) {
  return fetch(`${baseUrl}${routePath}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...headers },
    body: JSON.stringify(body)
  }).then(async (res) => ({
    status: res.status,
    body: await res.json()
  }));
}

/**
 * @issue reagent-sites-api
 * @ac AC1 - GET /api/reagent-sites returns all eight sites with positions unchanged,
 *     with or without a session.
 * @ac seed - apps/mock-api/data/reagent-sites.json holds exactly the eight sites in
 *     docs/specs/reagent-map-spec/acceptance-criteria.json's seed.reagentSites.
 */
test('GET /api/reagent-sites', async (t) => {
  const server = await startServer();
  const baseUrl = `http://localhost:${server.address().port}`;
  const seededSites = JSON.parse(fs.readFileSync(SEED_PATH, 'utf8')).seed.reagentSites;
  const usersSnapshot = fs.readFileSync(USERS_PATH, 'utf8');

  t.after(() => stopServer(server));
  t.after(() => fs.writeFileSync(USERS_PATH, usersSnapshot, 'utf8'));

  await t.test('returns all eight seeded sites with unchanged positions, no session', async () => {
    const { status, body } = await get(baseUrl, '/api/reagent-sites');

    assert.equal(status, 200);
    assert.equal(body.length, 8);
    assert.deepEqual(
      body.map(({ name, lat, lng }) => ({ name, lat, lng })),
      seededSites
    );
  });

  await t.test('returns the identical array with a valid session', async () => {
    const register = await fetch(`${baseUrl}/api/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: `qa-${Date.now()}@example.com`,
        password: 'password123',
        name: 'QA Tester'
      })
    });
    const { token } = await register.json();

    const withoutAuth = await get(baseUrl, '/api/reagent-sites');
    const withSession = await get(baseUrl, '/api/reagent-sites', {
      Authorization: `Bearer ${token}`
    });

    assert.equal(withSession.status, 200);
    assert.deepEqual(withSession.body, withoutAuth.body);
  });

  await t.test('returns the identical array with a garbage Authorization header', async () => {
    const withoutAuth = await get(baseUrl, '/api/reagent-sites');
    const withGarbageAuth = await get(baseUrl, '/api/reagent-sites', {
      Authorization: 'Bearer not-a-real-token'
    });

    assert.equal(withGarbageAuth.status, 200);
    assert.notEqual(withGarbageAuth.status, 401);
    assert.notEqual(withGarbageAuth.status, 403);
    assert.deepEqual(withGarbageAuth.body, withoutAuth.body);
  });

  await t.test('every site has an id, name, latitude, and longitude', () => {
    const sites = JSON.parse(
      fs.readFileSync(path.join(__dirname, 'data', 'reagent-sites.json'), 'utf8')
    );

    assert.equal(sites.length, 8);

    const ids = new Set();
    for (const site of sites) {
      assert.ok(site.id, `expected ${site.name} to have an id`);
      assert.equal(typeof site.name, 'string');
      assert.equal(typeof site.lat, 'number');
      assert.equal(typeof site.lng, 'number');
      ids.add(site.id);
    }
    assert.equal(ids.size, 8, 'expected every site id to be unique');
  });
});

// Mirrors the constants in server.js (not exported) so boundary tests can construct exact inputs.
const AURELIS_RADIUS_KM = 6000;
const MIN_EXPEDITION_DISTANCE_KM = 1;

/**
 * @issue expeditions-api
 * @ac AC3 - GET /api/expeditions with no session returns exactly the three seeded expeditions.
 * @ac AC4 - GET /api/expeditions with Elara's session returns the seeded expeditions and Elara's
 *     own, and never Borin's.
 * @ac AC4 (mirrored) - GET /api/expeditions with Borin's session returns the seeded expeditions
 *     and Borin's own, and never Elara's.
 * @ac resolveSessionUserId - GET /api/expeditions with an invalid/expired token behaves
 *     identically to no session at all (returns only the seeded expeditions).
 * @ac AC11 - POST /api/expeditions with a session stores the expedition for that adventurer
 *     (never an owner the body supplies), and it is returned by a later GET.
 * @ac AC12 - POST /api/expeditions with no session is refused with 401 and nothing is stored.
 * @ac AC13 - POST /api/expeditions with endpoints closer than 1 km apart is refused with 422
 *     naming a zero-length line, and nothing is stored.
 * @ac AC13 (boundary) - POST /api/expeditions with endpoints exactly 1 km apart is accepted,
 *     pinning the strict less-than semantics of the zero-length check.
 * @ac malformed body - POST /api/expeditions with a missing or non-numeric from/to is refused
 *     with 400 and nothing is stored.
 * @ac security-finding MEDIUM (isValidPoint) - POST /api/expeditions with a non-finite
 *     coordinate (NaN or Infinity) is refused with 400, not silently accepted as a zero-length
 *     line producing a corrupted {lat: null} record.
 * @ac security-finding MEDIUM (isValidPoint) - POST /api/expeditions with an out-of-range
 *     coordinate (latitude outside [-90, 90]) is refused with 400.
 */
test('expeditions endpoints', async (t) => {
  const server = await startServer();
  const baseUrl = `http://localhost:${server.address().port}`;
  const expeditionsSnapshot = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');
  const usersSnapshot = fs.readFileSync(USERS_PATH, 'utf8');

  t.after(() => stopServer(server));
  t.after(() => fs.writeFileSync(EXPEDITIONS_PATH, expeditionsSnapshot, 'utf8'));
  t.after(() => fs.writeFileSync(USERS_PATH, usersSnapshot, 'utf8'));

  async function loginAs(email) {
    const { body } = await post(baseUrl, '/api/auth/login', { email, password: 'password123' });
    return body.token;
  }

  async function userIdFor(token) {
    const { body } = await get(baseUrl, '/api/auth/me', { Authorization: `Bearer ${token}` });
    return body.id;
  }

  await t.test('GET with no session returns exactly the three seeded expeditions', async () => {
    const { status, body } = await get(baseUrl, '/api/expeditions');

    assert.equal(status, 200);
    assert.equal(body.length, 3);
    assert.ok(body.every((expedition) => expedition.owner === 'shop'));
  });

  await t.test("GET with Elara's session returns seeded plus her own, never Borin's", async () => {
    const elaraToken = await loginAs('elara@aurelis.example');
    const elaraId = await userIdFor(elaraToken);
    const borinId = await userIdFor(await loginAs('borin@aurelis.example'));

    const { status, body } = await get(baseUrl, '/api/expeditions', {
      Authorization: `Bearer ${elaraToken}`
    });

    assert.equal(status, 200);
    assert.equal(body.filter((expedition) => expedition.owner === 'shop').length, 3);
    assert.equal(body.filter((expedition) => expedition.owner === elaraId).length, 1);
    assert.equal(body.filter((expedition) => expedition.owner === borinId).length, 0);
  });

  await t.test("GET with Borin's session returns seeded plus his own, never Elara's", async () => {
    const borinToken = await loginAs('borin@aurelis.example');
    const borinId = await userIdFor(borinToken);
    const elaraId = await userIdFor(await loginAs('elara@aurelis.example'));

    const { status, body } = await get(baseUrl, '/api/expeditions', {
      Authorization: `Bearer ${borinToken}`
    });

    assert.equal(status, 200);
    assert.equal(body.filter((expedition) => expedition.owner === 'shop').length, 3);
    assert.equal(body.filter((expedition) => expedition.owner === borinId).length, 1);
    assert.equal(body.filter((expedition) => expedition.owner === elaraId).length, 0);
  });

  await t.test('GET with an invalid token returns exactly the three seeded expeditions', async () => {
    const { status, body } = await get(baseUrl, '/api/expeditions', {
      Authorization: 'Bearer not-a-real-token'
    });

    assert.equal(status, 200);
    assert.equal(body.length, 3);
    assert.ok(body.every((expedition) => expedition.owner === 'shop'));
  });

  await t.test('POST with a valid session stores the expedition and it is returned by GET', async () => {
    const token = await loginAs('elara@aurelis.example');

    const created = await post(
      baseUrl,
      '/api/expeditions',
      { from: { lat: 10, lng: 20 }, to: { lat: -10, lng: -20 } },
      { Authorization: `Bearer ${token}` }
    );
    const { body: afterPost } = await get(baseUrl, '/api/expeditions', {
      Authorization: `Bearer ${token}`
    });

    assert.equal(created.status, 200);
    assert.ok(afterPost.some((expedition) => expedition.id === created.body.id));
  });

  await t.test('POST ignores a forged owner field and always stores the session adventurer', async () => {
    const elaraToken = await loginAs('elara@aurelis.example');
    const elaraId = await userIdFor(elaraToken);
    const borinId = await userIdFor(await loginAs('borin@aurelis.example'));

    const created = await post(
      baseUrl,
      '/api/expeditions',
      { from: { lat: 10, lng: 20 }, to: { lat: -10, lng: -20 }, owner: borinId },
      { Authorization: `Bearer ${elaraToken}` }
    );

    assert.equal(created.body.owner, elaraId);
    assert.notEqual(created.body.owner, borinId);
  });

  await t.test('POST with no session is refused with 401 and nothing is stored', async () => {
    const before = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');

    const { status } = await post(baseUrl, '/api/expeditions', {
      from: { lat: 10, lng: 20 },
      to: { lat: -10, lng: -20 }
    });

    assert.equal(status, 401);
    assert.equal(fs.readFileSync(EXPEDITIONS_PATH, 'utf8'), before);
  });

  await t.test('POST with endpoints closer than 1 km apart is refused as a zero-length line', async () => {
    const token = await loginAs('elara@aurelis.example');
    const before = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');

    const { status, body } = await post(
      baseUrl,
      '/api/expeditions',
      { from: { lat: 10, lng: 20 }, to: { lat: 10.001, lng: 20 } },
      { Authorization: `Bearer ${token}` }
    );

    assert.equal(status, 422);
    assert.match(body.error, /zero-length/i);
    assert.equal(fs.readFileSync(EXPEDITIONS_PATH, 'utf8'), before);
  });

  await t.test('POST with endpoints just at the 1 km boundary is accepted (strict less-than)', async () => {
    const token = await loginAs('elara@aurelis.example');

    // A pure north-south displacement of MIN_EXPEDITION_DISTANCE_KM / AURELIS_RADIUS_KM radians
    // along a meridian is exactly MIN_EXPEDITION_DISTANCE_KM km of great-circle distance. Trig
    // floating-point error alone can land this a hair under the boundary (see the "closer than
    // 1 km" test above, whose 0.001-degree offset is nowhere near this close), so a tiny epsilon
    // is added to land reliably on the "at or just above" side the strict less-than check must accept.
    const dLatDegrees = (MIN_EXPEDITION_DISTANCE_KM / AURELIS_RADIUS_KM) * (180 / Math.PI) + 1e-9;
    const from = { lat: 10, lng: 20 };
    const to = { lat: 10 + dLatDegrees, lng: 20 };

    const distance = app.geodesicDistanceKm(from, to, AURELIS_RADIUS_KM);
    assert.ok(
      distance >= MIN_EXPEDITION_DISTANCE_KM && distance - MIN_EXPEDITION_DISTANCE_KM < 1e-5,
      'test setup must place endpoints at, not meaningfully past, the boundary'
    );

    const { status } = await post(baseUrl, '/api/expeditions', { from, to }, {
      Authorization: `Bearer ${token}`
    });

    assert.equal(status, 200);
  });

  await t.test('POST with a missing to field is refused with 400 and nothing is stored', async () => {
    const token = await loginAs('elara@aurelis.example');
    const before = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');

    const { status, body } = await post(
      baseUrl,
      '/api/expeditions',
      { from: { lat: 10, lng: 20 } },
      { Authorization: `Bearer ${token}` }
    );

    assert.equal(status, 400);
    assert.match(body.error, /from and to/i);
    assert.equal(fs.readFileSync(EXPEDITIONS_PATH, 'utf8'), before);
  });

  await t.test('POST with non-numeric coordinates is refused with 400 and nothing is stored', async () => {
    const token = await loginAs('elara@aurelis.example');
    const before = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');

    const { status, body } = await post(
      baseUrl,
      '/api/expeditions',
      { from: { lat: 10, lng: 20 }, to: { lat: '-10', lng: -20 } },
      { Authorization: `Bearer ${token}` }
    );

    assert.equal(status, 400);
    assert.match(body.error, /from and to/i);
    assert.equal(fs.readFileSync(EXPEDITIONS_PATH, 'utf8'), before);
  });

  async function postRawBody(rawJsonText, token) {
    const res = await fetch(`${baseUrl}/api/expeditions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: rawJsonText
    });
    return { status: res.status, body: await res.json() };
  }

  await t.test('POST with a coordinate that overflows to Infinity is refused with 400 and nothing is stored', async () => {
    const token = await loginAs('elara@aurelis.example');
    const before = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');

    // Literal "NaN"/"Infinity" tokens are not valid JSON, so body-parser rejects them before
    // our handler runs (a different, uninteresting failure). "1e400" IS valid JSON number
    // grammar but overflows double precision on parse, so the server's JSON.parse produces an
    // Infinity coordinate from a fully wire-valid request body — the exact non-finite value
    // isValidPoint's Number.isFinite check must catch. A JS numeric literal `1e400` in this
    // test's own source would overflow at parse time before ever reaching JSON.stringify, so
    // the payload is built as a raw JSON string instead, overflowing only on the server.
    const { status, body } = await postRawBody(
      '{"from":{"lat":10,"lng":20},"to":{"lat":1e400,"lng":21}}',
      token
    );

    assert.equal(status, 400);
    assert.match(body.error, /from and to/i);
    assert.equal(fs.readFileSync(EXPEDITIONS_PATH, 'utf8'), before);
  });

  await t.test('POST with a coordinate that overflows to -Infinity is refused with 400 and nothing is stored', async () => {
    const token = await loginAs('elara@aurelis.example');
    const before = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');

    const { status, body } = await postRawBody(
      '{"from":{"lat":10,"lng":20},"to":{"lat":-1e400,"lng":21}}',
      token
    );

    assert.equal(status, 400);
    assert.match(body.error, /from and to/i);
    assert.equal(fs.readFileSync(EXPEDITIONS_PATH, 'utf8'), before);
  });

  await t.test('POST with an out-of-range latitude is refused with 400 and nothing is stored', async () => {
    const token = await loginAs('elara@aurelis.example');
    const before = fs.readFileSync(EXPEDITIONS_PATH, 'utf8');

    const { status, body } = await post(
      baseUrl,
      '/api/expeditions',
      { from: { lat: 10, lng: 20 }, to: { lat: 999, lng: 20 } },
      { Authorization: `Bearer ${token}` }
    );

    assert.equal(status, 400);
    assert.match(body.error, /from and to/i);
    assert.equal(fs.readFileSync(EXPEDITIONS_PATH, 'utf8'), before);
  });
});

/**
 * @issue expeditions-api
 * @ac geodesicDistanceKm - the great-circle distance helper behind the 1 km zero-length check
 *     can fail independently of the HTTP layer.
 */
test('geodesicDistanceKm', () => {
  assert.equal(app.geodesicDistanceKm({ lat: 10, lng: 20 }, { lat: 10, lng: 20 }, 6000), 0);

  const quarterCircumference = (Math.PI * 6000) / 2;
  const distance = app.geodesicDistanceKm({ lat: 0, lng: 0 }, { lat: 90, lng: 0 }, 6000);
  assert.ok(Math.abs(distance - quarterCircumference) < 0.001);

  const closeDistance = app.geodesicDistanceKm({ lat: 10, lng: 20 }, { lat: 10.001, lng: 20 }, 6000);
  assert.ok(closeDistance < 1);
});
