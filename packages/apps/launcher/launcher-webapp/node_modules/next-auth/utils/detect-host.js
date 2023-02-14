"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.detectHost = detectHost;

function detectHost(forwardedHost) {
  var _process$env$VERCEL;

  if ((_process$env$VERCEL = process.env.VERCEL) !== null && _process$env$VERCEL !== void 0 ? _process$env$VERCEL : process.env.AUTH_TRUST_HOST) return forwardedHost;
  return process.env.NEXTAUTH_URL;
}