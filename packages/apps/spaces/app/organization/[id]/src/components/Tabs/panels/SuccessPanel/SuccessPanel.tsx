'use client';
import { PanelContainer } from './PanelContainer';
import { OnboardingStatus } from './OnboardingStatus';

export const SuccessPanel = () => {
  return (
    <PanelContainer title='Success'>
      <OnboardingStatus />
    </PanelContainer>
  );
};
