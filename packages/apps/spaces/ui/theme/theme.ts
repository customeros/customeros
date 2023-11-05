import { extendTheme } from '@chakra-ui/react';

import { Input } from '@ui/form/Input/Input.theme';
import { Button } from '@ui/form/Button/Button.theme';
import { Avatar } from '@ui/media/Avatar/Avatar.theme';
import { Checkbox } from '@ui/form/Checkbox/Checkbox.theme';
import { Textarea } from '@ui/form/Textarea/Textarea.theme';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.theme';
import { NumberInput } from '@ui/form/NumberInput/NumberInput.theme';

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
    Avatar,
    Button,
    Input,
    NumberInput,
    Textarea,
    Checkbox,
    Tooltip,
  },
  styles: {
    global: () => ({
      // Optionally set global CSS styles
      body: {
        '--chakra-colors-chakra-body-text': colors.gray['700'], // no idea how to change this variable
      },
    }),
  },
});
