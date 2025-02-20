import { editor, languages } from "monaco-editor";
import { Ref } from "vue";
import { getSQLiteKeywords } from "@/utils/keywords";

const keywords = getSQLiteKeywords();
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
        // {
        //   kind: languages.CompletionItemKind.File,
        //   label: tableId.value,
        //   insertText: tableId.value,
        //   range,
        // },
        ...keywords.map((keyword) => ({
          kind: languages.CompletionItemKind.Keyword,
          label: keyword,
          insertText: keyword,
          range,
        })),
      ],
    };
  },
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
