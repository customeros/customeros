"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = getAuthorizationUrl;

var _client = require("./client");

var _clientLegacy = require("./client-legacy");

var _stateHandler = require("./state-handler");

var _nonceHandler = require("./nonce-handler");

var _pkceHandler = require("./pkce-handler");

async function getAuthorizationUrl({
  options,
  query
}) {
  var _provider$version;

  const {
    logger,
    provider
  } = options;
  let params = {};

  if (typeof provider.authorization === "string") {
    const parsedUrl = new URL(provider.authorization);
    const parsedParams = Object.fromEntries(parsedUrl.searchParams);
    params = { ...params,
      ...parsedParams
    };
  } else {
    var _provider$authorizati;

    params = { ...params,
      ...((_provider$authorizati = provider.authorization) === null || _provider$authorizati === void 0 ? void 0 : _provider$authorizati.params)
    };
  }

  params = { ...params,
    ...query
  };

  if ((_provider$version = provider.version) !== null && _provider$version !== void 0 && _provider$version.startsWith("1.")) {
    var _provider$authorizati2;

    const client = (0, _clientLegacy.oAuth1Client)(options);
    const tokens = await client.getOAuthRequestToken(params);
    const url = `${(_provider$authorizati2 = provider.authorization) === null || _provider$authorizati2 === void 0 ? void 0 : _provider$authorizati2.url}?${new URLSearchParams({
      oauth_token: tokens.oauth_token,
      oauth_token_secret: tokens.oauth_token_secret,
      ...tokens.params
    })}`;
    logger.debug("GET_AUTHORIZATION_URL", {
      url,
      provider
    });
    return {
      redirect: url
    };
  }

  const client = await (0, _client.openidClient)(options);
  const authorizationParams = params;
  const cookies = [];
  const state = await (0, _stateHandler.createState)(options);

  if (state) {
    authorizationParams.state = state.value;
    cookies.push(state.cookie);
  }

  const nonce = await (0, _nonceHandler.createNonce)(options);

  if (nonce) {
    authorizationParams.nonce = nonce.value;
    cookies.push(nonce.cookie);
  }

  const pkce = await (0, _pkceHandler.createPKCE)(options);

  if (pkce) {
    authorizationParams.code_challenge = pkce.code_challenge;
    authorizationParams.code_challenge_method = pkce.code_challenge_method;
    cookies.push(pkce.cookie);
  }

  const url = client.authorizationUrl(authorizationParams);
  logger.debug("GET_AUTHORIZATION_URL", {
    url,
    cookies,
    provider
  });
  return {
    redirect: url,
    cookies
  };
}