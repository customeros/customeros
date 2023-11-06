import React from 'react';

import { Tag, TagLabel } from '@ui/presentation/Tag';

interface RoleTagProps {
  label: string;
}

export const getTagColorScheme = (label: string) => {
  switch (label) {
    case 'Decision Maker':
      return 'primary';
    case 'Influencer':
      return 'green';
    case 'User':
      return 'blue';
    case 'Stakeholder':
      return 'pink';
    case 'Gatekeeper':
      return 'orange';
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
      border='1px solid'
      borderColor={`${[colorScheme]}.200`}
      backgroundColor={`${[colorScheme]}.50`}
      color={`${[colorScheme]}.700`}
      boxShadow='none'
      fontWeight='normal'
      minHeight={6}
    >
      <TagLabel>{label}</TagLabel>
    </Tag>
  );
};
