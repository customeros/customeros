import { useState } from 'react';

import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn';
import { IconButton } from '@ui/form/IconButton';
import { MinusCircle } from '@ui/media/icons/MinusCircle';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

export const AdditionalDomainsCard = () => {
  const [expanded, setExpanded] = useState(false);

  const [selectedAdditionalDomains, setSelectedAdditionalDomains] =
    useLocalStorage<string[]>('selectedAdditionalDomains', []);

  const [_storedBrandName, _setStoredBrandName] = useLocalStorage<string[]>(
    'brandName',
    [],
  );
  const [storeUserName, _setStoreUserName] = useLocalStorage<string[]>(
    'userName',
    [],
  );

  return (
    <Card className='py-2 px-3 bg-white'>
      <CardHeader className='font-medium'>
        <div className='flex items-center justify-between text-sm'>
          Additional domains
          {selectedAdditionalDomains.length > 0 && (
            <span className='text-sm ml-2'>
              {`${selectedAdditionalDomains.length} x $18.99`}
            </span>
          )}
        </div>
      </CardHeader>

      <CardContent className='p-0'>
        {selectedAdditionalDomains.map((domain, index) => (
          <div
            key={`${domain}-${index}`}
            className='flex items-center justify-between mt-1 bg-gray-100 rounded-[4px] py-1 px-2'
          >
            <span className='text-sm'>{domain}</span>
            <IconButton
              size='xxs'
              variant='ghost'
              icon={<MinusCircle />}
              aria-label='remove-domain'
              onClick={() => {
                setSelectedAdditionalDomains((prev) => {
                  return prev.filter((item) => item !== domain);
                });
              }}
            />
          </div>
        ))}
        {selectedAdditionalDomains.length === 0 && (
          <p className='text-sm'> Add more domains at $18.99 each(53% off)</p>
        )}
        {storeUserName.length > 0 && selectedAdditionalDomains.length > 0 && (
          <CardHeader className='mt-2 flex justify-between items-center'>
            <span className='font-medium text-sm'>
              {storeUserName.length * selectedAdditionalDomains.length}{' '}
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
          selectedAdditionalDomains.reduce<JSX.Element[]>(
            (acc, brand, brandIndex) => {
              const brandUserCombos = storeUserName.map((user, userIndex) => {
                return (
                  <div
                    key={`${user}-${brand}-${userIndex}-${brandIndex}`}
                    className={cn(
                      'flex items-center justify-between ml-2',
                      (brandIndex * storeUserName.length + userIndex) % 2 === 1
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
};
