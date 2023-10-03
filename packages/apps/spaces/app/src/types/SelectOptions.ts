export type SelectOption<T = string> = {
  label: string;
  value: T;
};

export type GroupedOption<T = string> = {
  label: string;
  options: SelectOption<T>[];
};
