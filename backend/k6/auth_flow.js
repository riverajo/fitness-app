import http from 'k6/http';
import { check, sleep } from 'k6';


export const options = {
  vus: 1,
  duration: '10s',
};

const BASE_URL = __ENV.API_URL || 'http://localhost:8080';

export default function () {
  const cookieJar = http.cookieJar();

  // 1. Register a new user
  const email = `testuser_${Date.now()}@example.com`;
  const password = 'password123';

  const mutation = `
    mutation {
      register(input: {email: "${email}", password: "${password}"}) {
        token
        success
        message
        user {
          id
          email
        }
      }
    }
  `;

  const headers = { 'Content-Type': 'application/json' };
  let res = http.post(`${BASE_URL}/query`, JSON.stringify({ query: mutation }), { headers });

  check(res, {
    'register success': (r) => r.status === 200 && r.body.includes('register'),
    'has access token': (r) => r.body && JSON.parse(r.body).data?.register?.token !== '',
  });

  if (res.status !== 200) {
    console.error('Register failed:', res.status, res.body);
    return;
  }

  const accessToken = JSON.parse(res.body).data.register.token;

  // Verify cookies for refresh token are present in jar
  // k6 automatically handles cookies in the jar
  // We can't easily inspect HttpOnly cookies in k6 script logic without accessing the jar properties if exposed,
  // but subsequent requests will send them.

  sleep(1);

  // 2. Access protected resource with Access Token
  const query = `
    query {
      me {
        id
        email
      }
    }
  `;

  const authHeaders = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${accessToken}`
  };

  res = http.post(`${BASE_URL}/query`, JSON.stringify({ query: query }), { headers: authHeaders });

  check(res, {
    'me query success': (r) => r.status === 200 && r.body.includes(email),
  });

  // 3. Refresh Token
  // We call /auth/refresh. The cookie should be sent automatically by k6 since we are using the same VU.
  res = http.get(`${BASE_URL}/auth/refresh`);

  check(res, {
    'refresh success': (r) => r.status === 200,
    'new access token returned': (r) => JSON.parse(r.body).token !== '',
  });

  const newAccessToken = JSON.parse(res.body).token;
  check(res, {
    'tokens are different': () => newAccessToken !== accessToken
  })

  // 4. Logout
  const logoutMutation = `
    mutation {
      logout {
        success
      }
    }
  `;

  res = http.post(`${BASE_URL}/query`, JSON.stringify({ query: logoutMutation }), { headers: authHeaders });
  check(res, {
    'logout success': (r) => r.status === 200,
  });

  // 5. Try to Refresh again (Should fail)
  res = http.get(`${BASE_URL}/auth/refresh`);
  check(res, {
    'refresh fails after logout': (r) => r.status === 401,
  });
}
