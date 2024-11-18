import { useState } from 'react';

import { useLocalStorage } from 'usehooks-ts';

import { cn } from '@ui/utils/cn';
import { IconButton } from '@ui/form/IconButton';
import { MinusCircle } from '@ui/media/icons/MinusCircle';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

export const CheckoutCard = () => {
  const [expanded, setExpanded] = useState(false);
  const [storedBrandName, _setStoredBrandName] = useLocalStorage<string[]>(
    'brandName',
    [],
  );
  const [storeUserName, _setStoreUserName] = useLocalStorage<string[]>(
    'userName',
    [],
  );

  return (
    <Card className='py-2 px-3 bg-white'>
      <CardHeader className='flex items-center justify-between font-medium text-sm'>
        <span>Base Boundle</span>
        <span>$199.99</span>
      </CardHeader>
      <CardContent className='p-0'>
        <span>{storedBrandName.length} of 5 domains</span>

        {storedBrandName.map((name, index) => (
          <div
            key={`${name}-${index}`}
            className='flex items-center justify-between mt-1 bg-gray-100 rounded-[4px] py-1 px-2'
          >
            <span className='text-sm'>{name}</span>
            <IconButton
              size='xxs'
              variant='ghost'
              icon={<MinusCircle />}
              aria-label='remove-domain'
              onClick={() => {
                _setStoredBrandName((prev) => {
                  return prev.filter((item) => item !== name);
                });
              }}
            />
          </div>
        ))}
        {Array.from({ length: 5 - storedBrandName.length }).map((_, index) => (
          <div
            key={`placeholder-${index}`}
            className='flex items-center justify-between mt-1 border-dotted border border-gray-300 text-gray-400  rounded-[4px] py-1 px-2'
          >
            <span className='text-sm'>
              Domain {storedBrandName.length + index + 1}
            </span>
          </div>
        ))}
        {storeUserName.length > 0 && (
          <CardHeader className='mt-2 flex justify-between items-center'>
            <span className='font-medium'>
              {storeUserName.length * storedBrandName.length} of 10 mailboxes
            </span>
            <IconButton
              size='xxs'
              variant='ghost'
              aria-label='Expand'
              icon={<ChevronExpand />}
              onClick={() => setExpanded(!expanded)}
            />
          </CardHeader>
        )}
        {expanded &&
          storeUserName.map((user, userIndex) =>
            storedBrandName.map((brand, brandIndex) => {
              const isLastTwo =
                userIndex === storeUserName.length - 1 &&
                brandIndex >= storedBrandName.length - 2;

              return (
                <div
                  key={`${user}-${brand}-${userIndex}-${brandIndex}`}
                  className={cn(
                    'flex items-center justify-between mt-1',
                    !isLastTwo &&
                      (userIndex * storedBrandName.length + brandIndex) % 2 ===
                        1
                      ? 'mb-3'
                      : '',
                  )}
                >
                  <span className='text-sm'>{`${user}@${brand}`}</span>
                </div>
              );
            }),
          )}
      </CardContent>
    </Card>
  );
};
