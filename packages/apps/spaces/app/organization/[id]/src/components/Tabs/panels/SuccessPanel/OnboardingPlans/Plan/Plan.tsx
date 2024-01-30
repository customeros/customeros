import { useState } from 'react';
import { useParams } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { useDisclosure } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { pulseOpacity } from '@ui/utils/keyframes';
import { Collapse } from '@ui/transitions/Collapse';
import { ChevronExpand } from '@ui/media/icons/ChevronExpand';
import { ChevronCollapse } from '@ui/media/icons/ChevronCollapse';

import { PlanMenu } from './PlanMenu';
import { PlanDueDate } from './PlanDueDate';
import { MilestoneList } from './MilestoneList';
import { PlanDatum, MilestoneDatum } from '../types';
import { AddMilestoneModal } from './AddMilestoneModal';
import { usePlanMutations } from '../../hooks/usePlanMutations';
import { useMilestoneMutations } from '../../hooks/useMilestoneMutations';

interface PlanProps {
  plan: PlanDatum;
  isOpen?: boolean;
  onToggle?: (planId: string) => void;
}

export const Plan = ({ plan, isOpen, onToggle }: PlanProps) => {
  const organizationId = useParams()?.id as string;
  const addMilestoneModal = useDisclosure({ id: `${plan.id}-add-milestone` });
  const [openMilestoneId, setOpenMilestoneId] = useState<string | null>('');

  const { updateOnboardingPlan } = usePlanMutations({
    organizationId,
  });
  const { updateMilestone, addMilestone } = useMilestoneMutations({
    plan,
  });

  const isTemporary = plan.id.startsWith('temp');
  const activeMilestones = plan.milestones.filter((m) => !m.retired);
  const existingMilestoneNames = activeMilestones.map(
    (milestone) => milestone.name,
  );

  const handleTogglePlan = () => {
    onToggle?.(plan.id);
  };

  const handleToggleMilestone = (id: string) => {
    setOpenMilestoneId((prevId) => (prevId === id ? null : id));
  };

  const handleRemovePlan = () => {
    updateOnboardingPlan.mutate({
      input: {
        id: plan.id,
        retired: true,
        organizationId,
        statusDetails: plan.statusDetails,
      },
    });
  };

  const handleUpdateMilestone = (milestone: MilestoneDatum) => {
    updateMilestone.mutate({
      input: {
        ...milestone,
        organizationId,
        organizationPlanId: plan.id,
        updatedAt: new Date().toISOString(),
      },
    });
  };

  const handleRemoveMilestone = (id: string) => {
    const foundMilestone = plan.milestones.find((m) => m.id === id);
    if (!foundMilestone) return;

    updateMilestone.mutate({
      input: {
        ...foundMilestone,
        retired: true,
        organizationId,
        organizationPlanId: plan.id,
        updatedAt: new Date().toISOString(),
      },
    });
  };

  const handleAddMilestone = (name: string, dueDate: string) => {
    addMilestone.mutate({
      input: {
        name,
        dueDate,
        organizationId,
        organizationPlanId: plan.id,
        items: [],
        optional: false,
        createdAt: new Date().toISOString(),
        order: plan.milestones.length + 1,
      },
    });
  };

  const nextDueMilestone = plan.milestones.find(
    (m) => m.statusDetails?.status === 'NOT_STARTED' && m.retired === false,
  );

  return (
    <Flex
      p='4'
      pt='3'
      w='full'
      bg='gray.50'
      flexDir='column'
      borderRadius='lg'
      border='1px solid'
      borderColor='gray.200'
      animation={
        isTemporary ? `${pulseOpacity} 0.7s alternate ease-in-out` : undefined
      }
    >
      <Flex mb='2' align='center' flexDir='column' alignItems='flex-start'>
        <Flex align='center' justify='space-between' w='full'>
          <Text fontSize='sm' fontWeight='semibold' noOfLines={1}>
            {plan.name}
          </Text>

          <Flex align='center'>
            <IconButton
              size='xs'
              variant='ghost'
              color='gray.500'
              aria-label='Toggle plan'
              onClick={handleTogglePlan}
              icon={
                isOpen ? (
                  <ChevronCollapse color='gray.400' />
                ) : (
                  <ChevronExpand color='gray.400' />
                )
              }
            />
            <PlanMenu
              id={plan.id}
              onRemovePlan={handleRemovePlan}
              onAddMilestone={addMilestoneModal.onOpen}
            />
          </Flex>
        </Flex>

        <Flex gap='1'>
          {nextDueMilestone && (
            <Text fontSize='sm' color='gray.500' fontWeight='semibold'>
              {nextDueMilestone?.name}
            </Text>
          )}
          <PlanDueDate
            isDone={plan.statusDetails.status === 'DONE'}
            value={
              plan.statusDetails.status === 'DONE'
                ? plan.statusDetails.updatedAt
                : nextDueMilestone?.dueDate
            }
          />
        </Flex>
      </Flex>

      <Collapse in={isOpen}>
        <MilestoneList
          emptyText='Empty'
          title={plan.name}
          milestones={activeMilestones}
          openMilestoneId={openMilestoneId}
          onCreateMilestone={() => {}}
          onDuplicateMilestone={() => {}}
          onSyncMilestone={handleUpdateMilestone}
          onRemoveMilestone={handleRemoveMilestone}
          onToggleMilestone={handleToggleMilestone}
          onMakeMilestoneOptional={() => {}}
        />
      </Collapse>

      <AddMilestoneModal
        masterPlanId={plan.masterPlanId}
        isOpen={addMilestoneModal.isOpen}
        onClose={addMilestoneModal.onClose}
        onAddMilestone={handleAddMilestone}
        existingMilestoneNames={existingMilestoneNames}
      />
    </Flex>
  );
};
