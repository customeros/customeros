import React from 'react';
import { RouterProvider } from 'react-router-dom';

import ReactDOM from 'react-dom/client';

import { Providers } from '@shared/components/Providers/Providers';

import { router } from './routes/router';

import './styles/globals.scss';
import './styles/filepond.scss';
import './styles/date-picker.scss';
import './styles/normalization.scss';
import './styles/remirror-editor.scss';
import 'react-toastify/dist/ReactToastify.css';
import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.min.css';
import 'filepond/dist/filepond.min.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
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
  </Providers>,
);
