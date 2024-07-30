import React from 'react';
import { RouterProvider } from 'react-router-dom';

import ReactDOM from 'react-dom/client';

import { Providers } from '@shared/components/Providers/Providers';

import { router } from './routes/router';

import './styles/cmdk.scss';
import './styles/globals.scss';
import './styles/toastify.scss';
import './styles/normalization.scss';
import './styles/react-calendar.scss';
import './styles/remirror-editor.scss';
import 'react-toastify/dist/ReactToastify.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <Providers>
    <React.StrictMode>
      <RouterProvider router={router} />
    </React.StrictMode>
  </Providers>,
);
