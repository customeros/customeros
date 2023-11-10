import { Flex } from '@ui/layout/Flex';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';
import { Organization } from '@graphql/types';

interface OrganizationCellProps {
  organization: Organization;
  lastPositionParams?: string;
}

export const OrganizationCell = ({
  organization,
  lastPositionParams,
}: OrganizationCellProps) => {
  const href = getHref(organization.id, lastPositionParams);
  const hasParent = !!organization.subsidiaryOf?.length;
  const fullName = organization.name || 'Unnamed';
  const parentName = organization.subsidiaryOf?.[0]?.organization.name;

  return (
    <Flex isTruncated flexDir='column'>
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
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
