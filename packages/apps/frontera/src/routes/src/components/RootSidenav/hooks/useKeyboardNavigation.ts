import { useNavigate } from 'react-router-dom';

import { useSequentialShortcut } from '@shared/hooks';

interface KeyboardNavigationOptions {
  when?: boolean;
}

interface Presets {
  leadsPreset?: string;
  targetsPreset?: string;
  churnedPreset?: string;
  customersPreset?: string;
  addressBookPreset?: string;
  myPortfolioPreset?: string;
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
    'F',
    () => {
      if (!presets.churnedPreset) return;
      navigate(`/finder?preset=${presets.churnedPreset}`);
    },
    options,
  );
  useSequentialShortcut(
    'G',
    'A',
    () => {
      if (!presets.addressBookPreset) return;
      navigate(`/finder?preset=${presets.addressBookPreset}`);
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
    'P',
    () => {
      if (!presets.myPortfolioPreset) return;
      navigate(`/finder?preset=${presets.myPortfolioPreset}`);
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
};
