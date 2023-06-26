import React from 'react';
import { SingleSelect } from '@spaces/ui/form/select/components/single-select/SingleSelect';
import { HealthIndicatorSelectMenu } from '@spaces/organization/health-select/select/HealthIndicatorSelectMenu';
import { HealthIndicatorToggleButton } from '@spaces/organization/health-select/select/HealthIndicatorToggleButton';
import { SingleSelectWrapper } from '@spaces/ui/form/select/components/single-select/SingleSelectWrapper';
import {
  useHealthIndicators,
  useSetHealthIndicator,
} from '@spaces/hooks/useHealthIndicators';
import { HealthIndicator } from '@spaces/graphql';

interface HealthIndicatorSelectProps {
  organizationId: string;
  showIcon?: boolean;
  healthIndicator: HealthIndicator | undefined | null;
}

export const HealthIndicatorSelect: React.FC<HealthIndicatorSelectProps> = ({
  showIcon,
  organizationId,
  healthIndicator,
}) => {
  const { data } = useHealthIndicators();
  const { onSetHealthIndicator, onRemoveHealthIndicator, saving } =
    useSetHealthIndicator();
  return (
    <SingleSelect<string>
      onSelect={(newValue) => {
        if (healthIndicator?.id && !newValue) {
          onRemoveHealthIndicator({ variables: { organizationId } });
          return;
        }
        onSetHealthIndicator({
          variables: { organizationId, healthIndicatorId: newValue },
        });
      }}
      value={healthIndicator?.id}
      options={data ?? []}
    >
      <SingleSelectWrapper>
        <HealthIndicatorToggleButton
          placeholder='Health'
          showIcon={showIcon}
          saving={saving}
        />
        <HealthIndicatorSelectMenu />
      </SingleSelectWrapper>
    </SingleSelect>
  );
};
