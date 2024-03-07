import { cardAnatomy } from '@chakra-ui/anatomy';
import { createMultiStyleConfigHelpers } from '@chakra-ui/react';

const { definePartsStyle, defineMultiStyleConfig } =
  createMultiStyleConfigHelpers(cardAnatomy.keys);

// define custom styles for funky variant
const variants = {
  outlinedElevated: definePartsStyle({
    container: {
      boxShadow: 'xs',
      borderWidth: '1px',
      borderRadius: '8px',
      borderColor: 'gray.200',
      transition: 'all 0.2s ease-out',
      _hover: {
        boxShadow: 'md',
      },
    },
    body: {
      p: '4',
    },
  }),
};

// export variants in the component theme
export const Card = defineMultiStyleConfig({ variants });
