import { Link } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';

export const ContractCell = observer(
  ({
    organizationId,
    contractId,
  }: {
    contractId: string;
    organizationId: string;
  }) => {
    const store = useStore();
    const organization = store.organizations?.value?.get(organizationId)?.value;
    const contract = store.contracts?.value?.get(contractId)?.value;

    return (
      <div>
        <p className='text-xs text-gray-500'>{organization?.name}</p>
        <Link
          to={`/organization/${organizationId}?tab=account`}
          className='font-medium line-clamp-1 text-gray-700 no-underline hover:no-underline hover:text-gray-900 transition-colors'
        >
          {contract?.contractName}
        </Link>
      </div>
    );
  },
);
