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
