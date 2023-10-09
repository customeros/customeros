import { Metadata } from 'next';
import Script from 'next/script';
import localFont from 'next/font/local';
import { ToastContainer } from 'react-toastify';

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
  preload: false,
  display: 'swap',
  variable: '--font-barlow',
});

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang='en' className={barlow.variable} data-theme='light'>
      <Script
        async
        strategy='afterInteractive'
        id='openline-customer-os-clarity-script'
        dangerouslySetInnerHTML={{
          __html: `(function(c,l,a,r,i,t,y){
                        c[a]=c[a]||function(){(c[a].q=c[a].q||[]).push(arguments)};
                        t=l.createElement(r);t.async=1;t.src="https://www.clarity.ms/tag/"+i;
                        y=l.getElementsByTagName(r)[0];y.parentNode.insertBefore(t,y);
                    })(window, document, "clarity", "script", "fryzkewrjw");`,
        }}
      />
      {`${process.env.NEXT_PUBLIC_PRODUCTION}` === 'true' && (
        <Script
          async
          strategy='afterInteractive'
          id='openline-customer-os-prod-heap-script'
          dangerouslySetInnerHTML={{
            __html: `
            window.heap=window.heap||[],heap.load=function(e,t){window.heap.appid=e,window.heap.config=t=t||{};var r=document.createElement("script");r.type="text/javascript",r.async=!0,r.src="https://cdn.heapanalytics.com/js/heap-"+e+".js";var a=document.getElementsByTagName("script")[0];a.parentNode.insertBefore(r,a);for(var n=function(e){return function(){heap.push([e].concat(Array.prototype.slice.call(arguments,0)))}},p=["addEventProperties","addUserProperties","clearEventProperties","identify","resetIdentity","removeEventProperty","setEventProperties","track","unsetEventProperty"],o=0;o<p.length;o++)heap[p[o]]=n(p[o])};
            heap.load("1078792267");
            `,
          }}
        />
      )}

      {`${process.env.NEXT_PUBLIC_PRODUCTION}` !== 'true' && (
        <Script
          async
          strategy='afterInteractive'
          id='openline-customer-os-dev-heap-script'
          dangerouslySetInnerHTML={{
            __html: `
            window.heap=window.heap||[],heap.load=function(e,t){window.heap.appid=e,window.heap.config=t=t||{};var r=document.createElement("script");r.type="text/javascript",r.async=!0,r.src="https://cdn.heapanalytics.com/js/heap-"+e+".js";var a=document.getElementsByTagName("script")[0];a.parentNode.insertBefore(r,a);for(var n=function(e){return function(){heap.push([e].concat(Array.prototype.slice.call(arguments,0)))}},p=["addEventProperties","addUserProperties","clearEventProperties","identify","resetIdentity","removeEventProperty","setEventProperties","track","unsetEventProperty"],o=0;o<p.length;o++)heap[p[o]]=n(p[o])};
  heap.load("3563674186");
            `,
          }}
        />
      )}

      <body className='scrollbar'>
        <ThemeProvider>
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
