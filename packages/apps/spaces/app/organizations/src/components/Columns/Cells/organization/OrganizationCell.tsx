import { useState } from 'react';
import { useRouter } from 'next/navigation';

import { Organization } from '@graphql/types';
import { Avatar } from '@ui/media/Avatar';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';
import { Flex } from '@ui/layout/Flex';
import { IconButton } from '@ui/form/IconButton';
import { Tooltip } from '@ui/overlay/Tooltip';
import { Icons } from '@ui/media/Icon';
import {
  getExternalUrl,
  getFormattedLink,
} from '@spaces/utils/getExternalLink';

interface OrganizationCellProps {
  organization: Organization;
  lastPositionParams?: string;
}

export const OrganizationCell = ({
  organization,
  lastPositionParams,
}: OrganizationCellProps) => {
  const router = useRouter();
  const [isHovered, setIsHovered] = useState(false);
  const href = `/organization/${organization.id}?${
    lastPositionParams || 'tab=about'
  }`;
  const hasParent = !!organization.subsidiaryOf?.length;
  const fullName = hasParent
    ? organization.subsidiaryOf[0].organization.name
    : organization.name || 'Unnamed';

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
      <Flex display='inline-block' ml='3' isTruncated>
        <Link
          href={href}
          color='gray.700'
          fontWeight='semibold'
          _hover={{ textDecoration: 'none' }}
        >
          {fullName}
        </Link>
        <br />
        {organization.website && (
          <>
            <Text
              isTruncated
              color='gray.500'
              onMouseEnter={() => setIsHovered(true)}
            >
              {getFormattedLink(organization.website)}
            </Text>
            {isHovered && (
              <Flex
                position='absolute'
                bottom='8px'
                pl='1'
                ml='-5px'
                bg='white'
                borderRadius='md'
                border='1px solid'
                zIndex='overlay'
                borderColor='gray.200'
                onMouseLeave={() => setIsHovered(false)}
              >
                <Text color='gray.500' cursor='default' lineHeight='23px'>
                  {getFormattedLink(organization?.website)}
                </Text>
                <IconButton
                  ml='1'
                  variant='ghost'
                  size='xs'
                  borderRadius='5px'
                  onClick={() =>
                    window.open(
                      getExternalUrl(organization.website ?? '/'),
                      '_blank',
                      'noopener',
                    )
                  }
                  aria-label='organization website'
                  icon={<Icons.LinkExternal2 color='gray.500' />}
                />
              </Flex>
            )}
          </>
        )}
      </Flex>
    </Flex>
  );
};
