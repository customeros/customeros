import * as Toolbar from '@radix-ui/react-toolbar';

export const ToolbarRoot = (props: Toolbar.ToolbarProps) => {
  return <Toolbar.Root {...props} />;
};

export const ToolbarButton = (props: Toolbar.ToolbarButtonProps) => {
  return <Toolbar.Button {...props} />;
};

export const ToolbarLink = (props: Toolbar.ToolbarLinkProps) => {
  return <Toolbar.Link {...props} />;
};

export const ToolbarSeparator = (props: Toolbar.ToolbarSeparatorProps) => {
  return <Toolbar.Separator {...props} />;
};

export const ToolbarToggleGroup = (props: Toolbar.ToggleGroupProps) => {
  return <Toolbar.ToggleGroup {...props} />;
};

export const ToolbarToggleItem = (props: Toolbar.ToggleGroupItemProps) => {
  return <Toolbar.ToggleItem {...props} />;
};
