import { extendTheme } from '@chakra-ui/react';

import { Input } from '../form/Input/Input.theme';
import { Textarea } from '../form/Textarea/Textarea.theme';

import { colors } from './colors';

export const theme = extendTheme({
  fonts: {
    heading: 'var(--font-barlow)',
    body: 'var(--font-barlow)',
  },
  colors,
  components: {
    Input,
    Textarea,
  },
});
