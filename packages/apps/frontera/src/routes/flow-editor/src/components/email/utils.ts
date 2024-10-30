import { ReactNode, createElement } from 'react';

import parse, {
  Element,
  DOMNode,
  domToReact,
  HTMLReactParserOptions,
} from 'html-react-parser';

const getParserOptions = (): HTMLReactParserOptions => ({
  replace: (domNode) => {
    if (!(domNode instanceof Element)) {
      return;
    }

    const attributes: Record<string, unknown> = {};

    if (domNode.attribs) {
      Object.entries(domNode.attribs).forEach(([key, value]) => {
        if (key === 'style') {
          attributes.style = parseStyles(value);
        } else if (key === 'class') {
          attributes.className = value;
        } else {
          attributes[key] = value;
        }
      });
    }

    return createElement(
      domNode.name,
      attributes,
      domNode.children ? domToReact(domNode.children as DOMNode[]) : null,
    );
  },
});

const parseStyles = (styleString: string): Record<string, string> => {
  const styles: Record<string, string> = {};

  styleString.split(';').forEach((style) => {
    const [property, value] = style.split(':');

    if (property && value) {
      const formattedProperty = property.trim();

      styles[formattedProperty] = value.trim();
    }
  });

  return styles;
};

export const parseHtmlToReact = (htmlContent: string): ReactNode => {
  try {
    return parse(htmlContent, getParserOptions());
  } catch (error) {
    console.error('Error parsing HTML to React:', error);

    return null;
  }
};
