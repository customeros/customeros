import { Route, Routes, Navigate, useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { ProtectedRoute } from '@shared/components/ProtectedRoute';

import { Editor } from '../Editor';
import { AgentConfig } from '../AgentConfig/AgentConfig';

export const AgentSubRoutesWrapper = observer(() => {
  const store = useStore();
  const { id } = useParams<{ id: string }>();
  const agent = id ? store.agents.getById(id) : null;

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
            <div>People</div>
          </ProtectedRoute>
        }
      />
      <Route path='*' element={<Navigate replace to='setup' />} />
    </Routes>
  );
});
