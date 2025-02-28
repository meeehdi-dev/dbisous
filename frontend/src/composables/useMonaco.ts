import { editor, languages } from "monaco-editor";
import { getSQLiteKeywords } from "@/utils/keywords";
import { ref, Ref } from "vue";
import { createSharedComposable } from "@vueuse/core";

const keywords = getSQLiteKeywords();

interface Range {
  startLineNumber: number;
  endLineNumber: number;
  startColumn: number;
  endColumn: number;
}

interface Completion {
  kind: languages.CompletionItemKind;
  label: string;
  insertText: string;
  range: Range;
}

export const useCompletions = createSharedComposable(() => {
  const dispose = ref<() => void>();

  function register(columns: Record<string, Record<string, string[]>>) {
    if (dispose.value !== undefined) {
      dispose.value();
    }

    const disposable = languages.registerCompletionItemProvider("sql", {
      provideCompletionItems: function (model, position) {
        const word = model.getWordUntilPosition(position);
        const range: Range = {
          startLineNumber: position.lineNumber,
          endLineNumber: position.lineNumber,
          startColumn: word.startColumn,
          endColumn: word.endColumn,
        };

        const schemasCompletions: Completion[] = [];
        const tablesCompletions: Completion[] = [];
        const schemaTablesCompletions: Completion[] = [];
        const columnsCompletions: Completion[] = [];
        const tableColumnsCompletions: Completion[] = [];
        const schemaTableColumnsCompletions: Completion[] = [];

        Object.entries(columns).map(([schema, tables]) => {
          schemasCompletions.push({
            kind: languages.CompletionItemKind.File,
            label: schema,
            insertText: schema,
            range,
          });

          Object.entries(tables).map(([table, columns]) => {
            tablesCompletions.push({
              kind: languages.CompletionItemKind.File,
              label: table,
              insertText: table,
              range,
            });
            schemaTablesCompletions.push({
              kind: languages.CompletionItemKind.File,
              label: `${schema}.${table}`,
              insertText: `${schema}.${table}`,
              range,
            });

            columns.map((column) => {
              columnsCompletions.push({
                kind: languages.CompletionItemKind.Field,
                label: column,
                insertText: column,
                range,
              });
              tableColumnsCompletions.push({
                kind: languages.CompletionItemKind.Field,
                label: `${table}.${column}`,
                insertText: `${table}.${column}`,
                range,
              });
              schemaTableColumnsCompletions.push({
                kind: languages.CompletionItemKind.Field,
                label: `${schema}.${table}.${column}`,
                insertText: `${schema}.${table}.${column}`,
                range,
              });
            });
          });
        });

        return {
          suggestions: [
            ...keywords.map((keyword) => ({
              kind: languages.CompletionItemKind.Keyword,
              label: keyword,
              insertText: keyword,
              range,
            })),
            ...tablesCompletions,
            ...columnsCompletions,
            ...tableColumnsCompletions,
            ...schemasCompletions,
            ...schemaTablesCompletions,
            ...schemaTableColumnsCompletions,
          ],
        };
      },
    });

    // eslint-disable-next-line @typescript-eslint/unbound-method
    dispose.value = disposable.dispose;
  }

  return {
    register,
  };
});

export const useMonaco = () => {
  function create(el: HTMLDivElement, value: Ref<string>) {
    const e = editor.create(el, {
      language: "sql",
      minimap: { enabled: false },
      lineNumbers: "off",
      lineDecorationsWidth: 0,
      folding: false,
      contextmenu: false,
      value: value.value,
      theme: "vs-dark",
      scrollBeyondLastLine: false,
      wordWrap: "on",
    });
    e.onEndUpdate(() => {
      value.value = e.getValue();
    });

    return e;
  }

  return { create };
};
