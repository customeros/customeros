const fs = require('fs');
const file = require('../../theme/colors');

const buttonTypes = ['Link', 'Solid', 'Ghost', 'Outline'];
const sizes = ['sm', 'md', 'lg', 'xl', '2xl'];
const colors = Object.keys(file.colors).filter((color) => color !== 'white');
const variants = ['solid', 'outline', 'link', 'ghost'];

const solidButton = (color: string) => `
    ${color}: [
    'text-white',
    'border',
    'border-solid',
    'bg-${color}-600',
    'hover:bg-${color}-700',
    'focus:bg-${color}-700',
    'border-${color}-600',
    'hover:border-${color}-700',
    'focus:shadow-ringPrimary',
    'focus-visible:shadow-ringPrimary',
],`;

const outlineButton = (color: string) => `
    ${color}: [
    'text-${color}-600',
    'border',
    'border-solid',
    'border-${color}-300',
    'hover:bg-${color}-50',
    'hover:text-${color}-700',
    'focus:bg-${color}-50',
    'focus:shadow-ringPrimary',
    'focus-visible:shadow-ringPrimary',
],`;

const ghostButton = (color: string) => `
    ${color}: ${
  color === 'gray'
    ? `[
      'bg-transparent',
      'shadow-none',
      'text-${color}-500',
      'hover:text-${color}-700',
      'focus:text-${color}-700',
      'hover:bg-${color}-50',
      'focus:bg-${color}-50',
    ]`
    : `[
      'bg-transparent',
      'text-${color}-500',
      'shadow-none',
      'hover:text-${color}-700',
      'focus:text-${color}-700',
      'hover:bg-${color}-50',
      'focus:bg-${color}-50',
    ]`
},`;

const linkButton = (color: string) => `
    ${color}: ${
  color === 'gray'
    ? `[
      'text-${color}-500',
      'hover:text-${color}-700',
      'focus:text-${color}-700',
      'hover:underline',
      'focus:underline',
    ]`
    : `[
      'text-${color}-700',
      'hover:text-${color}-700',
      'focus:text-${color}-700',
      'hover:underline',
      'focus:underline',
    ]`
},`;

const buttonDefaultProp = `cva([
  'inline-flex',
  'items-center',
  'justify-center',
  'whitespace-nowrap',
  'gap-2',
  'text-sm',
  'font-semibold',
  'shadow-xs',
  'outline-none',
  'transition',
  'disabled:pointer-events-none',
  'disabled:opacity-50',
],`;

const genCompoundVariant = (
  size: string,
  variant: string,
  colorScheme: string,
) => {
  let iconSize = '';
  switch (size) {
    case 'sm':
    case 'md':
    case 'lg':
    case 'xl':
      iconSize = 'w-4 h-4';
    case '2xl':
      iconSize = 'w-6 h-6';
    default:
      iconSize = 'w-4 h-4';
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


${buttonTypes
  .flatMap(
    (buttonType) => `
  export const ${buttonType.toLowerCase()}Button = ${buttonDefaultProp} {
    variants: {
      colorScheme: {
        ${colors
          .map((color) => {
            switch (buttonType) {
              case 'Solid':
                return solidButton(color);
              case 'Outline':
                return outlineButton(color);
              case 'Ghost':
                return ghostButton(color);
              case 'Link':
                return linkButton(color);
              default:
                return '';
            }
          })
          .join('')}
      },
    },
    defaultVariants: {
      colorScheme: 'primary',
    },
  })
 
`,
  )
  .join('')}
`;

const filePath = './ui/form/Button/Button.variants.ts';

fs.writeFile(filePath, fileContent, () => {});
