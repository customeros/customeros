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
  }>(`customeros-player-last-position`, { root: 'renewals' });

  const lastPositionParams = tabs[id];
  const href = getHref(id, lastPositionParams);
  const fullName = name || 'Unnamed';

  return (
    <div className='flex items-center ml-[1px]'>
      <Tooltip
        side='bottom'
        align='center'
        asChild={false}
        label={fullName}
        className='font-normal'
      >
        <Avatar
          size='xs'
          textSize='xs'
          tabIndex={-1}
          name={fullName}
          src={src || undefined}
          variant='outlineSquare'
          onClick={() => {
            navigate(href);
          }}
          className='text-gray-700 cursor-pointer focus:outline-none'
        />
      </Tooltip>
    </div>
  );
};

function getHref(id: string, _lastPositionParams: string | undefined) {
  return `/organization/${id}?tab=account`;
}
