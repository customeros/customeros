import { Metadata } from 'next';
import Script from 'next/script';
import { Barlow } from 'next/font/google';
import { GlobalCache } from '@graphql/types';

import { PageLayout } from './components/PageLayout';
import { getGraphQLClient } from './util/getGraphQLClient';
import { Providers } from './components/Providers/Providers';
import { ThemeProvider } from './components/Providers/ThemeProvider';
import { GlobalCacheDocument } from './graphql/global_Cache.generated';

import 'react-date-picker/dist/DatePicker.css';
import 'react-calendar/dist/Calendar.css';
import 'react-toastify/dist/ReactToastify.css';
import './../styles/globals.scss';


const barlow = Barlow({
  weight: ['300', '400', '500'],
  style: ['normal'],
  subsets: ['latin'],
  display: 'swap',
  preload: true,
  variable: '--font-barlow',
});

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const graphqlClient = getGraphQLClient();
  let globalCache: GlobalCache | null = null;

  try {
    const { global_Cache } = await graphqlClient.request<{
      global_Cache: GlobalCache;
    }>(GlobalCacheDocument);

    globalCache = global_Cache;
  } catch (e) {
    // handle error
  }

  return (
    <html lang='en' className={barlow.className} data-theme='light'>
      <Script
        async
        strategy='afterInteractive'
        id='openline-spaces-clarity-script'
        dangerouslySetInnerHTML={{
          __html: `(function(c,l,a,r,i,t,y){
                        c[a]=c[a]||function(){(c[a].q=c[a].q||[]).push(arguments)};
                        t=l.createElement(r);t.async=1;t.src="https://www.clarity.ms/tag/"+i;
                        y=l.getElementsByTagName(r)[0];y.parentNode.insertBefore(t,y);
                    })(window, document, "clarity", "script", "fryzkewrjw");`,
        }}
      />
      <body className='scrollbar'>
        <ThemeProvider>
          <PageLayout isOwner={globalCache?.isOwner ?? false}>
            <Providers>{children}</Providers>
          </PageLayout>
        </ThemeProvider>
      </body>
    </html>
  );
}

export const metadata: Metadata = {
  title: 'Spaces',
  description: 'Customer OS',
  keywords: ['CustomerOS', 'Spaces', 'Openline'],
  viewport: 'width=device-width,initial-scale=1',
};
