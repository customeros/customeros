import { useParams } from 'react-router-dom';

import { PanelContainer } from './PanelContainer';
import { OnboardingStatus } from './OnboardingStatus';

export const SuccessPanel = () => {
  const id = useParams()?.id as string;

  return (
    <PanelContainer title='Success'>
      <OnboardingStatus id={id} />
    </PanelContainer>
  );
};
