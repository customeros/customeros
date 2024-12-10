import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { MinusCircle } from '@ui/media/icons/MinusCircle';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

interface BaseBundleCardProps {
  readonly?: boolean;
  defaultExpanded?: boolean;
}

export const BaseBundleCard = observer(
  ({ readonly, defaultExpanded = true }: BaseBundleCardProps) => {
    const store = useStore();
    const invalidDomains = store.mailboxes.invalidDomains;
    const [expanded, setExpanded] = useState(defaultExpanded);

    return (
      <Card className='py-2 px-3 bg-white'>
        <CardHeader>
          <div className='flex items-center justify-between font-medium text-sm mb-1'>
            <span>Starter bundle</span>
            <span>$199.99</span>
          </div>
          <p className='text-sm'>Unlock automated flows within CustomerOS</p>
        </CardHeader>
        <CardContent className='p-0'>
          <span className='text-sm font-medium'>
            {store.mailboxes.baseBundle.size} of 5 domains
          </span>

          {Array.from(store.mailboxes.baseBundle).map((name, index) => {
            const isInvalid = invalidDomains.includes(name);

            return (
              <div
                key={`${name}-${index}`}
                className={cn(
                  'flex items-center justify-between mt-1 rounded-[4px] py-1 px-2 border border-gray-100 bg-gray-100',
                  isInvalid && 'bg-error-50 border-error-50',
                )}
              >
                <span className='text-sm'>{name}</span>
                {!readonly && (
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    icon={<MinusCircle />}
                    aria-label='remove-domain'
                    onClick={() => {
                      store.mailboxes.removeDomain(name);
                    }}
                  />
                )}
              </div>
            );
          })}
          {Array.from({ length: 5 - store.mailboxes.baseBundle.size }).map(
            (_, index) => (
              <div
                key={index}
                className='flex items-center justify-between mt-1 border-dotted border border-gray-300 text-gray-400 rounded-[4px] py-1 px-2'
              >
                <span className='text-sm'>
                  Domain {store.mailboxes.baseBundle.size + index + 1}
                </span>
              </div>
            ),
          )}
          {store.mailboxes.hasUsernames && (
            <CardHeader className='mt-2 flex justify-between items-center'>
              <span className='font-medium text-sm'>
                {store.mailboxes.baseBundle.size *
                  store.mailboxes.usernamesCount}{' '}
                of 10 mailboxes
              </span>
              <IconButton
                size='xxs'
                variant='ghost'
                aria-label='Expand'
                onClick={() => setExpanded(!expanded)}
                icon={!expanded ? <ChevronExpand /> : <ChevronCollapse />}
              />
            </CardHeader>
          )}
          {expanded &&
            Array.from(store.mailboxes.baseBundle).reduce<JSX.Element[]>(
              (acc, brand, brandIndex) => {
                const brandUserCombos = store.mailboxes.usernames
                  .filter((v) => v !== '')
                  .map((user, userIndex) => {
                    return (
                      <div
                        key={`${user}-${brand}-${userIndex}-${brandIndex}`}
                        className={cn(
                          'flex items-center justify-between ml-2 ',
                          (brandIndex * store.mailboxes.usernamesCount +
                            userIndex) %
                            2 ===
                            1
                            ? 'mb-3'
                            : '',
                        )}
                      >
                        <span className='text-sm'>{`${user}@${brand}`}</span>
                      </div>
                    );
                  });

                return [...acc, ...brandUserCombos];
              },
              [],
            )}
        </CardContent>
      </Card>
    );
  },
);
