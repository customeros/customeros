export function convertPlainTextToHtml(plainText: string): string {
  // Escape HTML special characters to avoid XSS vulnerabilities
  const escapedText = plainText
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');

  // ensure variables are properly converted to nodes
  const textWithVariables = escapedText.replace(
    /\{\{([^}]+)\}\}/g,
    '<span data-lexical-variable>{{$1}}</span>',
  );

  // Replace newline characters with <br> tags
  return textWithVariables.replace(/\n/g, '<br>');
}
