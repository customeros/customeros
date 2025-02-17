import { useState, ReactNode } from 'react';

import { AgentScope } from '@graphql/types';
import { Icon, IconName } from '@ui/media/Icon';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';

export const Scope = ({ scope }: { scope: AgentScope }) => {
  const [isOpen, setOpen] = useState(false);

  return (
    <>
      <Tag
        tabIndex={0}
        role={'button'}
        variant='subtle'
        color='grayModern'
        className={'cursor-pointer'}
        onClick={() => setOpen(true)}
      >
        <TagLabel className='flex items-center gap-1'>
          <Icon className='size-3' name={scopeMap[scope].icon} />
          {scopeMap[scope].label}
        </TagLabel>
      </Tag>
      <InfoDialog
        isOpen={isOpen}
        confirmButtonLabel={'Got it'}
        onClose={() => setOpen(false)}
        onConfirm={() => setOpen(false)}
        body={scopeMap[scope].description}
        label={`${scopeMap[scope].label} agents`}
      />
    </>
  );
};
type ScopeMapType = {
  label: string;
  icon: IconName;
  description: ReactNode;
};

const scopeMap: Record<AgentScope, ScopeMapType> = {
  [AgentScope.Workspace]: {
    label: 'Workspace',
    icon: 'building-05',
    description: (
      <div className='text-sm'>
        <p className='mb-4'>
          This agent type manages team-wide information available to your entire
          team.
        </p>
        <p>
          For example, when someone visits your company website, the agent
          shares the details with everyone on your team, ensuring that all
          members have the same up-to-date information.
        </p>
      </div>
    ),
  },
  [AgentScope.Personal]: {
    label: 'Personal',
    icon: 'user-01',
    description: (
      <div className='text-sm'>
        <p className='mb-4'>
          This agent type works with your personal account to manage and share
          relevant information with your team.
        </p>
        <p>
          For example, when you're in a meeting, an agent picks out the key
          points and shares the details with everyone on your team, keeping
          everyone up to date.
        </p>
      </div>
    ),
  },
};
