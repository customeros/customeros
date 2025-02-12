import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

import { useStore } from '@shared/hooks/useStore';

/*
 * This hook is used to handle the slack oauth callback
 * and redirect to the agent page if the agent id is present.
 * add-slack-channel.usecase initiates a slack integration connection if none is present
 * slack redirects back to the /agents route with the agent id, capability id and a code.
 * The code is then used to authenticate the slack integration.
 */
export const useSlackOauthCallback = () => {
  const store = useStore();
  const navigate = useNavigate();
  const [queryParams] = useSearchParams();

  useEffect(() => {
    const state = queryParams.get('state');
    const slackCode = queryParams.get('code');
    const [id, cid] = decodeURIComponent(state ?? '').split(':');

    if (!store.session.isAuthenticated) return;

    if (slackCode) {
      store.settings.slack.oauthCallback(
        slackCode,
        `${import.meta.env.VITE_CLIENT_APP_URL}/agents`,
      );
    }

    if (id) {
      navigate(`/agents/${id}?${cid}`, { replace: true });
    }
  }, [store.session.isAuthenticated]);
};
