import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { MinusCircle } from '@ui/media/icons/MinusCircle';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

interface ExtendedBundleCardProps {
  readonly?: boolean;
  defaultExpanded?: boolean;
}

export const ExtendedBundleCard = observer(
  ({ readonly, defaultExpanded = true }: ExtendedBundleCardProps) => {
    const store = useStore();
    const invalidDomains = store.mailboxes.invalidDomains;
    const extendedDomains = store.mailboxes.extendedBundle;
    const [expanded, setExpanded] = useState(defaultExpanded);

    return (
      <Card className='py-2 px-3 bg-white'>
        <CardHeader className='font-medium'>
          <div className='flex items-center justify-between text-sm'>
            Additional domains
            {extendedDomains.size > 0 && (
              <span className='text-sm ml-2'>
                {`${extendedDomains.size} x $18.99`}
              </span>
            )}
          </div>
        </CardHeader>

        <CardContent className='p-0'>
          {Array.from(extendedDomains).map((domain, index) => {
            const isInvalid = invalidDomains.includes(domain);

            return (
              <div
                key={`${domain}-${index}`}
                className={cn(
                  'flex items-center justify-between mt-1 rounded-[4px] py-1 px-2 border border-gray-100 bg-gray-100',
                  isInvalid && 'bg-error-50 border-error-50',
                )}
              >
                <span className='text-sm'>{domain}</span>
                {!readonly && (
                  <IconButton
                    size='xxs'
                    variant='ghost'
                    icon={<MinusCircle />}
                    aria-label='remove-domain'
                    onClick={() => {
                      store.mailboxes.removeDomain(domain);
                    }}
                  />
                )}
              </div>
            );
          })}
          {extendedDomains.size === 0 && (
            <p className='text-sm'>
              {' '}
              Add more domains at $18.99 each (53% off)
            </p>
          )}
          {store.mailboxes.usernamesCount > 0 && extendedDomains.size > 0 && (
            <CardHeader className='mt-2 flex justify-between items-center'>
              <span className='font-medium text-sm'>
                {store.mailboxes.usernamesCount * extendedDomains.size}{' '}
                mailboxes
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
            Array.from(extendedDomains).reduce<JSX.Element[]>(
              (acc, brand, brandIndex) => {
                const brandUserCombos = store.mailboxes.usernames
                  .filter((v) => v !== '')
                  .map((user, userIndex) => {
                    return (
                      <div
                        key={`${user}-${brand}-${userIndex}-${brandIndex}`}
                        className={cn(
                          'flex items-center justify-between ml-2',
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
