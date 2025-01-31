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
    const agentId = queryParams.get('id');
    const capabilityId = queryParams.get('cid');
    const slackCode = queryParams.get('code');

    if (!store.session.isAuthenticated) return;

    if (slackCode) {
      store.settings.slack.oauthCallback(slackCode);
    }

    if (agentId) {
      navigate(
        `/agents/${agentId}${capabilityId ? `?cid=${capabilityId}` : ''}`,
      );
    }
  }, [store.session.isAuthenticated]);
};
