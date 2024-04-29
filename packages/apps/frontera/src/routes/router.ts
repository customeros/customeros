import { createBrowserRouter } from 'react-router-dom';

import { RootRoute } from './root';
import { AuthRoute } from './auth/route';

export const router = createBrowserRouter([RootRoute, AuthRoute]);
