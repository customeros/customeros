import { observer } from 'mobx-react-lite';

import { ContractStatus } from '@graphql/types';
import { Divider } from '@ui/presentation/Divider/Divider';

import { ProductsList } from './components/ProductsList';
import { AddNewProductMenu } from './components/AddNewProductMenu';

interface SubscriptionProductModalProps {
  id: string;
  currency?: string;
  contractStatus?: ContractStatus | null;
}

export const Products = observer(
  ({ id, currency, contractStatus }: SubscriptionProductModalProps) => {
    return (
      <>
        <div className='flex relative items-center h-8 '>
          <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
            Products
          </p>
          <Divider />
          <AddNewProductMenu contractId={id} />
        </div>

        <ProductsList
          id={id}
          currency={currency}
          contractStatus={contractStatus}
        />
      </>
    );
  },
);
