import { Props } from 'chakra-react-select';

export type SelectOption<T = string> = {
  value: T;
  label: string;
};

export type GroupedOption<T = string> = {
  label: string;
  options: SelectOption<T>[];
};

// Exhaustively typing this Props interface does not offer any benefit at this moment
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type chakraStyles = Props<any, any, any>['chakraStyles'];
