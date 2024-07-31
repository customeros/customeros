import { FC, PropsWithChildren } from 'react';

import { Button } from '@ui/form/Button/Button';
import { File02 } from '@ui/media/icons/File02';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { RelationshipButton } from '@organization/components/Tabs/panels/AccountPanel/RelationshipButton';
import { OrganizationPanel } from '@organization/components/Tabs/shared/OrganizationPanel/OrganizationPanel';

export const EmptyContracts: FC<
  PropsWithChildren<{ isPending: boolean; onCreate: () => void }>
> = ({ children, onCreate, isPending }) => {
  return (
    <OrganizationPanel title='Account' actionItem={<RelationshipButton />}>
      <article className='my-4 w-full flex flex-col items-center'>
        <FeaturedIcon size='lg' className='mb-4' colorScheme='primary'>
          <File02 className='size-4' />
        </FeaturedIcon>
        <h1 className='text-md font-semibold'>Draft a new contract</h1>

        <Button
          size='sm'
          variant='outline'
          onClick={onCreate}
          colorScheme='primary'
          isDisabled={isPending}
          className='text-sm mt-6 w-fit'
          data-Test='org-account-empty-new-contract'
        >
          {isPending ? 'Creating contract...' : 'New contract'}
        </Button>
      </article>
      {children}
    </OrganizationPanel>
  );
};
