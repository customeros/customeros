import { render } from '@react-email/render';
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

export const prepareEmailContent = async (
  bodyHtml: string,
): Promise<string> => {
  try {
    const emailHtml = await render(<EmailTemplate bodyHtml={bodyHtml} />, {
      pretty: true,
    });

    return emailHtml;
  } catch (error) {
    throw new Error('Unable to process email content');
  }
};
