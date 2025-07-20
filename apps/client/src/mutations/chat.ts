import { fetchPostRequest, MutationOptionsHelper } from ".";

import { useMutation } from "@tanstack/react-query";

export type ChatOutput = {
  content: string;
  role: "assistant";
  id: string;
};

export type Message = {
  id: string;
  content: string;
  role: "User" | "Assistant";
  jsonInput?: Record<string, unknown>;
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
