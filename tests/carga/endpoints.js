import http from 'k6/http';
import exec from 'k6/execution';
import { Rate } from 'k6/metrics';

const unexpectedFailures = new Rate('unexpected_failures');
const BASE_URL = validateBaseURL(__ENV.BASE_URL || 'http://127.0.0.1:8080');

export const options = {
  scenarios: {
    endpoints: {
      executor: 'constant-arrival-rate',
      rate: 100,
      timeUnit: '1s',
      duration: '60s',
      preAllocatedVUs: 20,
      maxVUs: 200,
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<300'],
    unexpected_failures: ['rate<0.01'],
  },
};

function validateBaseURL(raw) {
  const match = raw.match(/^http:\/\/(\[[^\]]+\]|[^/?#:@]+)(?::(\d+))?\/?$/i);
  if (!match) {
    throw new Error('BASE_URL must be credential-free http://localhost, 127.0.0.1, or ::1');
  }
  const hostname = match[1].replace(/^\[|\]$/g, '').toLowerCase();
  if (!['127.0.0.1', 'localhost', '::1'].includes(hostname)) {
    throw new Error('BASE_URL must target loopback only');
  }
  return raw.replace(/\/$/, '');
}

function identity() {
  const client = exec.scenario.iterationInTest % 300;
  return `2001:db8:${Math.floor(client / 256)}:${client % 256}::1`;
}

export default function () {
  const path = ((__VU + __ITER) % 2 === 0) ? '/salud' : '/api/leads';
  const response = http.get(`${BASE_URL}${path}`, {
    headers: { 'X-Forwarded-For': identity() },
    tags: { route: path },
  });
  const expected = response.status === 200;
  unexpectedFailures.add(!expected);
}
