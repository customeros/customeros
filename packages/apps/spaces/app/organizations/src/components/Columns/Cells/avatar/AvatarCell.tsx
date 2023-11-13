import { useRouter } from 'next/navigation';

import { useLocalStorage } from 'usehooks-ts';

import { Flex } from '@ui/layout/Flex';
import { Avatar } from '@ui/media/Avatar';
import { Tooltip } from '@ui/overlay/Tooltip';

interface AvatarCellProps {
  id: string;
  name: string;
}

export const AvatarCell = ({ name, id }: AvatarCellProps) => {
  const router = useRouter();
  const [tabs] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });

  const lastPositionParams = tabs[id];
  const href = getHref(id, lastPositionParams);
  const fullName = name || 'Unnamed';

  return (
    <Flex align='center'>
      <Tooltip label={fullName} fontWeight='normal'>
        <Avatar
          variant='outline'
          size='md'
          borderRadius='lg'
          name={fullName}
          cursor='pointer'
          onClick={() => router.push(href)}
        />
      </Tooltip>
    </Flex>
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
