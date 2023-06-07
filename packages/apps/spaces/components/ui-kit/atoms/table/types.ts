import { ReactNode } from 'react';

export type Column<T = unknown> = {
  id: string;
  width: string;
  label: string | ReactNode;
  subLabel?: string;
  template: (data: T) => JSX.Element;
  isLast?: boolean;
};
