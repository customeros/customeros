import { useRouter } from 'next/navigation';

import { Organization } from '@graphql/types';
import { Avatar } from '@ui/media/Avatar';
import { Link } from '@ui/navigation/Link';
import { Flex } from '@ui/layout/Flex';
import { Tooltip } from '@ui/overlay/Tooltip';
import {
  getExternalUrl,
  getFormattedLink,
} from '@spaces/utils/getExternalLink';

interface OrganizationTableCellProps {
  organization: Organization;
}

export const OrganizationTableCell = ({
  organization,
}: OrganizationTableCellProps) => {
  const router = useRouter();

  const href = `/organizations/${organization.id}?tab=about`;
  const hasParent = !!organization.subsidiaryOf?.length;
  const fullName = hasParent
    ? organization.subsidiaryOf[0].organization.name
    : organization.name || 'Unnamed';

  return (
    <Flex align='center'>
      <Tooltip label={fullName} fontWeight='normal'>
        <Avatar
          name={fullName}
          cursor='pointer'
          onClick={() => router.push(href)}
        />
      </Tooltip>
      <Flex flexDir='column' ml='2'>
        <Link
          href={href}
          color='gray.700'
          fontWeight='semibold'
          _hover={{ textDecoration: 'none' }}
        >
          {fullName}
        </Link>
        {organization.website && (
          <Link
            target='_blank'
            rel='noopener noreferrer'
            color='gray.500'
            href={getExternalUrl(organization.website)}
            transition='color 0.2s ease-in-out'
            _hover={{ textDecoration: 'none', color: 'gray.700' }}
          >
            {getFormattedLink(organization.website)}
          </Link>
        )}
      </Flex>
    </Flex>
  );
};
