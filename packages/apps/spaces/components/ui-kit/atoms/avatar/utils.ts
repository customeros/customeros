function hslToRgb(
  h: number,
  s: number,
  l: number,
): { r: number; g: number; b: number } {
  // Convert hue value to a value between 0 and 1
  const hue = h / 360;
  // Convert saturation and lightness values to a value between 0 and 1
  const saturation = s / 100;
  const lightness = l / 100;
  // Calculate intermediate values
  const c = (1 - Math.abs(2 * lightness - 1)) * saturation;
  const x = c * (1 - Math.abs(((hue * 6) % 2) - 1));
  const m = lightness - c / 2;
  // Calculate RGB values
  let r = 0,
    g = 0,
    b = 0;
  if (0 <= hue && hue < 1 / 6) {
    r = c;
    g = x;
  } else if (1 / 6 <= hue && hue < 2 / 6) {
    r = x;
    g = c;
  } else if (2 / 6 <= hue && hue < 3 / 6) {
    g = c;
    b = x;
  } else if (3 / 6 <= hue && hue < 4 / 6) {
    g = x;
    b = c;
  } else if (4 / 6 <= hue && hue < 5 / 6) {
    r = x;
    b = c;
  } else if (5 / 6 <= hue && hue < 1) {
    r = c;
    b = x;
  }
  // Apply the intermediate values
  r = Math.round((r + m) * 255);
  g = Math.round((g + m) * 255);
  b = Math.round((b + m) * 255);
  // Return the RGB color code as an object
  return { r, g, b };
}

function componentToHex(c: number): string {
  const hex = c.toString(16);
  return hex.length === 1 ? '0' + hex : hex;
}

export function getInitialsColor(initials: string): string {
  // Generate a unique hash code for the input string
  const hashCode = hashString(initials);
  // Calculate the hue value based on the hash code
  const hue = hashCode % 360;
  // Set the saturation and lightness values to produce a saturated color that can be used as background for white text
  const saturation = 75;
  const lightness = 50;
  // Convert the HSL color code to RGB color code
  const { r, g, b } = hslToRgb(hue, saturation, lightness);
  // Format the RGB color code as a hexadecimal string
  const colorCode = `#${componentToHex(r)}${componentToHex(g)}${componentToHex(
    b,
  )}`;
  return colorCode;
}

function hashString(str: string): number {
  let hash = 0;
  if (str.length === 0) {
    return hash;
  }
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash |= 0; // Convert to 32bit integer
  }
  return hash;
}
