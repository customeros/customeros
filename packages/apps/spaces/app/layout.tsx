import { Metadata } from 'next';
import Script from 'next/script';
import localFont from 'next/font/local';
import { ToastContainer } from 'react-toastify';

import { PageLayout } from './src/components/PageLayout';
import { Providers } from './src/components/Providers/Providers';
import { ThemeProvider } from './src/components/Providers/ThemeProvider';

import 'react-toastify/dist/ReactToastify.css';
import './../styles/globals.scss';
import './../styles/date-picker.scss';
import './../styles/remirror-editor.scss';

const barlow = localFont({
  src: [
    { path: './src/fonts/Barlow-Regular.woff', weight: '400', style: 'normal' },
    { path: './src/fonts/Barlow-Medium.woff', weight: '500', style: 'normal' },
    {
      path: './src/fonts/Barlow-SemiBold.woff',
      weight: '600',
      style: 'normal',
    },
  ],
  preload: true,
  display: 'swap',
  variable: '--font-barlow',
});

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
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
          <PageLayout>
            <Providers>
              {children}
              <ToastContainer
                position='bottom-right'
                autoClose={8000}
                limit={3}
                closeOnClick={true}
                hideProgressBar={true}
                theme='colored'
              />
            </Providers>
          </PageLayout>
        </ThemeProvider>
      </body>
    </html>
  );
}

export const metadata: Metadata = {
  title: 'Customer OS',
  description: 'Customer OS',
  keywords: ['CustomerOS', 'Spaces', 'Openline'],
  viewport: 'width=device-width,initial-scale=1',
};
