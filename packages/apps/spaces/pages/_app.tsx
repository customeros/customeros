import Head from 'next/head';
import { ApolloProvider } from '@apollo/client';
import { AppProps, NextWebVitalsMetric } from 'next/app';
import '../styles/normalization.css';
import '../styles/theme.css';
import client from './../apollo-client';
import React from 'react';
import { ToastContainer } from 'react-toastify';
import { MainPageWrapper } from '../src/components/layouts';

// Uncomment when adding google analitics
// export function reportWebVitals({ id, name, label, value } :NextWebVitalsMetric) {
//   // Use `window.gtag` if you initialized Google Analytics as this example:
//   // https://github.com/vercel/next.js/blob/canary/examples/with-google-analytics/pages/_app.js
//   window.gtag('event', name, {
//     event_category:
//         label === 'web-vital' ? 'Web Vitals' : 'Next.js custom metric',
//     value: Math.round(name === 'CLS' ? value * 1000 : value), // values must be integers
//     event_label: id, // id unique to current page load
//     non_interaction: true, // avoids affecting bounce rate.
//   })
// }

export default function MyApp({
  Component,
  pageProps: { session, ...pageProps },
}: AppProps) {
  return (
    <>
      <Head>
        <meta charSet='utf-8' />
        <meta httpEquiv='X-UA-Compatible' content='IE=edge' />
        <meta
          name='viewport'
          content='width=device-width,initial-scale=1,minimum-scale=1,maximum-scale=1,user-scalable=no'
        />
        <meta name='description' content='Description' />
        <meta name='keywords' content='Keywords' />
        <title>Spaces</title>

        <link rel='manifest' href='/manifest.json' />
      </Head>

      <MainPageWrapper>
        <ApolloProvider client={client}>
          <ToastContainer
            position='bottom-right'
            autoClose={3000}
            closeOnClick={true}
            hideProgressBar={true}
            theme='colored'
          />
          <Component {...pageProps} />
        </ApolloProvider>
      </MainPageWrapper>
    </>
  );
}
