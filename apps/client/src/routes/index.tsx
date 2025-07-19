import { ArrowUpIcon, WrenchIcon } from "@heroicons/react/16/solid";
import { createFileRoute } from "@tanstack/react-router";

import CodeMirror from "@uiw/react-codemirror";

import { EditorContent, useEditor } from "@tiptap/react";
import { Placeholder } from "@tiptap/extensions";

import { createTheme } from "@uiw/codemirror-themes";

import { TextStyleKit } from "@tiptap/extension-text-style";
import { StarterKit } from "@tiptap/starter-kit";
import { Code } from "@tiptap/extension-code";
import { CodeBlock } from "@tiptap/extension-code-block";

import { json } from "@codemirror/lang-json";

const theme = createTheme({
  theme: "dark",
  settings: {
    background: "#fff",
    lineHighlight: "#fff",
    gutterBackground: "#fff",
    gutterForeground: "#2e1065", // primary-100
    caret: "#fff",
  },
  styles: [],
});

export const Route = createFileRoute("/")({
  component: RouteComponent,
});

const extensions = [
  TextStyleKit,
  StarterKit,
  Placeholder.configure({
    placeholder: "Prompt your JSON...",
  }),
  Code,
  CodeBlock,
];

const JSONInput = () => {
  return (
    <CodeMirror
      value={`{"version": "9.99.99", "data": [1, 2, 3]}`}
      extensions={[json()]}
      theme={theme}
    />
  );
};

const InputBox = () => {
  const editor = useEditor({
    extensions,
  });

  return (
    <div className="rounded-xl ring-1 ring-zinc-100/50 bg-white shadow-sm min-h-40 p-4 relative flex flex-col justify-between">
      <EditorContent editor={editor} className="outline-none" />
      <div className="flex justify-between items-end mt-4">
        <div className="bg-violet-50/80 text-violet-900 px-2.5 py-1 font-semibold ring-1 ring-stone-100/40 rounded-md -translate-y-0.5">
          Claude Opus 4
        </div>

        <div className="flex items-center gap-4">
          <WrenchIcon className="w-5 h-5 opacity-40" />
          <div className="p-2 bg-violet-500 rounded-full text-white border-4 border-violet-50">
            <ArrowUpIcon className="w-5 h-5" />
          </div>
        </div>
      </div>
    </div>
  );
};

const Message = ({
  content,
  role,
  createdAt = "17:50 PM",
}: {
  content: string;
  role: "AI" | "User";
  createdAt?: string;
}) => {
  return (
    <div
      className={`rounded-xl p-5 flex flex-col gap-2 ${role === "User" ? "ring-1 ring-stone-200/60" : ""}`}
    >
      {role === "User" ? (
        <div className="font-semibold">{createdAt}</div>
      ) : null}
      <div className="text-zinc-700">{content}</div>
    </div>
  );
};

function RouteComponent() {
  return (
    <div className="flex p-4 h-screen">
      <div className="rounded-xl flex-1 bg-white ring-1 ring-zinc-200/[0.65] shadow-lg p-8 flex justify-between">
        <div className="text-violet-950">JSONkit</div>

        <div className="max-w-4xl flex flex-col justify-between w-full">
          <div className="flex-1 flex flex-col gap-2 px-8">
            <Message content="Hello, world!" role="User" />
            <Message
              content="Hello, User! How can I assist you today?"
              role="AI"
            />
            <Message content="Can you help me with this JSON?" role="User" />
            <Message content="Sure! Please provide the JSON data." role="AI" />
          </div>
          <div className="p-2 rounded-t-xl mb-2 ring-1 ring-zinc-100/50 text-zinc-600 shadow-sm mx-4 translate-y-2">
            <JSONInput />
          </div>
          <InputBox />
        </div>

        <div></div>
      </div>
    </div>
  );
}
