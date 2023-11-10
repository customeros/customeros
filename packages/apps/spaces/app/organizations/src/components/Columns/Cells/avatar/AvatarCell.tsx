import { useRouter } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { Avatar } from '@ui/media/Avatar';
import { Organization } from '@graphql/types';
import { Tooltip } from '@ui/overlay/Tooltip';

interface AvatarCellProps {
  organization: Organization;
  lastPositionParams?: string;
}

export const AvatarCell = ({
  organization,
  lastPositionParams,
}: AvatarCellProps) => {
  const router = useRouter();

  const href = getHref(organization.id, lastPositionParams);
  const fullName = organization.name || 'Unnamed';

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
