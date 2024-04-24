export enum HighlightColor {
  Yellow = 'yellow',
  Rose = 'rose',
  Purple = 'purple',
  BlueDark = 'blueDark',
  Cyan = 'cyan',
  Teal = 'teal',
  Moss = 'moss',
  Warning = 'warning',
  GrayWarm = 'grayWarm',
}

export function getColorByUUID(
  uuid: string,
  excludedColors: HighlightColor[] = [],
): HighlightColor {
  if (!uuid) return HighlightColor.Yellow;
  const hash = uuid
    .split('')
    .reduce((acc, char) => Math.imul(31, acc) + char.charCodeAt(0), 0);
  const availableColors = Object.values(HighlightColor).filter(
    (color) => !excludedColors.includes(color),
  );

  if (availableColors.length === 0) {
    return HighlightColor.Yellow;
  }

  const index = Math.abs(hash) % availableColors.length;

  return availableColors[index];
}

export function getVersionFromUUID(uuid: string): number {
  if (!uuid) return 1;

  const hash = uuid
    .split('')
    .reduce((acc, char) => Math.imul(31, acc) + char.charCodeAt(0), 0);

  return (Math.abs(hash) % 3) + 1;
}
