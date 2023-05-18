import { ReactNode, FC } from 'react';

export type Column = {
  id: string;
  width: string;
  label: string | ReactNode;
  subLabel?: string;
  template: (data: unknown) => JSX.Element;
  isLast?: boolean;
};
