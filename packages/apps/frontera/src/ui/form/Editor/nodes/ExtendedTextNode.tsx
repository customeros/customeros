import {
  NodeKey,
  TextNode,
  $isTextNode,
  LexicalNode,
  DOMConversion,
  DOMConversionMap,
  SerializedTextNode,
  DOMConversionOutput,
  $applyNodeReplacement,
} from 'lexical';

export class ExtendedTextNode extends TextNode {
  constructor(text: string, key?: NodeKey) {
    super(text, key);
  }

  static getType(): string {
    return 'extended-text';
  }

  static clone(node: ExtendedTextNode): ExtendedTextNode {
    return new ExtendedTextNode(node.__text, node.__key);
  }

  static importDOM(): DOMConversionMap | null {
    const importers = TextNode.importDOM();

    return {
      ...importers,
      code: () => ({
        conversion: patchStyleConversion(importers?.code),
        priority: 1,
      }),
      em: () => ({
        conversion: patchStyleConversion(importers?.em),
        priority: 1,
      }),
      span: () => ({
        conversion: patchStyleConversion(importers?.span),
        priority: 1,
      }),
      strong: () => ({
        conversion: patchStyleConversion(importers?.strong),
        priority: 1,
      }),
      sub: () => ({
        conversion: patchStyleConversion(importers?.sub),
        priority: 1,
      }),
      sup: () => ({
        conversion: patchStyleConversion(importers?.sup),
        priority: 1,
      }),
    };
  }

  static importJSON(serializedNode: SerializedTextNode): TextNode {
    return TextNode.importJSON(serializedNode);
  }

  isSimpleText() {
    return this.__type === 'extended-text' && this.__mode === 0;
  }

  exportJSON(): SerializedTextNode {
    return {
      ...super.exportJSON(),
      type: 'extended-text',
      version: 1,
    };
  }
}

export function $createExtendedTextNode(text: string): ExtendedTextNode {
  return $applyNodeReplacement(new ExtendedTextNode(text));
}

export function $isExtendedTextNode(
  node: LexicalNode | null | undefined,
): node is ExtendedTextNode {
  return node instanceof ExtendedTextNode;
}

function patchStyleConversion(
  originalDOMConverter?: (node: HTMLElement) => DOMConversion | null,
): (node: HTMLElement) => DOMConversionOutput | null {
  return (node) => {
    const original = originalDOMConverter?.(node);

    if (!original) {
      return null;
    }
    const originalOutput = original.conversion(node);

    if (!originalOutput) {
      return originalOutput;
    }

    const backgroundColor = node.style.backgroundColor;
    const color = node.style.color;
    const fontFamily = node.style.fontFamily;
    const fontWeight = node.style.fontWeight;
    const fontSize = node.style.fontSize;
    const textDecoration = node.style.textDecoration;

    return {
      ...originalOutput,
      forChild: (lexicalNode, parent) => {
        const originalForChild = originalOutput?.forChild ?? ((x) => x);
        const result = originalForChild(lexicalNode, parent);

        if ($isTextNode(result)) {
          const style = [
            backgroundColor ? `background-color: ${backgroundColor}` : null,
            color ? `color: ${color}` : null,
            fontFamily ? `font-family: ${fontFamily}` : null,
            fontWeight ? `font-weight: ${fontWeight}` : null,
            fontSize ? `font-size: ${fontSize}` : null,
            textDecoration ? `text-decoration: ${textDecoration}` : null,
          ]
            .filter((value) => value != null)
            .join('; ');

          if (style.length) {
            return result.setStyle(style);
          }
        }

        return result;
      },
    };
  };
}
