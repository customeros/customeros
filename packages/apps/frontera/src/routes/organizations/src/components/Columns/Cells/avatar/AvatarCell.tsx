import { useNavigate } from 'react-router-dom';

import { useLocalStorage } from 'usehooks-ts';

import { Avatar } from '@ui/media/Avatar/Avatar';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';

interface AvatarCellProps {
  id: string;
  name: string;
  src?: string | null;
}

export const AvatarCell = ({ name, id, src }: AvatarCellProps) => {
  const navigate = useNavigate();
  const [tabs] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });

  const lastPositionParams = tabs[id];
  const href = getHref(id, lastPositionParams);
  const fullName = name || 'Unnamed';

  return (
    <div className='items-center'>
      <Tooltip align='center' side='bottom' label={fullName} asChild={false}>
        <Avatar
          className='rounded-lg text-primary-700 cursor-pointer'
          textSizes={'lg'}
          variant='outline'
          size='md'
          src={src || undefined}
          name={fullName}
          onClick={() => {
            navigate(href);
          }}
        />
      </Tooltip>
    </div>
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
