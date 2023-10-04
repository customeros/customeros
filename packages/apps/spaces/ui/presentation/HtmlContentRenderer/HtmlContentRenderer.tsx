import React from 'react';
import parse, {
  HTMLReactParserOptions,
  domToReact,
  Element,
} from 'html-react-parser';
import linkifyHtml from 'linkify-html';

import { Flex } from '@ui/layout/Flex';
import { ImageAttachment } from './ImageAttachment';
import { ChakraProps } from '@chakra-ui/react';
import { InteractivityProps } from '@chakra-ui/styled-system';

interface HtmlContentRendererProps extends InteractivityProps, ChakraProps {
  htmlContent: string;
  showAsInlineText?: boolean;
}

export const HtmlContentRenderer: React.FC<HtmlContentRendererProps> = ({
  htmlContent,
  noOfLines,
  pointerEvents,
  showAsInlineText,
  ...rest
}) => {
  const linkifiedContent = linkifyHtml(htmlContent, {
    defaultProtocol: 'https',
    rel: 'noopener noreferrer',
  }).replace(/<\/?body>|<\/?html>|<\/?head>/g, '');

  const parseOptions: HTMLReactParserOptions = {
    replace: (domNode) => {
      if (domNode instanceof Element) {
        if (domNode.attribs && domNode.attribs.style) {
          delete domNode.attribs.style;
        }

        if (domNode.children.length === 0 && domNode.name !== 'img') {
          return <React.Fragment />;
        }

        switch (domNode.name) {
          case 'td': {
            return (
              <Flex flexDir='column' noOfLines={noOfLines}>
                {domToReact(domNode.children)}
              </Flex>
            );
          }
          case 'img': {
            return <ImageAttachment {...domNode.attribs} />;
          }
          default:
            return;
        }
      }
    },
  };
  const parsedContent = parse(linkifiedContent, parseOptions);
  return (
    <Flex
      flexDir='column'
      pointerEvents={pointerEvents}
      noOfLines={noOfLines}
      {...rest}
      sx={{
        '& ol, ul': {
          pl: showAsInlineText ? 0 : '5',
        },
        '& pre': {
          whiteSpace: 'normal',
          fontSize: '12px',
          color: 'gray.700',
          border: '1px solid',
          borderColor: 'gray.300',
          borderRadius: '4',
          p: '2',
          py: '1',
          my: '2',
        },
        '& blockquote': {
          position: 'relative',
          pl: '3',
          borderRadius: 0,
          verticalAlign: 'bottom',

          '&:before': {
            content: '""',
            position: 'absolute',
            left: 0,
            background: 'gray.300',
            width: '3px',
            height: '100%',
            borderRadius: '8px',
            bottom: 0,
          },
          '& p': {
            color: 'gray.500',
          },
          '& .customeros-tag': {
            color: 'gray.700',
            fontWeight: 'medium',

            '&:before': {
              content: '"#"',
            },
          },
        },
        ...(showAsInlineText
          ? {
              '&': {
                display: 'inline',
              },
              '& *': {
                display: 'inline',
              },
            }
          : {}),
      }}
    >
      {parsedContent}
    </Flex>
  );
};
