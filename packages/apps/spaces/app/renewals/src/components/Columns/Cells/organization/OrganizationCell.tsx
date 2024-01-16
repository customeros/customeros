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
  const href = `/organization/${id}?tab=account`;
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
