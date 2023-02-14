"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.createNonce = createNonce;
exports.useNonce = useNonce;

var jwt = _interopRequireWildcard(require("../../../jwt"));

var _openidClient = require("openid-client");

function _getRequireWildcardCache(nodeInterop) { if (typeof WeakMap !== "function") return null; var cacheBabelInterop = new WeakMap(); var cacheNodeInterop = new WeakMap(); return (_getRequireWildcardCache = function (nodeInterop) { return nodeInterop ? cacheNodeInterop : cacheBabelInterop; })(nodeInterop); }

function _interopRequireWildcard(obj, nodeInterop) { if (!nodeInterop && obj && obj.__esModule) { return obj; } if (obj === null || typeof obj !== "object" && typeof obj !== "function") { return { default: obj }; } var cache = _getRequireWildcardCache(nodeInterop); if (cache && cache.has(obj)) { return cache.get(obj); } var newObj = {}; var hasPropertyDescriptor = Object.defineProperty && Object.getOwnPropertyDescriptor; for (var key in obj) { if (key !== "default" && Object.prototype.hasOwnProperty.call(obj, key)) { var desc = hasPropertyDescriptor ? Object.getOwnPropertyDescriptor(obj, key) : null; if (desc && (desc.get || desc.set)) { Object.defineProperty(newObj, key, desc); } else { newObj[key] = obj[key]; } } } newObj.default = obj; if (cache) { cache.set(obj, newObj); } return newObj; }

const NONCE_MAX_AGE = 60 * 15;

async function createNonce(options) {
  var _provider$checks;

  const {
    cookies,
    logger,
    provider
  } = options;

  if (!((_provider$checks = provider.checks) !== null && _provider$checks !== void 0 && _provider$checks.includes("nonce"))) {
    return;
  }

  const nonce = _openidClient.generators.nonce();

  const expires = new Date();
  expires.setTime(expires.getTime() + NONCE_MAX_AGE * 1000);
  const encryptedNonce = await jwt.encode({ ...options.jwt,
    maxAge: NONCE_MAX_AGE,
    token: {
      nonce
    }
  });
  logger.debug("CREATE_ENCRYPTED_NONCE", {
    nonce,
    maxAge: NONCE_MAX_AGE
  });
  return {
    cookie: {
      name: cookies.nonce.name,
      value: encryptedNonce,
      options: { ...cookies.nonce.options,
        expires
      }
    },
    value: nonce
  };
}

async function useNonce(nonce, options) {
  var _provider$checks2, _value$nonce;

  const {
    cookies,
    provider
  } = options;

  if (!(provider !== null && provider !== void 0 && (_provider$checks2 = provider.checks) !== null && _provider$checks2 !== void 0 && _provider$checks2.includes("nonce")) || !nonce) {
    return;
  }

  const value = await jwt.decode({ ...options.jwt,
    token: nonce
  });
  return {
    value: (_value$nonce = value === null || value === void 0 ? void 0 : value.nonce) !== null && _value$nonce !== void 0 ? _value$nonce : undefined,
    cookie: {
      name: cookies.nonce.name,
      value: "",
      options: { ...cookies.nonce.options,
        maxAge: 0
      }
    }
  };
}