import fs from 'fs';
const { format } = require('prettier');

const file = require('../../theme/colors');

const prettierConfig = JSON.parse(
  fs.readFileSync(process.cwd() + '/.prettierrc', 'utf8'),
);

const sizes = ['xxs', 'xs', 'sm', 'md', 'lg'];
const colors = Object.keys(file.colors).filter((color) => color !== 'white');
const variants = ['solid', 'outline', 'ghost'];

const genCompoundVariant = (
  size: string,
  variant: string,
  colorScheme: string,
) => {
  let iconSize = '';

  switch (size) {
    case 'xxs':
      iconSize = 'w-3 h-3';
      break;
    case 'xs':
      iconSize = 'w-4 h-4';
      break;
    case 'sm':
      iconSize = 'w-5 h-5';
      break;
    case 'md':
      iconSize = 'w-5 h-5';
      break;
    case 'lg':
      iconSize = 'w-6 h-6';
      break;
    default:
      iconSize = 'w-5 h-5';
      break;
  }

  let iconColor = '';

  switch (variant) {
    case 'solid':
      iconColor = 'text-white';
      break;
    case 'ghost':
      iconColor = `text-${colorScheme}-600`;
      break;
    case 'link':
      iconColor = `text-${colorScheme}-700`;
      break;
    case 'outline':
      iconColor = `text-${colorScheme}-600`;
      break;
    default:
      break;
  }

  return {
    size,
    variant,
    colorScheme,
    className: [iconSize, iconColor],
  };
};

interface CompoundVariant {
  size: string;
  variant: string;
  colorScheme: string;
  className: string[];
}

function generateIconVariant(
  variants: string[],
  sizes: string[],
  colors: string[],
) {
  const compoundVariants: CompoundVariant[] = [];

  sizes.forEach((size) => {
    variants.forEach((variant) => {
      colors.forEach((colorScheme) => {
        compoundVariants.push(genCompoundVariant(size, variant, colorScheme));
      });
    });
  });

  return `const iconVariant = cva('', {
  variants: {
    size: {
      ${sizes.map((size) => `"${size}": [],`).join('\n      ')}
    },
    variant: {
      ${variants.map((variant) => `${variant}: [],`).join('\n      ')}
    },
    colorScheme: {
      ${colors.map((colorScheme) => `${colorScheme}: [],`).join('\n      ')}
    },
  },
  compoundVariants: [
    ${compoundVariants
      .map(
        (variant) =>
          `{
      size: '${variant.size}',
      variant: '${variant.variant}',
      colorScheme: '${variant.colorScheme}',
      className: ${JSON.stringify(variant.className)}
    },`,
      )
      .join('\n    ')}
  ]
});`;
}

const fileContent = `
import { cva } from 'class-variance-authority';
export ${generateIconVariant(variants, sizes, colors)}
`;

const formattedContent = format(fileContent, {
  ...prettierConfig,
  parser: 'babel',
});

const filePath = process.cwd() + '/ui/form/IconButton/IconButton.variants.ts';

fs.writeFile(filePath, formattedContent, () => {});
