import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { Input } from '@ui/form/Input';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { RefreshCw01 } from '@ui/media/icons/RefreshCw01';
import { ShoppingCartAdd } from '@ui/media/icons/ShopingCartAdd';
import {
  Card,
  CardFooter,
  CardHeader,
  CardContent,
} from '@ui/presentation/Card/Card';

export const AddDomainsCard = observer(() => {
  const store = useStore();
  const [isHovered, setIsHovered] = useState<number | null>(null);
  const [showSecondHalf, setShowSecondHalf] = useState(false);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    store.mailboxes.setDomainName(e.target.value);
  };

  const handleInputBlur = async () => {
    if (store.mailboxes.domain.trim() !== '') {
      store.mailboxes.getDomainSuggestions();
      setShowSecondHalf(false);
    }
  };

  const toggleDisplay = () => {
    setShowSecondHalf((prev) => !prev);
  };

  const displayedDomains = showSecondHalf
    ? store.mailboxes.domainSuggestions.slice(10, 20)
    : store.mailboxes.domainSuggestions.slice(0, 10);

  return (
    <Card className='py-2 px-3 bg-white'>
      <CardHeader className='flex items-end font-medium text-sm gap-1'>
        Add outbound domains
        <IconButton
          size='xxs'
          variant='ghost'
          aria-label='info'
          icon={<InfoCircle />}
        />
      </CardHeader>
      <CardContent className='p-0 text-sm'>
        Search and add your ideal outbound domains based on your brand
        <CardFooter className='w-full px-0 mb-0 py-0 mt-2 relative'>
          <Input
            size='sm'
            variant='outline'
            placeholder='Brand name'
            onBlur={handleInputBlur}
            onChange={handleInputChange}
            value={store.mailboxes.domain}
          />
          <SearchSm className='absolute right-2 text-gray-500' />
        </CardFooter>
        <div className='mt-2'>
          {displayedDomains.map((domain, index) => (
            <div
              onMouseLeave={() => setIsHovered(null)}
              key={`${domain}-${crypto.randomUUID()}`}
              onMouseEnter={() => setIsHovered(index)}
              className='flex items-center justify-between py-1 ml-[9px]'
            >
              <span className='text-sm'>{domain}</span>
              {isHovered === index && (
                <IconButton
                  size='xxs'
                  variant='ghost'
                  aria-label='add to cart'
                  icon={<ShoppingCartAdd className='text-primary-700' />}
                  onClick={() => {
                    store.mailboxes.selectDomain(domain);
                  }}
                />
              )}
            </div>
          ))}
        </div>
        {store.mailboxes.domainSuggestions.length > 10 && (
          <Button
            size='xxs'
            variant='ghost'
            colorScheme='primary'
            onClick={toggleDisplay}
            className='ml-[3px] mt-1'
            leftIcon={<RefreshCw01 />}
          >
            Suggest new
          </Button>
        )}
      </CardContent>
    </Card>
  );
});
