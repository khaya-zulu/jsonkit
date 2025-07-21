import { fetchPostRequest, MutationOptionsHelper } from ".";

import { useMutation } from "@tanstack/react-query";

type ToolMessage = {
  id: string;
  input: {
    input: Record<string, unknown> | Array<any>;
    query: string;
    description: string;
  };
  name: "process_json";
  result: {
    content: [{ text: string; type: "text" }];
    is_error: boolean;
    tool_use_id: string;
    type: "tool_use";
  };
};

export type ChatOutput = {
  content: string;
  role: "assistant";
  id: string;
  toolCalls: ToolMessage[];
};

export type Message = {
  id: string;
  content: string;
  role: "User" | "Assistant";
  jsonInput?: Record<string, unknown>;
  toolCalls?: ToolMessage[];
};

export type ChatInput = {
  content: string;
  jsonInput?: Record<string, unknown>;
  chatId: string;
  messages: Array<Message>;
};

export const useChatMutation = (
  options?: MutationOptionsHelper<ChatOutput, ChatInput>
) => {
  const path = "/chat";

  return useMutation({
    ...options,
    mutationKey: [path],
    mutationFn: async (input) => {
      const response = await fetchPostRequest<ChatOutput>(path, input);
      return response;
    },
  });
};
