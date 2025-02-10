import { Route, Routes, Navigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { ProtectedRoute } from '@shared/components/ProtectedRoute';

import { Editor } from '../Editor';
import { List } from '../List/List';
import { AgentConfig } from '../AgentConfig/AgentConfig';

export const AgentSubRoutesWrapper = observer(() => {
  // todo change condition
  return (
    <Routes>
      <Route path='setup' element={<AgentConfig />} />
      <Route
        path='editor'
        element={
          <ProtectedRoute condition={true} fallback={'../setup'}>
            <Editor />
          </ProtectedRoute>
        }
      />
      <Route
        path='list'
        element={
          <ProtectedRoute condition={true} fallback={'../setup'}>
            <List />
          </ProtectedRoute>
        }
      />
      <Route index element={<Navigate replace to={`setup`} />} />
    </Routes>
  );
});
