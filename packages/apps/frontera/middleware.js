import cors from 'cors';
import express from 'express';
import jwt from 'jsonwebtoken';
import { google } from 'googleapis';
import { createProxyMiddleware } from 'http-proxy-middleware';

import 'dotenv/config';

const PUBLIC_PATHS = ['/google-auth', '/callback/google-auth'];

const jwtMiddleware = (req, res, next) => {
  if (PUBLIC_PATHS.includes(req.path)) {
    return next();
  }

  const authorizationHeader = req.headers['authorization'];

  if (!authorizationHeader) {
    return res.status(400).json({
      message: 'missing authorization header',
    });
  }

  const sessionToken = authorizationHeader.split(' ')[1];

  if (!sessionToken) {
    return res.status(400).json({
      message: 'invalid token format',
    });
  }

  try {
    const session = jwt.verify(sessionToken, process.env.NEXTAUTH_SECRET);
    req.session = session;
    next();
  } catch (err) {
    return res.status(401).json({
      message: 'invalid authorization token',
    });
  }
};

const oauth2Client = new google.auth.OAuth2(
  process.env.GMAIL_CLIENT_ID,
  process.env.GMAIL_CLIENT_SECRET,
  `${process.env.VITE_MIDDLEWARE_API_URL}/callback/google-auth`,
);

async function customerOsSignIn(
  payload = { email: '', provider: '', oAuthToken: {} },
) {
  try {
    await fetch(`${process.env.USER_ADMIN_API_URL}/signin`, {
      method: 'POST',
      headers: {
        'X-Openline-API-KEY': process.env.USER_ADMIN_API_KEY,
      },
      body: JSON.stringify(payload),
    });
  } catch (err) {
    console.error(err);
  }
}

function fetchTenant(email) {
  return fetch(`${process.env.CUSTOMER_OS_API_PATH + '/query'}`, {
    method: 'POST',
    headers: {
      'X-Openline-API-KEY': process.env.CUSTOMER_OS_API_KEY,
      'X-Openline-USERNAME': email,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      operationName: 'tenant',
      query: `query tenant {
        tenant
      }`,
    }),
  });
}

function createIntegrationAppToken(tenant) {
  const WORKSPACE_KEY = process.env.INTEGRATION_APP_WORKSPACE_KEY;
  const PRIVATE_KEY_VALUE = process.env.INTEGRATION_APP_PRIVATE_KEY_VALUE;

  const tokenData = {
    id: tenant,
    name: tenant,
  };

  const token = jwt.sign(tokenData, PRIVATE_KEY_VALUE, {
    issuer: WORKSPACE_KEY,
    expiresIn: '30d',
    algorithm: 'ES256',
  });

  return token;
}

async function createServer() {
  const app = express();
  app.use(cors());
  app.use(jwtMiddleware);

  const customerOsApiProxy = createProxyMiddleware({
    pathFilter: '/customer-os-api',
    pathRewrite: { '^/customer-os-api': '' },
    target: process.env.CUSTOMER_OS_API_PATH + '/query',
    changeOrigin: true,
    headers: {
      'X-Openline-API-KEY': process.env.CUSTOMER_OS_API_KEY,
    },
    logger: console,
    preserveHeaderKeyCase: true,
    followRedirects: true,
  });

  const settingsApiProxy = createProxyMiddleware({
    pathFilter: '/sa',
    pathRewrite: { '^/sa': '' },
    target: process.env.SETTINGS_API_PATH,
    changeOrigin: true,
    headers: {
      'X-Openline-API-KEY': process.env.SETTINGS_API_KEY,
    },
    logger: console,
    preserveHeaderKeyCase: true,
    followRedirects: true,
  });

  const userAdminApiProxy = createProxyMiddleware({
    pathFilter: '/ua',
    pathRewrite: { '^/ua': '' },
    target: process.env.USER_ADMIN_API_URL,
    changeOrigin: true,
    headers: {
      'X-Openline-API-KEY': process.env.USER_ADMIN_API_KEY,
    },
    logger: console,
    preserveHeaderKeyCase: true,
    followRedirects: true,
  });

  app.use(customerOsApiProxy);
  app.use(settingsApiProxy);
  app.use(userAdminApiProxy);

  app.use('/google-auth', (_req, res) => {
    const scopes = ['openid', 'email', 'profile'];

    const url = oauth2Client.generateAuthUrl({
      access_type: 'offline',
      scope: scopes,
    });

    res.json({ url });
  });

  app.use('/callback/google-auth', async (req, res) => {
    const { code, state = '/organizations' } = req.query;

    try {
      const { tokens } = await oauth2Client.getToken(code);
      oauth2Client.setCredentials(tokens);

      const { access_token, refresh_token, expiry_date, scope } = tokens;

      const profileRes = await google
        .oauth2({
          auth: oauth2Client,
          version: 'v2',
        })
        .userinfo.get();

      await customerOsSignIn({
        accessToken: access_token,
        refreshToken: refresh_token,
        expiresAt: expiry_date ? new Date(expiry_date * 1000) : new Date(),
        scope,
        providerAccountId: profileRes.data.id,
        idToken: tokens.id_token,
      });

      const tenantReq = await fetchTenant(profileRes.data.email);
      const tenantRes = await tenantReq.json();
      const tenant = tenantRes?.data?.tenant ?? '';

      const integrations_token = createIntegrationAppToken(tenant);

      const sessionToken = jwt.sign(
        {
          tenant,
          access_token,
          refresh_token,
          integrations_token,
          profile: profileRes.data,
        },
        process.env.NEXTAUTH_SECRET,
        {
          expiresIn: '30d',
        },
      );

      res.redirect(
        `http://localhost:5173/auth/success?sessionToken=${sessionToken}&origin=${state}`,
      );
    } catch (err) {
      console.error(err);
      res.redirect(`http://localhost:5173/auth/failure?message=${err.message}`);
    }
  });

  app.use('/session', (req, res) => {
    res.json({ session: req?.session ?? null });
  });

  app.use('/enable/google-sync', (req, res) => {
    const scopes = [
      'openid',
      'email',
      'profile',
      'https://www.googleapis.com/auth/gmail.readonly',
      'https://www.googleapis.com/auth/gmail.send',
      'https://www.googleapis.com/auth/calendar.readonly',
    ];

    const url = oauth2Client.generateAuthUrl({
      access_type: 'offline',
      scope: scopes,
      state: req.query.origin,
    });

    res.json({ url });
  });

  app.listen(5174);
  console.info('Middleware server started on port 5174');
}

createServer();
