import { useLocalStorage } from 'usehooks-ts';

import { Flex } from '@ui/layout/Flex';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';

interface OrganizationCellProps {
  id: string;
  name: string;
  isSubsidiary: boolean;
  parentOrganizationName: string;
}

export const OrganizationCell = ({
  id,
  name,
  isSubsidiary,
  parentOrganizationName,
}: OrganizationCellProps) => {
  const [tabs] = useLocalStorage<{
    [key: string]: string;
  }>(`customeros-player-last-position`, { root: 'organization' });

  const lastPositionParams = tabs[id];
  const href = getHref(id, lastPositionParams);
  const fullName = name || 'Unnamed';

  return (
    <Flex isTruncated flexDir='column'>
      {isSubsidiary && (
        <Text fontSize='xs' color='gray.500'>
          {parentOrganizationName}
        </Text>
      )}
      <Link
        href={href}
        color='gray.700'
        fontWeight='semibold'
        _hover={{ textDecoration: 'none' }}
      >
        {fullName}
      </Link>
    </Flex>
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
