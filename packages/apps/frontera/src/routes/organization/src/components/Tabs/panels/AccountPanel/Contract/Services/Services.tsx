import { Plus } from '@ui/media/icons/Plus';
import { ServiceLineItem } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { ServicesList } from '@organization/components/Tabs/panels/AccountPanel/Contract/Services/ServicesList';

interface ServicesProps {
  id: string;
  onModalOpen: () => void;
  currency?: string | null;
  data?: Array<ServiceLineItem> | null;
}

export const Services = ({
  id,
  data,
  currency,
  onModalOpen,
}: ServicesProps) => {
  return (
    <>
      <p className='w-full flex items-center justify-between'>
        {!data?.length && (
          <span className='text-sm font-semibold mt-2'>No services</span>
        )}

        {!data?.length && (
          <IconButton
            size='xs'
            variant='ghost'
            colorScheme='gray'
            aria-label={'Add services'}
            data-test='org-account-add-services'
            icon={<Plus className='text-gray-400' />}
            onClick={() => {
              onModalOpen();
            }}
          />
        )}
      </p>

      {!!data?.length && (
        <ServicesList id={id} currency={currency} onModalOpen={onModalOpen} />
      )}
    </>
  );
};
