/** @type {import('next').NextConfig} */

const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});

const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  images: {
    minimumCacheTTL: 31536000,
  },
  env: {
    SSR_PUBLIC_PATH: process.env.SSR_PUBLIC_PATH,
    FILE_STORAGE_PUBLIC_URL: process.env.FILE_STORAGE_PUBLIC_URL,
    COMMS_MAIL_API_KEY: process.env.COMMS_MAIL_API_KEY,
    GOOGLE_MAPS_API_KEY: process.env.GOOGLE_MAPS_API_KEY,
    NEXT_PUBLIC_JUNE_ENABLED: process.env.NEXT_PUBLIC_JUNE_ENABLED,
    NEXT_PUBLIC_PRODUCTION: process.env.NEXT_PUBLIC_PRODUCTION,
    NEXT_PUBLIC_NOTIFICATION_PROD_APP_IDENTIFIER:
      process.env.NEXT_PUBLIC_NOTIFICATION_PROD_APP_IDENTIFIER,
    NEXT_PUBLIC_NOTIFICATION_TEST_APP_IDENTIFIER:
      process.env.NEXT_PUBLIC_NOTIFICATION_TEST_APP_IDENTIFIER,
  },
  i18n: {
    locales: ['en'],
    defaultLocale: 'en',
  },
  output: 'standalone',
};

module.exports = withBundleAnalyzer(nextConfig);
