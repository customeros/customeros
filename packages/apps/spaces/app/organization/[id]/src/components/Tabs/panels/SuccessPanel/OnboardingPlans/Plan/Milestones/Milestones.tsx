import { VStack } from '@ui/layout/Stack';

import { Milestone } from './Milestone';
import { TaskDatum, MilestoneDatum } from '../../types';

interface MilestonesProps {
  openMilestoneId: string | null;
  onToggleMilestone: (id: string) => void;
  onRemoveMilestone: (id: string) => void;
  onSyncMilestone: (milestone: MilestoneDatum) => void;
  milestones: (MilestoneDatum & { items: TaskDatum[] })[];
}

export const Milestones = ({
  milestones,
  openMilestoneId,
  onSyncMilestone,
  onRemoveMilestone,
  onToggleMilestone,
}: MilestonesProps) => {
  return (
    <VStack mb='2' mt='3'>
      {milestones?.map((milestone, idx, arr) => (
        <Milestone
          key={milestone.id}
          milestone={milestone}
          onSync={onSyncMilestone}
          onToggle={onToggleMilestone}
          onRemove={onRemoveMilestone}
          isLast={idx === arr.length - 1}
          isOpen={openMilestoneId === milestone.id}
        />
      ))}
    </VStack>
  );
};
