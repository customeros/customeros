import { useEffect } from 'react';

import { TextNode } from 'lexical';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';

export const setInlineStyle = (node: TextNode) => {
  const oldAttr = node.getStyle();
  const oldAttrs = oldAttr.split(';');
  // const newStylesMap: Set<string> = new Set<string>();
  const newStylesMap: Map<string, string> = new Map<string, string>();

  for (const a of oldAttrs) {
    const keyValue = a.split(':');
    const key = keyValue[0];
    const value = keyValue[1];

    key && value && newStylesMap.set(key.trim(), value.trim());
  }

  if (node.hasFormat('bold')) {
    newStylesMap.set('font-weight', 'bold');
  } else {
    newStylesMap.delete('font-weight');
  }

  if (node.hasFormat('italic')) {
    newStylesMap.set('font-style', 'italic');
  } else {
    newStylesMap.delete('font-style');
  }

  const hasUnderline = node.hasFormat('underline');
  const hasStrikeThrough = node.hasFormat('strikethrough');

  if (hasUnderline && hasStrikeThrough) {
    newStylesMap.set('text-decoration', 'underline line-through');
  } else if (hasUnderline) {
    newStylesMap.set('text-decoration', 'underline');
  } else if (hasStrikeThrough) {
    newStylesMap.set('text-decoration', 'line-through');
  } else {
    newStylesMap.delete('text-decoration');
  }

  const attr = Array.from(newStylesMap.entries())
    .map(([k, v]) => `${k}: ${v}`)
    .join(';');

  oldAttr != attr && node.setStyle(attr);
};

const TextNodeTransformer = () => {
  const [editor] = useLexicalComposerContext();

  useEffect(() => {
    editor.registerNodeTransform(TextNode, (node) => {
      setInlineStyle(node);
    });
  }, [editor]);

  return null;
};
export default TextNodeTransformer;
