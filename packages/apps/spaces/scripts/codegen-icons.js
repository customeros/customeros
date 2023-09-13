const { readFileSync, readdirSync, writeFileSync } = require('fs');
const { format } = require('prettier');

const makeIconComponent = (name, content) => `
import { Icon, IconProps } from "@ui/media/Icon";

export const ${name} = (props: IconProps) => (
  <Icon viewBox="0 0 24 24" {...props}>
    ${content}
  </Icon>
);
`;

const files = readdirSync(process.cwd() + '/public/icons/new');
const prettierConfig = JSON.parse(
  readFileSync(process.cwd() + '/.prettierrc', 'utf8'),
);

files.forEach((name) => {
  try {
    const file = readFileSync(
      process.cwd() + '/public/icons/new/' + name,
      'utf8',
    );
    const lines = file.split('\n');
    const svgInnerContent = lines
      .slice(1, lines.length - 2)
      .join('\n')
      .replace('stroke="black"', 'stroke="currentColor"')
      .replace('stroke-width', 'strokeWidth')
      .replace('stroke-linecap', 'strokeLinecap')
      .replace('stroke-linejoin', 'strokeLinejoin');

    const componentName = camelize(name.split('.')[0]);
    const outFileName = `${componentName}.tsx`;
    const outContent = makeIconComponent(componentName, svgInnerContent);

    const formattedOutContent = format(outContent, {
      ...prettierConfig,
      parser: 'babel',
    });

    const filePath = process.cwd() + '/ui/media/icons/' + outFileName;

    writeFileSync(filePath, formattedOutContent);
  } catch (e) {}
});

function camelize(str) {
  let arr = str.split('-');
  let capital = arr.map(
    (item) => item.charAt(0).toUpperCase() + item.slice(1).toLowerCase(),
  );
  let capitalString = capital.join('');

  return capitalString;
}
