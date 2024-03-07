import { useRef, useState, useEffect } from 'react';

import { useLocalStorage } from 'usehooks-ts';

import { Flex } from '@ui/layout/Flex';
import { Link } from '@ui/navigation/Link';
import { Text } from '@ui/typography/Text';
import { Tooltip } from '@ui/overlay/Tooltip';

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
  const linkRef = useRef<HTMLAnchorElement>(null);
  const [isOverflowing, setIsOverflowing] = useState(false);

  const lastPositionParams = tabs[id];
  const href = getHref(id, lastPositionParams);
  const fullName = name || 'Unnamed';
  useEffect(() => {
    const element = linkRef.current;
    if (element) {
      const isOverflow = element.scrollWidth > element.clientWidth;
      setIsOverflowing(isOverflow);
    }
  }, [linkRef]);

  return (
    <Tooltip
      hasArrow
      placement='bottom-start'
      fontWeight='normal'
      label={isOverflowing ? fullName : ''}
    >
      <Flex isTruncated flexDir='column'>
        {isSubsidiary && (
          <Text fontSize='xs' color='gray.500'>
            {parentOrganizationName}
          </Text>
        )}
        <Link
          ref={linkRef}
          href={href}
          color='gray.700'
          fontWeight='semibold'
          overflow='hidden'
          textOverflow='ellipsis'
          _hover={{ textDecoration: 'none' }}
        >
          {fullName}
        </Link>
      </Flex>
    </Tooltip>
  );
};

function getHref(id: string, lastPositionParams: string | undefined) {
  return `/organization/${id}?${lastPositionParams || 'tab=about'}`;
}
