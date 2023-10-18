// This file sets a custom webpack configuration to use your Next.js app
// with Sentry.
// https://nextjs.org/docs/api-reference/next.config.js/introduction
// https://docs.sentry.io/platforms/javascript/guides/nextjs/manual-setup/
const { withSentryConfig } = require('@sentry/nextjs');
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});

/** @type {import('next').NextConfig} */
const webpack = require('webpack');

const config = {
  reactStrictMode: true,
  swcMinify: true,
  images: {
    minimumCacheTTL: 31536000,
  },
  env: {
    SSR_PUBLIC_PATH: process.env.SSR_PUBLIC_PATH,
    FILE_STORAGE_PUBLIC_URL: process.env.FILE_STORAGE_PUBLIC_URL,
    WEB_CHAT_API_KEY: process.env.WEB_CHAT_API_KEY,
    WEB_CHAT_HTTP_PATH: process.env.WEB_CHAT_HTTP_PATH,
    WEB_CHAT_WS_PATH: process.env.WEB_CHAT_WS_PATH,
    NEXT_PUBLIC_WEBSOCKET_PATH: process.env.NEXT_PUBLIC_WEBSOCKET_PATH,
    WEB_CHAT_TRACKER_ENABLED: process.env.WEB_CHAT_TRACKER_ENABLED,
    WEB_CHAT_TRACKER_APP_ID: process.env.WEB_CHAT_TRACKER_APP_ID,
    WEB_CHAT_TRACKER_ID: process.env.WEB_CHAT_TRACKER_ID,
    WEB_CHAT_TRACKER_COLLECTOR_URL: process.env.WEB_CHAT_TRACKER_COLLECTOR_URL,
    WEB_CHAT_TRACKER_BUFFER_SIZE: process.env.WEB_CHAT_TRACKER_BUFFER_SIZE,
    WEB_CHAT_TRACKER_MINIMUM_VISIT_LENGTH:
    process.env.WEB_CHAT_TRACKER_MINIMUM_VISIT_LENGTH,
    WEB_CHAT_TRACKER_HEARTBEAT_DELAY:
    process.env.WEB_CHAT_TRACKER_HEARTBEAT_DELAY,
    COMMS_MAIL_API_KEY: process.env.COMMS_MAIL_API_KEY,
    GOOGLE_MAPS_API_KEY: process.env.GOOGLE_MAPS_API_KEY,
    NEXT_PUBLIC_JUNE_ENABLED: process.env.NEXT_PUBLIC_JUNE_ENABLED,
    NEXT_PUBLIC_PRODUCTION: process.env.NEXT_PUBLIC_PRODUCTION,
  },
  i18n: {
    locales: ['en'],
    defaultLocale: 'en',
  },
  sentry: {
    hideSourceMaps: true,
  },
  output: 'standalone',
};

module.exports = withBundleAnalyzer(
  withSentryConfig(
    {
      ...config,
      swcMinify: true,
      webpack(config) {
        config.module.rules.push({
          test: /\.svg$/i,
          issuer: /\.[jt]sx?$/,
          use: ['@svgr/webpack'],
        });
        config.plugins.push(
          new webpack.DefinePlugin({
            __SENTRY_DEBUG__: true,
            __SENTRY_TRACING__: false,
          }),
        );
        // return the modified config
        return {
          ...config,
        };
      },
    },
    {},
  ),
);
