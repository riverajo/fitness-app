import http from 'k6/http';
import { check } from 'k6';
import { config } from './config.js';

export const options = {
  iterations: 1,
};

export default function () {
  const email = `k6-test-${Date.now()}@example.com`;
  const password = 'Password123!';

  const mutation = `
    mutation {
      register(input: {
        email: "${email}",
        password: "${password}"
      }) {
        success
        message
        user {
          id
          email
        }
      }
    }
  `;

  const payload = JSON.stringify({
    query: mutation,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post(`${config.BASE_URL}/query`, payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'registration successful': (r) => {
      const body = r.json();
      return body.data && body.data.register && body.data.register.success === true;
    },
    'no errors': (r) => r.json('errors') === undefined,
  });
}
