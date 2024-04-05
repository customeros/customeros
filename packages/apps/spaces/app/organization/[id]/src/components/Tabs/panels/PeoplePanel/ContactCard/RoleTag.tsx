import React from 'react';

import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';

interface RoleTagProps {
  label: string;
}

export const getTagColorScheme = (label: string) => {
  switch (label) {
    case 'Decision Maker':
      return 'primary';
    case 'Influencer':
      return 'greenLight';
    case 'User':
      return 'blueDark';
    case 'Stakeholder':
      return 'rose';
    case 'Gatekeeper':
      return 'warning';
    case 'Champion':
      return 'error';
    default:
      return 'gray';
  }
};

export const RoleTag: React.FC<RoleTagProps> = ({ label }) => {
  const colorScheme = (() => getTagColorScheme(label))();

  return (
    <Tag
      size='sm'
      variant='outline'
      colorScheme={colorScheme}
      className='min-h-6 '
    >
      <TagLabel>{label}</TagLabel>
    </Tag>
  );
};
