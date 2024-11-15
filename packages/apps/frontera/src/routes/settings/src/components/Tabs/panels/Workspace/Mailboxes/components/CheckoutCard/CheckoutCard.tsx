import { useLocalStorage } from 'usehooks-ts';

import { IconButton } from '@ui/form/IconButton';
import { MinusCircle } from '@ui/media/icons/MinusCircle';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

export const CheckoutCard = () => {
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
      <CardContent>
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
          <CardHeader>
            <span className='font-medium'>
              {storeUserName.length} of 10 mailboxes
            </span>
          </CardHeader>
        )}
        {storeUserName.map((user, userIndex) =>
          storedBrandName.map((brand, brandIndex) => (
            <div
              key={`${user}-${brand}-${userIndex}-${brandIndex}`}
              className='flex items-center justify-between mt-1 bg-gray-100 rounded-[4px] py-1 px-2'
            >
              <span className='text-sm'>{`${user}@${brand}`}</span>
            </div>
          )),
        )}
      </CardContent>
    </Card>
  );
};
