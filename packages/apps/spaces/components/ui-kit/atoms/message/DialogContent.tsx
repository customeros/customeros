import * as React from 'react';
import styles from './message.module.scss';
import sanitizeHtml from 'sanitize-html';
import linkifyHtml from 'linkify-html';

interface Content {
  type?: string;
  mimetype: string;
  body: string;
}

interface DialogContentProps {
  dialog: Content;
}

export const DialogContent: React.FC<DialogContentProps> = ({ dialog }) => {
  if (dialog.mimetype === 'text/plain') {
    return (
      <>
        {linkifyHtml(dialog.body, {
          defaultProtocol: 'https',
          rel: 'noopener noreferrer',
        })}
      </>
    );
  }
  return dialog.mimetype === 'text/html' ? (
    <div
      className={`text-overflow-ellipsis ${styles.emailContent}`}
      dangerouslySetInnerHTML={{
        __html: sanitizeHtml(
          linkifyHtml(dialog.body, {
            defaultProtocol: 'https',
            rel: 'noopener noreferrer',
          }),
        ),
      }}
    ></div>
  ) : (
    <>Error</>
  );
};
