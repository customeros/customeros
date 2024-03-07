import { useSearchParams } from 'next/navigation';
import { useRef, useMemo, useState } from 'react';

import { useDeepCompareEffect } from 'rooks';

import { MilestoneDatum } from '../types';
import { useMilestonesMutations } from './useMilestonesMutations';

export const useMilestonesMethods = (options: {
  milestones: MilestoneDatum[];
}) => {
  const isBulkAction = useRef(false);
  const searchParams = useSearchParams();
  const masterPlanId = searchParams?.get('planId') ?? '';
  const [allMilestones, setAllMilestones] = useState<MilestoneDatum[]>(
    options.milestones,
  );

  const {
    createMilestone,
    updateMilestone,
    updateMilestones,
    duplicateMilestone,
  } = useMilestonesMutations();

  const defaultMilestones = useMemo(
    () => allMilestones.filter((m) => !m.optional),
    [allMilestones],
  );

  const optionalMilestones = useMemo(
    () => allMilestones.filter((m) => m.optional),
    [allMilestones],
  );

  const onSyncMilestone = (milestone: MilestoneDatum) => {
    updateMilestone.mutate({
      input: {
        masterPlanId,
        ...milestone,
      },
    });
  };

  const onMilestonesChange = (milestones: MilestoneDatum[]) => {
    setAllMilestones(milestones);
    isBulkAction.current = true;
  };

  const onCreateMilestone = (options?: { optional?: boolean }) => {
    const prevUnnamedCount = allMilestones.filter((m) =>
      m.name.startsWith('Unnamed milestone'),
    ).length;

    createMilestone.mutate({
      input: {
        masterPlanId,
        durationHours: 24,
        items: [],
        optional: options?.optional ?? false,
        order: allMilestones.length + 1,
        name: `Unnamed milestone (${prevUnnamedCount + 1})`,
      },
    });
  };

  const onCreateOptionalMilestone = () => {
    onCreateMilestone({ optional: true });
  };

  const onRemoveMilestone = (id: string) => {
    const milestone = allMilestones.find((m) => m.id === id);
    if (!milestone) return;

    onSyncMilestone({
      ...milestone,
      retired: true,
    });
  };

  const onDuplicateMilestone = (id: string) => {
    duplicateMilestone.mutate({
      id,
      masterPlanId,
    });
  };

  const onMakeMilestoneOptional = (id: string) => {
    const foundMilestone = allMilestones.find((m) => m.id === id);
    if (!foundMilestone) return;

    onSyncMilestone({
      ...foundMilestone,
      optional: true,
    });
  };

  const onMakeMilestoneDefault = (id: string) => {
    const foundMilestone = allMilestones.find((m) => m.id === id);
    if (!foundMilestone) return;

    onSyncMilestone({
      ...foundMilestone,
      optional: false,
    });
  };

  useDeepCompareEffect(() => {
    setAllMilestones(options.milestones);
  }, [options.milestones]);

  useDeepCompareEffect(() => {
    if (!isBulkAction.current) return;

    updateMilestones.mutate({
      input: allMilestones.map((m, idx) => ({
        ...m,
        order: idx,
        masterPlanId,
      })),
    });

    isBulkAction.current = false;
  }, [allMilestones?.map((m) => m.id)]);

  return {
    allMilestones,
    onSyncMilestone,
    defaultMilestones,
    onCreateMilestone,
    onRemoveMilestone,
    optionalMilestones,
    onMilestonesChange,
    onDuplicateMilestone,
    onMakeMilestoneDefault,
    onMakeMilestoneOptional,
    onCreateOptionalMilestone,
  };
};
