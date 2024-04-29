import React from 'react';
import { RouterProvider } from 'react-router-dom';

import ReactDOM from 'react-dom/client';
import { GoogleOAuthProvider } from '@react-oauth/google';

import { Providers } from '@shared/components/Providers/Providers';

import { router } from './routes/router';

import './styles/globals.scss';
import './styles/filepond.scss';
import './styles/date-picker.scss';
import './styles/normalization.scss';
import './styles/remirror-editor.scss';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <GoogleOAuthProvider clientId='17459587911-41bnga38ph1o6ksnbe798k2cu65v1gdd.apps.googleusercontent.com'>
    <Providers
      env={{
        PRODUCTION: 'false',
        NOTIFICATION_PROD_APP_IDENTIFIER: 'B9Ctz-VBB6MN',
        NOTIFICATION_TEST_APP_IDENTIFIER: 'aq3ddgqJSmmv',
        REALTIME_WS_API_KEY: '87e9561b-fd73-4024-ad3f-0e8c7bb28856',
        REALTIME_WS_PATH: 'ws://127.0.0.1:4000',
      }}
    >
      <React.StrictMode>
        <RouterProvider router={router} />
      </React.StrictMode>
    </Providers>
  </GoogleOAuthProvider>,
);
