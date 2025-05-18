import type { LexicalEditor } from "lexical";
import type { Binding, Provider } from "@lexical/yjs";

import { createHeadlessEditor } from "@lexical/headless";
import {
  Doc,
  XmlText,
  type YEvent,
  Transaction,
  AbstractType,
  encodeStateAsUpdate,
} from "yjs";
import {
  createBinding,
  syncLexicalUpdateToYjs,
  syncYjsChangesToLexical,
} from "@lexical/yjs";

import { nodes } from "./nodes";

console.log = () => {};
console.info = () => {};
console.error = () => {};
console.warn = () => {};

/**
 * Parses a Lexical editor state (as JSON string) and applies it to a blank Y.doc.
 * Returns the Y.doc with the applied state.
 *
 * @param lexicalState The Lexical editor state as a JSON string
 * @returns The Y.doc with the applied Lexical state and the binary state of the Y.doc
 */
export function syncLexicalToYDoc(
  id: string,
  lexicalState: string
): {
  doc: Doc;
  binary: Uint8Array;
} {
  // Create a new blank Y.doc
  const doc = new Doc();

  // Get the root XML text from the doc
  doc.get("root", XmlText);

  // Create a headless editor with minimal configuration
  const editor = createHeadlessEditor({
    nodes,
    onError: (error) => {
      console.error("Error in headless editor:", error);
    },
  });

  // Create a binding between the editor and the Y.doc
  // Create a map for docMap parameter
  const provider = createNoOpProvider();
  const docMap = new Map([[id, doc]]);
  const binding = createBinding(editor, provider, id, doc, docMap);
  const unsubscribe = registerCollaborationListeners(editor, provider, binding);

  editor.update(
    () => {
      const parsedState = editor.parseEditorState(lexicalState);
      editor.setEditorState(parsedState);
    },
    { discrete: true }
  );

  unsubscribe();

  // Get the binary representation of the Y.doc
  // currently used for debugging and comparison with what the realtime server stores in the database
  const binary = encodeStateAsUpdate(doc);

  return { doc, binary };
}

function registerCollaborationListeners(
  editor: LexicalEditor,
  provider: Provider,
  binding: Binding
): () => void {
  const unsubscribeUpdateListener = editor.registerUpdateListener(
    ({
      dirtyElements,
      dirtyLeaves,
      editorState,
      normalizedNodes,
      prevEditorState,
      tags,
    }) => {
      if (tags.has("skip-collab") === false) {
        syncLexicalUpdateToYjs(
          binding,
          provider,
          prevEditorState,
          editorState,
          dirtyElements,
          dirtyLeaves,
          normalizedNodes,
          tags
        );
      }
    }
  );

  // Use a more generic type for YEvent to avoid type incompatibilities
  const observer = (
    events: Array<YEvent<AbstractType<unknown>>>,
    transaction: Transaction
  ) => {
    if (transaction.origin !== binding) {
      syncYjsChangesToLexical(binding, provider, events as any, false);
    }
  };

  binding.root.getSharedType().observeDeep(observer);

  return () => {
    unsubscribeUpdateListener();
    binding.root.getSharedType().unobserveDeep(observer);
  };
}

function createNoOpProvider(): Provider {
  const emptyFunction = () => {};

  return {
    awareness: {
      getLocalState: () => null,
      getStates: () => new Map(),
      off: emptyFunction,
      on: emptyFunction,
      setLocalState: emptyFunction,
      setLocalStateField: emptyFunction,
    },
    connect: emptyFunction,
    disconnect: emptyFunction,
    off: emptyFunction,
    on: emptyFunction,
  };
}

/**
 * Main CLI function that processes command line arguments and outputs the result
 */
async function main() {
  try {
    // Get command line arguments
    const args = process.argv.slice(2);

    // Check if we have the required arguments
    if (args.length < 2) {
      process.stderr.write("Usage: bun run convert_lexical_to_yjs.ts <id> <lexicalState>\n");
      process.exit(1);
    }

    const id = args[0];
    // The lexical state might be passed as a file path or directly as a string
    let lexicalState = args[1];

    // If the lexical state is a file path, read the file
    if (lexicalState.startsWith("@")) {
      const filePath = lexicalState.substring(1);
      try {
        lexicalState = await Bun.file(filePath).text();
      } catch (error) {
        process.stderr.write(`Error reading file ${filePath}: ${error}\n`);
        process.exit(1);
      }
    }

    try {
      // Process the lexical state
      const result = syncLexicalToYDoc(id, lexicalState);
      // Output the binary directly to stdout
      process.stdout.write(Buffer.from(result.binary));
    } catch (err) {}
  } catch (error) {
    // Handle any errors
    process.stderr.write(
      `Error: ${error instanceof Error ? error.message : String(error)}\n`
    );
    process.exit(1);
  }
}

main();
