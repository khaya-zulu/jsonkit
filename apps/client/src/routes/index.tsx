import { useState } from "react";

import { ArrowUpIcon, WrenchIcon } from "@heroicons/react/16/solid";
import { createFileRoute } from "@tanstack/react-router";

import CodeMirror from "@uiw/react-codemirror";
import Turndown from "turndown";

import { EditorContent, useEditor } from "@tiptap/react";
import { Placeholder } from "@tiptap/extensions";

import { createTheme } from "@uiw/codemirror-themes";

import { TextStyleKit } from "@tiptap/extension-text-style";
import { StarterKit } from "@tiptap/starter-kit";
import { Code } from "@tiptap/extension-code";
import { CodeBlock } from "@tiptap/extension-code-block";

import { json } from "@codemirror/lang-json";
import { type Message, useChatMutation } from "../mutations/chat";

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

const InputBox = ({ onSubmit }: { onSubmit: (value: string) => void }) => {
  const editor = useEditor({
    extensions,
  });

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (editor) {
      const content = editor.getHTML();
      const markdown = new Turndown().turndown(content);

      onSubmit(markdown);
      editor.commands.clearContent();
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="rounded-xl ring-1 ring-zinc-100/50 bg-white shadow-sm min-h-40 p-4 relative flex flex-col justify-between"
    >
      <EditorContent editor={editor} className="outline-none" />
      <div className="flex justify-between items-end mt-4">
        <div className="bg-violet-50/80 text-violet-900 px-2.5 py-1 font-semibold ring-1 ring-stone-100/40 rounded-md -translate-y-0.5">
          Claude Opus 4
        </div>

        <button type="submit" className="flex items-center gap-4">
          <WrenchIcon className="w-5 h-5 opacity-40" />
          <div className="p-2 bg-violet-500 rounded-full text-white border-4 border-violet-50">
            <ArrowUpIcon className="w-5 h-5" />
          </div>
        </button>
      </div>
    </form>
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
  const [messages, setMessages] = useState<Array<Message>>([]);

  const [chatId] = useState<string>(`${Date.now()}`);

  const chatMutation = useChatMutation({});

  const handleSubmit = async (value: string) => {
    setMessages((prev) => [
      ...prev,
      {
        content: value,
        role: "user",
        id: Date.now().toString(),
        jsonInput: {},
      },
    ]);

    const chatInput = await chatMutation.mutateAsync({
      content: value,
      jsonInput: { version: "9.99.99", data: [1, 2, 3] }, // Example JSON input
      chatId,
      messages,
    });

    setMessages((prev) => [
      ...prev,
      {
        content: chatInput.content,
        role: "assistant",
        id: chatInput.id,
        jsonInput: {},
      },
    ]);
  };

  return (
    <div className="flex p-4 h-screen">
      <div className="rounded-xl flex-1 bg-white ring-1 ring-zinc-200/[0.65] shadow-lg p-8 flex justify-between">
        <div className="text-violet-950">JSONkit</div>

        <div className="max-w-4xl flex flex-col justify-between w-full">
          <div className="flex-1 flex flex-col gap-2 px-8 overflow-y-auto">
            {messages.map((message, index) => (
              <Message
                key={index}
                content={message.content}
                role={message.role === "user" ? "User" : "AI"}
                createdAt={new Date().toString()}
              />
            ))}
          </div>
          <div className="p-2 rounded-t-xl mb-2 ring-1 ring-zinc-100/50 text-zinc-600 shadow-sm mx-4 translate-y-2">
            <JSONInput />
          </div>
          <InputBox onSubmit={handleSubmit} />
        </div>

        <div></div>
      </div>
    </div>
  );
}
