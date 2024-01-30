import { VStack } from '@ui/layout/Stack';

import { Milestone } from './Milestone';
import { MilestoneDatum } from '../../types';

interface MilestoneListProps {
  title: string;
  emptyText: string;
  milestones: MilestoneDatum[];
  onCreateMilestone: () => void;
  openMilestoneId: string | null;
  onToggleMilestone: (id: string) => void;
  onRemoveMilestone: (id: string) => void;
  onDuplicateMilestone: (id: string) => void;
  onMakeMilestoneOptional?: (id: string) => void;
  onSyncMilestone: (milestone: MilestoneDatum) => void;
}

export const MilestoneList = ({
  title,
  emptyText,
  milestones,
  openMilestoneId,
  onSyncMilestone,
  onCreateMilestone,
  onRemoveMilestone,
  onToggleMilestone,
  onDuplicateMilestone,
  onMakeMilestoneOptional,
}: MilestoneListProps) => {
  return (
    <VStack>
      {milestones?.map((milestone, idx, arr) => (
        <Milestone
          key={milestone.id}
          milestone={milestone}
          onSync={onSyncMilestone}
          onToggle={onToggleMilestone}
          onRemove={onRemoveMilestone}
          isLast={idx === arr.length - 1}
          onDuplicate={onDuplicateMilestone}
          onMakeOptional={onMakeMilestoneOptional}
          isOpen={openMilestoneId === milestone.id}
        />
      ))}
    </VStack>
  );
};
