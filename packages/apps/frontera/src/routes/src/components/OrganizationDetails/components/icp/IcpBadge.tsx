import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn.ts';
import { IcpFit } from '@graphql/types';
import { Icon, IconName } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { Target04 } from '@ui/media/icons/Target04.tsx';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';
import { Tag, TagLabel, TagLeftIcon, TagRightIcon } from '@ui/presentation/Tag';

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
    icon: 'target-04',
    colorScheme: 'gray',
  },
};

interface IcpBadgeProps {
  id: string;
}

export const IcpBadge = observer(({ id }: IcpBadgeProps) => {
  const store = useStore();
  const [open, setOpen] = useState(false);
  const icpAgent = store.agents.icpQualificationAgent;
  const navigate = useNavigate();
  const icpAgentActive = icpAgent?.value.isActive;

  const organization = store.organizations.getById(id);

  if (!organization) return null;

  const data = icpData?.[organization.value.icpFit];
  const icpFitReasons = organization.value.icpFitReasons;
  const icpFitUpdatedAt = organization.value.icpFitUpdatedAt;

  if (!data) return null;

  const noAgentConfigured =
    organization.value.icpFit === IcpFit.IcpNotSet &&
    !icpFitReasons.length &&
    !icpAgentActive;
  const icpProfilingInProgress =
    organization.value.icpFit === IcpFit.IcpNotSet &&
    icpAgentActive &&
    !icpFitUpdatedAt &&
    !icpFitReasons.length;

  return (
    <Popover open={open} modal={true} onOpenChange={(value) => setOpen(value)}>
      <PopoverTrigger>
        {noAgentConfigured ? (
          <Button
            size='xxs'
            variant='ghost'
            colorScheme='primary'
            leftIcon={<Target04 />}
          >
            ICP fit?
          </Button>
        ) : (
          <Tag className='ml-4' variant='subtle' colorScheme={data.colorScheme}>
            {data.icon && (
              <TagLeftIcon className='mr-1'>
                <div>
                  <Icon
                    width={12}
                    height={12}
                    name={data.icon}
                    className={cn('size-3 text-gray-500', {
                      [`text-${data?.colorScheme}-500`]: true,
                    })}
                  />
                </div>
              </TagLeftIcon>
            )}

            <TagLabel className='flex items-center'>{data.label}</TagLabel>
            <TagRightIcon className='ml-1'>
              <div>
                <Icon
                  width={12}
                  height={12}
                  name='chevron-down'
                  className={cn('size-3 text-gray-500', {
                    [`text-${data?.colorScheme}-500`]: true,
                  })}
                />
              </div>
            </TagRightIcon>
          </Tag>
        )}
      </PopoverTrigger>
      <PopoverContent align='end' side='bottom' className='text-sm'>
        <div className='max-w-[295px]'>
          {icpAgent?.value.isActive &&
            organization.value.icpFit !== IcpFit.IcpNotSet &&
            organization.value.icpFitReasons.length && (
              <>
                <p>
                  The
                  <span className='mx-1 font-medium underline underline-offset-1'>
                    ICP qualifier
                  </span>
                  agent determined that this company{' '}
                  {organization.value.icpFit === IcpFit.IcpFit
                    ? 'fits'
                    : 'does not fit'}{' '}
                  your ideal customer profile.
                </p>

                <div className='pt-3'>
                  <span>Here's why:</span>
                  <ol className='list-decimal pl-5'>
                    {icpFitReasons.map((reason, index) => (
                      <li key={`${index}-reason`}>{reason}</li>
                    ))}
                  </ol>
                </div>
                <div className='p-2 px-3 mt-2 bg-grayModern-50 flex items-center'>
                  <Icon
                    name='message-question-circle'
                    className='text-grayModern-500'
                  />
                  <p className='mx-2'>Is this qualification correct?</p>
                  <IconButton
                    size='xxs'
                    aria-label={''}
                    variant='ghost'
                    icon={
                      <Icon
                        name='thumbs-down'
                        className='text-grayModern-500 hover:text-grayModern-700'
                      />
                    }
                  />
                </div>
              </>
            )}
          {icpProfilingInProgress && (
            <>
              <p>
                The
                <span className='mx-1 font-medium underline underline-offset-1'>
                  ICP qualifier
                </span>
                agent is busy determining whether this company fits your ideal
                customer profile or not
              </p>
            </>
          )}

          {noAgentConfigured && (
            <>
              <p>
                To determine whether this company fits your ideal customer
                profile, configure and enable the
                <span className='mx-1 font-medium underline underline-offset-1'>
                  ICP qualifier
                </span>
                agent.
              </p>
              <Button
                size='xs'
                variant='outline'
                colorScheme='primary'
                className={'w-full mt-4'}
                onClick={() => {
                  navigate(`/agents/${icpAgent?.id}`);
                }}
              >
                Go to ICP qualifier
              </Button>
            </>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
});
