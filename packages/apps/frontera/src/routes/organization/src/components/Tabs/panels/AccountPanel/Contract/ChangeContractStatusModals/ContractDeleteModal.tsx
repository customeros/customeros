import React from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Trash01 } from '@ui/media/icons/Trash01.tsx';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';

interface ContractStartModalProps {
  contractId: string;
  onClose: () => void;
}

export const ContractDeleteModal = observer(
  ({ onClose, contractId }: ContractStartModalProps) => {
    const store = useStore();
    const organizationId = useParams()?.id as string;

    const handleDeleteContract = () => {
      const contractsStore = store.contracts;

      contractsStore.delete(contractId, organizationId);

      onClose();
    };

    return (
      <>
        <div
          className={
            'rounded-2xl max-w-[600px] h-full flex flex-col justify-between'
          }
        >
          <div>
            <div>
              <FeaturedIcon size='lg' colorScheme='error'>
                <Trash01 className='text-error-600' />
              </FeaturedIcon>

              <h1 className={cn('text-lg font-semibold  mt-5 mb-1')}>
                Delete this contract?
              </h1>
            </div>
            <div className='flex flex-col'>
              <p className='text-sm mt-3'>
                Are you sure you want to delete this contract?
              </p>
            </div>
          </div>

          <div className='mt-6 flex'>
            <Button
              size='lg'
              variant='outline'
              onClick={onClose}
              className='w-full'
            >
              Cancel
            </Button>
            <Button
              size='lg'
              variant='outline'
              colorScheme='error'
              className='ml-3 w-full'
              onClick={handleDeleteContract}
              dataTest='contract-card-confirm-contract-deletion'
            >
              Delete contract
            </Button>
          </div>
        </div>
      </>
    );
  },
);
