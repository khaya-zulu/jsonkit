import { useState } from "react";

import {
  AdjustmentsHorizontalIcon,
  ArrowUpIcon,
} from "@heroicons/react/16/solid";
import {
  ClipboardIcon,
  CodeBracketSquareIcon,
  FolderIcon,
  SparklesIcon,
} from "@heroicons/react/24/outline";
import { createFileRoute } from "@tanstack/react-router";

import CodeMirror from "@uiw/react-codemirror";
import Turndown from "turndown";

import { EditorContent, useEditor } from "@tiptap/react";
import { Placeholder } from "@tiptap/extensions";
import Markdown from "react-markdown";

import { createTheme } from "@uiw/codemirror-themes";
import { formatDate } from "date-fns";

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

const JsonInput = ({
  value,
  isFormatted,
}: {
  value: string;
  isFormatted?: boolean;
}) => {
  const [jsonValue, setJsonValue] = useState(() => {
    try {
      return isFormatted ? JSON.stringify(JSON.parse(value), null, 2) : value;
    } catch (error) {
      return value;
    }
  });

  return <CodeMirror value={jsonValue} extensions={[json()]} theme={theme} />;
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
      <EditorContent
        editor={editor}
        className="outline-none"
        onKeyDown={(ev) => {
          if (ev.key === "Enter" && !ev.shiftKey) {
            ev.preventDefault();
            handleSubmit(ev);
          }
        }}
      />
      <div className="flex justify-between items-end mt-4">
        <div className="bg-violet-50/80 text-violet-900 px-2.5 py-1 font-semibold ring-1 ring-stone-100/40 rounded-md -translate-y-0.5">
          Claude Opus 4
        </div>

        <button type="submit" className="flex items-center gap-4">
          <AdjustmentsHorizontalIcon className="size-5" />
          <div className="p-2 bg-violet-700 rounded-full text-white border-4 border-violet-50">
            <ArrowUpIcon className="size-5" />
          </div>
        </button>
      </div>
    </form>
  );
};

