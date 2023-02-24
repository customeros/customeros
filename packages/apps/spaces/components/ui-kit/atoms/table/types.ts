export type Column = {
  width: string;
  label: string;
  subLabel?: string;
  template: (data: unknown) => JSX.Element;
};
