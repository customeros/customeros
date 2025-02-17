import { useMemo, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { CreateAgentUsecase } from '@domain/usecases/agents/create-agent.usecase';

import { cn } from '@ui/utils/cn';
import { Spinner } from '@ui/feedback/Spinner';
import { Icon, IconName } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { IcpFit, AgentType } from '@graphql/types';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';
import {
  Tag,
  TagLabel,
  TagLeftIcon,
  TagRightButton,
} from '@ui/presentation/Tag';

const icpData: Record<
  IcpFit,
  {
    label: string;
    icon?: IconName;
    colorScheme: 'success' | 'warning' | 'gray';
  }
> = {
  [IcpFit.IcpFit]: {
    label: 'ICP match',
    icon: 'check-verified-02',
    colorScheme: 'success',
  },
  [IcpFit.IcpNotFit]: {
    label: 'Not ICP',
    icon: undefined,
    colorScheme: 'warning',
  },
  [IcpFit.IcpNotSet]: {
    label: 'ICP not set',
    icon: undefined,
    colorScheme: 'gray',
  },
};

type PopoverState =
  | { isMatch: boolean; reasons: string[]; type: 'has_reasons' }
  | { type: 'agent_error' }
  | { type: 'agent_inactive' }
  | { type: 'profiling' }
  | { type: 'manual_change' };

interface IcpBadgeProps {
  id: string;
}

export const IcpBadge = observer(({ id }: IcpBadgeProps) => {
  const store = useStore();
  const [open, setOpen] = useState(false);

  const icpAgent = store.agents.getFirstAgentByType(AgentType.IcpQualifier);
  const icpAgentActive = icpAgent?.value?.isActive;
  const icpAgentHasError =
    icpAgent?.value?.error !== null || !icpAgent?.value?.isConfigured;
  const organization = store.organizations.getById(id);
  const navigate = useNavigate();

  const navigateToAgent = (id: string) => {
    navigate(`/agents/${id}`);
  };
  const usecase = useMemo(() => new CreateAgentUsecase(navigateToAgent), []);

  if (!organization) return null;

  const data = organization.value.icpFit
    ? icpData[organization.value.icpFit]
    : icpData[IcpFit.IcpNotSet];

  const handleAgentNavigation = () => {
    if (!icpAgent) {
      return usecase.execute(AgentType.IcpQualifier);
    }

    return navigate(`/agents/${icpAgent.id}`);
  };

  const getPopoverState = (): PopoverState => {
    const { icpFitReasons, icpFit, icpFitUpdatedAt } = organization.value;

    if (icpFitReasons.length > 0) {
      return {
        type: 'has_reasons',
        reasons: icpFitReasons,
        isMatch: icpFit === IcpFit.IcpFit,
      };
    }

    if (icpAgentHasError) {
      return { type: 'agent_error' };
    }

    if (!icpAgentActive || !icpAgent) {
      return { type: 'agent_inactive' };
    }

    const isProfilingInProgress =
      icpFit === IcpFit.IcpNotSet && !icpFitUpdatedAt;

    if (isProfilingInProgress) {
      return { type: 'profiling' };
    }

    return { type: 'manual_change' };
  };

  const popoverState = getPopoverState();
  const showSpinner = popoverState.type === 'profiling';

  return (
    <Popover open={open} modal={true} onOpenChange={setOpen}>
      <PopoverTrigger>
        <Tag className='ml-4' variant='subtle' colorScheme={data.colorScheme}>
          {data.icon && (
            <TagLeftIcon className='mr-1'>
              <div>
                {showSpinner ? (
                  <Spinner
                    size='xs'
                    label='icp profiling'
                    className='text-gray-300 fill-gray-400'
                  />
                ) : (
                  <Icon
                    width={12}
                    height={12}
                    name={data.icon}
                    className={cn('size-3 text-gray-500', {
                      [`text-${data?.colorScheme}-500`]: true,
                    })}
                  />
                )}
              </div>
            </TagLeftIcon>
          )}

          <TagLabel className='flex items-center whitespace-nowrap'>
            {data.label}
          </TagLabel>
          <TagRightButton>
            <Icon
              name={open ? 'chevron-up' : 'chevron-down'}
              className={cn('text-gray-500 ml-0', {
                [`text-${data?.colorScheme}-500`]: true,
              })}
            />
          </TagRightButton>
        </Tag>
      </PopoverTrigger>
      <PopoverContent align='end' side='bottom' className='text-sm'>
        <div className='max-w-[295px]'>
          <PopoverContents
            state={popoverState}
            agentId={icpAgent?.id}
            onNavigate={handleAgentNavigation}
          />
        </div>
      </PopoverContent>
    </Popover>
  );
});

const AgentLink = ({ agentId }: { agentId?: string }) => (
  <Link
    to={`/agents/${agentId}`}
    className='mx-1 font-medium underline underline-offset-1 cursor-pointer'
  >
    ICP qualifier
  </Link>
);

const QualifierButton = ({ onClick }: { onClick: () => void }) => (
  <Button
    size='xs'
    variant='outline'
    onClick={onClick}
    colorScheme='primary'
    className='w-full mt-4'
  >
    Go to ICP qualifier
  </Button>
);

const PopoverContents = ({
  state,
  agentId,
  onNavigate,
}: {
  agentId?: string;
  state: PopoverState;
  onNavigate: () => void;
}) => {
  switch (state.type) {
    case 'has_reasons':
      return (
        <>
          <p>
            The <AgentLink agentId={agentId} /> agent determined that this
            company {state.isMatch ? 'fits' : 'does not fit'} your ideal
            customer profile.
          </p>
          <div className='pt-3'>
            <span>Here's why:</span>
            <ol className='list-decimal pl-5'>
              {state.reasons.map((reason, index) => (
                <li key={`${index}-reason`}>{reason}</li>
              ))}
            </ol>
          </div>
        </>
      );

    case 'agent_error':
      return (
        <>
          <p>
            To determine whether this company fits your ideal customer profile,
            ensure the ICP qualifier agent is configured and enabled without
            errors.
          </p>
          <QualifierButton onClick={onNavigate} />
        </>
      );

    case 'agent_inactive':
      return (
        <>
          <p>
            To determine whether this company fits your ideal customer profile,
            configure and enable the
            <span className='mx-1 font-medium'>ICP qualifier</span>
            agent.
          </p>
          <QualifierButton onClick={onNavigate} />
        </>
      );

    case 'profiling':
      return (
        <p>
          The <AgentLink agentId={agentId} /> agent is busy determining whether
          this company fits your ideal customer profile or not
        </p>
      );

    case 'manual_change':
      return (
        <p>
          A user changed this company's ICP status by updating its relationship
          and stage.
        </p>
      );
  }
};
