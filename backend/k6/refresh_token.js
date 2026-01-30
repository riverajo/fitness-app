import http from 'k6/http';
import { check, sleep } from 'k6';
import { config } from './config.js';

export const options = {
    scenarios: {
        refresh_flow: {
            executor: 'per-vu-iterations',
            exec: 'refreshTokenTest',
            vus: 1,
            iterations: 1,
            maxDuration: '30s',
        },
    },
};

export function refreshTokenTest() {
    const email = `k6-refresh-${Date.now()}@example.com`;
    const password = 'Password123!';

    // 1. Register
    const registerMutation = `
    mutation {
      register(input: {
        email: "${email}",
        password: "${password}"
      }) {
        success
        token
        refreshToken
        user {
          id
          email
        }
      }
    }
  `;

    const headers = {
        'Content-Type': 'application/json',
    };

    let res = http.post(`${config.BASE_URL}/query`, JSON.stringify({ query: registerMutation }), { headers });

    check(res, {
        'register status is 200': (r) => r.status === 200,
        'register success': (r) => r.json('data.register.success') === true,
        'register has access token': (r) => r.json('data.register.token') !== undefined,
        'register has refresh token': (r) => r.json('data.register.refreshToken') !== undefined,
    });

    const accessToken = res.json('data.register.token');
    const refreshToken = res.json('data.register.refreshToken');
    const userId = res.json('data.register.user.id');

    console.log(`Registered user ${userId}`);

    // 2. Call Protected Endpoint (Me) with Access Token
    const meQuery = `
    query {
      me {
        id
        email
      }
    }
  `;

    res = http.post(`${config.BASE_URL}/query`, JSON.stringify({ query: meQuery }), {
        headers: { ...headers, 'Authorization': `Bearer ${accessToken}` }
    });

    check(res, {
        'me status is 200': (r) => r.status === 200,
        'me has correct id': (r) => r.json('data.me.id') === userId,
    });

    // Short sleep to simulate time passing (though jti handles uniqueness, strict time check might be same second)
    sleep(1);

    // 3. Refresh Token
    const refreshMutation = `
    mutation {
      refreshToken(refreshToken: "${refreshToken}") {
        success
        token
        refreshToken
        user {
          id
        }
      }
    }
  `;

    res = http.post(`${config.BASE_URL}/query`, JSON.stringify({ query: refreshMutation }), { headers });

    check(res, {
        'refresh mutation status is 200': (r) => r.status === 200,
        'refresh success': (r) => r.json('data.refreshToken.success') === true,
        'refresh returns new access token': (r) => r.json('data.refreshToken.token') !== undefined,
        'refresh returns new refresh token': (r) => r.json('data.refreshToken.refreshToken') !== undefined,
        'refresh token is rotated': (r) => r.json('data.refreshToken.refreshToken') !== refreshToken,
    });

    const newAccessToken = res.json('data.refreshToken.token');
    const newRefreshToken = res.json('data.refreshToken.refreshToken');

    // 4. Verify Access with NEW Access Token
    res = http.post(`${config.BASE_URL}/query`, JSON.stringify({ query: meQuery }), {
        headers: { ...headers, 'Authorization': `Bearer ${newAccessToken}` }
    });

    check(res, {
        'me (new token) status is 200': (r) => r.status === 200,
        'me (new token) has correct id': (r) => r.json('data.me.id') === userId,
    });
}
