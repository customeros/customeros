import { useState } from 'react';

import { runInAction } from 'mobx';
import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Input } from '@ui/form/Input';
import { Spinner } from '@ui/feedback/Spinner';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { SearchSm } from '@ui/media/icons/SearchSm';
import { Highlight } from '@ui/presentation/Highlight';
import { InfoCircle } from '@ui/media/icons/InfoCircle';
import { RefreshCw01 } from '@ui/media/icons/RefreshCw01';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { ShoppingCartAdd } from '@ui/media/icons/ShopingCartAdd';
import {
  Card,
  CardFooter,
  CardHeader,
  CardContent,
} from '@ui/presentation/Card/Card';

export const AddDomainsCard = observer(() => {
  const store = useStore();
  const [showSecondHalf, setShowSecondHalf] = useState(false);
  const { open, onOpen, onClose } = useDisclosure();

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    store.mailboxes.setDomainName(e.target.value.toLowerCase());

    if (store.mailboxes.domain === '') {
      runInAction(() => {
        store.mailboxes.domainSuggestions = [];
      });
    }

    store.mailboxes.validateDomain();
  };

  const handleInputBlur = async () => {
    if (store.mailboxes.domain.trim() !== '') {
      if (!store.mailboxes.invalidDomain) {
        store.mailboxes.getDomainSuggestions();
        setShowSecondHalf(false);
      }
    }
  };

  const toggleDisplay = () => {
    setShowSecondHalf((prev) => !prev);
  };

  const displayedDomains = showSecondHalf
    ? store.mailboxes.domainSuggestions.slice(10, 20)
    : store.mailboxes.domainSuggestions.slice(0, 10);

  return (
    <>
      <Card className='py-2 px-3 bg-white'>
        <CardHeader className='flex items-end font-medium text-sm gap-1'>
          Add outbound domains
          <IconButton
            size='xxs'
            variant='ghost'
            onClick={onOpen}
            aria-label='info'
            icon={<InfoCircle />}
          />
        </CardHeader>
        <CardContent className='p-0 text-sm'>
          Use your brand to find your ideal domains
          <CardFooter className='w-full px-0 mb-0 py-0 mt-2 relative'>
            <Input
              size='sm'
              variant='outline'
              placeholder='Brand name'
              onBlur={handleInputBlur}
              onChange={handleInputChange}
              value={store.mailboxes.domain}
              dataTest='mailboxes-domain-input'
              autoFocus={store.mailboxes.domain.length === 0}
              invalid={store.mailboxes.invalidDomain.length > 0}
              onKeyDown={(e) => {
                e.stopPropagation();

                if (e.key === 'Enter') {
                  e.currentTarget.blur();
                }
              }}
            />
            {store.mailboxes.isLoading ? (
              <Spinner
                size='sm'
                label='loading'
                className='absolute right-2 text-gray-300 fill-gray-700'
              />
            ) : (
              <IconButton
                size='xs'
                variant='ghost'
                aria-label='search'
                icon={<SearchSm />}
                className='absolute right-[2px] rounded-[2px]'
              />
            )}
          </CardFooter>
          <span
            className={cn(
              'text-[12px] ml-[9px] text-error-400',
              store.mailboxes.invalidDomain.length === 0 && 'opacity-0',
              displayedDomains.length > 0 && 'hidden',
            )}
          >
            {store.mailboxes.invalidDomain ?? '_'}
          </span>
          <div
            data-test='mailboxes-domain-list'
            className={cn(displayedDomains.length > 0 && 'mt-2')}
          >
            {displayedDomains.map((domain) => (
              <div
                key={`${domain}-${crypto.randomUUID()}`}
                className='group/item flex items-center justify-between py-1 ml-[9px] relative'
              >
                <span className='text-sm'>
                  <Highlight
                    className='font-semibold'
                    term={store.mailboxes.domain}
                  >
                    {domain}
                  </Highlight>
                </span>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='add to cart'
                  dataTest='mailboxes-domains-add-to-cart'
                  icon={<ShoppingCartAdd className='text-primary-600' />}
                  className='absolute right-0 invisible group-hover/item:visible'
                  onClick={() => {
                    store.mailboxes.selectDomain(domain);
                  }}
                />
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
              dataTest='mailboxes-domains-suggest-new'
            >
              Suggest new
            </Button>
          )}
        </CardContent>
      </Card>
      <InfoDialog
        isOpen={open}
        onClose={onClose}
        onConfirm={onClose}
        confirmButtonLabel='Got it'
        label='Email domain best practices'
        body={
          <div className='space-y-4'>
            <div>
              <p className='text-sm font-semibold'>
                Protect your primary domain
              </p>
              <p className='text-sm'>
                Keep your main and sub-domains safe—never use them for email
                outbound.
              </p>
            </div>

            <div>
              <p className='text-sm font-semibold'>.com is king</p>
              <p className='text-sm'>
                Domains with a .com tend to get better reply rates than any
                other extension.
              </p>
            </div>

            <div>
              <p className='text-sm font-semibold'>Brand recognition matters</p>
              <p className='text-sm'>
                Your brand’s name in the domain is key but minor misspellings
                won’t hurt.
              </p>
            </div>

            <div>
              <p className='text-sm font-semibold'>Smart rotation</p>
              <p className='text-sm'>
                For optimized sending, we auto rotate mailboxes and limit each
                to 40 emails a day. We recommend starting with a min of 5
                domains and adding more one by one.
              </p>
            </div>
          </div>
        }
      />
    </>
  );
});
