import { Body, Html, Tailwind } from '@react-email/components';

import { parseHtmlToReact } from './utils.ts';

export const EmailTemplate = ({ bodyHtml }: { bodyHtml: string }) => {
  const reactEmailContent = parseHtmlToReact(bodyHtml);

  return (
    <Html>
      <Tailwind>
        <Body
          style={{
            margin: 'auto',
            fontFamily: 'ui-sans-serif, system-ui, sans-serif',
          }}
        >
          {reactEmailContent}
        </Body>
      </Tailwind>
    </Html>
  );
};
