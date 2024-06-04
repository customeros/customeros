import { Link, useNavigate } from 'react-router-dom';

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

  const navigate = useNavigate();
  const organization = store.organizations?.value?.get(id);
  if (!organization) return null;

  const subsidiaries = organization.value.subsidiaries;

  const findOrgSubsidiaries = subsidiaries?.map((i) => {
    if (i.organization?.metadata?.id) {
      return store.organizations.value.get(i.organization.metadata?.id);
    }
  });

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
              store.organizations.create(undefined, {
                onSucces(serverId) {
                  const findOrg = store.organizations.value.get(serverId);
                  if (!findOrg) return;
                  setTimeout(() => {
                    organization.update((org) => {
                      org.subsidiaries = [{ organization: findOrg.value }];

                      return org;
                    });
                  }, 100);

                  navigate(`/organization/${serverId}?tab=about`);
                },
              });
            }}
            icon={<Plus className='size-4' />}
          />
        )}
      </CardHeader>
      <CardContent className='flex flex-col p-0 pt-0 gap-2 items-baseline'>
        {findOrgSubsidiaries?.map((i) =>
          i?.value.metadata.id ? (
            <Link
              className='line-clamp-1 break-keep text-gray-700 hover:text-primary-600 no-underline hover:underline'
              to={`/organization/${i.value.metadata.id}?tab=about`}
              key={i.value.metadata.id}
            >
              {i.value.name}
            </Link>
          ) : null,
        )}
      </CardContent>
    </Card>
  );
});
