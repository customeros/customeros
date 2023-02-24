import { ChangeEvent } from 'react';

export type Props = Omit<DOMRect, 'toJSON'>;
export type Result = {
  onFileChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
};
