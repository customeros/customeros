import React from 'react';

import linkifyHtml from 'linkify-html';
import sanitizeHtml from 'sanitize-html';
import { InteractivityProps } from '@chakra-ui/styled-system';
import parse, {
  Element,
  domToReact,
  HTMLReactParserOptions,
} from 'html-react-parser';

import { cn } from '@ui/utils/cn';

import { ImageAttachment } from './ImageAttachment';

interface HtmlContentRendererProps
  extends InteractivityProps,
    React.HTMLAttributes<HTMLDivElement> {
  noOfLines?: number;
  className?: string;
  htmlContent: string;
  showAsInlineText?: boolean;
}

export const HtmlContentRenderer: React.FC<HtmlContentRendererProps> = ({
  htmlContent,
  noOfLines,
  className,
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
              <div
                className='flex flex-col'
                style={{
                  lineClamp: `${noOfLines}`,
                  WebkitLineClamp: `${noOfLines}`,
                }}
              >
                {children}
              </div>
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

  const textRendererClass = showAsInlineText
    ? 'inline-text-renderer'
    : 'block-text-renderer';

  return (
    <div
      className={cn(textRendererClass, className, 'flex flex-col')}
      style={{
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        pointerEvents: pointerEvents as any,
        WebkitLineClamp: `${noOfLines}`,
      }}
      {...rest}
    >
      {parsedContent}
    </div>
  );
};
