import fs from 'fs';
const { format } = require('prettier');

const file = require('../../theme/colors');
const prettierConfig = JSON.parse(
  fs.readFileSync(process.cwd() + '/.prettierrc', 'utf8'),
);
const colors = Object.keys(file.colors).filter((color) => color !== 'white');
type compoundVariants = {
  colorScheme: string;
  className: string[];
}[];

const compoundVariants: compoundVariants = [];

colors.forEach((colorScheme) => {
  const bgColor = `bg-${colorScheme}-100`;
  const ringColor = `ring-${colorScheme}-50`;
  const ringOffsetColor = `ring-offset-${colorScheme}-100`;
  const textColor = `text-${colorScheme}-600`;

  const className = [
    `${bgColor} ${ringColor} ${ringOffsetColor} ${textColor}`,
  ].filter(Boolean) as string[];

  compoundVariants.push({
    colorScheme,
    className,
  });
});

const fileContent = `
import { cva } from 'class-variance-authority';

export const featureIconVariant = cva(
  ['flex', 'justify-center', 'items-center', 'rounded-full', 'overflow-visible'],
  {
    variants: {
     
      colorScheme: {
        primary: [],
        gray: [],
        grayBlue: [],
        warm: [],
        error: [],
        rose: [],
        warning: [],
        blueDark: [],
        teal: [],
        success: [],
        moss: [],
        greenLight: [],
        violet: [],
        fuchsia: [],
        blue:[],
        yellow: [],
        orangeDark: [],
      },
    },
    compoundVariants: ${JSON.stringify(compoundVariants, null, 2)}
  },
);
`;

const formattedContent = format(fileContent, {
  ...prettierConfig,
  parser: 'babel',
});

const filePath = process.cwd() + '/ui/media/Icon/Icon.variants.ts';

fs.writeFile(filePath, formattedContent, () => {});
