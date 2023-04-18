import React from 'react';
import Head from 'next/head';
import Script from 'next/script';
import { AppProps } from 'next/app';
import { RecoilRoot } from 'recoil';
import { ToastContainer } from 'react-toastify';
import 'primereact/resources/primereact.min.css';
import 'primeflex/primeflex.css';
import 'primeicons/primeicons.css';
import '../styles/normalization.scss';
import '../styles/theme.css';
import '../styles/globals.css';
import 'react-toastify/dist/ReactToastify.css';
import { MainPageWrapper } from '../components/ui-kit/layouts';

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

      <Script
        id='openline-spaces-clarity-script'
        strategy='afterInteractive'
        dangerouslySetInnerHTML={{
          __html: `(function(c,l,a,r,i,t,y){
                        c[a]=c[a]||function(){(c[a].q=c[a].q||[]).push(arguments)};
                        t=l.createElement(r);t.async=1;t.src="https://www.clarity.ms/tag/"+i;
                        y=l.getElementsByTagName(r)[0];y.parentNode.insertBefore(t,y);
                    })(window, document, "clarity", "script", "fryzkewrjw");`,
        }}
      />
      <Script
        id='openline-spaces-june-script'
        strategy='afterInteractive'
        dangerouslySetInnerHTML={{
          __html: `
          window.analytics = {};
  function juneify(writeKey) { 
      window.analytics._writeKey = writeKey;
      var script = document.createElement("script");
      script.type = "application/javascript";
      script.onload = function () {
          window.analytics.page();
      }
      script.src = "https://cdn.jsdelivr.net/npm/@june-so/analytics-next@2.0.0/dist/cjs/index.min.js";
      var first = document.getElementsByTagName('script')[0];
      first.parentNode.insertBefore(script, first);
  }
  juneify("M2QnaR2vqHiuu3W2");
          `,
        }}
      />

      <RecoilRoot>
        <MainPageWrapper>
          <Component {...pageProps} />
        </MainPageWrapper>
      </RecoilRoot>

      <ToastContainer
        position='bottom-right'
        autoClose={3000}
        limit={3}
        closeOnClick={true}
        hideProgressBar={true}
        theme='colored'
      />
    </>
  );
}
