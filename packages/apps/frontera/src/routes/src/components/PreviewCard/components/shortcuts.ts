type ModifierIcon = 'command' | 'arrow-block-up';
type KeyIcon = 'arrow-block-up' | 'arrow-block-down' | 'delete';

interface ShortcutWithSection extends BaseShortcut {
  section: string;
}
interface BaseShortcut {
  key?: string;
  label: string;
  keyIcon?: KeyIcon;
  description?: string;
  modifierIcon?: ModifierIcon;
  sequence?: [string, string];
}

type Shortcut = BaseShortcut;

export type Section = {
  title: string;
  shortcuts: Shortcut[];
};

interface ShortcutData {
  sections: Section[];
}

export type SearchResult = {
  item: ShortcutWithSection;
};

export const data: ShortcutData = {
  sections: [
    {
      title: 'General',
      shortcuts: [
        {
          label: 'Open command menu',
          key: 'K',
          modifierIcon: 'command',
        },
        {
          label: 'Add or Search companies',
          key: '/',
        },
        {
          label: 'Cancel',
          key: 'Esc',
        },
        {
          label: 'Move up',
          keyIcon: 'arrow-block-up',
        },
        {
          label: 'Move down',
          keyIcon: 'arrow-block-down',
        },
      ],
    },
    {
      title: 'Navigation',
      shortcuts: [
        {
          label: 'Go to Companies',
          sequence: ['G', 'C'],
        },
        {
          label: 'Go to Contacts',
          sequence: ['G', 'N'],
          description: 'View and manage all contacts',
        },
        {
          label: 'Go to Opportunities',
          sequence: ['G', 'O'],
        },
        {
          label: 'Go to Invoices',
          sequence: ['G', 'I'],
        },
        {
          label: 'Go to Contracts',
          sequence: ['G', 'R'],
        },
        {
          label: 'Go to Agents',
          sequence: ['G', 'A'],
        },
        {
          label: 'Go to Settings',
          sequence: ['G', 'S'],
        },
        {
          label: 'Go to Workspace switcher',
          sequence: ['G', 'W'],
        },
        {
          label: 'Go to Customer map',
          sequence: ['G', 'D'],
        },
      ],
    },
    {
      title: 'Companies',
      shortcuts: [
        {
          label: 'Add contacts',
          key: 'C',
        },
        {
          label: 'Change or add tags',
          key: 'T',
          modifierIcon: 'arrow-block-up',
        },
        {
          label: 'Assign owner',
          key: 'O',
          modifierIcon: 'arrow-block-up',
        },
        {
          label: 'Archive company',
          keyIcon: 'delete',
          modifierIcon: 'command',
        },
        {
          label: 'Preview company',
          key: 'Space',
        },
      ],
    },
    {
      title: 'Contacts',
      shortcuts: [
        {
          label: 'Change or add tags',
          key: 'T',
          modifierIcon: 'arrow-block-up',
        },
        {
          label: 'Add to flow',
          key: 'A',
          modifierIcon: 'arrow-block-up',
        },
        {
          label: 'Add email',
          key: 'E',
          modifierIcon: 'arrow-block-up',
        },
        {
          label: 'Edit name',
          key: 'R',
          modifierIcon: 'arrow-block-up',
        },
        {
          label: 'Archive contact',
          keyIcon: 'delete',
          modifierIcon: 'arrow-block-up',
        },
        {
          label: 'Preview contact',
          key: 'Space',
        },
      ],
    },
    {
      title: 'Opportunities',
      shortcuts: [
        {
          label: 'Change stage',
          key: 'S',
          modifierIcon: 'arrow-block-up',
        },
        {
          label: 'Rename opportunity',
          key: 'R',
          modifierIcon: 'arrow-block-up',
        },
        {
          label: 'Assign owner',
          key: 'O',
          modifierIcon: 'arrow-block-up',
        },
        {
          label: 'Archive opportunity',
          keyIcon: 'delete',
          modifierIcon: 'arrow-block-up',
        },
      ],
    },
  ],
};

export const getAllShortcuts = (sections: Section[]): Shortcut[] => {
  return sections.flatMap((section: Section) =>
    section.shortcuts.map((shortcut: Shortcut) => ({
      ...shortcut,
      section: section.title,
    })),
  );
};
