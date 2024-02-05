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
import {
  OnboardingPlanStatus,
  OnboardingPlanMilestoneStatus,
} from '@shared/types/__generated__/graphql.types';

import { PlanMenu } from './PlanMenu';
import { PlanDueDate } from './PlanDueDate';
import { MilestoneList } from './MilestoneList';
import { AddMilestoneModal } from './AddMilestoneModal';
import { ProgressCompletion } from './ProgressCompletion';
import { usePlanMutations } from '../../hooks/usePlanMutations';
import { PlanDatum, MilestoneDatum, NewMilestoneInput } from '../types';
import { useMilestoneMutations } from '../../hooks/useMilestoneMutations';

interface PlanProps {
  plan: PlanDatum;
  isOpen?: boolean;
  onToggle?: (planId: string) => void;
}

export const Plan = ({ plan, isOpen, onToggle }: PlanProps) => {
  const organizationId = useParams()?.id as string;
  const [isHovered, setIsHovered] = useState(false);
  const planMenu = useDisclosure({ id: `${plan.id}-menu` });
  const addMilestoneModal = useDisclosure({ id: `${plan.id}-add-milestone` });
  const [openMilestoneId, setOpenMilestoneId] = useState<string | null>('');

  const { updateOnboardingPlan } = usePlanMutations({
    organizationId,
  });
  const { updateMilestone, addMilestone } = useMilestoneMutations({
    plan,
  });

  const isTemporary = plan.id.startsWith('temp');
  const activeMilestones = plan.milestones
    .filter((m) => !m.retired)
    .sort((a, b) => a.order - b.order);

  const hasMilestones = activeMilestones.length > 0;
  const hasOneMilestone = activeMilestones.length === 1;
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

  const handleAddMilestone = (input: NewMilestoneInput) => {
    addMilestone.mutate({
      input: {
        ...input,
        organizationId,
        adhoc: false,
        optional: false,
        organizationPlanId: plan.id,
        createdAt: new Date().toISOString(),
      },
    });
  };

  const nextDueMilestone = plan.milestones.find(
    (m) =>
      [
        OnboardingPlanMilestoneStatus.Started,
        OnboardingPlanMilestoneStatus.StartedLate,
        OnboardingPlanMilestoneStatus.NotStarted,
        OnboardingPlanMilestoneStatus.NotStartedLate,
      ].includes(m.statusDetails?.status) && m.retired === false,
  );

  const isPlanDone = [
    OnboardingPlanStatus.Done,
    OnboardingPlanStatus.DoneLate,
  ].includes(plan.statusDetails.status);

  return (
    <Flex
      px='3'
      pb='2'
      pt='3'
      w='full'
      bg='gray.50'
      id='plm'
      flexDir='column'
      borderRadius='lg'
      border='1px solid'
      borderColor='gray.200'
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => !isOpen && setIsHovered(false)}
      animation={
        isTemporary ? `${pulseOpacity} 0.7s alternate ease-in-out` : undefined
      }
    >
      <Flex
        mx='1'
        align='center'
        flexDir='column'
        alignItems='flex-start'
        // mb={isPlanDone ? '0' : '3'}
      >
        <Flex align='center' justify='space-between' w='full'>
          <Text fontSize='sm' fontWeight='semibold' noOfLines={1}>
            {plan.name}
          </Text>

          <Flex
            align='center'
            opacity={isHovered || planMenu.isOpen ? '1' : '0'}
            transition='opacity 0.2s ease-out'
          >
            {((hasMilestones && !hasOneMilestone) || isPlanDone) && (
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
            )}
            <PlanMenu
              id={plan.masterPlanId}
              isOpen={planMenu.isOpen}
              onOpen={planMenu.onOpen}
              onClose={planMenu.onClose}
              onRemovePlan={handleRemovePlan}
              onAddMilestone={addMilestoneModal.onOpen}
            />
          </Flex>
        </Flex>

        <Flex gap='1'>
          {hasMilestones ? (
            <>
              {nextDueMilestone && (
                <Text fontSize='sm' color='gray.500' fontWeight='semibold'>
                  {nextDueMilestone?.name}
                </Text>
              )}
              <PlanDueDate
                isDone={isPlanDone}
                value={
                  isPlanDone
                    ? plan.statusDetails.updatedAt
                    : nextDueMilestone?.dueDate
                }
              />
            </>
          ) : (
            <Text fontSize='sm' color='gray.500'>
              No milestones added yet
            </Text>
          )}
        </Flex>
      </Flex>

      <Collapse in={(hasOneMilestone && !isPlanDone) || isOpen}>
        <MilestoneList
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

      <Collapse in={!isOpen}>
        <ProgressCompletion plan={plan} />
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
