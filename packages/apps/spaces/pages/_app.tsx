import React, { useEffect, useState } from 'react';
import Head from 'next/head';
import Script from 'next/script';
import dynamic from 'next/dynamic';
import { Barlow } from 'next/font/google';
import { AppProps } from 'next/app';
import { RecoilRoot } from 'recoil';
import 'remirror/styles/all.css';
import '../styles/overwrite.scss';
import '../styles/normalization.scss';
import '../styles/theme.css';
import '../styles/globals.scss';
import 'react-date-picker/dist/DatePicker.css';
import 'react-calendar/dist/Calendar.css';
import 'react-toastify/dist/ReactToastify.css';
import { useRouter } from 'next/router';
import { PageSkeleton } from '../components/shared/page-skeleton/PageSkeleton';
import { PageContentLayout } from '@spaces/layouts/page-content-layout';
import { SessionProvider } from 'next-auth/react';
import { ChakraProvider } from '@chakra-ui/react';
import { theme } from '@ui/theme/theme';

const ToastContainer = dynamic(
  () => import('react-toastify').then((res) => res.ToastContainer),
  { ssr: true },
);

const MainPageWrapper = dynamic(
  () =>
    import('../components/ui-kit/layouts').then((res) => res.MainPageWrapper),
  { ssr: false },
);
const barlow = Barlow({
  weight: ['300', '400', '500'],
  style: ['normal'],
  subsets: ['latin'],
  display: 'swap',
  preload: true,
  variable: '--font-main',
});

export default function MyApp({
  Component,
  pageProps: { session, ...pageProps },
}: AppProps) {
  const [loading, setLoading] = useState(false);
  const [loadingUrl, setLoadingUrl] = useState('');
  const router = useRouter();

  useEffect(() => {
    router.events.on('routeChangeStart', (url) => {
      setLoadingUrl(url);
      setLoading(true);
    });

    router.events.on('routeChangeComplete', (url) => {
      setLoading(false);
    });

    router.events.on('routeChangeError', (url) => {
      setLoading(false);
    });
  }, [router]);

  return (
    <>
      <Head>
        <meta charSet='utf-8' />
        <meta httpEquiv='X-UA-Compatible' content='IE=edge' />
        <meta name='viewport' content='width=device-width,initial-scale=1' />
        <meta name='description' content='Description' />
        <meta name='keywords' content='customerOS' />
        <title>Spaces</title>
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
        <div className={`global_container`}>
          <ChakraProvider theme={theme}>
            <SessionProvider session={session}>
              <MainPageWrapper>
                {loading ? (
                  <PageContentLayout>
                    <PageSkeleton loadingUrl={loadingUrl} />
                  </PageContentLayout>
                ) : (
                  <Component {...pageProps} />
                )}
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
      />
    </>
  );
}
