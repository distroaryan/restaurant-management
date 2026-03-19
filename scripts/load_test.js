import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 20 }, // Simulate ramp-up to 20 users
    { duration: '1m', target: 20 },  // Stay at 20 users for 1 minute
    { duration: '10s', target: 0 },  // Ramp-down to 0 users
  ],
};

export default function () {
  const url = 'http://localhost:8000/api/v1/menus'; // Adjust port if needed
  const res = http.get(url);
  
  check(res, {
    'is status 200': (r) => r.status === 200,
  });
  
  sleep(1);
}
