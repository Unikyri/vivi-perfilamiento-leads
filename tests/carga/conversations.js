import http from 'k6/http';
import { Rate, Trend } from 'k6/metrics';

const conversationFailures = new Rate('conversation_failures');
const conversationOverhead = new Trend('conversation_overhead', true);
const BASE_URL = validateBaseURL(__ENV.BASE_URL || 'http://127.0.0.1:8080');

export const options = {
  scenarios: {
    conversations: {
      executor: 'per-vu-iterations',
      vus: 20,
      iterations: 1,
      maxDuration: '30s',
    },
  },
  thresholds: {
    conversation_overhead: ['p(95)<500'],
    conversation_failures: ['rate<0.01'],
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

export default function () {
  const started = Date.now();
  const clientID = `conversation-${__VU}`;
  const response = http.post(
    `${BASE_URL}/api/conversations`,
    JSON.stringify({ client_id: clientID, message: 'solicitud de prueba local' }),
    { headers: { 'Content-Type': 'application/json', 'X-Forwarded-For': `2001:db8:1:${__VU}::1` } },
  );
  conversationOverhead.add(Date.now() - started);
  conversationFailures.add(response.status !== 200);
}
