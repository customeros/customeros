import React, { useMemo } from 'react';

import linkifyHtml from 'linkify-html';
import sanitizeHtml from 'sanitize-html';
import { ChakraProps } from '@chakra-ui/react';
import { InteractivityProps } from '@chakra-ui/styled-system';
import parse, {
  Element,
  domToReact,
  HTMLReactParserOptions,
} from 'html-react-parser';

import { Flex } from '@ui/layout/Flex';
import { getTextRendererStyles } from '@ui/theme/textRendererStyles';

import { ImageAttachment } from './ImageAttachment';

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
  const linkifiedContent = sanitizeHtml(
    linkifyHtml(htmlContent.replace(/&zwnj;/g, ''), {
      defaultProtocol: 'https',
      rel: 'noopener noreferrer',
    }),
    {
      ...sanitizeHtml.defaults,
      allowedAttributes: {
        a: ['href', 'name', 'target'],
        // We don't currently allow img itself by default, but
        // these attributes would make sense if we did.
        img: ['src', 'srcset', 'alt', 'title', 'width', 'height', 'loading'],
        '*': ['class', 'aria-hidden'],
      },
      allowedClasses: {
        '*': ['*'],
      },
    },
  );

  const parseOptions: HTMLReactParserOptions = {
    replace: (domNode) => {
      if (domNode instanceof Element) {
        if (domNode.tagName === 'style') {
          return <React.Fragment />;
        }

        if (domNode.attribs && domNode.attribs.style) {
          delete domNode.attribs.style;
        }

        if (domNode.children.length === 0 && domNode.name !== 'img') {
          return <React.Fragment />;
        }

        let newAttribs = {};
        if (domNode.attribs) {
          newAttribs = Object.keys(domNode.attribs).reduce(
            (result: Record<string, string>, key) => {
              if (key !== 'style') {
                result[key] = domNode.attribs[key];
              }

              return result;
            },
            {},
          );
        }
        const children = domToReact(domNode.children, parseOptions);
        switch (domNode.name) {
          case 'td': {
            return (
              <Flex flexDir='column' noOfLines={noOfLines}>
                {children}
              </Flex>
            );
          }
          case 'img': {
            return <ImageAttachment {...domNode.attribs} />;
          }
          default:
            return React.createElement(domNode.name, newAttribs, children);
        }
      }
    },
  };
  const textRendererStyles = useMemo(
    () => getTextRendererStyles(showAsInlineText),
    [showAsInlineText],
  );
  const parsedContent = parse(linkifiedContent, parseOptions);

  return (
    <Flex
      flexDir='column'
      pointerEvents={pointerEvents}
      noOfLines={noOfLines}
      {...rest}
      sx={textRendererStyles}
    >
      {parsedContent}
    </Flex>
  );
};
