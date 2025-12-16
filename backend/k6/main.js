import http from 'k6/http';
import { check, sleep } from 'k6';
import { config } from './config.js';

export const options = {
  discardResponseBodies: false,
  scenarios: {
    register: {
      executor: 'per-vu-iterations',
      exec: 'registerTest',
      vus: 1,
      iterations: 1,
      maxDuration: '30s',
    },
    login: {
      executor: 'per-vu-iterations',
      exec: 'loginTest',
      vus: 1,
      iterations: 1,
      startTime: '1s',
      maxDuration: '30s',
    },
  },
};

export function setup() {
  // Create a dedicated user for login tests
  const email = `k6-setup-${Date.now()}@example.com`;
  const password = 'Password123!';

  const mutation = `
    mutation {
      register(input: {
        email: "${email}",
        password: "${password}"
      }) {
        success
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
    'setup: user created': (r) => r.status === 200,
    'setup: registration successful': (r) => {
      const body = r.json();
      return body.data && body.data.register && body.data.register.success === true;
    }
  });

  return { email, password };
}

export function registerTest() {
  const email = `k6-reg-${__VU}-${__ITER}-${Date.now()}@example.com`;
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
    'register: status is 200': (r) => r.status === 200,
    'register: success': (r) => {
      const body = r.json();
      return body.data && body.data.register && body.data.register.success === true;
    },
    'register: no errors': (r) => r.json('errors') === undefined,
  });

  sleep(1);
}

export function loginTest(data) {
  const { email, password } = data;

  const mutation = `
    mutation {
      login(input: {
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

  const checkRes = check(res, {
    'login: status is 200': (r) => r.status === 200,
    'login: success': (r) => {
      const body = r.json();
      return body.data && body.data.login && body.data.login.success === true;
    },
    'login: no errors': (r) => r.json('errors') === undefined,
  });

  if (!checkRes) {
    console.log(`Login failed. Status: ${res.status}, Body: ${res.body}, Email: ${email}`);
  }

  sleep(1);
}
