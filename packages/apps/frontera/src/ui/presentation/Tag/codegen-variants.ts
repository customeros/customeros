import fs from 'fs';
const { format } = require('prettier');

const file = require('../../theme/colors');

const prettierConfig = JSON.parse(
  fs.readFileSync(process.cwd() + '/.prettierrc', 'utf8'),
);

const capitalize = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);

const sizes = ['sm', 'md', 'lg'];
const variants = ['subtle', 'solid', 'outline'];
const colors = Object.keys(file.colors).filter((color) => color !== 'white');

const subtle = (color: string) => `
${color}: [
  'text-${color}-700',
  'bg-${color}-100',
],`;

const solid = (color: string) => `
${color}: [
  'bg-${color}-500',
  'text-white',
],`;

const outline = (color: string) => `
${color}: [
  'bg-${color}-50',
  'text-${color}-700',
  'border',
  'border-solid',
  'border-${color}-200'
],`;

const variantsClasses: Record<string, (color: string) => string> = {
  subtle,
  solid,
  outline,
};

const sizeClasses: Record<string, string> = {
  sm: 'px-2 text-xs',
  md: 'px-2 text-sm',
  lg: 'px-3 text-base',
};

const tagCommonClasses = `[
  'w-fit',
  'flex',
  'items-center',
  'rounded-[4px]',
  'leading-none',
]`;

const sizeVariant = `
export const tagSizeVariant = cva('', {
  variants: {
    size: {
      ${sizes.map((size) => `${size}: '${sizeClasses[size]}',`).join('')}
    }
  },
  defaultVariants: {
    size: 'md',
  }
})
`;

const makeColorVariant = (variant: string) => `
export const tag${capitalize(variant)}Variant = cva(${tagCommonClasses}, {
  variants: {
    colorScheme: {
      ${colors.map((color) => `${variantsClasses[variant](color)}`).join('')}
    },
  },
  defaultVariants: {
    colorScheme: 'gray',
  },
});
`;

const fileContent = `
import { cva } from 'class-variance-authority';

${variants.map((variant) => makeColorVariant(variant)).join('\n')}
${sizeVariant}
`;

const formattedContent = format(fileContent, {
  ...prettierConfig,
  parser: 'babel',
});

const filePath = process.cwd() + '/ui/presentation/Tag/Tag.variants.ts';

fs.writeFile(filePath, formattedContent, () => {});
