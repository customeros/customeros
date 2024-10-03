import { QuoteNode } from '@lexical/rich-text';
import { LexicalNode, EditorConfig, $applyNodeReplacement } from 'lexical';

export class ExtendedQuoteNode extends QuoteNode {
  static getType() {
    return 'custom-quote';
  }

  static clone(node: ExtendedQuoteNode): ExtendedQuoteNode {
    return new ExtendedQuoteNode(node.__key);
  }

  createDOM(config: EditorConfig) {
    const element = super.createDOM(config);

    element.style.borderLeft = '2px solid #D0D5DD';
    element.style.paddingLeft = '12px';

    return element;
  }
}

export function $createExtendedQuoteNode(text: string): ExtendedQuoteNode {
  return $applyNodeReplacement(new ExtendedQuoteNode(text));
}

export function $isExtendedQuoteNode(
  node: LexicalNode | null | undefined,
): node is ExtendedQuoteNode {
  return node?.getType() === 'custom-quote';
}
