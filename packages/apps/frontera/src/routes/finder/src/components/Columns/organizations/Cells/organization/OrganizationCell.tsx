import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useLocalStorage } from 'usehooks-ts';

import { Eye } from '@ui/media/icons/Eye';
import { useStore } from '@shared/hooks/useStore';

interface OrganizationCellProps {
  id: string;
}

export const OrganizationCell = observer(({ id }: OrganizationCellProps) => {
  const store = useStore();
  const org = store.organizations.getById(id);
  const [previewCard, setPreviewCard] = useLocalStorage('previewCard', false);

  const name = org?.value?.name;
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
    <div className='flex items-center gap-2 group'>
      <p
        onClick={handleNavigate}
        data-test='organization-name-in-all-orgs-table'
        className='overflow-ellipsis overflow-hidden font-medium no-underline hover:no-underline cursor-pointer'
      >
        {fullName}
      </p>
      <Eye
        className='opacity-0 group-hover:opacity-100 text-gray-500'
        onClick={() => {
          if (previewCard === true && store.ui.focusRow === id) {
            setPreviewCard(false);
          } else {
            store.ui.setFocusRow(id);
            setPreviewCard(true);
          }
        }}
      />
    </div>
  );
});

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
