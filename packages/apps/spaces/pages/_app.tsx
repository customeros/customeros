import React from 'react';
import Head from 'next/head';
import Script from 'next/script';
import dynamic from 'next/dynamic';
import { AppProps } from 'next/app';
import { RecoilRoot } from 'recoil';
import 'remirror/styles/all.css';
import '../styles/overwrite.scss';
import '../styles/normalization.scss';
import '../styles/theme.css';
import '../styles/globals.css';
import 'react-date-picker/dist/DatePicker.css';
import 'react-calendar/dist/Calendar.css';
import 'react-toastify/dist/ReactToastify.css';

const ToastContainer = dynamic(
  () => import('react-toastify').then((res) => res.ToastContainer),
  { ssr: true },
);
const MainPageWrapper = dynamic(
  () =>
    import('../components/ui-kit/layouts').then((res) => res.MainPageWrapper),
  { ssr: true },
);
export default function MyApp({
  Component,
  pageProps: { session, ...pageProps },
}: AppProps) {
  // if (process.env.NODE_ENV === 'development') {
  //   require('../mocks');
  // }
  return (
    <>
      <Head>
        <meta charSet='utf-8' />
        <meta httpEquiv='X-UA-Compatible' content='IE=edge' />
        <meta name='viewport' content='width=device-width,initial-scale=1' />
        <meta name='description' content='Description' />
        <meta name='keywords' content='Keywords' />
        <title>Spaces</title>

        <link rel='manifest' href='/manifest.json' />
      </Head>

      <Script
        id='openline-spaces-clarity-script'
        strategy='afterInteractive'
        async
        dangerouslySetInnerHTML={{
          __html: `(function(c,l,a,r,i,t,y){
                        c[a]=c[a]||function(){(c[a].q=c[a].q||[]).push(arguments)};
                        t=l.createElement(r);t.async=1;t.src="https://www.clarity.ms/tag/"+i;
                        y=l.getElementsByTagName(r)[0];y.parentNode.insertBefore(t,y);
                    })(window, document, "clarity", "script", "fryzkewrjw");`,
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
