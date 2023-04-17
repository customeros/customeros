import * as React from 'react';
import { ReactNode } from 'react';

import { Message } from './Message';

interface TranscriptElement {
  party: any;
  text: string;
  file_id?: string;
}

interface TranscriptContentProps {
  messages: Array<TranscriptElement>;
  children?: ReactNode;
  firstIndex: {
    received: number | null;
    send: number | null;
  };
  contentType?: string;
}

export const TranscriptContent: React.FC<TranscriptContentProps> = ({
  messages = [],
  children,
  firstIndex,
  contentType,
}) => {
  return (
    <>
      {messages?.map((transcriptElement: TranscriptElement, index: number) => {
        return (
          <Message key={`message-item-${index}`}
            transcriptElement={transcriptElement}
            index={index}
            contentType={contentType}
            firstIndex={firstIndex}
          >
            {children}
          </Message>
        );
      })}
    </>
  );
};
