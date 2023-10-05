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
import sanitizeHtml from 'sanitize-html';

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
    linkifyHtml(htmlContent, {
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
          '& .customeros-mention': {
            color: 'gray.700',
            fontWeight: 'medium',

            '&:before': {
              content: '"@"',
            },
          },
        },
        "[aria-hidden='true']": {
          display: 'none',
        },

        // code to nicely present google meeting email notifications
        '& h2.primary-text': {
          color: 'gray.700',
          fontWeight: 'medium',
        },
        '& a.primary-button-text': {
          paddingY: 1,
          paddingX: 2,
          mb: 2,
          border: '1px solid',
          borderColor: 'primary.200',
          color: 'primary.700',
          background: 'primary.50',
          borderRadius: 'lg',
          width: 'fit-content',
          '&:hover': {
            textDecoration: 'none',
            bg: `primary.100`,
            color: `primary.700`,
            borderColor: `primary.200`,
          },
          '&:focus-visible': {
            textDecoration: 'none',
            bg: `primary.100`,
            color: `primary.700`,
            borderColor: `primary.200`,
            boxShadow: `0 0 0 4px var(--chakra-colors-primary-100)`,
          },
          '&:active': {
            textDecoration: 'none',
            bg: `primary.100`,
            color: `primary.700`,
            borderColor: `primary.200`,
          },
        },
        '& .body-container': {
          mt: 3,
          padding: 4,
          display: 'block',
          border: '1px solid',
          borderColor: 'gray.300',
          borderRadius: 'md',

          '& tr': {
            mr: 2,
            display: 'flex',
          },
          '& tbody': {
            mb: 2,
          },

          '& table': {
            marginInlineStart: 0,
            marginInlineEnd: 0,
          },
        },
        '& .main-column-table-ltr': {
          my: 3,
        },
        '& .grey-button-text:not(a)': {
          color: 'gray.700',
          width: 'fit-content',
          fontWeight: 'medium',
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
