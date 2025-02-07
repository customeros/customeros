import { useNavigate } from 'react-router-dom';

import { useStore } from '@shared/hooks/useStore';
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
  isPlatformOwner: boolean,
) => {
  const navigate = useNavigate();
  const store = useStore();

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
    'F',
    () => {
      if (!presets.flowSequencesPreset) return;
      navigate(`/finder?preset=${presets.flowSequencesPreset}`);
    },
    options,
  );

  useSequentialShortcut(
    'G',
    'C',
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
    'W',
    () => {
      store.ui.commandMenu.setType('SwitchWorkspace');
      store.ui.commandMenu.setOpen(true);
    },
    {
      when: options.when && isPlatformOwner,
    },
  );

  // useSequentialShortcut(
  //   'G',
  //   'D',
  //   () => {
  //     navigate('/customer-map');
  //   },
  //   options,
  // );

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
