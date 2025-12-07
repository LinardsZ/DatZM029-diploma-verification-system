/* eslint-disable no-param-reassign -- for status code reassignment */
const initActiveSession = {
  sid: '01J884W8WH62XWXA136WGG3Q4V',
  active: true,
  st: 'authorized',
  code: '12119912245',
  given_name: 'Ādams',
  family_name: 'Kalniņš',
  role: 'GAMER',
  scope: ['game/list:read'],
  secondsToLive: 597,
  secondsToCountdown: 300,
};

const activeSession = { ...initActiveSession };

let authorized = false;
/**
 * @typedef {Object} Middleware
 * @property {string} pattern
 * @property {string} method
 * @property {(req: import('node:http').IncomingMessage, res: import('node:http').ServerResponse) => void} handle
 */

/** @type {Middleware[]} */
export default [
  {
    pattern: '/idauth/api/1.0/session/keep-alive',
    method: 'GET',
    handle: (req, res) => {
      res.setHeader('Content-Type', 'application/json');
      if (req.headers.authorization) {
        const token = req.headers?.authorization?.split(' ')?.[1];
        authorized = token?.length > 1;
      }
      if (authorized) {
        res.end(JSON.stringify(activeSession));
      }
      res.statusCode = 401;
      res.end();
    },
  },
  {
    pattern: '/idauth/api/1.0/session',
    method: 'GET',
    handle: (req, res) => {
      res.setHeader('Content-Type', 'application/json');
      if (req.headers.authorization) {
        const token = req.headers?.authorization?.split(' ')?.[1];
        authorized = token?.length > 1;
      }
      if (authorized) {
        res.end(JSON.stringify(activeSession));
      }
      res.statusCode = 401;
      res.end();
    },
  },
  {
    pattern: '/idauth/api/1.0/session',
    method: 'DELETE',
    handle: (req, res) => {
      res.setHeader('Content-Type', 'application/json');
      authorized = false;
      res.end();
    },
  },
  {
    pattern: '/idauth/authorize',
    method: 'GET',
    handle: (req, res) => {
      res.setHeader('Content-Type', 'application/json');
      authorized = true;
      const url = new URL(req.url, 'http://localhost');
      const redirectUri = url.searchParams.get('redirect_uri');
      if (redirectUri) {
        res.statusCode = 302;
        res.setHeader('Location', `${redirectUri}?code=123456789`);
        res.end();
      }
      res.end();
    },
  },
];
