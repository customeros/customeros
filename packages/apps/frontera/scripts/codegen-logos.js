const { format } = require('prettier');
const { readFileSync, readdirSync, writeFileSync } = require('fs');

const makeIconComponent = (name, content, viewBox) => `
import React from 'react';

import { twMerge } from 'tailwind-merge';

interface IconProps extends React.SVGAttributes<SVGElement> {
  className?: string;
}


export const ${name} = ({ className, ...props }: IconProps) => (
  <svg viewBox='${viewBox}' fill='none' {...props} className={twMerge('inline-block size-4', className)}>
    ${content}
  </svg>
);
`;
const files = readdirSync(process.cwd() + '/public/icons/logos');
const prettierConfig = JSON.parse(
  readFileSync(process.cwd() + '/.prettierrc', 'utf8'),
);

function getSvgViewBox(svgString) {
  const match = svgString.match(/viewBox="([^"]*)"/);

  return match ? match[1] : '0 0 24 24'; // return matched viewBox value or null if not found
}

files.forEach((name) => {
  try {
    const file = readFileSync(
      process.cwd() + '/public/icons/logos/' + name,
      'utf8',
    );
    const lines = file.split('\n');
    const svgInnerContent = lines
      .slice(1, lines.length - 2)
      .join('\n')
      .replaceAll('stroke-width', 'strokeWidth')
      .replaceAll('stroke-linecap', 'strokeLinecap')
      .replaceAll('stroke-linejoin', 'strokeLinejoin')
      .replaceAll('fill-rule', 'fillRule')
      .replaceAll('stop-color', 'stopColor')
      .replaceAll('clip-path', 'clipPath')
      .replaceAll('clip-rule', 'clipRule')
      .replaceAll('stop-opacity', 'stopOpacity');

    const componentName = camelize(name.split('.')[0]);
    const outFileName = `${componentName}.tsx`;

    const viewBox = getSvgViewBox(file);

    const outContent = makeIconComponent(
      componentName,
      svgInnerContent,
      viewBox,
    );

    const formattedOutContent = format(outContent, {
      ...prettierConfig,
      parser: 'babel',
    });

    const filePath = process.cwd() + '/src/ui/media/logos/' + outFileName;

    writeFileSync(filePath, formattedOutContent);
  } catch (e) {
    // handle error
  }
});

function camelize(str) {
  let arr = str.split('-');
  let capital = arr.map(
    (item) => item.charAt(0).toUpperCase() + item.slice(1).toLowerCase(),
  );
  let capitalString = capital.join('');

  return capitalString;
}
