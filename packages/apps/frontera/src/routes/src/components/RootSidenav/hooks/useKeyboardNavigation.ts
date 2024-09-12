import { useNavigate } from 'react-router-dom';

import { useSequentialShortcut } from '@shared/hooks';

interface KeyboardNavigationOptions {
  when?: boolean;
}

interface Presets {
  targetsPreset?: string;
  contactsPreset?: string;
  customersPreset?: string;
  contractsPreset?: string;
  organizationsPreset?: string;
  flowSequencesPreset?: string;
  upcomingInvoicesPreset?: string;
}

export const useKeyboardNavigation = (
  presets: Presets,
  options: KeyboardNavigationOptions = { when: true },
) => {
  const navigate = useNavigate();

  useSequentialShortcut(
    'G',
    'T',
    () => {
      if (!presets.targetsPreset) return;
      navigate(`/finder?preset=${presets.targetsPreset}`);
    },
    options,
  );
  useSequentialShortcut(
    'G',
    'O',
    () => {
      navigate(`/prospects`);
    },
    options,
  );
  useSequentialShortcut(
    'G',
    'C',
    () => {
      if (!presets.customersPreset) return;
      navigate(`/finder?preset=${presets.customersPreset}`);
    },
    options,
  );
  useSequentialShortcut(
    'G',
    'Q',
    () => {
      if (!presets.flowSequencesPreset) return;
      navigate(`/finder?preset=${presets.flowSequencesPreset}`);
    },
    options,
  );

  useSequentialShortcut(
    'G',
    'Z',
    () => {
      if (!presets.organizationsPreset) return;
      navigate(`/finder?preset=${presets.organizationsPreset}`);
    },
    options,
  );

  useSequentialShortcut(
    'G',
    'N',
    () => {
      if (!presets.contactsPreset) return;
      navigate(`/finder?preset=${presets.contactsPreset}`);
    },
    options,
  );
  useSequentialShortcut(
    'G',
    'I',
    () => {
      if (!presets.upcomingInvoicesPreset) return;
      navigate(`/finder?preset=${presets.upcomingInvoicesPreset}`);
    },
    options,
  );
  useSequentialShortcut(
    'G',
    'S',
    () => {
      navigate('/settings');
    },
    options,
  );

  useSequentialShortcut(
    'G',
    'D',
    () => {
      navigate('/customer-map');
    },
    options,
  );

  useSequentialShortcut(
    'G',
    'R',
    () => {
      if (!presets.contractsPreset) return;
      navigate(`/finder?preset=${presets.contractsPreset}`);
    },
    options,
  );
};
