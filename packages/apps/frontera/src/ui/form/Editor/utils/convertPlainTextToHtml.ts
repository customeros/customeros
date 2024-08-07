export function convertPlainTextToHtml(plainText: string): string {
  // Escape HTML special characters to avoid XSS vulnerabilities
  const escapedText = plainText
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');

  // Replace newline characters with <br> tags
  const htmlText = escapedText.replace(/\n/g, '<br>');

  return htmlText;
}
