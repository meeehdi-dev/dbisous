import { editor, languages } from "monaco-editor";
import { Ref } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";

export const useMonaco = (value: Ref<string>, _columns: Array<unknown>) => {
  const { tableId } = useUrlParams();

  languages.registerCompletionItemProvider("sql", {
    provideCompletionItems: function (model, position) {
      const word = model.getWordUntilPosition(position);
      const range = {
        startLineNumber: position.lineNumber,
        endLineNumber: position.lineNumber,
        startColumn: word.startColumn,
        endColumn: word.endColumn,
      };

      return {
        suggestions: [
          // ...columns.map((c) => ({
          //   kind: languages.CompletionItemKind.Field,
          //   label: c.name,
          //   insertText: `${c.name}`,
          //   range,
          // })),
          {
            kind: languages.CompletionItemKind.File,
            label: tableId.value,
            insertText: tableId.value,
            range,
          },
          {
            kind: languages.CompletionItemKind.Keyword,
            label: "SELECT",
            insertText: "SELECT",
            range,
          },
          {
            kind: languages.CompletionItemKind.Keyword,
            label: "FROM",
            insertText: "FROM",
            range,
          },
        ],
      };
    },
  });

  function init(el: HTMLDivElement) {
    const e = editor.create(el, {
      language: "sql",
      minimap: { enabled: false },
      lineNumbers: "off",
      lineDecorationsWidth: 0,
      contextmenu: false,
      value: `SELECT * FROM ${tableId.value};`,
    });
    e.onEndUpdate(() => {
      value.value = e.getValue();
    });
  }

  return { init, value };
};
