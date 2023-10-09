import { useState } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { LinkExternal02 } from '@ui/media/icons/LinkExternal02';
import {
  getExternalUrl,
  getFormattedLink,
} from '@spaces/utils/getExternalLink';

interface WebsiteCellProps {
  website?: string | null;
}

export const WebsiteCell = ({ website }: WebsiteCellProps) => {
  const [isHovered, setIsHovered] = useState(false);

  if (!website) return <Text color='gray.400'>Unknown</Text>;

  return (
    <Flex
      align='center'
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <Text isTruncated color='gray.700' cursor='default'>
        {getFormattedLink(website)}
      </Text>
      {isHovered && (
        <IconButton
          ml='1'
          variant='ghost'
          size='xs'
          borderRadius='5px'
          onClick={() =>
            window.open(getExternalUrl(website ?? '/'), '_blank', 'noopener')
          }
          aria-label='organization website'
          icon={<LinkExternal02 color='gray.500' />}
        />
      )}
    </Flex>
  );
};