const ChatMessage = ({
  message,
  createdAt,
}: {
  message: Message;
  createdAt: Date;
}) => {
  return (
    <div
      className={`rounded-xl flex flex-col ${message.role === "User" ? "bg-stone-100 shadow-2xs" : ""}`}
    >
      {message.role === "User" ? (
        <div className="font-semibold px-4 pt-2 text-slate-950 flex justify-between">
          User <div>• {formatDate(createdAt, "hh:mm a")}</div>
        </div>
      ) : null}

      <div className="p-1.5">
        <div
          className={
            message.role === "User" ? "bg-white rounded-lg shadow-xs p-5" : ""
          }
        >
          <Markdown>{message.content}</Markdown>
          {message.toolCalls?.map((t) => {
            return (
              <div className="rounded-lg mt-4 ring-1 ring-violet-950/10 overflow-hidden shadow-xs">
                <div className="py-2 px-4 flex justify-between bg-violet-50/50 text-[#1b0746]">
                  <div className="truncate">{t.input.query}</div>
                  <div className="text-violet-950 truncate">
                    {t.input.description}
                  </div>
                </div>
                <div className="p-4 text-[#1b0746]">
                  {t.result.content[0].text}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

function RouteComponent() {
  const [messages, setMessages] = useState<Array<Message>>([]);
  const [jsonInputVal, setJsonInputVal] = useState<string>(() =>
    JSON.stringify({
      logs: [
        {
          context: {
            email: "john.doe@example.com",
            ip: "192.168.1.105",
            userAgent:
              "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
            userId: "usr_123456",
          },
          duration: 145,
          host: "api-server-01",
          level: "INFO",
          message: "User authentication successful",
          requestId: "req_789xyz",
          service: "user-api",
          timestamp: "2024-01-15T10:23:45.123Z",
        },
        {
          context: {
            amount: 129.99,
            currency: "USD",
            orderId: "ord_abc123",
            paymentMethod: "credit_card",
            userId: "usr_789012",
          },
          error: {
            code: "CARD_DECLINED",
            message: "Card declined by issuer",
            stackTrace:
              "at PaymentProcessor.charge (payment.js:123)\n at OrderService.complete (order.js:456)",
            type: "PaymentGatewayError",
          },
          host: "api-server-02",
          level: "ERROR",
          message: "Payment processing failed",
          requestId: "req_456def",
          service: "payment-service",
          timestamp: "2024-01-15T10:23:46.234Z",
        },
        {
          context: {
            currentStock: 5,
            productId: "prod_987654",
            productName: "Wireless Headphones",
            threshold: 10,
            warehouseId: "wh_east_01",
          },
          host: "api-server-01",
          level: "WARN",
          message: "Low stock warning",
          requestId: "req_111aaa",
          service: "inventory-service",
          timestamp: "2024-01-15T10:23:47.345Z",
        },
        {
          context: {
            fallback: "database",
            key: "user:preferences:usr_123456",
            latency: 0.5,
            operation: "GET",
          },
          host: "cache-server-01",
          level: "DEBUG",
          message: "Cache miss for key",
          requestId: "req_222bbb",
          service: "cache-service",
          timestamp: "2024-01-15T10:23:48.456Z",
        },
        {
          context: {
            apiVersion: "v1",
            rateLimitRemaining: 95,
            userId: "usr_345678",
          },
          host: "gateway-01",
          http: {
            contentLength: 1456,
            method: "POST",
            path: "/api/v1/orders",
            responseTime: 234,
            statusCode: 201,
          },
          level: "INFO",
          message: "API request completed",
          requestId: "req_333ccc",
          service: "api-gateway",
          timestamp: "2024-01-15T10:23:49.567Z",
        },
      ],
      metadata: {
        environment: "production",
        region: "us-east-1",
        version: "2.3.1",
      },
    })
  );

  const [chatId] = useState<string>(`${Date.now()}`);

  const chatMutation = useChatMutation({});

  const handleSubmit = async (value: string) => {
    try {
      const jsonInput = JSON.parse(jsonInputVal);
      setMessages((prev) => [
        ...prev,
        {
          content: value,
          role: "User",
          id: Date.now().toString(),
          jsonInput,
        },
      ]);

      const chatInput = await chatMutation.mutateAsync({
        content: value,
        jsonInput,
        chatId,
        messages,
      });

      setMessages((prev) => [
        ...prev,
        {
          content: chatInput.content,
          role: "Assistant",
          id: chatInput.id,
          jsonInput: {},
          toolCalls: chatInput.toolCalls,
        },
      ]);
    } catch (error) {
      console.error("Error submitting chat message:", error);
      return;
    }
  };

  return (
    <div className="flex p-4 h-screen bg-gradient-to-b from-stone-50/40 to-violet-50/80 gap-4">
      <div className="rounded-xl flex-1 ring-1 ring-stone-200/40 shadow-lg flex flex-col bg-white">
        <div className="border-b border-stone-100 flex justify-between items-center">
          <div className="px-8 py-3 border-r border-stone-100 flex-1">
            <input
              className="outline-none placeholder:text-stone-400"
              placeholder="jqlang query"
            />
          </div>
          <div className="pr-8 pl-4 py-3 flex gap-4 items-center">
            <SparklesIcon className="size-5" />
            <CodeBracketSquareIcon className="size-6" />
            <FolderIcon className="size-5" />
            <ClipboardIcon className="size-5" />
          </div>
        </div>
        <div className="h-12 border-b border-stone-100"></div>
        <div className="flex-1 overflow-auto px-4 py-2 hide-scrollbar overscroll-contain">
          <JsonInput value={jsonInputVal} isFormatted />
        </div>
        <div className="px-8 py-3.5 border-t border-stone-100 text-stone-500">
          Temporary storage: tmp/test.json
        </div>
      </div>

      <div className="rounded-xl w-[35rem] shadow-lg ring-1 ring-stone-200/40 p-4 bg-white overflow-hidden flex flex-col">
        <div className="flex-1"></div>
        <InputBox onSubmit={handleSubmit} />
      </div>
    </div>
  );
}
