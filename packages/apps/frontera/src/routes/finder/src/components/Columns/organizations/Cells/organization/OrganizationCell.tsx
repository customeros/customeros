import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { useStore } from '@shared/hooks/useStore';

interface OrganizationCellProps {
  id: string;
}

export const OrganizationCell = observer(({ id }: OrganizationCellProps) => {
  const store = useStore();
  const org = store.organizations.getById(id);

  const name = org?.value?.name;
  const isSubsidiary = org?.value?.subsidiaries?.length > 0;
  const parentOrganizationName = org?.value?.parentName;
  const isEnriching = org?.isEnriching;

  const [tabs] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });
  const navigate = useNavigate();

  const fullName = name || 'Unnamed';

  const handleNavigate = () => {
    const lastPositionParams = tabs[id];
    const href = getHref(id, lastPositionParams);

    if (!href) return;

    navigate(href);
  };

  if (isEnriching) {
    return <p className='text-gray-400'>Enriching...</p>;
  }

  if (!org) return <p className='text-gray-400'>Not set</p>;

  return (
    <span className='inline'>
      {isSubsidiary && (
        <span className='text-xs text-gray-500'>{parentOrganizationName}</span>
      )}
      <p
        role='button'
        onClick={handleNavigate}
        data-test='organization-name-in-all-orgs-table'
        className='overflow-ellipsis overflow-hidden font-medium no-underline hover:no-underline cursor-pointer'
      >
        {fullName}
      </p>
    </span>
  );
});

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
