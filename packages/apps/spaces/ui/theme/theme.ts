import { extendTheme } from '@chakra-ui/react';

import { Input } from '@ui/form/Input/Input.theme';
import { Textarea } from '@ui/form/Textarea/Textarea.theme';
import { Checkbox } from '@ui/form/Checkbox/Checkbox.theme';

import { colors } from './colors';
import { shadows } from './shadows';

export const theme = extendTheme({
  fonts: {
    heading: 'var(--font-barlow)',
    body: 'var(--font-barlow)',
  },
  colors,
  shadows,
  components: {
    Input,
    Textarea,
    Checkbox,
  },
});
