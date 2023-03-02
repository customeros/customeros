import { ReactNode } from 'react';

export type Column = {
  width: string;
  label: string | ReactNode;
  subLabel?: string;
  template: (data: unknown) => JSX.Element;
};
