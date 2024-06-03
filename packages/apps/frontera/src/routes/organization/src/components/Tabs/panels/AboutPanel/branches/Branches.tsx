import { Link } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Plus } from '@ui/media/icons/Plus';
import { Organization } from '@graphql/types';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Card, CardHeader, CardContent } from '@ui/presentation/Card/Card';

interface BranchesProps {
  id: string;
  isReadOnly?: boolean;
  branches?: Organization['subsidiaries'];
}

export const Branches = observer(({ id, isReadOnly }: BranchesProps) => {
  const store = useStore();
  const organization = store.organizations?.value?.get(id);

  const subsidiaries = organization?.value.subsidiaries;

  const parentOrgId =
    organization?.value.parentCompanies[0]?.organization.metadata.id;

  if (!parentOrgId) return;

  return (
    <Card className='w-full mt-2 p-4 bg-white rounded-md border-1 shadow-lg'>
      <CardHeader className='flex mb-4 items-center justify-between'>
        <h2 className='text-base'>Branches</h2>
        {!isReadOnly && (
          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Add'
            onClick={() => {
              organization?.create();
            }}
            icon={<Plus className='size-4' />}
          />
        )}
      </CardHeader>
      <CardContent className='flex flex-col p-0 pt-0 gap-2 items-baseline'>
        {/* {subsidiaries &&
          subsidiaries[0].organization.metadata &&
          subsidiaries?.map((i) =>
            i.metadata.id ? (
              <Link
                className='line-clamp-1 break-keep text-gray-700 hover:text-primary-600 no-underline hover:underline'
                to={`/organization/${i.metadata?.id}?tab=about`}
                key={i.metadata?.id}
              >
                {i.name}
              </Link>
            ) : null,
          )} */}
        {subsidiaries && subsidiaries[0]?.organization.metadata.id && (
          <Link
            to={`/organization/${subsidiaries[0]?.organization.metadata.id}?tab=about`}
          >
            {subsidiaries[0]?.organization.name}
          </Link>
        )}
      </CardContent>
    </Card>
  );
});
