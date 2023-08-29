import React from 'react';
import Head from 'next/head';
import localFont from 'next/font/local';
import Script from 'next/script';
import dynamic from 'next/dynamic';
import { AppProps } from 'next/app';
import { RecoilRoot } from 'recoil';
import { SessionProvider } from 'next-auth/react';
import { ChakraProvider } from '@chakra-ui/react';
import { theme } from '@ui/theme/theme';
import '../styles/overwrite.scss';
import '../styles/normalization.scss';
import '../styles/theme.css';
import '../styles/globals.scss';
import './../styles/remirror-editor.scss';
import 'react-date-picker/dist/DatePicker.css';
import 'react-calendar/dist/Calendar.css';
import 'react-toastify/dist/ReactToastify.css';
import Times from '@spaces/atoms/icons/Times';

const ToastContainer = dynamic(
  () => import('react-toastify').then((res) => res.ToastContainer),
  { ssr: true },
);

const MainPageWrapper = dynamic(
  () =>
    import('../components/ui-kit/layouts').then((res) => res.MainPageWrapper),
  { ssr: false },
);

const barlow = localFont({
  src: [
    {
      path: '../app/fonts/Barlow-Regular.woff',
      weight: '500',
      style: 'normal',
    },
    {
      path: '../app/fonts/Barlow-SemiBold.woff',
      weight: '600',
      style: 'normal',
    },
  ],
  preload: true,
  display: 'swap',
  variable: '--font-barlow',
});

export default function MyApp({
  Component,
  pageProps: { session, ...pageProps },
}: AppProps) {
  return (
    <>
      <Head>
        <meta charSet='utf-8' />
        <meta httpEquiv='X-UA-Compatible' content='IE=edge' />
        <meta name='viewport' content='width=device-width,initial-scale=1' />
        <meta name='description' content='Description' />
        <meta name='keywords' content='customerOS' />
        <title>Customer OS</title>
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
        <div className={barlow.className}>
          <ChakraProvider theme={theme}>
            <SessionProvider session={session}>
              <MainPageWrapper>
                <Component {...pageProps} />
              </MainPageWrapper>
            </SessionProvider>
          </ChakraProvider>
        </div>
      </RecoilRoot>

      <ToastContainer
        position='bottom-right'
        autoClose={3000}
        limit={3}
        closeOnClick={true}
        hideProgressBar={true}
        theme='colored'
        closeButton={({ closeToast }) => (
          <div onClick={closeToast}>
            <Times height={30} width={30} color='#17B26A' />
          </div>
        )}
      />
    </>
  );
}
