import React, { useMemo } from 'react';
import { Flex } from '@ui/layout/Flex';
import { ChakraProps } from '@chakra-ui/react';
import { InteractivityProps } from '@chakra-ui/styled-system';
import ReactMarkdown from 'react-markdown';
import { getTextRendererStyles } from '@ui/theme/textRendererStyles';

interface MarkdownContentRendererProps extends InteractivityProps, ChakraProps {
  markdownContent: string;
  showAsInlineText?: boolean;
}

export const MarkdownContentRenderer: React.FC<
  MarkdownContentRendererProps
> = ({
  markdownContent,
  noOfLines,
  pointerEvents,
  showAsInlineText,
  ...rest
}) => {
  const textRendererStyles = useMemo(
    () => getTextRendererStyles(showAsInlineText),
    [showAsInlineText],
  );
  return (
    <Flex
      as={ReactMarkdown}
      flexDir='column'
      pointerEvents={pointerEvents}
      noOfLines={noOfLines}
      {...rest}
      sx={textRendererStyles}
    >
      {markdownContent}
    </Flex>
  );
};
