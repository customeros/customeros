import express from 'express';
import { createProxyMiddleware } from 'http-proxy-middleware';

import 'dotenv/config';

async function createServer() {
  const app = express();

  const customerOsApiProxy = createProxyMiddleware({
    target: process.env.CUSTOMER_OS_API_PATH + '/query',
    changeOrigin: true,
    headers: {
      'X-Openline-API-KEY': process.env.CUSTOMER_OS_API_KEY,
      'X-Openline-USERNAME': 'customerostestuser@gmail.com',
    },
    logger: console,
    preserveHeaderKeyCase: true,
    followRedirects: true,
  });

  app.use('/customer-os-api', customerOsApiProxy);

  app.listen(5174);
}

createServer();
