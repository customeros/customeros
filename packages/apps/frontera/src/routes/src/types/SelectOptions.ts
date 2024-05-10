export type SelectOption<T = string> = {
  value: T;
  label: string;
};

export type GroupedOption<T = string> = {
  label: string;
  options: SelectOption<T>[];
};
