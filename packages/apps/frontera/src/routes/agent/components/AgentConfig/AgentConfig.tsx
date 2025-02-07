import { useMemo, useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';

import get from 'lodash/get';
import { observer } from 'mobx-react-lite';
import { AgentConfigViewUsecase } from '@domain/usecases/agents/agent-config-view.usecase.ts';

import { cn } from '@ui/utils/cn.ts';
import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';

import { capabilities } from './Capabilities';

const goalMap = {
  identify_web_visitor:
    'Identify website visitors and add them as enriched leads to CustomerOS',
  evaluate_icp_fit:
    'Qualify new leads based on whether they match your ideal customer profile or not',
  receive_reply: 'Receive replies from leads and forward them to your team',
};

const agentData = {
  __typename: 'Query',
  value: {
    __typename: 'Agent',
    id: 'campaign-manager-1',
    type: 'campaign_manager',
    name: 'Campaign manager',
    goal: 'receive_reply',
    isActive: true,
    flowId: null,
    visible: true,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    error: null,
    color: '#000000',
    icon: '',
    capabilities: [
      {
        __typename: 'Capability',
        id: 'cap-1',
        type: 'enrich_email_address',
        name: 'Enrich Email Address',
        action: 'enrich_email_address',
        active: true,
        config: '',
        errors: null,
      },
      {
        __typename: 'Capability',
        id: 'cap-2',
        type: 'validate_email_deliverability',
        name: 'Validate Email Deliverability',
        action: 'validate_email_deliverability',
        active: true,
        config: '',
        errors: null,
      },
      {
        __typename: 'Capability',
        id: 'cap-4',
        type: 'select_optimal_sending_mailbox',
        name: 'Select Optimal Sending Mailbox',
        action: 'select_optimal_sending_mailbox',
        active: true,
        config: '{}',
        errors: null,
      },
      {
        __typename: 'Capability',
        id: 'cap-3',
        type: 'manage_campaign_execution',
        name: 'Manage Campaign Execution',
        action: 'manage_campaign_execution',
        active: true,
        config: '{}',
        errors: null,
      },

      {
        __typename: 'Capability',
        id: 'cap-5',
        type: 'manage_email_delivery_failures',
        name: 'Manage Email Delivery Failures',
        action: 'manage_email_delivery_failures',
        active: true,
        config: '',
        errors: null,
      },
      {
        __typename: 'Capability',
        id: 'cap-6',
        type: 'forward_email_replies',
        name: 'Forward Email Replies',
        action: 'forward_email_replies',
        active: true,
        config: '{}',
        errors: null,
      },
    ],
  },
};
export const AgentConfig = observer(() => {
  const store = useStore();

  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [queryParams, setQueryParams] = useSearchParams();

  const agent = agentData;

  const usecase = useMemo(
    () => new AgentConfigViewUsecase(id ?? '', queryParams.get('cid')),
    [id],
  );

  const ActiveCapability = useMemo(
    () =>
      usecase.activeCapability
        ? capabilities[usecase.activeCapability.type]
        : () => null,
    [usecase?.activeCapability?.type],
  );

  console.log(
    '🏷️ -----   usecase.activeCapability: ',
    usecase.activeCapability,
  );
  useEffect(() => {
    if (!queryParams.get('cid')) {
      setQueryParams((params) => {
        // if (!usecase.activeCapability) return params;
        console.log(
          '🏷️ ----- : R',
          agent.value.capabilities.find((e) => e.config.length)?.id,
        );
        params.set(
          'cid',
          agent.value.capabilities.find((e) => e.config.length)?.id || 'cap-3',
        );

        return params;
      });
    }
  }, []);

  if (!id) {
    throw new Error('No id provided');
  }

  return (
    <div className='flex h-screen'>
      <div className='w-[448px] border-r border-r-grayModern-200 px-4 py-3'>
        <div className='mb-2'>
          <h2 className='font-medium mb-1'>Goal</h2>
          <p className='pb-2 text-sm'>
            {get(goalMap, agent?.value.goal ?? '', 'Unknown')}
          </p>
        </div>

        <h2 className='mb-2 font-medium text-sm'>This agent will:</h2>

        <ul className='space-y-1'>
          {agent?.value.capabilities.map((capability) => (
            <div
              key={capability.id}
              onClick={() => {
                if (capability.config.length) {
                  usecase.setActiveCapability(capability);
                  navigate(`?cid=${capability.id}`);
                }
              }}
              className={cn(
                'flex items-center px-2 py-1 justify-between rounded-lg select-none',
                capability.config.length &&
                  'hover:bg-grayModern-200 cursor-pointer',
                capability.id === usecase?.activeCapability?.id &&
                  'bg-grayModern-100 hover:bg-grayModern-100',
              )}
            >
              <div className='flex items-center gap-2'>
                <Icon
                  stroke={capability.errors ? 'currentColor' : 'none'}
                  name={capability.errors ? 'radio-dot' : 'dot-single'}
                  className={
                    capability.errors ? 'text-error-500' : 'text-grayModern-500'
                  }
                />
                <p className='text-sm'>{capability.name ?? 'Unknown'}</p>
              </div>

              {capability.config.length > 0 && (
                <Icon name='settings-02' className='text-grayModern-500' />
              )}
            </div>
          ))}
        </ul>
      </div>

      {usecase.activeCapability && (
        <div className='w-[418px] border-r border-r-grayModern-200 px-4 py-3'>
          <ActiveCapability />
        </div>
      )}
    </div>
  );
});
