import { useRouter } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { Avatar } from '@ui/media/Avatar';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';
import { Organization } from '@graphql/types';
import { Tooltip } from '@ui/overlay/Tooltip';

interface OrganizationCellProps {
  organization: Organization;
  lastPositionParams?: string;
}

export const OrganizationCell = ({
  organization,
  lastPositionParams,
}: OrganizationCellProps) => {
  const router = useRouter();

  const href = getHref(organization.id, lastPositionParams);
  const hasParent = !!organization.subsidiaryOf?.length;
  const fullName = organization.name || 'Unnamed';
  const parentName = organization.subsidiaryOf?.[0]?.organization.name;

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
      <Flex ml='3' isTruncated flexDir='column'>
        {hasParent && (
          <Text fontSize='xs' color='gray.500'>
            {parentName}
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
    </Flex>
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
